package application

import (
	"errors"

	"github.com/Antibrag/gcalc-server/pkg/calc"
)

var (
	ErrInternalError       error = errors.New("internal error")
	ErrRequestBodyEmpty    error = errors.New("request body empty")
	ErrUnsupportedBodyType error = errors.New("unsupported request body type")

	CalculatorErrors []*error = []*error{
		&calc.ErrDivideByZero,
		&calc.ErrExpressionEmpty,
		&calc.ErrOperationWithoutValue,
		&calc.ErrBracketsNotFound,
		&calc.ErrExpressionIncorrect,
	}
)
