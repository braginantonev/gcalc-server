package agent

import (
	"encoding/json"
	"net/http"

	"github.com/Antibrag/gcalc-server/pkg/calc"
)

const (
	TIME_ADDITION_MS       = 100
	TIME_SUBTRACTION_MS    = 150
	TIME_MULTIPLICATION_MS = 300
	TIME_DIVISIONS_MS      = 350
	COMPUTING_POWER        = 1
)

func SolveExpressions() error {
	resp, err := http.Get("localhost/internal/task")
	if err != nil {
		return err
	}

	resp_json := make([]byte, 1024)
	n, err := resp.Body.Read(resp_json)
	if err != nil {
		return err
	}
	resp_json = resp_json[n:]

	var example calc.Example
	if err = json.Unmarshal(resp_json, example); err != nil {
		return err
	}


}

func SolveExample(ex calc.Example) (float64, error) {
	if ex.SecondArgument.Value == 0 && ex.Operation == calc.Division {
		return 0, calc.ErrDivideByZero
	}

	switch ex.Operation {
	case calc.Plus:
		return ex.FirstArgument.Value + ex.SecondArgument.Value, nil
	case calc.Minus:
		return ex.FirstArgument.Value - ex.SecondArgument.Value, nil
	case calc.Multiply:
		return ex.FirstArgument.Value * ex.SecondArgument.Value, nil
	case calc.Division:
		return ex.FirstArgument.Value / ex.SecondArgument.Value, nil
	case calc.Equals:
		return ex.FirstArgument.Value, nil
	}
	return 0, calc.ErrExpressionIncorrect
}
