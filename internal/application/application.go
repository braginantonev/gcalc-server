package application

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"os"

	"github.com/Antibrag/gcalc-server/pkg/calc"
)

var (
	InternalError    error    = errors.New("Internal error")
	RequestBodyEmpty error    = errors.New("Request body empty")
	CalculatorErrors []*error = []*error{
		&calc.DivideByZero,
		&calc.ExpressionEmpty,
		&calc.OperationWithoutValue,
		&calc.BracketsNotFound,
		&calc.ParseError,
	}
)

// * -------------------- Config --------------------
type Config struct {
	Port string
}

func NewConfig() *Config {
	cfg := new(Config)
	cfg.Port = os.Getenv("PORT")
	if cfg.Port == "" {
		cfg.Port = "8080"
		slog.Warn("env: \"PORT\" not found. ")
	}
	slog.Info("Server has been configured successfully")
	return cfg
}

// * ------------------- Application --------------------
type Application struct {
	cfg *Config
}

func NewApplication() *Application {
	return &Application{cfg: NewConfig()}
}

func (app Application) Run() error {
	slog.Info("Start server", slog.String("port", app.cfg.Port))

	http.HandleFunc("/api/v1/calculate", RequestEmpty(CalcHandler))
	err := http.ListenAndServe(":"+app.cfg.Port, nil)
	if err != nil {
		slog.Error("Failed to start server")
		return err
	}

	return nil
}

// * ------------------- HTTP Server --------------------
type Request struct {
	Expression string `json:"expression"`
}

type Response struct {
	Result float64 `json:"result,omitempty"`
	Error  string  `json:"error,omitempty"`
}

var expression []byte

func CalcHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	resq := Request{}
	err := json.Unmarshal(expression, &resq)
	if err != nil {
		slog.Error("Failed unmarshal expression json.", slog.String("error", err.Error()))
	}

	result, err := calc.Calc(resq.Expression)
	if err != nil {
		slog.Error("Failed calculate expression.", slog.String("error", err.Error()))

		isInternalError := true
		for _, i := range CalculatorErrors {
			if errors.Is(err, *i) {
				isInternalError = false
				break
			}
		}

		var resp_json []byte
		if isInternalError {
			w.WriteHeader(500)
			resp_json, _ = json.Marshal(Response{Error: InternalError.Error()})
		} else {
			w.WriteHeader(422)
			resp_json, _ = json.Marshal(Response{Error: err.Error()})
		}
		w.Write(resp_json)
		return
	}

	slog.Info("Calculation success")
	resp_json, _ := json.Marshal(Response{Result: result})
	w.WriteHeader(200)
	w.Write(resp_json)
}

func RequestEmpty(fn http.HandlerFunc) http.HandlerFunc {
	log_failed_conv := func(resp_json string, w http.ResponseWriter) {
		w.WriteHeader(500)
		slog.Error("Failed convert error response to json")
		slog.Debug("expression:", string(expression), "\nresponse json:", (resp_json))
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error
		expression, err = io.ReadAll(r.Body)

		if err != nil {
			slog.Error("Failed read request body")

			resp_json, err := json.Marshal(Response{Error: err.Error()})
			if err != nil {
				log_failed_conv(string(resp_json), w)
				return
			}

			w.WriteHeader(500)
			w.Write(resp_json)
			return
		}

		if len(expression) == 0 {
			slog.Error("Request body empty")

			resp_json, err := json.Marshal(Response{Error: RequestBodyEmpty.Error()})
			if err != nil {
				log_failed_conv(string(resp_json), w)
				return
			}

			w.WriteHeader(422)
			w.Write(resp_json)
			return
		}

		fn.ServeHTTP(w, r)
	})
}
