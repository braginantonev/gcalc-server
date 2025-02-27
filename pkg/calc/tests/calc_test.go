package calc_test

import (
	"errors"
	"testing"

	"github.com/Antibrag/gcalc-server/pkg/calc"
)

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
		// {
		// 	name:           "addition with negative value",
		// 	expression:     "-3+1",
		// 	expected_value: -2,
		// 	expected_err:   nil,
		// },
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
			expected_err:   calc.ErrExpressionIncorrect,
		},
		{
			name:           "unkown operator 2",
			expression:     "1+&1",
			expected_value: 0,
			expected_err:   calc.ErrExpressionIncorrect,
		},
		{
			name:           "operation without value",
			expression:     "1+1*",
			expected_value: 0,
			expected_err:   calc.ErrOperationWithoutValue,
		},
		{
			name:           "operation without value 2",
			expression:     "2+2**2",
			expected_value: 0,
			expected_err:   calc.ErrExpressionIncorrect,
		},
		{
			name:           "without closed bracket",
			expression:     "((2+2-*(2",
			expected_value: 0,
			expected_err:   calc.ErrBracketsNotFound,
		},
		{
			name:           "nothing",
			expression:     "",
			expected_value: 0,
			expected_err:   calc.ErrExpressionEmpty,
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
