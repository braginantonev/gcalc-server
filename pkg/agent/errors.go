package agent

import "errors"

var (
	ErrDivideByZero error = errors.New("divide by zero")
	DHT             error = errors.New("don't have task")
)
