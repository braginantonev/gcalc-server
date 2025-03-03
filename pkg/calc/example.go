package calc

import "time"

type Status string

const (
	StatusAnalyze         Status = "analyze"
	StatusBacklog         Status = "backlog"
	StatusIsWaitingValues Status = "in waiting values"
	StatusInProgress      Status = "in progress"
	StatusComplete        Status = "complete"
)

type Operator rune

const (
	Plus     Operator = '+'
	Minus    Operator = '-'
	Multiply Operator = '*'
	Division Operator = '/'
	Equals   Operator = '='
)

type Argument struct {
	Value    float64
	Expected string // Id example; wait while status - not complete
}

type Example struct {
	Id             string `json:"id"`
	FirstArgument  Argument `json:"arg1"`
	SecondArgument Argument `json:"arg2"`
	Operation      Operator `json:"operation"`
	OperationTime  time.Duration `json:"operation_time"`
	Status         Status
	String         string
	Answer         float64
}
