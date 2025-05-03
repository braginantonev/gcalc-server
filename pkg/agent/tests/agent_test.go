package agent_test

import (
	"errors"
	"testing"

	"github.com/braginantonev/gcalc-server/pkg/agent"
	"github.com/braginantonev/gcalc-server/pkg/orchestrator"
	pb "github.com/braginantonev/gcalc-server/proto/orchestrator"
)

func TestSolve(t *testing.T) {
	cases := []struct {
		name         string
		task         *pb.Task
		expected     float64
		expected_err error
	}{
		{
			name:         "1 + 1",
			task:         &pb.Task{FirstArgument: &pb.Argument{Value: 1}, SecondArgument: &pb.Argument{Value: 1}, Operation: orchestrator.Plus.ToString()},
			expected:     2,
			expected_err: nil,
		},
		{
			name:         "1 * 1",
			task:         &pb.Task{FirstArgument: &pb.Argument{Value: 1}, SecondArgument: &pb.Argument{Value: 1}, Operation: orchestrator.Multiply.ToString()},
			expected:     1,
			expected_err: nil,
		},
		{
			name:         "divide by zero",
			task:         &pb.Task{FirstArgument: &pb.Argument{Value: 0}, SecondArgument: &pb.Argument{Value: 0}, Operation: orchestrator.Division.ToString()},
			expected:     0,
			expected_err: agent.ErrDivideByZero,
		},
		{
			name:         "unknown operator",
			task:         &pb.Task{FirstArgument: &pb.Argument{Value: 1}, SecondArgument: &pb.Argument{Value: 0}, Operation: "&"},
			expected:     0,
			expected_err: orchestrator.ErrExpressionIncorrect,
		},
		{
			name:         "123 + 10",
			task:         &pb.Task{FirstArgument: &pb.Argument{Value: 123}, SecondArgument: &pb.Argument{Value: 10}, Operation: orchestrator.Plus.ToString()},
			expected:     133,
			expected_err: nil,
		},
		// {
		// 	name:         "equal(1)",
		// 	example:      calc.Example{FirstArgument: calc.Argument{Value: 1}, Operation: calc.Equals},
		// 	expected:     1,
		// 	expected_err: nil,
		// },
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			err := agent.Solve(0, test.task)
			if test.task.Answer != test.expected || !errors.Is(err, test.expected_err) {
				t.Errorf("SolveExample(%#v) = (%f, %q), but expected: (%f, %q)", test.task, test.task.Answer, err, test.expected, test.expected_err)
			}
		})
	}
}
