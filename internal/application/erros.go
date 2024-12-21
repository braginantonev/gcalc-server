package application

import (
	"errors"

	"github.com/Antibrag/gcalc-server/pkg/calc"
)

var (
	InternalError       error = errors.New("Internal error")
	RequestBodyEmpty    error = errors.New("Request body empty")
	UnsupportedBodyType error = errors.New("Unsupported request body type")

	CalculatorErrors []*error = []*error{
		&calc.DivideByZero,
		&calc.ExpressionEmpty,
		&calc.OperationWithoutValue,
		&calc.BracketsNotFound,
		&calc.ExpressionIncorrect,
	}
)
