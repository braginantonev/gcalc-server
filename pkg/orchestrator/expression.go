package orchestrator

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	pb "github.com/braginantonev/gcalc-server/proto/orchestrator"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

//Todo: Исправить тесты для использования grpc
//Todo: Протестировать сервер
//Todo: Удалить пакет calc

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

func (s *Server) GetTask(ctx context.Context, id *wrapperspb.StringValue) (*pb.Task, error) {
	id_str := id.GetValue()
	if id_str == "" {
		expression_local, err := s.GetExpression(ctx, id)
		if err != nil {
			return nil, DHT
		}

		expId, err := strconv.Atoi(expression_local.Id)
		if err != nil {
			return nil, err
		}

		//! При возникновении багов - обратить внимание
		/* Я заменил сохраняемый тип в очереди
		Теперь сохраняется именно указатель на выражение,
		соответственно вся вот эта херня, где я получал ссылку на выражение была убрана.
		Я хер знает как оно теперь будет себя вести, поэтому надо тестить */
		p_expression := expressionsQueue[expId]
		for i := range p_expression.TasksQueue {
			p_task := p_expression.TasksQueue[i]
			if p_task.Status == pb.ETStatus_Backlog {
				p_task.Status = pb.ETStatus_InProgress
				p_expression.Status = pb.ETStatus_InProgress
				return p_task, nil
			}
		}
		return nil, DHT
	}

	for _, p_expression := range expressionsQueue {
		for _, p_task := range p_expression.TasksQueue {
			if p_task.Id == id_str {
				return p_task, nil
			}
		}
	}

	return nil, ErrTaskNotFound
}

func (s *Server) SaveTaskResult(ctx context.Context, result *pb.TaskResult) (*emptypb.Empty, error) {
	p_task, err := s.GetTask(ctx, &wrapperspb.StringValue{Value: result.Id})
	if err != nil {
		return nil, err
	}

	if p_task.Status == pb.ETStatus_Complete {
		log.Println("task", result.Id, "already complete")
		return nil, nil
	}

	task_low_line_idx := strings.IndexRune(p_task.Id, '_')
	p_expression, err := s.GetExpression(ctx, &wrapperspb.StringValue{Value: p_task.Id[:task_low_line_idx]})
	if err != nil {
		return nil, err
	}

	task_id, err := strconv.Atoi(p_task.Id[task_low_line_idx+1:])
	if err != nil {
		return nil, err
	}

	p_task.Answer = result.GetResult()
	p_task.Status = pb.ETStatus_Complete

	if task_id == len(p_expression.TasksQueue)-1 {
		p_expression.Result = result.GetResult()
		p_expression.Status = pb.ETStatus_Complete
		return nil, nil
	}

	// Return true, if example result expected
	delExpectation := func(arg *pb.Argument) {
		if arg.Expected == p_task.Id {
			arg.Value = result.GetResult()
			arg.Expected = ""
		}
	}

	for _, p_task_local := range p_expression.TasksQueue {
		if p_task.Id == p_task_local.Id {
			continue
		}

		delExpectation(p_task_local.FirstArgument)
		delExpectation(p_task_local.SecondArgument)

		if p_task_local.FirstArgument.Expected == "" && p_task_local.SecondArgument.Expected == "" {
			p_task_local.Status = pb.ETStatus_Backlog
		}
	}

	return nil, nil
}

// Return expression id and error
func (s *Server) AddExpression(ctx context.Context, expression *wrapperspb.StringValue) (*wrapperspb.StringValue, error) {
	expression_str := expression.GetValue()
	if expression_str == "" {
		return nil, ErrExpressionEmpty
	}

	ex := pb.Expression{
		//!!! Если придётся реализовывать удаление выражения, то нужно изменить систему выдачи индекса !!!
		//!!! При удалении элемента, длина уменьшается, следовательно следующее добавленное выражение, будет иметь такой же индекс, что и предпоследний !!!
		Id:     fmt.Sprint(len(expressionsQueue)),
		Status: pb.ETStatus_Analyze,
		Str:    expression_str,
	}

	if err := setTasksQueue(&ex); err != nil {
		return nil, err
	}

	expressionsQueue = append(expressionsQueue, &ex)
	return &wrapperspb.StringValue{Value: ex.Id}, nil
}

func (s *Server) GetExpressions(ctx context.Context, empty *emptypb.Empty) (*pb.Expressions, error) {
	return &pb.Expressions{Queue: expressionsQueue}, nil
}

func (s *Server) GetExpression(ctx context.Context, id *wrapperspb.StringValue) (*pb.Expression, error) {
	id_str := id.GetValue()
	if id_str == "" {
		for _, expression := range expressionsQueue {
			if expression.Status == pb.ETStatus_InProgress || expression.Status == pb.ETStatus_Backlog {
				return expression, nil
			}
		}
		return nil, DHT
	}

	for _, expression := range expressionsQueue {
		if expression.Id == id_str {
			return expression, nil
		}
	}

	return nil, ErrExpressionNotFound
}

