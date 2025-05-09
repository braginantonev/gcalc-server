package database

import "errors"

var (
	ErrDBPathIsEmpty         error = errors.New("database path is empty")
	ErrUnexpectedRequestType error = errors.New("unexpected database request type")
	ErrBadArguments          error = errors.New("bad arguments in database request")
)
