package application

import (
	"log/slog"
	"net/http"
	"os"
	"time"

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
	//Todo: Добавить мьюксы для хендлеров

	mux := http.NewServeMux()

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
	Error  string      `json:"error,omitempty"`
}

type ResponseExpressionsQueue struct {
	Expressions []ResponseExpression `json:"expressions"`
}

type ResponseTask struct {
	Id             string        `json:"id"`
	FirstArgument  calc.Argument `json:"arg1"`
	SecondArgument calc.Argument `json:"arg2"`
	Operation      calc.Operator `json:"operation"`
	OperationTime  time.Time     `json:"operation_time"`
}

var expression []byte
