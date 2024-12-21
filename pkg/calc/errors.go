package calc

import "errors"

var (
	DivideByZero          error = errors.New("divide by zero")
	ExpressionEmpty       error = errors.New("expression empty")
	OperationWithoutValue error = errors.New("operation don't have a value")
	BracketsNotFound      error = errors.New("not found opened or closed bracket")
	ExpressionIncorrect   error = errors.New("expression incorrect")
)
