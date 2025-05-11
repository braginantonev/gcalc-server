package orchestrator

import (
	"context"
	"fmt"
	"log"

	"github.com/braginantonev/gcalc-server/pkg/database"
	dbreq "github.com/braginantonev/gcalc-server/pkg/database/requests-types"
	"github.com/braginantonev/gcalc-server/pkg/orchestrator/orchreq"

	pb "github.com/braginantonev/gcalc-server/proto/orchestrator"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

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

	END_STR        string = "end"
	TASK_ID_FORMAT string = "%s-%d-%d" // User-ExpressionID-TaskID
)

type Server struct {
	pb.OrchestratorServiceServer
	db                           *database.DataBase
	unAuthorizedExpressionsQueue *pb.Expressions
	tasks                        map[string]*pb.Task
}

func NewServer(server_db *database.DataBase) *Server {
	return &Server{
		db:                           server_db,
		unAuthorizedExpressionsQueue: &pb.Expressions{Queue: make([]*pb.Expression, 0)},
		tasks:                        make(map[string]*pb.Task, 0),
	}
}

// For test only!
func (s *Server) GetTasksQueue() map[string]*pb.Task {
	return s.tasks
}

// In task set required expression id and task id. Use taskID with internal id = -1, to get task to solve (internal).
func (s *Server) GetTask(ctx context.Context, task_id *pb.TaskID) (*pb.Task, error) {
	//log.Println("[Debug] GetTask() - get task with id:", task_id)
	if task_id == nil || task_id.Expression.Internal == -1 || task_id.Internal == -1 {
		for _, p_task := range s.tasks {
			if p_task.Status == pb.ETStatus_Backlog {
				p_task.Status = pb.ETStatus_InProgress

				expression, err := s.GetExpression(ctx, p_task.Id.Expression)
				if err != nil {
					return nil, err
				}

				expression.Status = pb.ETStatus_InProgress
				if task_id.Expression.User != "" {
					err = s.db.Update(ctx, dbreq.NewDBRequest(orchreq.DBRequest_UPDATE_Expression, int32(expression.Status), 0.0, expression.Id.User, expression.Id.Internal))
					if err != nil {
						return nil, err
					}
				}

				return p_task, nil
			}
		}
		return nil, DHT
	}

	//log.Println("[Debug] GetTask() - find task with id", task_id)

	if len(s.tasks) == 0 {
		return nil, ErrTaskNotFound
	}

	task, ok := s.tasks[fmt.Sprintf(TASK_ID_FORMAT, task_id.Expression.User, task_id.Expression.Internal, task_id.Internal)]
	if ok {
		//log.Println("[Debug] GetTask() - task found. Task:", task)
		return task, nil
	}

	return nil, ErrTaskNotFound
}

func (s *Server) SaveTaskResult(ctx context.Context, result *pb.TaskResult) (*emptypb.Empty, error) {
	//log.Println("[Debug] SaveTaskResult() - get task to save with taskId =", result.TaskID)
	task, err := s.GetTask(ctx, result.TaskID)
	if err != nil {
		return nil, err
	}

	if task.GetStatus() == pb.ETStatus_Complete {
		log.Println("task", result.TaskID, "already complete")
		return nil, nil
	}

	//! if error - check this
	expression, err := s.GetExpression(ctx, task.Id.Expression)
	if err != nil {
		return nil, err
	}

	task.Answer = result.GetResult()
	task.Status = pb.ETStatus_Complete

	log.Println("Save -", task)

	if task.IsLast {
		expression.Result = result.GetResult()
		expression.Status = pb.ETStatus_Complete

		if expression.Id.User != "" {
			err := s.db.Update(ctx, dbreq.NewDBRequest(orchreq.DBRequest_UPDATE_Expression, int32(expression.Status), expression.Result, expression.Id.User, expression.Id.Internal))
			if err != nil {
				return nil, err
			}
		}
		return nil, nil
	}

	// Return true, if task result expected
	delExpectation := func(arg *pb.Argument) {
		if arg.GetExpected() == task.Id.Internal {
			arg.Value = result.GetResult()
			arg.Expected = -1
		}
	}

	for i := 0; ; i++ {
		//log.Println("[Debug] SaveTaskResult() - del expectation to task with id:", result.TaskID.Expression.User, result.TaskID.Expression.Internal, int32(i))
		local_task, err := s.GetTask(ctx, pb.NewTaskIDWithValues(result.TaskID.Expression.User, result.TaskID.Expression.Internal, int32(i)))
		if err != nil {
			return nil, err
		}

		//log.Println("[Debug] SaveTaskResult() - got task:", local_task)

		if task.GetId() == local_task.GetId() {
			continue
		}

		if local_task.Status != pb.ETStatus_IsWaitingValues {
			continue
		}

		delExpectation(local_task.GetFirstArgument())
		delExpectation(local_task.GetSecondArgument())

		if local_task.FirstArgument.GetExpected() == -1 && local_task.SecondArgument.GetExpected() == -1 {
			local_task.Status = pb.ETStatus_Backlog
			//log.Println("[Debug] Expectation del succesfull. Task:", local_task, s.tasks)
		}

		if local_task.IsLast {
			break
		}
	}

	return nil, nil
}

