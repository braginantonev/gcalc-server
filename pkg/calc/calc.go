package calc

import (
	"fmt"
)

type Status string

const (
	StatusBacklog         Status = "backlog"
	StatusIsWaitingValues Status = "in waiting values"
	StatusInProgress      Status = "in progress"
	StatusComplete        Status = "complete"
)

type Argument struct {
	Value    float64
	Expected int // Id example; wait while status - not complete
}

type Example struct {
	Id             int
	FirstArgument  Argument
	SecondArgument Argument
	Operation      Operator
	Status         Status
	Answer         float64
}

func (ex Example) ToString() string {
	return fmt.Sprint(ex.FirstArgument.Value, ex.Operation, ex.SecondArgument.Value)
}

func Calc(expression string) (result float64, err error) {
	if expression == "" {
		return 0, ErrExpressionEmpty
	}

	// for {
	// 	ex_str, pri_idx, example, err := GetExample(expression)
	// 	if err != nil {
	// 		return 0, err
	// 	}

	// 	result, err = SolveExample(example)
	// 	if err != nil {
	// 		return 0, err
	// 	}

	// 	if ex_str == "end" {
	// 		break
	// 	}

	// 	expression = EraseExample(expression, ex_str, pri_idx, result)
	// }
	return
}
