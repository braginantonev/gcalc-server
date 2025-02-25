package orchestrator

import "errors"

var (
	ErrEOQ                error = errors.New("end of queue")
	ErrExpressionNotFound error = errors.New("expression not found")
	ErrTaskNotFound       error = errors.New("task not found")
	ErrExpectation        error = errors.New("expectation error")
)
