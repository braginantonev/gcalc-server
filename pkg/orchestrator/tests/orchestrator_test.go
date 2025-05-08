package orchestrator_test

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/braginantonev/gcalc-server/pkg/orchestrator"
	pb "github.com/braginantonev/gcalc-server/proto/orchestrator"
	"google.golang.org/protobuf/types/known/emptypb"
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
			expected_str: "1+1-(id:0)",
		},
		{
			name:         "1+1-(1+1)-1+1",
			example:      "1+1-(1+1)-1+1",
			expected_str: "1+1-(id:0)-1+1",
		},
		{
			name:         "(1)",
			example:      "(1)",
			expected_str: "id:0",
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			gotExample, pri_idx, err := orchestrator.GetTask(test.example)
			if err != nil {
				t.Error(err)
			}

			got := orchestrator.EraseExample(test.example, gotExample.Str, pri_idx, gotExample.Id)
			if got != test.expected_str {
				t.Errorf("EraseExample(%q, %q, %d, %d) = %q, but expected: %q", test.example, gotExample.Str, pri_idx, gotExample.Id, got, test.expected_str)
			}
		})
	}
}

// Todo: Исправить в тесте статусы
func TestSetTasksQueue(t *testing.T) {
	ctx := context.Background()
	s := orchestrator.Server{}

	tests := []struct {
		name           string
		expression     string
		expected_queue []*pb.Task
		expected_err   error
	}{
		{
			name:       "1+1",
			expression: "1+1",
			expected_queue: []*pb.Task{
				{
					ExpressionId:   0,
					Id:             0,
					FirstArgument:  &pb.Argument{Value: 1, Expected: -1},
					SecondArgument: &pb.Argument{Value: 1, Expected: -1},
					Operation:      orchestrator.Plus.ToString(),
					Status:         pb.ETStatus_Backlog,
					Str:            "1+1",
				},
			},
			expected_err: nil,
		},
		{
			name:       "1+1-1",
			expression: "1+1-1",
			expected_queue: []*pb.Task{
				{
					ExpressionId:   1,
					Id:             0,
					FirstArgument:  &pb.Argument{Value: 1, Expected: -1},
					SecondArgument: &pb.Argument{Value: 1, Expected: -1},
					Operation:      orchestrator.Plus.ToString(),
					Status:         pb.ETStatus_Backlog,
					Str:            "1+1",
				},
				{
					ExpressionId:   1,
					Id:             1,
					FirstArgument:  &pb.Argument{Expected: 0},
					SecondArgument: &pb.Argument{Value: 1, Expected: -1},
					Operation:      orchestrator.Minus.ToString(),
					Status:         pb.ETStatus_IsWaitingValues,
					Str:            "id:0-1",
				},
			},
		},
		{
			name:       "1*(2*1)+1",
			expression: "1*(2*1)+1",
			expected_queue: []*pb.Task{
				{
					ExpressionId:   2,
					Id:             0,
					FirstArgument:  &pb.Argument{Value: 2, Expected: -1},
					SecondArgument: &pb.Argument{Value: 1, Expected: -1},
					Operation:      orchestrator.Multiply.ToString(),
					Status:         pb.ETStatus_Backlog,
					Str:            "2*1",
				},
				{
					ExpressionId:   2,
					Id:             1,
					FirstArgument:  &pb.Argument{Value: 1, Expected: -1},
					SecondArgument: &pb.Argument{Expected: 0},
					Operation:      orchestrator.Multiply.ToString(),
					Status:         pb.ETStatus_IsWaitingValues,
					Str:            "1*id:0",
				},
				{
					ExpressionId:   2,
					Id:             2,
					FirstArgument:  &pb.Argument{Expected: 1},
					SecondArgument: &pb.Argument{Value: 1, Expected: -1},
					Operation:      orchestrator.Plus.ToString(),
					Status:         pb.ETStatus_IsWaitingValues,
					Str:            "id:1+1",
				},
			},
		},
		{
			name:       "1/0",
			expression: "1/0",
			expected_queue: []*pb.Task{
				{
					ExpressionId:   3,
					Id:             0,
					FirstArgument:  &pb.Argument{Value: 1, Expected: -1},
					SecondArgument: &pb.Argument{Expected: -1},
					Operation:      orchestrator.Division.ToString(),
					Status:         pb.ETStatus_Backlog,
					Str:            "1/0",
				},
			},
		},
	}

	for i, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, gotErr := s.AddExpression(ctx, &wrapperspb.StringValue{Value: test.expression})
			expression, err := s.GetExpression(ctx, wrapperspb.Int32(int32(i)))
			if err != nil {
				t.Error(err)
			}

			gotQueue := expression.GetTasksQueue()

			if !errors.Is(gotErr, test.expected_err) {
				t.Log("got:", gotQueue)
				t.Error("expected:", test.expected_err, "but got:", gotErr)
				return
			}

			if !reflect.DeepEqual(test.expected_queue, gotQueue) {
				t.Error("got:", gotQueue, "but expected:", test.expected_queue)
			}
		})
	}
}

