package calc

import (
	"fmt"
)

type Argument struct {
	Value      float64
	IsExpected bool
}

type Example struct {
	FirstArgument  Argument
	SecondArgument Argument
	Operation      Operator
}

func (ex Example) ToString() string {
	return fmt.Sprint(ex.FirstArgument.Value, ex.Operation, ex.SecondArgument.Value)
}

var ExamplesQueue []Example

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
