package calc

type Status string

const (
	StatusBacklog         Status = "backlog"
	StatusIsWaitingValues Status = "in waiting values"
	StatusInProgress      Status = "in progress"
	StatusComplete        Status = "complete"
)

type Argument struct {
	Value    float64
	Expected string // Id example; wait while status - not complete
}

type Example struct {
	Id             string
	FirstArgument  Argument
	SecondArgument Argument
	Operation      Operator
	Status         Status
	String         string
	Answer         float64
	AnswerChannel  chan float64
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
