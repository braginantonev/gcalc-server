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

	"github.com/braginantonev/gcalc-server/pkg/agent"
	"github.com/braginantonev/gcalc-server/pkg/orchestrator"
	orch_pb "github.com/braginantonev/gcalc-server/proto/orchestrator"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// * -------------------- Config --------------------

type Config struct {
	Port              string
	GRPCPort          string
	GRPCServerAddress string
	ComputingPower    int
}

func NewConfig() *Config {
	cfg := new(Config)
	cfg.Port = os.Getenv("PORT")
	if cfg.Port == "" {
		cfg.Port = "8080"
		slog.Warn("env: \"PORT\" not found. Set default - 8080")
	}
	slog.Info("Server has been configured successfully")

	cfg.GRPCPort = os.Getenv("GRPCPort")
	if cfg.GRPCPort == "" {
		cfg.GRPCPort = "5000"
		slog.Warn("env: \"GRPCPort\" not found. Set default - 5000")
	}

	cfg.GRPCServerAddress = "localhost:" + cfg.GRPCPort
	slog.Info("Orchestrator will be started", slog.String("address", cfg.GRPCServerAddress))

	var err error
	cfg.ComputingPower, err = strconv.Atoi(os.Getenv("COMPUTING_POWER"))
	if err != nil {
		cfg.ComputingPower = 1
		slog.Warn("env: \"COMPUTING_POWER\" not found or not integer")
	}
	slog.Info("Set", slog.String("Computing power", fmt.Sprint(cfg.ComputingPower)))
	return cfg
}

// * ----------------- Services clients -----------------

var (
	GRPCConnectionClient      *grpc.ClientConn
	OrchestratorServiceClient orch_pb.OrchestratorServiceClient
	//LogRegClient lr_pb.LogRegServiceClient
)

// * ------------------- Application --------------------

type Application struct {
	cfg *Config
}

func NewApplication() *Application {
	return &Application{cfg: NewConfig()}
}

func (app Application) Run(grpc_server *grpc.Server) error {
	mux := http.NewServeMux()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mux.HandleFunc("/api/v1/calculate", RequestEmpty(AddExpressionHandler))
	mux.HandleFunc("/api/v1/expressions", GetExpressionsQueueHandler)
	mux.HandleFunc("/api/v1/expressions/", GetExpressionHandler)

	app.EnableOrchestratorService(grpc_server)

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

	// Create orchestrator client
	OrchestratorServiceClient = orch_pb.NewOrchestratorServiceClient(GRPCConnectionClient)

	agent.Enable(ctx, OrchestratorServiceClient, app.cfg.ComputingPower)

	slog.Info("Start server", slog.String("port", app.cfg.Port))
	err = http.ListenAndServe(":"+app.cfg.Port, mux)
	if err != nil {
		slog.Error("Failed to start server")
		return err
	}

	return nil
}

func (app Application) EnableOrchestratorService(grpc_server *grpc.Server) {
	orchestratorServiceServer := orchestrator.NewServer()
	orch_pb.RegisterOrchestratorServiceServer(grpc_server, orchestratorServiceServer)
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
