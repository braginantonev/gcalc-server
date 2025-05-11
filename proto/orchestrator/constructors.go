package orchestrator

func NewTask() *Task {
	return &Task{
		Id:             NewTaskID(),
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
		Expression: NewExpressionID(),
		Internal:   -1,
	}
}

func NewTaskIDWithValues(user string, expression_id, task_id int32) *TaskID {
	return &TaskID{
		Expression: NewExpressionIDWithValues(user, expression_id),
		Internal:   task_id,
	}
}

func NewExpressionID() *ExpressionID {
	return &ExpressionID{
		Internal: -1,
	}
}

func NewExpressionIDWithValues(user string, id int32) *ExpressionID {
	return &ExpressionID{
		User:     user,
		Internal: id,
	}
}
