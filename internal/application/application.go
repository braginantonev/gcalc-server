package application

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/braginantonev/gcalc-server/pkg/database"

	"github.com/braginantonev/gcalc-server/pkg/agent"
	"github.com/braginantonev/gcalc-server/pkg/logreg"
	"github.com/braginantonev/gcalc-server/pkg/orchestrator"
	lr_pb "github.com/braginantonev/gcalc-server/proto/logreg"
	orch_pb "github.com/braginantonev/gcalc-server/proto/orchestrator"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// * -------------------- Config --------------------

type Config struct {
	Port               string
	GRPCPort           string
	GRPCServerAddress  string
	JWTSecretSignature string
	ComputingPower     int
}

func NewConfig() *Config {
	cfg := new(Config)
	var loaded bool

	cfg.Port, loaded = os.LookupEnv("PORT")
	if !loaded {
		cfg.Port = "8080"
		slog.Warn("env: \"PORT\" not found. Set default - 8080")
	}
	slog.Info("Server has been configured successfully")

	cfg.GRPCPort, loaded = os.LookupEnv("GRPCPort")
	if !loaded {
		cfg.GRPCPort = "5000"
		slog.Warn("env: \"GRPCPort\" not found. Set default - 5000")
	}

	cfg.GRPCServerAddress = "localhost:" + cfg.GRPCPort
	slog.Info("Orchestrator will be started", slog.String("address", cfg.GRPCServerAddress))

	var err error
	got_compower, loaded := os.LookupEnv("COMPUTING_POWER")
	if !loaded {
		cfg.ComputingPower = 1
		slog.Warn("env: \"COMPUTING_POWER\" not found. Set default - 1")
	}

	cfg.ComputingPower, err = strconv.Atoi(got_compower)
	if err != nil {
		cfg.ComputingPower = 1
		slog.Warn("env: \"COMPUTING_POWER\" not integer")
	}
	slog.Info("Set", slog.String("COMPUTING_POWER", fmt.Sprint(cfg.ComputingPower)))

	cfg.JWTSecretSignature, loaded = os.LookupEnv("JWTSecretSignature")
	if !loaded {
		panic(`
		!!! Attention !!!
		JWT signature in env JWTSecretSignature not found.
		Please go to README.md - Usage, and follow the instruction!
		
		!!! Внимание !!!
		JWT сигнатура в переменной окружения JWTSecretSignature не найдена.
		Пожалуйста перейдите в README.md - Usage и проследуйте инструкции!
		`)
	}

	JWTSignature = cfg.JWTSecretSignature
	return cfg
}

// * ----------------- Services clients -----------------

var (
	GRPCConnectionClient      *grpc.ClientConn
	OrchestratorServiceClient orch_pb.OrchestratorServiceClient
	LogRegClient              lr_pb.LogRegServiceClient
	JWTSignature              string
)

// * ------------------- Application --------------------

type Application struct {
	cfg *Config
}

func NewApplication() *Application {
	return &Application{cfg: NewConfig()}
}

func (app Application) Run(grpc_server *grpc.Server) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	//Todo: Добавить path в env app
	expressions_db, err := database.NewDataBase(ctx, "expressions.db")
	if err != nil {
		return err
	}

	users_db, err := database.NewDataBase(ctx, "users.db")
	if err != nil {
		return err
	}

	if err := logreg.RegisterServer(ctx, grpc_server, users_db, app.cfg.JWTSecretSignature); err != nil {
		return err
	}

	if err := orchestrator.RegisterServer(ctx, grpc_server, expressions_db); err != nil {
		return err
	}

	//* Start gRPC server
	lis, err := net.Listen("tcp", app.cfg.GRPCServerAddress)
	if err != nil {
		slog.Info(app.cfg.GRPCServerAddress)
		panic("Failed to enable orchestrator process. " + err.Error())
	}

	slog.Info("orchestrator tcp listener started", slog.String("addr", app.cfg.GRPCServerAddress))

	go func() {
		slog.Info("Start grpc server", slog.String("port", app.cfg.GRPCPort))
		if err := grpc_server.Serve(lis); err != nil {
			panic("error serving grpc: " + err.Error())
		}
	}()

	<-time.After(time.Second)

	GRPCConnectionClient, err = grpc.NewClient(app.cfg.GRPCServerAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}

	LogRegClient = lr_pb.NewLogRegServiceClient(GRPCConnectionClient)
	OrchestratorServiceClient = orch_pb.NewOrchestratorServiceClient(GRPCConnectionClient)
	agent.Enable(ctx, OrchestratorServiceClient, app.cfg.ComputingPower)

	//* Start REST API main server
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/calculate", RequestEmpty(AddExpressionHandler))
	mux.HandleFunc("/api/v1/expressions", GetExpressionsQueueHandler)
	mux.HandleFunc("/api/v1/expressions/", GetExpressionHandler)
	mux.HandleFunc("/api/v1/register", RegisterHandler)
	mux.HandleFunc("/api/v1/login", LoginHandler)

	slog.Info("Start server", slog.String("port", app.cfg.Port))
	err = http.ListenAndServe(":"+app.cfg.Port, mux)
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
	Id     int32   `json:"id"`
	Status string  `json:"status,omitempty"`
	Result float64 `json:"result,omitempty"`
}

type User struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}
