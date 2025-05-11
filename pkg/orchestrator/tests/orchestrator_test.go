package orchestrator_test

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/braginantonev/gcalc-server/pkg/orchestrator"
	pb "github.com/braginantonev/gcalc-server/proto/orchestrator"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func TestGetExample(t *testing.T) {
	cases := []struct {
		name         string
		example      string
		expected_str string
		expected_err bool
	}{
		{
			name:         "1+1",
			example:      "1+1",
			expected_str: "1+1",
		},
		{
			name:         "1+145",
			example:      "1+145",
			expected_str: "1+145",
		},
		{
			name:         "10+10+10+10",
			example:      "10+10",
			expected_str: "10+10",
		},
		{
			name:         "1+145*10",
			example:      "1+145*10",
			expected_str: "145*10",
		},
		{
			name:         "(1)",
			example:      "(1)",
			expected_str: "(1)",
		},
		{
			name:         "(112+1)+145*10",
			example:      "(112+1)+145*10",
			expected_str: "112+1",
		},
		{
			name:         "(10+30+(20+1*10))+1",
			example:      "(10+30+(20+1*10))+1",
			expected_str: "1*10",
		},
		{
			name:         "1+1-(2.000000)",
			example:      "1+1-(2.000000)",
			expected_str: "(2.000000)",
		},
		{
			name:         "52",
			example:      "52",
			expected_str: "end",
		},
		{
			name:         "with id",
			example:      "1+id:1",
			expected_str: "1+id:1",
		},
		{
			name:         "with id and brackets",
			example:      "1+(13+id:3)-2",
			expected_str: "13+id:3",
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			gotExample, _, err := orchestrator.GetTask(test.example)
			if err != nil {
				t.Error(err)
			}

			if gotExample.Str != test.expected_str {
				t.Errorf("GetExample(%s).ex - expected %s, but got %s", test.example, test.expected_str, gotExample.Str)
			}
		})
	}
}

func TestEraseExample(t *testing.T) {
	cases := []struct {
		name         string
		example      string
		expected_str string
	}{
		{
			name:         "1+1-(1+1)",
			example:      "1+1-(1+1)",
			expected_str: "1+1-(id:-1)",
		},
		{
			name:         "1+1-(1+1)-1+1",
			example:      "1+1-(1+1)-1+1",
			expected_str: "1+1-(id:-1)-1+1",
		},
		{
			name:         "(1)",
			example:      "(1)",
			expected_str: "id:-1",
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			gotTask, pri_idx, err := orchestrator.GetTask(test.example)
			if err != nil {
				t.Error(err)
			}

			got := orchestrator.EraseExample(test.example, gotTask.Str, pri_idx, gotTask.Id.Internal)
			if got != test.expected_str {
				t.Errorf("EraseExample(%q, %q, %d, %d) = %q, but expected: %q", test.example, gotTask.Str, pri_idx, gotTask.Id.Internal, got, test.expected_str)
			}
		})
	}
}