// Return expression id and error
func (s *Server) AddExpression(ctx context.Context, added_expression *pb.AddedExpression) (*wrapperspb.Int32Value, error) {
	expression_str := added_expression.Str
	if expression_str == "" {
		return nil, ErrExpressionEmpty
	}

	var ex *pb.Expression
	if added_expression.User == "" { //Not authorized
		ex = &pb.Expression{
			Id:     &pb.ExpressionID{Internal: int32(len(s.unAuthorizedExpressionsQueue.Queue))},
			Status: pb.ETStatus_Analyze,
			Str:    expression_str,
		}

		if err := s.setTasksQueue(ex); err != nil {
			return nil, err
		}

		s.unAuthorizedExpressionsQueue.Queue = append(s.unAuthorizedExpressionsQueue.Queue, ex)
	} else {
		expressions, err := s.GetExpressions(ctx, wrapperspb.String(added_expression.User))
		if err != nil {
			return nil, err
		}

		ex = &pb.Expression{
			Id:     pb.NewExpressionIDWithValues(added_expression.User, int32(len(expressions.Queue))),
			Status: pb.ETStatus_Analyze,
			Str:    expression_str,
		}

		if err := s.setTasksQueue(ex); err != nil {
			return nil, err
		}

		err = s.db.Add(ctx, dbreq.NewDBRequest(orchreq.DBRequest_INSERT_Expression, ex.Id.User, ex.Id.Internal, expression_str, int32(ex.Status), 0.0))
		if err != nil {
			return nil, err
		}
	}

	return wrapperspb.Int32(ex.Id.Internal), nil
}

func (s *Server) GetExpressions(ctx context.Context, user *wrapperspb.StringValue) (*pb.Expressions, error) {
	if user.Value == "" { //Not authorized
		return s.unAuthorizedExpressionsQueue, nil
	} else {
		db_resp, err := s.db.Get(ctx, dbreq.NewDBRequest(orchreq.DBRequest_SELECT_Expressions, user.Value))
		if err != nil {
			return nil, err
		}
		return db_resp.(*pb.Expressions), nil
	}
}

/*
In expression set required id.

Return: expression
*/
func (s *Server) GetExpression(ctx context.Context, expression_id *pb.ExpressionID) (*pb.Expression, error) {
	if expression_id == nil || expression_id.Internal == -1 {
		return nil, ErrExpressionNotFound
	}

	if expression_id.User != "" {
		got, err := s.db.Get(ctx, dbreq.NewDBRequest(orchreq.DBRequest_SELECT_Expression, expression_id.User, expression_id.Internal))
		if err != nil {
			return nil, err
		}
		return got.(*pb.Expression), nil
	}

	for _, local_expression := range s.unAuthorizedExpressionsQueue.Queue {
		if local_expression.Id.User == expression_id.User && local_expression.Id.Internal == expression_id.Internal {
			return local_expression, nil
		}
	}
	return nil, ErrExpressionNotFound
}

func (s *Server) setTasksQueue(expression *pb.Expression) error {
	expression_str := expression.GetStr()
	var last_task *pb.Task
	var counter int32 = 0

	for {
		task, priority_idx, err := GetTask(expression_str)
		if err != nil {
			return err
		}

		task_str := task.GetStr()
		if task_str == END_STR {
			expression.Status = pb.ETStatus_Backlog
			last_task.IsLast = true
			log.Println(s.tasks)
			return nil
		}

		if task.GetOperation() == Equals.ToString() {
			expression_str = EraseExample(expression_str, task_str, priority_idx, counter-1)
			continue
		}

		task.Id = pb.NewTaskIDWithValues(expression.GetId().User, expression.GetId().Internal, counter)

		if task.GetStatus() != pb.ETStatus_IsWaitingValues {
			task.Status = pb.ETStatus_Backlog
		}

		s.tasks[fmt.Sprintf(TASK_ID_FORMAT, expression.Id.User, expression.Id.Internal, task.Id.Internal)] = task
		last_task = task

		expression_str = EraseExample(expression_str, task.GetStr(), priority_idx, task.GetId().Internal)
		counter++
		log.Println("task", counter, task)
	}
}

func Register(ctx context.Context, grpcServer *grpc.Server, server_db *database.DataBase) error {
	if server_db == nil {
		return database.ErrDBNotInit
	}

	if err := server_db.Create(ctx, dbreq.NewDBRequest(orchreq.DBRequest_CREATE_Orchestrator_Table)); err != nil {
		return err
	}

	pb.RegisterOrchestratorServiceServer(grpcServer, NewServer(server_db))
	return nil
}
