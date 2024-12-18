package application

import (
	"log/slog"
	"net/http"
	"os"
)

var InternalError error

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

}

func RequestEmpty(fn http.HandlerFunc) http.HandlerFunc {

}
