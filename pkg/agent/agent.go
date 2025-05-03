package agent

import (
	"context"
	"fmt"
	"log/slog"
	"slices"
	"sync"
	"time"

	"github.com/braginantonev/gcalc-server/pkg/orchestrator"
	pb "github.com/braginantonev/gcalc-server/proto/orchestrator"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

//! Оставить как минимум один поток для выполнения анализа

const (
	TASK_WAIT_TIME_MS      = 250
	TIME_ADDITION_MS       = 1000
	TIME_SUBTRACTION_MS    = 2000
	TIME_MULTIPLICATION_MS = 3000
	TIME_DIVISIONS_MS      = 4000
)

var (
	tasks           []string
	orchClient      pb.OrchestratorServiceClient
	COMPUTING_POWER int
	TIMES_MS        = map[string]int{
		"+": TIME_ADDITION_MS,
		"-": TIME_SUBTRACTION_MS,
		"*": TIME_MULTIPLICATION_MS,
		"/": TIME_DIVISIONS_MS,
	}
)

// Return true, if task has been appended, else - false
func appendTask(task_id string) bool {
	if slices.Contains(tasks, task_id) {
		return false
	}

	tasks = append(tasks, task_id)
	return true
}

func Enable(ctx context.Context, orch_client pb.OrchestratorServiceClient, comp_power int) {
	mux := sync.Mutex{}

	orchClient = orch_client
	COMPUTING_POWER = comp_power

	for i := range COMPUTING_POWER {
		go func(id int, ctx context.Context) {
			for {
				select {
				case <-ctx.Done():
					return

				default:
					var err_message string
					task, got_err := orchClient.GetTask(ctx, wrapperspb.String(""))
					if st, ok := status.FromError(got_err); ok {
						err_message = st.Message()
					}

					if err_message == orchestrator.DHT.Error() {
						<-time.After(TASK_WAIT_TIME_MS * time.Millisecond)
						continue
					}

					if got_err != nil {
						fmt.Println("enter", got_err)
						SendRequest(task, got_err)
						continue
					}

					mux.Lock()
					if !appendTask(task.Id) {
						mux.Unlock()
						continue
					}
					mux.Unlock()

					if got_err = Solve(i, task); got_err != nil {
						SendRequest(task, got_err)
					}

					SendRequest(task, nil)
				}
			}
		}(i, ctx)
	}
}

func SendRequest(task *pb.Task, err error) {
	task_res := &pb.TaskResult{}
	if err != nil {
		task_res.Id = task.GetId()
		task_res.Error = err.Error()
		//log.Printf("DEBUG: save task %s result with error %s", task_res.Id, task_res.Error)
	} else {
		task_res.Id = task.GetId()
		task_res.Result = task.GetAnswer()
		//log.Printf("DEBUG: save task %s result - %f", task_res.Id, task_res.Result)
	}

	_, err = orchClient.SaveTaskResult(context.TODO(), task_res)
	if err != nil {
		slog.Error("[Agent] Failed to send task result", slog.String("error", err.Error()))
	}
}

func Solve(agent_id int, task *pb.Task) error {
	//log.Printf("[Agent %d] DEBUG: solve - %s (%s)", agent_id, task.GetId(), task.GetStr())
	if task.GetSecondArgument().Value == 0 && task.GetOperation() == orchestrator.Division.ToString() {
		return ErrDivideByZero
	}

	<-time.After(time.Duration(TIMES_MS[(task.GetOperation())]) * time.Millisecond)

	switch task.GetOperation() {
	case orchestrator.Plus.ToString():
		task.Answer = task.GetFirstArgument().Value + task.GetSecondArgument().Value
		return nil

	case orchestrator.Minus.ToString():
		task.Answer = task.GetFirstArgument().Value - task.GetSecondArgument().Value
		return nil

	case orchestrator.Multiply.ToString():
		task.Answer = task.GetFirstArgument().Value * task.GetSecondArgument().Value
		return nil

	case orchestrator.Division.ToString():
		task.Answer = task.GetFirstArgument().Value / task.GetSecondArgument().Value
		return nil
	}
	return orchestrator.ErrExpressionIncorrect
}
