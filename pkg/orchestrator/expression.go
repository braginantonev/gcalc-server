package orchestrator

import (
	"fmt"
	"strings"

	"github.com/Antibrag/gcalc-server/pkg/calc"
)

//! Оставить как минимум один поток для выполнения анализа

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

		if example.Status != calc.StatusIsWaitingValues {
			example.Status = calc.StatusBacklog
		}

		expression.TasksQueue = append(expression.TasksQueue, example)

		expressionStr = EraseExample(expressionStr, example.String, priority_idx, example.Id)
	}

	return nil
}

func AddExpression(expression string) error {
	if expression == "" {
		return ErrExpressionEmpty
	}

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
			if ex.Status == calc.StatusInProgress || ex.Status == calc.StatusBacklog {
				return ex, nil
			}
		}
		return Expression{}, nil
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

// TODO: Написать тесты
func GetTask(id string) (calc.Example, error) {
	if id == "" {
		exp, err := GetExpression("")
		if err != nil {
			return calc.Example{}, err
		}

		for _, example := range exp.TasksQueue {
			if example.Status == calc.StatusBacklog {
				example.Status = calc.StatusInProgress
				return example, nil
			}
		}
		return calc.Example{}, nil
	}

	for _, exp := range expressionsQueue {
		for _, example := range exp.TasksQueue {
			if example.Id == id {
				return example, nil
			}
		}
	}

	return calc.Example{}, ErrTaskNotFound
}

// TODO: Добавить тесты
func SetExampleResult(id string, result float64) error {
	example, err := GetTask(id)
	if err != nil {
		return err
	}

	example.Answer = result
	example.Status = calc.StatusComplete

	low_line_idx := strings.IndexRune(example.Id, '_')
	exp, err := GetExpression(example.Id[:low_line_idx])
	if err != nil {
		return err
	}

	var exampleIdx int
	for i, example := range exp.TasksQueue {
		if example.Id == id {
			exampleIdx = i
			break
		}
	}

	if exampleIdx == len(exp.TasksQueue)-1 {
		exp.Result = result
		exp.Status = calc.StatusComplete
		exp.TasksQueue[exampleIdx] = example
		return nil
	}

	// Return true - if argument excepted value
	delExpectation := func(arg *calc.Argument) bool {
		if arg.Expected == example.Id {
			arg.Value = result
			arg.Expected = ""
			exp.Status = calc.StatusBacklog
			return true
		}
		return false
	}

	if isExpected := delExpectation(&exp.TasksQueue[exampleIdx+1].FirstArgument); !isExpected {
		if isExpected = delExpectation(&exp.TasksQueue[exampleIdx+1].SecondArgument); !isExpected {
			return ErrExpectation
		}
	}

	exp.TasksQueue[exampleIdx] = example

	return nil
}

// TODO: Написать функцию принятия результата вычисления задания
// TODO: При получении результата сменить статус примера, который ожидает данные, на StatusBacklog
// TODO: Также проверяется, если все задачи примера выполнены, то примеру выставляется StatusComplete
