package orchestrator

import "errors"

var (
	ErrEOQ                error = errors.New("end of queue")
	ErrExpressionNotFound error = errors.New("expression not found")
)
