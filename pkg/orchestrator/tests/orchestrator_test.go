package orchestrator_test

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/Antibrag/gcalc-server/pkg/calc"
	"github.com/Antibrag/gcalc-server/pkg/orchestrator"
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
			gotExample, _, _ := orchestrator.GetExample(test.example)
			if gotExample.String != test.expected_str {
				t.Errorf("GetExample(%s).ex - expected %s, but got %s", test.example, test.expected_str, gotExample.String)
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
			expected_str: "1+1-(id:)",
		},
		{
			name:         "1+1-(1+1)-1+1",
			example:      "1+1-(1+1)-1+1",
			expected_str: "1+1-(id:)-1+1",
		},
		{
			name:         "(1)",
			example:      "(1)",
			expected_str: "id:",
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			gotExample, pri_idx, _ := orchestrator.GetExample(test.example)
			got := orchestrator.EraseExample(test.example, gotExample.String, pri_idx, gotExample.Id)
			if got != test.expected_str {
				t.Errorf("EraseExample(%q, %q, %d, %s) = %q, but expected: %q", test.example, gotExample.String, pri_idx, gotExample.Id, got, test.expected_str)
			}
		})
	}
}

func TestSetTasksQueue(t *testing.T) {
	tests := []struct {
		name           string
		expression     string
		expected_queue []calc.Example
		expected_err   error
	}{
		{
			name:       "1+1",
			expression: "1+1",
			expected_queue: []calc.Example{
				{
					Id:             "0_0",
					FirstArgument:  calc.Argument{Value: 1},
					SecondArgument: calc.Argument{Value: 1},
					Operation:      calc.Plus,
					Status:         calc.StatusBacklog,
					String:         "1+1",
				},
			},
			expected_err: nil,
		},
		{
			name:       "1+1-1",
			expression: "1+1-1",
			expected_queue: []calc.Example{
				{
					Id:             "1_0",
					FirstArgument:  calc.Argument{Value: 1},
					SecondArgument: calc.Argument{Value: 1},
					Operation:      calc.Plus,
					Status:         calc.StatusBacklog,
					String:         "1+1",
				},
				{
					Id:             "1_1",
					FirstArgument:  calc.Argument{Expected: "1_0"},
					SecondArgument: calc.Argument{Value: 1},
					Operation:      calc.Minus,
					Status:         calc.StatusBacklog,
					String:         "id:1_0-1",
				},
			},
		},
		{
			name:       "1*(2*1)+1",
			expression: "1*(2*1)+1",
			expected_queue: []calc.Example{
				{
					Id:             "2_0",
					FirstArgument:  calc.Argument{Value: 2},
					SecondArgument: calc.Argument{Value: 1},
					Operation:      calc.Multiply,
					Status:         calc.StatusBacklog,
					String:         "2*1",
				},
				{
					Id:             "2_1",
					FirstArgument:  calc.Argument{Value: 1},
					SecondArgument: calc.Argument{Expected: "2_0"},
					Operation:      calc.Multiply,
					Status:         calc.StatusBacklog,
					String:         "1*id:2_0",
				},
				{
					Id:             "2_2",
					FirstArgument:  calc.Argument{Expected: "2_1"},
					SecondArgument: calc.Argument{Value: 1},
					Operation:      calc.Plus,
					Status:         calc.StatusBacklog,
					String:         "id:2_1+1",
				},
			},
		},
		{
			name:       "1/0",
			expression: "1/0",
			expected_queue: []calc.Example{
				{
					Id:             "3_0",
					FirstArgument:  calc.Argument{Value: 1},
					SecondArgument: calc.Argument{Value: 0},
					Operation:      calc.Division,
					Status:         calc.StatusBacklog,
					String:         "1/0",
				},
			},
		},
	}

	for i, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, gotErr := orchestrator.AddExpression(test.expression)
			ex, err := orchestrator.GetExpression(fmt.Sprint(i))
			if err != nil {
				t.Error(err)
			}
			t.Log(ex)

			gotQueue := ex.TasksQueue

			if !errors.Is(gotErr, test.expected_err) {
				t.Log("got:", gotQueue)
				t.Error("expected:", test.expected_err, "but got:", gotErr)
				return
			}

			if !reflect.DeepEqual(test.expected_queue, gotQueue) {
				t.Log("got:", gotQueue, "but expected:", test.expected_queue)
				t.Error()
			}
		})
	}
}

