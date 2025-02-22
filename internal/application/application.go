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
		resp, _ := json.Marshal(Response{Error: ErrUnsupportedBodyType.Error()})

		w.WriteHeader(http.StatusUnsupportedMediaType)
		w.Write(resp)
		return
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
			w.WriteHeader(http.StatusInternalServerError)
			resp_json, _ = json.Marshal(Response{Error: ErrInternalError.Error()})
		} else {
			w.WriteHeader(http.StatusUnprocessableEntity)
			resp_json, _ = json.Marshal(Response{Error: err.Error()})
		}
		w.Write(resp_json)
		return
	}

	slog.Info("Calculation success")

	resp := Response{Result: result}
	var resp_json []byte
	if resp.Result != 0 {
		resp_json, _ = json.Marshal(Response{Result: result})
	} else {
		resp_json = []byte("{\"result\":0}")
	}

	w.WriteHeader(http.StatusOK)
	w.Write(resp_json)
}

func RequestEmpty(fn http.HandlerFunc) http.HandlerFunc {
	log_failed_conv := func(resp_json string, w http.ResponseWriter) {
		w.WriteHeader(http.StatusInternalServerError)
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

			w.WriteHeader(http.StatusInternalServerError)
			w.Write(resp_json)
			return
		}

		if len(expression) == 0 {
			slog.Error("Request body empty")

			resp_json, err := json.Marshal(Response{Error: ErrRequestBodyEmpty.Error()})
			if err != nil {
				log_failed_conv(string(resp_json), w)
				return
			}

			w.WriteHeader(http.StatusBadRequest)
			w.Write(resp_json)
			return
		}

		fn.ServeHTTP(w, r)
	})
}
