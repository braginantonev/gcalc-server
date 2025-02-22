package calc

import (
	"fmt"
)

type Example struct {
	First_value  float64
	Second_value float64
	Operation    Operator
}

func (ex Example) ToString() string {
	return fmt.Sprint(ex.First_value, ex.Operation, ex.Second_value)
}

func Calc(expression string) (result float64, err error) {
	if expression == "" {
		return 0, ErrExpressionEmpty
	}

	for {
		ex_str, pri_idx, example, err := GetExample(expression)
		if err != nil {
			return 0, err
		}

		result, err = SolveExample(example)
		if err != nil {
			return 0, err
		}

		if ex_str == "end" {
			break
		}

		expression = EraseExample(expression, ex_str, pri_idx, result)
	}
	return
}
