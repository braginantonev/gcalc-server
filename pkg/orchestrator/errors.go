package orchestrator

import "errors"

var (
	ErrExpressionEmpty       error = errors.New("expression empty")
	ErrOperationWithoutValue error = errors.New("operation don't have a value")
	ErrBracketsNotFound      error = errors.New("not found opened or closed bracket")

	DHT error = errors.New("don't have task")

	ErrExpressionNotFound  error = errors.New("expression not found")
	ErrTaskNotFound        error = errors.New("task not found")
	ErrExpectation         error = errors.New("expectation error")
	ErrExpressionIncorrect error = errors.New("expression incorrect")
)
