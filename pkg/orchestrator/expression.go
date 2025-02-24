package orchestrator

import (
	"fmt"

	"github.com/Antibrag/gcalc-server/pkg/calc"
)

const END_STR = "end"

var expressionsQueue []Expression

type Expression struct {
	Id         string      `json:"id"`
	Status     calc.Status `json:"status"`
	String     string
	TasksQueue []calc.Example
	Result     float64 `json:"result"`
}

func (expression *Expression) setTasksQueue() error {
	if expression.String == "" {
		return calc.ErrExpressionEmpty
	}

	expressionStr := expression.String

	for range 5 {
		example, priority_idx, err := GetExample(expressionStr)
		if err != nil {
			return err
		}

		if example.String == END_STR {
			expression.Status = calc.StatusBacklog
			return nil
		}

		if example.Operation == calc.Equals {
			expressionStr = EraseExample(expressionStr, example.String, priority_idx, expression.TasksQueue[len(expression.TasksQueue)-1].Id)
			continue
		}

		example.Id = expression.Id + "_" + fmt.Sprint(len(expression.TasksQueue))
		example.Status = calc.StatusBacklog
		expression.TasksQueue = append(expression.TasksQueue, example)

		expressionStr = EraseExample(expressionStr, example.String, priority_idx, example.Id)
	}

	return nil
}

func AddExpression(expression string) error {
	ex := Expression{
		Id:     fmt.Sprint(len(expressionsQueue)),
		Status: calc.StatusAnalyze,
		String: expression,
	}

	if err := ex.setTasksQueue(); err != nil {
		return err
	}

	expressionsQueue = append(expressionsQueue, ex)
	return nil
}

func GetExpression(id string) (Expression, error) {
	if id == "" {
		for _, ex := range expressionsQueue {
			if ex.Status == calc.StatusBacklog {
				return ex, nil
			}
		}
		return Expression{}, ErrEOQ
	}

	for _, ex := range expressionsQueue {
		if ex.Id == id {
			return ex, nil
		}
	}

	return Expression{}, ErrExpressionNotFound
}

func GetExpressionsQueue() []Expression {
	return expressionsQueue
}
