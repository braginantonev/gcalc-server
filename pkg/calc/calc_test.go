package calc_test

import (
	"errors"
	"testing"

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
			example:      calc.Example{First_value: 1, Second_value: 1, Operation: calc.Plus},
			expected:     2,
			expected_err: nil,
		},
		{
			name:         "1 * 1",
			example:      calc.Example{First_value: 1, Second_value: 1, Operation: calc.Multiply},
			expected:     1,
			expected_err: nil,
		},
		{
			name:         "divide by zero",
			example:      calc.Example{First_value: 0, Second_value: 0, Operation: calc.Division},
			expected:     0,
			expected_err: calc.DivideByZero,
		},
		{
			name:         "unknown operator",
			example:      calc.Example{First_value: 1, Second_value: 1, Operation: '&'},
			expected:     0,
			expected_err: calc.ParseError,
		},
		{
			name:         "123 + 10",
			example:      calc.Example{First_value: 123, Second_value: 10, Operation: calc.Plus},
			expected:     133,
			expected_err: nil,
		},
		{
			name:         "equal(1)",
			example:      calc.Example{First_value: 1, Second_value: 52, Operation: calc.Equals},
			expected:     1,
			expected_err: nil,
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			got, err := calc.SolveExample(test.example)
			if got != test.expected || !errors.Is(err, test.expected_err) {
				t.Errorf("SolveExample(%#v) = (%f, %q), but expected: (%f, %q)", test.example, got, err, test.expected, test.expected_err)
			}
		})
	}
}

func TestGetExample(t *testing.T) {
	cases := []struct {
		name           string
		example        string
		expected_str   string
		expected_value float64
		expected_err   bool
	}{
		{
			name:           "1+1",
			example:        "1+1",
			expected_str:   "1+1",
			expected_value: 2,
		},
		{
			name:           "1+145",
			example:        "1+145",
			expected_str:   "1+145",
			expected_value: 146,
		},
		{
			name:           "10+10+10+10",
			example:        "10+10",
			expected_str:   "10+10",
			expected_value: 20,
		},
		{
			name:           "1+145*10",
			example:        "1+145*10",
			expected_str:   "145*10",
			expected_value: 1450,
		},
		{
			name:           "(1)",
			example:        "(1)",
			expected_str:   "(1)",
			expected_value: 1,
		},
		{
			name:           "(112+1)+145*10",
			example:        "(112+1)+145*10",
			expected_str:   "112+1",
			expected_value: 113,
		},
		{
			name:           "(10+30+(20+1*10))+1",
			example:        "(10+30+(20+1*10))+1",
			expected_str:   "1*10",
			expected_value: 10,
		},
		{
			name:           "1+1-(2.000000)",
			example:        "1+1-(2.000000)",
			expected_str:   "(2.000000)",
			expected_value: 2,
		},
		{
			name:           "52",
			example:        "52",
			expected_str:   "end",
			expected_value: 52,
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			got, _, ex, _ := calc.GetExample(test.example)
			value, _ := calc.SolveExample(ex)
			if got != test.expected_str || value != test.expected_value {
				t.Logf("GetExample(%q).ex = %#v", test.example, ex)
				t.Errorf("GetExample(%q) = (%q, %f), but expected: (%q, %f)", test.example, got, value, test.expected_str, test.expected_value)
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
			expected_str: "1+1-(2.000000)",
		},
		{
			name:         "1+1-(1+1)-1+1",
			example:      "1+1-(1+1)-1+1",
			expected_str: "1+1-(2.000000)-1+1",
		},
		{
			name:         "(1)",
			example:      "(1)",
			expected_str: "1.000000",
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			erase_str, pri_idx, ex, _ := calc.GetExample(test.example)
			answ, _ := calc.SolveExample(ex)
			got := calc.EraseExample(test.example, erase_str, pri_idx, answ)
			if got != test.expected_str {
				t.Errorf("EraseExample(%q, %q, %d, %f) = %q, but expected: %q", test.example, erase_str, pri_idx, answ, got, test.expected_str)
			}
		})
	}
}

func TestCalc(t *testing.T) {
	cases := []struct {
		name           string
		expression     string
		expected_value float64
		expected_err   error
	}{
		{
			name:           "simple addition",
			expression:     "1+1",
			expected_value: 2,
			expected_err:   nil,
		},
		{
			name:           "addition with negative value",
			expression:     "-3+1",
			expected_value: -2,
			expected_err:   nil,
		},
		{
			name:           "addition with 3 values",
			expression:     "1+1+1",
			expected_value: 3,
			expected_err:   nil,
		},
		{
			name:           "simple multiply",
			expression:     "1*1",
			expected_value: 1,
			expected_err:   nil,
		},
		{
			name:           "simple division",
			expression:     "1/1",
			expected_value: 1,
			expected_err:   nil,
		},
		{
			name:           "division with addition",
			expression:     "2+1*1",
			expected_value: 3,
			expected_err:   nil,
		},
		{
			name:           "hard example 1",
			expression:     "2+1*1+10/2",
			expected_value: 8,
			expected_err:   nil,
		},
		{
			name:           "brackets",
			expression:     "(1+1)/(1+1)",
			expected_value: 1,
			expected_err:   nil,
		},
		{
			name:           "hard example with brackets",
			expression:     "(1+10*(23-3)/2)-12",
			expected_value: 89,
			expected_err:   nil,
		},
		{
			name:           "unkown operator 1",
			expression:     "1&1",
			expected_value: 0,
			expected_err:   calc.ParseError,
		},
		{
			name:           "unkown operator 2",
			expression:     "1+&1",
			expected_value: 0,
			expected_err:   calc.ParseError,
		},
		{
			name:           "operation without value",
			expression:     "1+1*",
			expected_value: 0,
			expected_err:   calc.OperationWithoutValue,
		},
		{
			name:           "operation without value 2",
			expression:     "2+2**2",
			expected_value: 0,
			expected_err:   calc.ParseError,
		},
		{
			name:           "without closed bracket",
			expression:     "((2+2-*(2",
			expected_value: 0,
			expected_err:   calc.BracketsNotFound,
		},
		{
			name:           "nothing",
			expression:     "",
			expected_value: 0,
			expected_err:   calc.ExpressionEmpty,
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			got, err := calc.Calc(test.expression)

			if !errors.Is(err, test.expected_err) {
				t.Errorf("Calc(%s) got %q, but expected %q", test.expression, err, test.expected_err)
				return
			}

			if got != test.expected_value {
				t.Errorf("Calc(%q) = %f, but expected - %f", test.expression, got, test.expected_value)
			}
		})
	}
}