func TestGetExpressionsQueue(t *testing.T) {
	tests := []struct {
		name           string
		expression     string
		expected_queue []orchestrator.Expression
		expected_err   error
	}{
		{
			name:         "error",
			expression:   "",
			expected_err: orchestrator.ErrExpressionEmpty,
		},
		{
			name:       "1 expression",
			expression: "1+1",
			expected_queue: []orchestrator.Expression{
				{
					Id:     "0",
					Status: calc.StatusBacklog,
					String: "1+1",
					TasksQueue: []calc.Example{
						{
							Id:             "0_0",
							FirstArgument:  calc.Argument{Value: 1},
							SecondArgument: calc.Argument{Value: 1},
							Operation:      calc.Plus,
							String:         "1+1",
							Status:         calc.StatusBacklog,
						},
					},
				},
			},
		},
		{
			name:       "2 expressions",
			expression: "6-5",
			expected_queue: []orchestrator.Expression{
				{
					Id:     "0",
					Status: calc.StatusBacklog,
					String: "1+1",
					TasksQueue: []calc.Example{
						{
							Id:             "0_0",
							FirstArgument:  calc.Argument{Value: 1},
							SecondArgument: calc.Argument{Value: 1},
							Operation:      calc.Plus,
							String:         "1+1",
							Status:         calc.StatusBacklog,
						},
					},
				},
				{
					Id:     "1",
					Status: calc.StatusBacklog,
					String: "6-5",
					TasksQueue: []calc.Example{
						{
							Id:             "1_0",
							FirstArgument:  calc.Argument{Value: 6},
							SecondArgument: calc.Argument{Value: 5},
							Operation:      calc.Minus,
							String:         "6-5",
							Status:         calc.StatusBacklog,
						},
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, gotErr := orchestrator.AddExpression(test.expression)
			if !errors.Is(gotErr, test.expected_err) {
				t.Error("expected error:", test.expected_err, "but expected:", gotErr)
			}

			gotQueue := orchestrator.GetExpressionsQueue()
			if !reflect.DeepEqual(test.expected_queue, gotQueue) {
				t.Error("got:", gotQueue, "but expected:", test.expected_queue)
			}
		})
	}
}

func TestSetExampleResult(t *testing.T) {
	tests := []struct {
		name          string
		expression    string
		result        float64
		example_id    string
		expected_task calc.Example
		expected_err  error
	}{
		{
			name:       "1+1",
			expression: "1+1",
			result:     2,
			example_id: "0_0",
			expected_task: calc.Example{
				Id:             "0_0",
				FirstArgument:  calc.Argument{Value: 1},
				SecondArgument: calc.Argument{Value: 1},
				Operation:      calc.Plus,
				String:         "1+1",
				Answer:         2,
				Status:         calc.StatusComplete,
			},
		},
		{
			name:       "1+1+1",
			expression: "1+1+1",
			result:     2,
			example_id: "1_0",
			expected_task: calc.Example{
				Id:             "1_0",
				FirstArgument:  calc.Argument{Value: 1},
				SecondArgument: calc.Argument{Value: 1},
				Operation:      calc.Plus,
				String:         "1+1",
				Answer:         2,
				Status:         calc.StatusComplete,
			},
		},
		{
			name:       "previous id:1_0+1",
			expression: "1+1+1",
			example_id: "1_1",
			result:     3,
			expected_task: calc.Example{
				Id:             "1_1",
				FirstArgument:  calc.Argument{Value: 2},
				SecondArgument: calc.Argument{Value: 1},
				Operation:      calc.Plus,
				String:         "id:1_0+1",
				Answer:         3,
				Status:         calc.StatusComplete,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			orchestrator.AddExpression(test.expression)
			gotErr := orchestrator.SetExampleResult(test.example_id, test.result)
			if !errors.Is(gotErr, test.expected_err) {
				t.Error("expected error:", test.expected_err, "but got:", gotErr)
			}

			gotTask, err := orchestrator.GetTask(test.example_id)
			if err != nil {
				t.Error(err)
			}

			if !reflect.DeepEqual(test.expected_task, gotTask) {
				t.Error("got:", gotTask, "but expected:", test.expected_task)
			}
		})
	}
}

func TestGetTask(t *testing.T) {
	tests := []struct {
		name          string
		expression    string
		id            string
		result        float64
		expected_task calc.Example
		expected_err  error
	}{
		{
			name:       "2*4 complete",
			expression: "2*4",
			id:         "",
			result:     8,
			expected_task: calc.Example{
				Id:             "0_0",
				FirstArgument:  calc.Argument{Value: 2},
				SecondArgument: calc.Argument{Value: 4},
				Operation:      calc.Multiply,
				Status:         calc.StatusBacklog,
				String:         "2*4",
			},
		},
		{
			name:       "1+1+1 in progress",
			expression: "1+1+1",
			id:         "",
			result:     2,
			expected_task: calc.Example{
				Id:             "1_0",
				FirstArgument:  calc.Argument{Value: 1},
				SecondArgument: calc.Argument{Value: 1},
				Operation:      calc.Plus,
				Status:         calc.StatusBacklog,
				String:         "1+1",
			},
		},
		{
			name:       "1+1+1 complete",
			expression: "1+1+1",
			id:         "",
			result:     3,
			expected_task: calc.Example{
				Id:             "1_1",
				FirstArgument:  calc.Argument{Value: 1},
				SecondArgument: calc.Argument{Expected: "1_0"},
				Operation:      calc.Plus,
				Status:         calc.StatusIsWaitingValues,
				String:         "1+id:1_0",
			},
		},
		{
			name:       "find by id",
			expression: "1-1",
			id:         "0_0",
			expected_task: calc.Example{
				Id:             "0_0",
				FirstArgument:  calc.Argument{Value: 2},
				SecondArgument: calc.Argument{Value: 4},
				Operation:      calc.Multiply,
				Status:         calc.StatusBacklog,
				String:         "2*4",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, gotErr := orchestrator.AddExpression(test.expression)
			if !errors.Is(gotErr, test.expected_err) {
				t.Error("expected error:", test.expected_err, "but got:", gotErr)
			}

			gotTask, gotErr := orchestrator.GetTask(test.id)
			if !errors.Is(gotErr, test.expected_err) {
				t.Error("expected error:", test.expected_err, "but got:", gotErr)
			}

			if !reflect.DeepEqual(test.expected_task, gotTask) {
				t.Error("got:", gotTask, "but expected:", test.expected_task)
			}

			orchestrator.SetExampleResult(test.expected_task.Id, test.result)
		})
	}
}
