package orchestrator

import (
	"context"
	"log"

	pb "github.com/braginantonev/gcalc-server/proto/orchestrator"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

//Todo: Протестировать сервер

type Operator string

func (op Operator) ToString() string {
	return string(op)
}

const (
	Plus     Operator = "+"
	Minus    Operator = "-"
	Multiply Operator = "*"
	Division Operator = "/"
	Equals   Operator = "="
)

const END_STR = "end"

var expressionsQueue []*pb.Expression

type Server struct {
	pb.OrchestratorServiceServer
}

func NewServer() *Server {
	return &Server{}
}

/*
In task set required expression id and task id. Use nil task to get task to solve (internal).

Return: task; if nil - task to solve (internal)
*/
func (s *Server) GetTask(ctx context.Context, task *pb.TaskID) (*pb.Task, error) {
	if task == nil || task.GetExpression() == -1 || task.GetInternal() == -1 {
		expression, err := s.GetExpression(ctx, nil)
		if err != nil {
			return nil, err
		}

		for i := range expression.TasksQueue {
			p_task := expression.TasksQueue[i]
			if p_task.GetStatus() == pb.ETStatus_Backlog {
				p_task.Status = pb.ETStatus_InProgress
				expression.Status = pb.ETStatus_InProgress
				return p_task, nil
			}
		}
		return nil, DHT
	}

	if len(expressionsQueue) == 0 {
		return nil, nil
	}

	req_task := expressionsQueue[task.Expression].TasksQueue[task.Internal]
	if req_task == nil {
		return nil, ErrTaskNotFound
	}

	return req_task, nil
}

func (s *Server) SaveTaskResult(ctx context.Context, result *pb.TaskResult) (*emptypb.Empty, error) {
	task, err := s.GetTask(ctx, result.TaskID)
	if err != nil {
		return nil, err
	}

	if task.GetStatus() == pb.ETStatus_Complete {
		log.Println("task", result.TaskID, "already complete")
		return nil, nil
	}

	//! if error - check this
	expression, err := s.GetExpression(ctx, wrapperspb.Int32(task.GetExpressionId()))
	if err != nil {
		return nil, err
	}

	task.Answer = result.GetResult()
	task.Status = pb.ETStatus_Complete

	if task.GetId() == int32(len(expression.GetTasksQueue())-1) {
		expression.Result = result.GetResult()
		expression.Status = pb.ETStatus_Complete
		return nil, nil
	}

	// Return true, if example result expected
	delExpectation := func(arg *pb.Argument) {
		if arg.GetExpected() == task.Id {
			arg.Value = result.GetResult()
			arg.Expected = -1
		}
	}

	for _, p_task_local := range expression.GetTasksQueue() {
		if task.GetId() == p_task_local.GetId() {
			continue
		}

		delExpectation(p_task_local.GetFirstArgument())
		delExpectation(p_task_local.GetSecondArgument())

		if p_task_local.FirstArgument.GetExpected() == -1 && p_task_local.SecondArgument.GetExpected() == -1 {
			p_task_local.Status = pb.ETStatus_Backlog
		}
	}

	return nil, nil
}

// Return expression id and error
func (s *Server) AddExpression(ctx context.Context, expression *wrapperspb.StringValue) (*wrapperspb.Int32Value, error) {
	expression_str := expression.GetValue()
	if expression_str == "" {
		return nil, ErrExpressionEmpty
	}

	ex := pb.Expression{
		//!!! Если придётся реализовывать удаление выражения, то нужно изменить систему выдачи индекса !!!
		//!!! При удалении элемента, длина уменьшается, следовательно следующее добавленное выражение, будет иметь такой же индекс, что и предпоследний !!!
		Id:     int32(len(expressionsQueue)),
		Status: pb.ETStatus_Analyze,
		Str:    expression_str,
	}

	if err := setTasksQueue(&ex); err != nil {
		return nil, err
	}

	expressionsQueue = append(expressionsQueue, &ex)
	return &wrapperspb.Int32Value{Value: ex.GetId()}, nil
}

func (s *Server) GetExpressions(ctx context.Context, empty *emptypb.Empty) (*pb.Expressions, error) {
	return &pb.Expressions{Queue: expressionsQueue}, nil
}

/*
In expression set required id. Use nil expression to get in working or backlog expression (internal).

Return: expression; if nil - in working or backlog expression (internal)
*/
func (s *Server) GetExpression(ctx context.Context, expression *wrapperspb.Int32Value) (*pb.Expression, error) {
	if expression == nil {
		for _, expression := range expressionsQueue {
			expression_status := expression.GetStatus()
			if expression_status == pb.ETStatus_InProgress || expression_status == pb.ETStatus_Backlog {
				return expression, nil
			}
		}
		return nil, DHT
	}

	for _, local_expression := range expressionsQueue {
		if local_expression.GetId() == expression.GetValue() {
			return local_expression, nil
		}
	}

	return nil, ErrExpressionNotFound
}

func setTasksQueue(expression *pb.Expression) error {
	expression_str := expression.GetStr()

	for {
		task, priority_idx, err := GetTask(expression_str)
		if err != nil {
			return err
		}

		task_str := task.GetStr()
		if task_str == END_STR {
			expression.Status = pb.ETStatus_Backlog
			return nil
		}

		if task.GetOperation() == Equals.ToString() {
			expression_str = EraseExample(expression_str, task_str, priority_idx, expression.GetTasksQueue()[len(expression.TasksQueue)-1].Id)
			continue
		}

		task.Id = int32(len(expression.GetTasksQueue()))
		task.ExpressionId = expression.GetId()

		if task.GetStatus() != pb.ETStatus_IsWaitingValues {
			task.Status = pb.ETStatus_Backlog
		}

		expression.TasksQueue = append(expression.GetTasksQueue(), task)

		expression_str = EraseExample(expression_str, task.GetStr(), priority_idx, task.GetId())
	}
}

func Register(ctx context.Context, grpcServer *grpc.Server) error {
	log.Println("Orchestrator: tcp listener started at port:")

	orchestartorServiceServer := NewServer()
	pb.RegisterOrchestratorServiceServer(grpcServer, orchestartorServiceServer)

	return nil
}
