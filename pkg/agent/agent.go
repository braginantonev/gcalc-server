package agent

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"slices"
	"sync"
	"time"

	"github.com/braginantonev/gcalc-server/pkg/orchestrator"
	pb "github.com/braginantonev/gcalc-server/proto/orchestrator"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

//! Оставить как минимум один поток для выполнения анализа

const (
	COMPUTING_POWER   = 5
	TASK_WAIT_TIME_MS = 250
)

var (
	tasks           []string
	main_server_url string
	conn            *grpc.ClientConn
	orchClient      pb.OrchestratorServiceClient
)

// Return true, if task has been appended, else - false
func appendTask(task_id string) bool {
	if slices.Contains(tasks, task_id) {
		return false
	}

	tasks = append(tasks, task_id)
	return true
}

func Enable(ctx context.Context, orchestrator_addr, main_server_port string) {
	mux := sync.Mutex{}

	main_server_url = fmt.Sprintf("http://localhost:%s/internal/task", main_server_port)

	//Wait enable server
	<-time.After(1 * time.Second)

	var err error
	conn, err = grpc.NewClient(orchestrator_addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic("Agents: could not connect to grpc server:" + err.Error())
	}
	defer conn.Close()

	orchClient = pb.NewOrchestratorServiceClient(conn)

	//! Для сервера нужно сделать очередь из запросов, для избежания получения повторных примеров
	//Todo: Исправить наслаивание потоков

	for range COMPUTING_POWER {
		go func() {
			for {
				select {
				case <-ctx.Done():
					return

				default:
					task, err := orchClient.GetTask(ctx, &wrapperspb.StringValue{Value: ""})
					if err == DHT {
						<-time.After(TASK_WAIT_TIME_MS * time.Millisecond)
						continue
					}

					if err != nil {
						SendRequest(task, err)
					}

					mux.Lock()
					if !appendTask(task.Id) {
						mux.Unlock()
						continue
					}
					mux.Unlock()

					if err = Solve(task); err != nil {
						SendRequest(task, err)
					}

					SendRequest(task, nil)
				}
			}
		}()
	}
}

func SendRequest(task *pb.Task, err error) {
	req := &pb.TaskResult{}
	if err != nil {
		req.Id = task.GetId()
		req.Error = err.Error()
	} else {
		req.Id = task.GetId()
		req.Result = task.GetAnswer()
	}

	req_json, err := protojson.Marshal(req)
	if err != nil {
		log.Printf("Agents error: SendRequest() - %s", err.Error())
	}

	_, err = http.Post(main_server_url, "application/json", bytes.NewReader(req_json))
	if err != nil {
		log.Printf("Agents error: SendRequest(): %s", err.Error())
	}
}

func Solve(task *pb.Task) error {
	log.Println("AGENT DEBUG: solve -", task.Str)
	if task.SecondArgument.Value == 0 && task.Operation == orchestrator.Division.ToString() {
		return ErrDivideByZero
	}

	<-time.After(time.Second * time.Duration(task.OperationTimeSeconds))

	switch task.Operation {
	case orchestrator.Plus.ToString():
		task.Answer = task.FirstArgument.Value + task.SecondArgument.Value
		return nil

	case orchestrator.Minus.ToString():
		task.Answer = task.FirstArgument.Value - task.SecondArgument.Value
		return nil

	case orchestrator.Multiply.ToString():
		task.Answer = task.FirstArgument.Value * task.SecondArgument.Value
		return nil

	case orchestrator.Division.ToString():
		task.Answer = task.FirstArgument.Value / task.SecondArgument.Value
		return nil
	}
	return orchestrator.ErrExpressionIncorrect
}
