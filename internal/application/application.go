package application

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"os"

	"github.com/Antibrag/gcalc-server/pkg/calc"
)

var (
	InternalError    error = errors.New("Internal error")
	RequsetBodyEmpty error = errors.New("Request body empty")
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
		slog.Info("env: \"PORT\" not found. ")
	}
	slog.Info("Configure port to ", cfg.Port)
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
	slog.Info("Start server...")

	http.HandleFunc("/api/v1/calculate", RequestEmpty(CalcHandler))
	err := http.ListenAndServe(":"+app.cfg.Port, nil)
	if err != nil {
		slog.Error("Failed to start server")
		return err
	}

	slog.Info("Server started with port", app.cfg.Port)
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

	err := json.Unmarshal()

	result, err := calc.Calc()
}

func RequestEmpty(fn http.HandlerFunc) http.HandlerFunc {
	log_failed_conv := func(resp_json string, w http.ResponseWriter) {
		w.WriteHeader(500)
		slog.Error("Failed convert error response to json")
		slog.Debug("expression:", string(expression), "\nresponse json:", string(resp_json))
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n, err := r.Body.Read(expression)
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

		if n == 0 {
			slog.Info("Request body is empty")
			resp_json, err := json.Marshal(Response{Error: RequsetBodyEmpty.Error()})
			if err != nil {
				log_failed_conv(string(resp_json), w)
				return
			}

			w.WriteHeader(400)
			w.Write(resp_json)
			return
		}

		fn.ServeHTTP(w, r)
	})
}
