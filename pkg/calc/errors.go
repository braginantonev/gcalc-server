package calc

import "errors"

var (
	ErrDivideByZero          error = errors.New("divide by zero")
	ErrExpressionEmpty       error = errors.New("expression empty")
	ErrOperationWithoutValue error = errors.New("operation don't have a value")
	ErrBracketsNotFound      error = errors.New("not found opened or closed bracket")
	ErrExpressionIncorrect   error = errors.New("expression incorrect")
)
