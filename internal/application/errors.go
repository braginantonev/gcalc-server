package application

import (
	"errors"

	"github.com/Antibrag/gcalc-server/pkg/agent"
	"github.com/Antibrag/gcalc-server/pkg/calc"
	"github.com/Antibrag/gcalc-server/pkg/orchestrator"
)

var (
	ErrInternalError       error = errors.New("internal error")
	ErrRequestBodyEmpty    error = errors.New("request body empty")
	ErrUnsupportedBodyType error = errors.New("unsupported request body type")

	OrchestratorErrors []*error = []*error{
		&orchestrator.ErrExpressionEmpty,
		&orchestrator.ErrOperationWithoutValue,
		&orchestrator.ErrBracketsNotFound,
		&orchestrator.ErrExpressionNotFound,
		&orchestrator.ErrTaskNotFound,
		&orchestrator.ErrExpectation,
		&calc.ErrExpressionIncorrect,
	}

	AgentErrors []*error = []*error{
		&calc.ErrExpressionIncorrect,
		&agent.ErrDivideByZero,
	}
)
