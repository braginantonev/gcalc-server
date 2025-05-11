package application

import (
	"errors"

	"github.com/braginantonev/gcalc-server/pkg/agent"
	"github.com/braginantonev/gcalc-server/pkg/orchestrator"
)

var (
	ErrInternalError       error = errors.New("internal error")
	ErrRequestBodyEmpty    error = errors.New("request body empty")
	ErrUnsupportedBodyType error = errors.New("unsupported request body type")
	ErrJWTTokenNotValid    error = errors.New("jwt token not valid -> relogin")

	OrchestratorErrors []*error = []*error{
		&orchestrator.ErrExpressionEmpty,
		&orchestrator.ErrOperationWithoutValue,
		&orchestrator.ErrBracketsNotFound,
		&orchestrator.ErrExpressionNotFound,
		&orchestrator.ErrTaskNotFound,
		&orchestrator.ErrExpectation,
		&orchestrator.ErrExpressionIncorrect,
	}

	AgentErrors []*error = []*error{
		&orchestrator.ErrExpressionIncorrect,
		&agent.ErrDivideByZero,
	}
)