func TestSetTasksQueue(t *testing.T) {
	ctx := context.Background()
	s := orchestrator.NewServer(nil)

	tests := []struct {
		name         string
		expression   string
		queue        []*pb.Task
		expected_err error
	}{
		{
			name:       "1+1",
			expression: "1+1",
			queue: []*pb.Task{
				{
					Id:             pb.NewTaskIDWithValues("", 0, 0),
					FirstArgument:  &pb.Argument{Value: 1, Expected: -1},
					SecondArgument: &pb.Argument{Value: 1, Expected: -1},
					Operation:      orchestrator.Plus.ToString(),
					Status:         pb.ETStatus_Backlog,
					Str:            "1+1",
					IsLast:         true,
				},
			},
			expected_err: nil,
		},
		{
			name:       "1+1-1",
			expression: "1+1-1",
			queue: []*pb.Task{
				{
					Id:             pb.NewTaskIDWithValues("", 1, 0),
					FirstArgument:  &pb.Argument{Value: 1, Expected: -1},
					SecondArgument: &pb.Argument{Value: 1, Expected: -1},
					Operation:      orchestrator.Plus.ToString(),
					Status:         pb.ETStatus_Backlog,
					Str:            "1+1",
				},
				{
					Id:             pb.NewTaskIDWithValues("", 1, 1),
					FirstArgument:  &pb.Argument{Expected: 0},
					SecondArgument: &pb.Argument{Value: 1, Expected: -1},
					Operation:      orchestrator.Minus.ToString(),
					Status:         pb.ETStatus_IsWaitingValues,
					Str:            "id:0-1",
					IsLast:         true,
				},
			},
		},
		{
			name:       "1*(2*1)+1",
			expression: "1*(2*1)+1",
			queue: []*pb.Task{
				{
					Id:             pb.NewTaskIDWithValues("", 2, 0),
					FirstArgument:  &pb.Argument{Value: 2, Expected: -1},
					SecondArgument: &pb.Argument{Value: 1, Expected: -1},
					Operation:      orchestrator.Multiply.ToString(),
					Status:         pb.ETStatus_Backlog,
					Str:            "2*1",
				},
				{
					Id:             pb.NewTaskIDWithValues("", 2, 1),
					FirstArgument:  &pb.Argument{Value: 1, Expected: -1},
					SecondArgument: &pb.Argument{Expected: 0},
					Operation:      orchestrator.Multiply.ToString(),
					Status:         pb.ETStatus_IsWaitingValues,
					Str:            "1*id:0",
				},
				{
					Id:             pb.NewTaskIDWithValues("", 2, 2),
					FirstArgument:  &pb.Argument{Expected: 1},
					SecondArgument: &pb.Argument{Value: 1, Expected: -1},
					Operation:      orchestrator.Plus.ToString(),
					Status:         pb.ETStatus_IsWaitingValues,
					Str:            "id:1+1",
					IsLast:         true,
				},
			},
		},
		{
			name:       "1/0",
			expression: "1/0",
			queue: []*pb.Task{
				{
					Id:             pb.NewTaskIDWithValues("", 3, 0),
					FirstArgument:  &pb.Argument{Value: 1, Expected: -1},
					SecondArgument: &pb.Argument{Expected: -1},
					Operation:      orchestrator.Division.ToString(),
					Status:         pb.ETStatus_Backlog,
					Str:            "1/0",
					IsLast:         true,
				},
			},
		},
	}

	expected_queue := make(map[string]*pb.Task, 0)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, gotErr := s.AddExpression(ctx, &pb.AddedExpression{Str: test.expression})
			if !errors.Is(gotErr, test.expected_err) {
				t.Error("expected:", test.expected_err, "but got:", gotErr)
				return
			}

			for _, task := range test.queue {
				expected_queue[fmt.Sprintf(orchestrator.TASK_ID_FORMAT, task.Id.Expression.User, task.Id.Expression.Internal, task.Id.Internal)] = task
			}

			got_queue := s.GetTasksQueue()
			if !reflect.DeepEqual(expected_queue, got_queue) {
				t.Error("got:", got_queue, "but expected:", test.queue)
			}
		})
	}
}

