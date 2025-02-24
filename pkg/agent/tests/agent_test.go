package agent_test

import (
	"errors"
	"testing"

	"github.com/Antibrag/gcalc-server/pkg/agent"
	"github.com/Antibrag/gcalc-server/pkg/calc"
)

func TestSolveExample(t *testing.T) {
	cases := []struct {
		name         string
		example      calc.Example
		expected     float64
		expected_err error
	}{
		{
			name:         "1 + 1",
			example:      calc.Example{FirstArgument: calc.Argument{Value: 1}, SecondArgument: calc.Argument{Value: 1}, Operation: calc.Plus},
			expected:     2,
			expected_err: nil,
		},
		{
			name:         "1 * 1",
			example:      calc.Example{FirstArgument: calc.Argument{Value: 1}, SecondArgument: calc.Argument{Value: 1}, Operation: calc.Multiply},
			expected:     1,
			expected_err: nil,
		},
		{
			name:         "divide by zero",
			example:      calc.Example{FirstArgument: calc.Argument{Value: 0}, SecondArgument: calc.Argument{Value: 0}, Operation: calc.Division},
			expected:     0,
			expected_err: calc.ErrDivideByZero,
		},
		{
			name:         "unknown operator",
			example:      calc.Example{FirstArgument: calc.Argument{Value: 1}, SecondArgument: calc.Argument{Value: 0}, Operation: '&'},
			expected:     0,
			expected_err: calc.ErrExpressionIncorrect,
		},
		{
			name:         "123 + 10",
			example:      calc.Example{FirstArgument: calc.Argument{Value: 123}, SecondArgument: calc.Argument{Value: 10}, Operation: calc.Plus},
			expected:     133,
			expected_err: nil,
		},
		{
			name:         "equal(1)",
			example:      calc.Example{FirstArgument: calc.Argument{Value: 1}, Operation: calc.Equals},
			expected:     1,
			expected_err: nil,
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			got, err := agent.SolveExample(test.example)
			if got != test.expected || !errors.Is(err, test.expected_err) {
				t.Errorf("SolveExample(%#v) = (%f, %q), but expected: (%f, %q)", test.example, got, err, test.expected, test.expected_err)
			}
		})
	}
}
