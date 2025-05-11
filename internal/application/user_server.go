package application

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"slices"
	"strconv"
	"strings"

	lr_pb "github.com/braginantonev/gcalc-server/proto/logreg"
	orch_pb "github.com/braginantonev/gcalc-server/proto/orchestrator"
	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

var expression []byte

func getUsername(r *http.Request) (string, bool) {
	token := r.Header.Get("Authorization")
	if token == "" {
		return "", true
	}

	if token[:6] == "Bearer" {
		token = token[7:]
	}

	tokenFromString, err := jwt.Parse(token, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(JWTSignature), nil
	})

	if err != nil {
		fmt.Println(err)
		return "", false
	}

	if claims, ok := tokenFromString.Claims.(jwt.MapClaims); ok {
		return claims["name"].(string), true
	} else {
		return "", false
	}
}

func logFailedConvert(handler_name, resp_json string, w *http.ResponseWriter) {
	(*w).WriteHeader(http.StatusInternalServerError)
	slog.Error("Failed convert response to JSON", slog.String("handler_name", handler_name))
	slog.Debug("expression:", string(expression), "\nresponse json:", resp_json)
}

func RequestEmpty(fn http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Info("request - new expression")

		var err error
		expression, err = io.ReadAll(r.Body)

		if err != nil {
			slog.Error("Failed read request body")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if len(expression) == 0 {
			slog.Error("Request body empty")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		fn.ServeHTTP(w, r)
	})
}

func AddExpressionHandler(w http.ResponseWriter, r *http.Request) {
	req := RequestExpression{}
	err := json.Unmarshal(expression, &req)
	if err != nil {
		slog.Error("Failed unmarshal expression json.", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}

	//Todo: Создать расшифровку jwt токена, для получения имени пользователя
	username, ok := getUsername(r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(ErrJWTTokenNotValid.Error()))
		return
	}

	id, err := OrchestratorServiceClient.AddExpression(context.TODO(), &orch_pb.AddedExpression{User: username, Str: req.Expression})
	slog.Info("add expression to queue. ", slog.String("id", fmt.Sprint(id.GetValue())))

	if err != nil {
		slog.Error("Failed add expression. ", slog.String("error", err.Error()))

		if slices.Contains(OrchestratorErrors, &err) {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusUnprocessableEntity)
		}
		return
	}

	resp := ResponseExpression{
		Id: id.GetValue(),
	}

	resp_json, err := json.Marshal(resp)
	if err != nil {
		logFailedConvert("AddExpressionHandler()", string(resp_json), &w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(resp_json)
}

func GetExpressionsQueueHandler(w http.ResponseWriter, r *http.Request) {
	slog.Info("request - get expressions queue")

	username, ok := getUsername(r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(ErrJWTTokenNotValid.Error()))
		return
	}

	expressions, err := OrchestratorServiceClient.GetExpressions(context.TODO(), wrapperspb.String(username))
	if err != nil {
		slog.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var resp struct {
		Expressions []ResponseExpression `json:"expressions"`
	}

	for _, exp := range expressions.GetQueue() {
		resp.Expressions = append(resp.Expressions, ResponseExpression{
			Id:     exp.Id.Internal,
			Status: exp.Status.String(),
			Result: exp.Result,
		})
	}

	resp_json, err := json.Marshal(resp)
	if err != nil {
		logFailedConvert("GetExpressionsQueue()", string(resp_json), &w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp_json)
}

func GetExpressionHandler(w http.ResponseWriter, r *http.Request) {
	paths := strings.Split(r.URL.Path, "/")
	id, err := strconv.ParseInt(paths[len(paths)-1], 10, 32)
	if err != nil {
		slog.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	slog.Info("request - get expression.", slog.String("id", fmt.Sprint(id)))

	username, ok := getUsername(r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(ErrJWTTokenNotValid.Error()))
		return
	}

	exp, err := OrchestratorServiceClient.GetExpression(context.TODO(), orch_pb.NewExpressionIDWithValues(username, int32(id)))
	if err != nil {
		slog.Error("expression not found", slog.String("id", fmt.Sprint(id)), slog.String("err", err.Error()))
		w.WriteHeader(http.StatusNotFound)
		return
	}

	resp := struct {
		Expression ResponseExpression `json:"expression"`
	}{
		Expression: ResponseExpression{
			Id:     exp.Id.Internal,
			Status: exp.Status.String(),
			Result: exp.Result,
		},
	}

	resp_json, err := json.Marshal(resp)
	if err != nil {
		logFailedConvert("getExpressionHandler()", string(resp_json), &w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp_json)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	slog.Info("Request - Login")

	readed, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	}

	if len(readed) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user := lr_pb.User{}
	if err := json.Unmarshal(readed, &user); err != nil {
		slog.Error(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	jwt, err := LogRegClient.Login(context.TODO(), &user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if st, ok := status.FromError(err); ok {
			w.Header().Set("Content-Type", "text/plain")
			w.Write([]byte(st.Message()))
		}
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(jwt.Token))
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	slog.Info("Request - Register")

	readed, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	}

	if len(readed) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user := lr_pb.User{}
	if err := json.Unmarshal(readed, &user); err != nil {
		slog.Error(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	jwt, err := LogRegClient.Register(context.TODO(), &user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if st, ok := status.FromError(err); ok {
			w.Header().Set("Content-Type", "text/plain")
			w.Write([]byte(st.Message()))
		}
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(jwt.Token))
}