func TestGetExpressionsQueue(t *testing.T) {
	ctx := context.Background()
	s := orchestrator.NewServer(nil)

	tests := []struct {
		name           string
		expression     string
		user           string
		expected_queue []*pb.Expression
		expected_err   error
	}{
		{
			name:           "empty",
			expression:     "",
			user:           "",
			expected_queue: make([]*pb.Expression, 0),
			expected_err:   orchestrator.ErrExpressionEmpty,
		},
		{
			name:       "1 expression",
			expression: "1+1",
			user:       "",
			expected_queue: []*pb.Expression{
				{
					Id:     pb.NewExpressionIDWithValues("", 0),
					Status: pb.ETStatus_Backlog,
					Str:    "1+1",
				},
			},
		},
		{
			name:       "2 expressions",
			expression: "6-5",
			user:       "",
			expected_queue: []*pb.Expression{
				{
					Id:     pb.NewExpressionIDWithValues("", 0),
					Status: pb.ETStatus_Backlog,
					Str:    "1+1",
				},
				{
					Id:     pb.NewExpressionIDWithValues("", 1),
					Status: pb.ETStatus_Backlog,
					Str:    "6-5",
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, gotErr := s.AddExpression(ctx, &pb.AddedExpression{User: test.user, Str: test.expression})
			if !errors.Is(gotErr, test.expected_err) {
				t.Error("expected error:", test.expected_err, "but expected:", gotErr)
			}

			gotQueue, err := s.GetExpressions(ctx, wrapperspb.String(test.user))
			if err != nil {
				t.Error(err)
			}

			if !reflect.DeepEqual(test.expected_queue, gotQueue.GetQueue()) {
				t.Error("got:", gotQueue, "but expected:", test.expected_queue)
			}
		})
	}
}

func TestSetTaskResult(t *testing.T) {
	ctx := context.Background()
	s := orchestrator.NewServer(nil)

	tests := []struct {
		name          string
		expression    *pb.AddedExpression
		result        *pb.TaskResult
		expected_task *pb.Task
		expected_err  error
	}{
		{
			name: "1+1",
			expression: &pb.AddedExpression{
				User: "",
				Str:  "1+1",
			},
			result: &pb.TaskResult{
				TaskID: pb.NewTaskIDWithValues("", 0, 0),
				Result: 2,
			},
			expected_task: &pb.Task{
				Id:             pb.NewTaskIDWithValues("", 0, 0),
				FirstArgument:  &pb.Argument{Value: 1, Expected: -1},
				SecondArgument: &pb.Argument{Value: 1, Expected: -1},
				Operation:      orchestrator.Plus.ToString(),
				Str:            "1+1",
				Answer:         2,
				Status:         pb.ETStatus_Complete,
				IsLast:         true,
			},
		},
		{
			name: "1+1+1",
			expression: &pb.AddedExpression{
				User: "",
				Str:  "1+1+1",
			},
			result: &pb.TaskResult{
				TaskID: pb.NewTaskIDWithValues("", 1, 0),
				Result: 2,
			},
			expected_task: &pb.Task{
				Id:             pb.NewTaskIDWithValues("", 1, 0),
				FirstArgument:  &pb.Argument{Value: 1, Expected: -1},
				SecondArgument: &pb.Argument{Value: 1, Expected: -1},
				Operation:      orchestrator.Plus.ToString(),
				Str:            "1+1",
				Answer:         2,
				Status:         pb.ETStatus_Complete,
			},
		},
		{
			name: "previous id:1_0+1",
			expression: &pb.AddedExpression{
				User: "",
				Str:  "1+1+1",
			},
			result: &pb.TaskResult{
				TaskID: pb.NewTaskIDWithValues("", 1, 1),
				Result: 3,
			},
			expected_task: &pb.Task{
				Id:             pb.NewTaskIDWithValues("", 1, 1),
				FirstArgument:  &pb.Argument{Value: 2, Expected: -1},
				SecondArgument: &pb.Argument{Value: 1, Expected: -1},
				Operation:      orchestrator.Plus.ToString(),
				Str:            "id:0+1",
				Answer:         3,
				Status:         pb.ETStatus_Complete,
				IsLast:         true,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s.AddExpression(ctx, test.expression)
			_, gotErr := s.SaveTaskResult(ctx, test.result)
			if !errors.Is(gotErr, test.expected_err) {
				t.Error("expected error:", test.expected_err, "but got:", gotErr)
			}

			gotTask, err := s.GetTask(ctx, pb.NewTaskIDWithValues(test.expression.User, test.result.TaskID.Expression.Internal, test.result.TaskID.Internal))
			if err != nil {
				t.Error(err)
			}

			if !reflect.DeepEqual(test.expected_task, gotTask) {
				t.Error("got:", gotTask, "but expected:", test.expected_task)
			}
		})
	}
}
