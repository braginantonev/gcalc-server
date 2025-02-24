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

func TestAddExpression(t *testing.T) {

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
			gotErr := orchestrator.AddExpression(test.expression)
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
