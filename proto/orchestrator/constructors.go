package orchestrator

func NewTask() *Task {
	return &Task{
		FirstArgument:  &Argument{},
		SecondArgument: &Argument{},
	}
}
