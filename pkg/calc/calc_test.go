package calc_test

import (
	"testing"

	"github.com/Antibrag/gcalc-server/pkg/calc"
)

func TestCalculate(t *testing.T) {
	cases := []struct {
		name     string
		example  calc.Example
		expected float64
		err_str  string
	}{
		{
			name:     "1 + 1",
			example:  calc.Example{First_value: 1, Second_value: 1, Operation: calc.Plus},
			expected: 2,
			err_str:  "",
		},
		{
			name:     "1 * 1",
			example:  calc.Example{First_value: 1, Second_value: 1, Operation: calc.Multiply},
			expected: 1,
			err_str:  "",
		},
		{
			name:     "divide by zero",
			example:  calc.Example{First_value: 0, Second_value: 0, Operation: calc.Division},
			expected: 0,
			err_str:  calc.Err_DivideByZero,
		},
		{
			name:     "unkown operator",
			example:  calc.Example{First_value: 1, Second_value: 1, Operation: '&'},
			expected: 0,
			err_str:  calc.Err_UnkownOperator,
		},
		{
			name:     "123 + 10",
			example:  calc.Example{First_value: 123, Second_value: 10, Operation: calc.Plus},
			expected: 133,
			err_str:  "",
		},
		{
			name:     "equal(1)",
			example:  calc.Example{First_value: 1, Second_value: 52, Operation: calc.Equals},
			expected: 1,
			err_str:  "",
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			got, err := calc.Calculate(test.example)
			if got != test.expected || err != test.err_str {
				t.Errorf("Calculate(%#v) = (%f, %q), but expected: (%f, %q)", test.example, got, err, test.expected, test.err_str)
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
			value, _ := calc.Calculate(ex)
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
			answ, _ := calc.Calculate(ex)
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
	}{
		{
			name:           "1+1",
			expression:     "1+1",
			expected_value: 2,
		},
		{
			name:           "-3+1",
			expression:     "1+1",
			expected_value: -2,
		},
		{
			name:           "1+1+1",
			expression:     "1+1+1",
			expected_value: 3,
		},
		{
			name:           "1*1",
			expression:     "1*1",
			expected_value: 1,
		},
		{
			name:           "1/1",
			expression:     "1/1",
			expected_value: 1,
		},
		{
			name:           "2+1*1",
			expression:     "2+1*1",
			expected_value: 3,
		},
		{
			name:           "2+1*1+10/2",
			expression:     "2+1*1+10/2",
			expected_value: 8,
		},
		{
			name:           "(1+1)/(1+1)",
			expression:     "(1+1)/(1+1)",
			expected_value: 1,
		},
		{
			name:           "(1+10*(23-3)/2)-12",
			expression:     "(1+10*(23-3)/2)-12",
			expected_value: 89,
		},
		{
			name:           "1&1",
			expression:     "1&1",
			expected_value: 0,
		},
		{
			name:           "1+&1",
			expression:     "1+&1",
			expected_value: 0,
		},
		{
			name:           "1+1*",
			expression:     "1+1*",
			expected_value: 0,
		},
		{
			name:           "2+2**2",
			expression:     "2+2**2",
			expected_value: 0,
		},
		{
			name:           "((2+2-*(2",
			expression:     "((2+2-*(2",
			expected_value: 0,
		},
		{
			name:           "nothing",
			expression:     "",
			expected_value: 0,
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			got, err := calc.Calc(test.expression)

			if err != nil {
				t.Error(err)
			}

			if got != test.expected_value {
				t.Errorf("Calc(%q) = %f, but expected - %f", test.expression, got, test.expected_value)
			}
		})
	}
}
