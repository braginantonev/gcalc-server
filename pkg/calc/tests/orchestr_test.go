package calc_test

import (
	"testing"

	"github.com/Antibrag/gcalc-server/pkg/calc"
)

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
			_, gotExample, _ := calc.GetExample(test.example)
			value, _ := calc.SolveExample(gotExample)
			if gotExample.String != test.expected_str || value != test.expected_value {
				t.Logf("GetExample(%q).ex = %#v", test.example, gotExample)
				t.Errorf("GetExample(%q) = (%q, %f), but expected: (%q, %f)", test.example, gotExample.String, value, test.expected_str, test.expected_value)
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
			pri_idx, gotExample, _ := calc.GetExample(test.example)
			//answ, _ := calc.SolveExample(ex)
			got := calc.EraseExample(test.example, gotExample.String, pri_idx, gotExample.Id)
			if got != test.expected_str {
				t.Errorf("EraseExample(%q, %q, %d, %d) = %q, but expected: %q", test.example, gotExample.String, pri_idx, gotExample.Id, got, test.expected_str)
			}
		})
	}
}
