package orchestrator

func NewTask() *Task {
	return &Task{
		FirstArgument:  &Argument{Expected: -1},
		SecondArgument: &Argument{Expected: -1},
	}
}

func NewTaskResult() *TaskResult {
	return &TaskResult{
		TaskID: NewTaskID(),
	}
}

func NewTaskID() *TaskID {
	return &TaskID{
		Expression: -1,
		Internal:   -1,
	}
}

func NewTaskIDWithValues(expression_id, task_id int32) *TaskID {
	return &TaskID{
		Expression: expression_id,
		Internal:   task_id,
	}
}
