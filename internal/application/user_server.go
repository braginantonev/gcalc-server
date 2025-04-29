package application

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"slices"
	"strings"

	"github.com/braginantonev/gcalc-server/pkg/orchestrator"
)

func logFailedConvert(handler_name, resp_json string, w *http.ResponseWriter) {
	(*w).WriteHeader(http.StatusInternalServerError)
	slog.Error("Failed convert response to JSON", slog.String("handler_name", handler_name))
	slog.Debug("expression:", string(expression), "\nresponse json:", resp_json)
}

func AddExpressionHandler(w http.ResponseWriter, r *http.Request) {
	req := RequestExpression{}
	err := json.Unmarshal(expression, &req)
	if err != nil {
		slog.Error("Failed unmarshal expression json.", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}

	id, err := orchestrator.AddExpression(req.Expression)
	slog.Info("add expression to queue. ", slog.String("id", id))

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
		Id: id,
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

func GetExpressionsQueueHandler(w http.ResponseWriter, r *http.Request) {
	slog.Info("request - get expressions queue")

	expressions := orchestrator.GetExpressionsQueue()

	var resp struct {
		Expressions []ResponseExpression `json:"expressions"`
	}

	for _, exp := range expressions {
		resp.Expressions = append(resp.Expressions, ResponseExpression{
			Id:     exp.Id,
			Status: exp.Status,
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
	id := paths[len(paths)-1]

	slog.Info("request - get expression.", slog.String("id", id))

	exp, err := orchestrator.GetExpression(id)
	if err != nil {
		slog.Error("expression not found", slog.String("id", id))
		w.WriteHeader(http.StatusNotFound)
		return
	}

	resp := struct {
		Expression ResponseExpression `json:"expression"`
	}{
		Expression: ResponseExpression{
			Id:     exp.Id,
			Status: exp.Status,
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