func TestGetExpressionsQueue(t *testing.T) {
	ctx := context.Background()
	s := orchestrator.NewServer()

	tests := []struct {
		name           string
		expression     string
		expected_queue []*pb.Expression
		expected_err   error
	}{
		{
			name:         "empty",
			expression:   "",
			expected_err: orchestrator.ErrExpressionEmpty,
		},
		{
			name:       "1 expression",
			expression: "1+1",
			expected_queue: []*pb.Expression{
				{
					Id:     0,
					Status: pb.ETStatus_Backlog,
					Str:    "1+1",
					TasksQueue: []*pb.Task{
						{
							ExpressionId:   0,
							Id:             0,
							FirstArgument:  &pb.Argument{Value: 1, Expected: -1},
							SecondArgument: &pb.Argument{Value: 1, Expected: -1},
							Operation:      orchestrator.Plus.ToString(),
							Str:            "1+1",
							Status:         pb.ETStatus_Backlog,
						},
					},
				},
			},
		},
		{
			name:       "2 expressions",
			expression: "6-5",
			expected_queue: []*pb.Expression{
				{
					Id:     0,
					Status: pb.ETStatus_Backlog,
					Str:    "1+1",
					TasksQueue: []*pb.Task{
						{
							ExpressionId:   0,
							Id:             0,
							FirstArgument:  &pb.Argument{Value: 1, Expected: -1},
							SecondArgument: &pb.Argument{Value: 1, Expected: -1},
							Operation:      orchestrator.Plus.ToString(),
							Str:            "1+1",
							Status:         pb.ETStatus_Backlog,
						},
					},
				},
				{
					Id:     1,
					Status: pb.ETStatus_Backlog,
					Str:    "6-5",
					TasksQueue: []*pb.Task{
						{
							ExpressionId:   1,
							Id:             0,
							FirstArgument:  &pb.Argument{Value: 6, Expected: -1},
							SecondArgument: &pb.Argument{Value: 5, Expected: -1},
							Operation:      orchestrator.Minus.ToString(),
							Str:            "6-5",
							Status:         pb.ETStatus_Backlog,
						},
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, gotErr := s.AddExpression(ctx, &wrapperspb.StringValue{Value: test.expression})
			if !errors.Is(gotErr, test.expected_err) {
				t.Error("expected error:", test.expected_err, "but expected:", gotErr)
			}

			gotQueue, err := s.GetExpressions(ctx, &emptypb.Empty{})
			if err != nil {
				t.Error(err)
			}

			if !reflect.DeepEqual(test.expected_queue, gotQueue.GetQueue()) {
				t.Error("got:", gotQueue, "but expected:", test.expected_queue)
			}
		})
	}
}

func TestSetExampleResult(t *testing.T) {
	ctx := context.Background()
	s := orchestrator.NewServer()

	tests := []struct {
		name          string
		expression    string
		result        float64
		task_id       int32
		expression_id int32
		expected_task *pb.Task
		expected_err  error
	}{
		{
			name:          "1+1",
			expression:    "1+1",
			result:        2,
			expression_id: 0,
			task_id:       0,
			expected_task: &pb.Task{
				ExpressionId:   0,
				Id:             0,
				FirstArgument:  &pb.Argument{Value: 1, Expected: -1},
				SecondArgument: &pb.Argument{Value: 1, Expected: -1},
				Operation:      orchestrator.Plus.ToString(),
				Str:            "1+1",
				Answer:         2,
				Status:         pb.ETStatus_Complete,
			},
		},
		{
			name:          "1+1+1",
			expression:    "1+1+1",
			result:        2,
			expression_id: 1,
			task_id:       0,
			expected_task: &pb.Task{
				ExpressionId:   1,
				Id:             0,
				FirstArgument:  &pb.Argument{Value: 1, Expected: -1},
				SecondArgument: &pb.Argument{Value: 1, Expected: -1},
				Operation:      orchestrator.Plus.ToString(),
				Str:            "1+1",
				Answer:         2,
				Status:         pb.ETStatus_Complete,
			},
		},
		{
			name:          "previous id:1_0+1",
			expression:    "1+1+1",
			expression_id: 1,
			task_id:       1,
			result:        3,
			expected_task: &pb.Task{
				ExpressionId:   1,
				Id:             1,
				FirstArgument:  &pb.Argument{Value: 2, Expected: -1},
				SecondArgument: &pb.Argument{Value: 1, Expected: -1},
				Operation:      orchestrator.Plus.ToString(),
				Str:            "id:0+1",
				Answer:         3,
				Status:         pb.ETStatus_Complete,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s.AddExpression(ctx, &wrapperspb.StringValue{Value: test.expression})
			_, gotErr := s.SaveTaskResult(ctx, &pb.TaskResult{TaskID: pb.NewTaskIDWithValues(test.expression_id, test.task_id), Result: test.result})
			if !errors.Is(gotErr, test.expected_err) {
				t.Error("expected error:", test.expected_err, "but got:", gotErr)
			}

			gotTask, err := s.GetTask(ctx, pb.NewTaskIDWithValues(test.expression_id, test.task_id))
			if err != nil {
				t.Error(err)
			}

			if !reflect.DeepEqual(test.expected_task, gotTask) {
				t.Error("got:", gotTask, "but expected:", test.expected_task)
			}
		})
	}
}

//Todo: Fix this test
/*
func TestGetTask(t *testing.T) {
	ctx := context.Background()
	s := orchestrator.NewServer()

	tests := []struct {
		name          string
		expression    string
		expression_id int32
		task_id       int32
		result        float64
		expected_task *pb.Task
		expected_err  error
	}{
		{
			name:          "2*4 complete",
			expression:    "2*4",
			expression_id: -1,
			result:        8,
			expected_task: &pb.Task{
				ExpressionId:   0,
				Id:             0,
				FirstArgument:  &pb.Argument{Value: 2, Expected: -1},
				SecondArgument: &pb.Argument{Value: 4, Expected: -1},
				Operation:      orchestrator.Multiply.ToString(),
				Status:         pb.ETStatus_InProgress,
				Str:            "2*4",
			},
		},
		{
			name:          "1+1+1 in progress",
			expression:    "1+1+1",
			expression_id: -1,
			result:        2,
			expected_task: &pb.Task{
				ExpressionId:   1,
				Id:             0,
				FirstArgument:  &pb.Argument{Value: 1, Expected: -1},
				SecondArgument: &pb.Argument{Value: 1, Expected: -1},
				Operation:      orchestrator.Plus.ToString(),
				Status:         pb.ETStatus_InProgress,
				Str:            "1+1",
			},
		},
		{
			name:          "1+1+1 complete",
			expression:    "1+1+1",
			expression_id: -1,
			result:        3,
			expected_task: &pb.Task{
				ExpressionId:   1,
				Id:             1,
				FirstArgument:  &pb.Argument{Value: 2, Expected: -1},
				SecondArgument: &pb.Argument{Value: 1, Expected: -1},
				Operation:      orchestrator.Plus.ToString(),
				Status:         pb.ETStatus_InProgress,
				Str:            "id:0+1",
			},
		},
		{
			name:          "find by id",
			expression:    "1-1",
			expression_id: 0,
			task_id:       0,
			result:        8,
			expected_task: &pb.Task{
				ExpressionId:   0,
				Id:             0,
				FirstArgument:  &pb.Argument{Value: 2, Expected: -1},
				SecondArgument: &pb.Argument{Value: 4, Expected: -1},
				Operation:      orchestrator.Multiply.ToString(),
				Status:         pb.ETStatus_Complete,
				Answer:         8,
				Str:            "2*4",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, gotErr := s.AddExpression(ctx, &wrapperspb.StringValue{Value: test.expression})
			if !errors.Is(gotErr, test.expected_err) {
				t.Error("expected error:", test.expected_err, "but got:", gotErr)
			}

			var task *pb.Task
			if test.expression_id != -1 {
				task = pb.NewTask()
				task.ExpressionId = test.expression_id
				task.Id = test.task_id
			}

			gotTask, gotErr := s.GetTask(ctx, task)
			if !errors.Is(gotErr, test.expected_err) {
				t.Error("expected error:", test.expected_err, "but got:", gotErr)
			}

			if !reflect.DeepEqual(test.expected_task, gotTask) {
				t.Error("got:", gotTask, "but expected:", test.expected_task)
			}

			s.SaveTaskResult(ctx, &pb.TaskResult{ExpressionId: task.ExpressionId, Id: test.expected_task.Id, Result: test.result})
		})
	}
}
*/
