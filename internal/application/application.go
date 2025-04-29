package application

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/braginantonev/gcalc-server/pkg/agent"
	"github.com/braginantonev/gcalc-server/pkg/calc"
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
	mux := http.NewServeMux()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mux.HandleFunc("/api/v1/calculate", RequestEmpty(AddExpressionHandler))
	mux.HandleFunc("/api/v1/expressions", GetExpressionsQueueHandler)
	mux.HandleFunc("/api/v1/expressions/", GetExpressionHandler)
	mux.HandleFunc("/internal/task", ResultOrGet(ResultHandler, GetTaskHandler))

	agent.Enable(ctx)

	slog.Info("Start server", slog.String("port", app.cfg.Port))
	err := http.ListenAndServe(":"+app.cfg.Port, mux)
	if err != nil {
		slog.Error("Failed to start server")
		return err
	}

	return nil
}

// * ------------------- HTTP Server --------------------

type RequestExpression struct {
	Expression string `json:"expression"`
}

type RequestTask struct {
	Id     string  `json:"id"`
	Result float64 `json:"result"`
}

type ResponseExpression struct {
	Id     string      `json:"id"`
	Status calc.Status `json:"status,omitempty"`
	Result float64     `json:"result,omitempty"`
}

var expression []byte
