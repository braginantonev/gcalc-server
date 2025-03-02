package calc

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
	Id             string
	FirstArgument  Argument
	SecondArgument Argument
	Operation      Operator
	Status         Status
	String         string
	Answer         float64
	//AnswerChannel  chan float64
}