func setTasksQueue(expression *pb.Expression) error {
	expression_str := expression.Str

	for {
		example, priority_idx, err := GetExample(expression_str)
		if err != nil {
			return err
		}

		if example.Str == END_STR {
			expression.Status = pb.ETStatus_Backlog
			return nil
		}

		if example.Operation == Equals.ToString() {
			expression_str = EraseExample(expression_str, example.Str, priority_idx, expression.TasksQueue[len(expression.TasksQueue)-1].Id)
			continue
		}

		example.Id = expression.Id + "_" + fmt.Sprint(len(expression.TasksQueue))

		if example.Status != pb.ETStatus_IsWaitingValues {
			example.Status = pb.ETStatus_Backlog
		}

		expression.TasksQueue = append(expression.TasksQueue, example)

		expression_str = EraseExample(expression_str, example.Str, priority_idx, example.Id)
	}
}

// Старый вариант - на всякий случай
// Если ничего не заработает верну отсюда, хоть я и работаю в отдельной ветке
//! Перед слиянием с основной веткой - удалить этот рудимент
/*

type Expression struct {
	Id         string      `json:"id"`
	Status     calc.Status `json:"status"`
	String     string
	TasksQueue []calc.Example
	Result     float64 `json:"result"`
}


// Return id expression and error
func AddExpression(expression string) (string, error) {
	if expression == "" {
		return "", ErrExpressionEmpty
	}

	ex := Expression{
		//!!! Если придётся реализовывать удаление выражения, то нужно изменить систему выдачи индекса !!!
		//!!! При удалении элемента, длина уменьшается, следовательно следующее добавленное выражение, будет иметь такой же индекс, что и предпоследний !!!
		Id:     fmt.Sprint(len(expressionsQueue)),
		Status: calc.StatusAnalyze,
		String: expression,
	}

	if err := ex.setTasksQueue(); err != nil {
		return "", err
	}

	expressionsQueue = append(expressionsQueue, ex)
	return ex.Id, nil
}

func GetExpression(id string) (Expression, error) {
	if id == "" {
		for _, ex := range expressionsQueue {
			if ex.Status == calc.StatusInProgress || ex.Status == calc.StatusBacklog {
				return ex, nil
			}
		}
		return Expression{}, DHT
	}

	for _, ex := range expressionsQueue {
		if ex.Id == id {
			return ex, nil
		}
	}

	return Expression{}, ErrExpressionNotFound
}

func GetExpressionsQueue() []Expression {
	return expressionsQueue
}

func GetTask(id string) (calc.Example, error) {
	if id == "" {
		expression_local, err := GetExpression("")
		if err != nil {
			return calc.Example{}, DHT
		}

		expId, err := strconv.Atoi(expression_local.Id)
		if err != nil {
			return calc.Example{}, err
		}

		p_expression := &expressionsQueue[expId]
		for i := range p_expression.TasksQueue {
			p_example := &p_expression.TasksQueue[i]
			if p_example.Status == calc.StatusBacklog {
				p_example.Status = calc.StatusInProgress
				p_expression.Status = calc.StatusInProgress
				return *p_example, nil
			}
		}
		return calc.Example{}, DHT
	}

	for _, exp := range expressionsQueue {
		for _, example := range exp.TasksQueue {
			if example.Id == id {
				return example, nil
			}
		}
	}

	return calc.Example{}, ErrTaskNotFound
}

func SetExampleResult(id string, result float64) error {
	example_local, err := GetTask(id)
	if err != nil {
		return err
	}

	if example_local.Status == calc.StatusComplete {
		fmt.Println("example", id, "already complete")
		return nil
	}

	low_line_idx := strings.IndexRune(example_local.Id, '_')
	exp_local, err := GetExpression(example_local.Id[:low_line_idx])
	if err != nil {
		return err
	}

	expressionId_int, err := strconv.Atoi(exp_local.Id)
	if err != nil {
		return err
	}

	exampleId_int, err := strconv.Atoi(example_local.Id[low_line_idx+1:])
	if err != nil {
		return err
	}

	p_expression := &expressionsQueue[expressionId_int]
	p_example := &p_expression.TasksQueue[exampleId_int]

	p_example.Answer = result
	p_example.Status = calc.StatusComplete

	if exampleId_int == len(p_expression.TasksQueue)-1 {
		p_expression.Result = result
		p_expression.Status = calc.StatusComplete
		return nil
	}

	// Return true, if example result expected
	delExpectation := func(arg *calc.Argument) {
		if arg.Expected == p_example.Id {
			arg.Value = result
			arg.Expected = ""
		}
	}

	for i := range p_expression.TasksQueue {
		p_local_example := &p_expression.TasksQueue[i]
		if p_example.Id == p_local_example.Id {
			continue
		}

		delExpectation(&p_local_example.FirstArgument)
		delExpectation(&p_local_example.SecondArgument)

		if p_local_example.FirstArgument.Expected == "" && p_local_example.SecondArgument.Expected == "" {
			p_local_example.Status = calc.StatusBacklog
		}
	}

	return nil
}
*/
