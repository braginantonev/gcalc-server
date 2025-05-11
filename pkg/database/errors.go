package database

import "errors"

var (
	ErrDBNotInit             error = errors.New("database not initialized")
	ErrDBPathIsEmpty         error = errors.New("database path is empty")
	ErrUnexpectedRequestType error = errors.New("unexpected database request type")
	ErrBadArguments          error = errors.New("bad arguments in database request")
	ErrCacheUpdateFailed     error = errors.New("update cache value failed. key not found")
)
