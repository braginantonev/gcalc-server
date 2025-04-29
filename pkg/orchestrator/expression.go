package orchestrator

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/braginantonev/gcalc-server/pkg/calc"
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

	for {
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
}

// Return id expression and error
func AddExpression(expression string) (string, error) {
	if expression == "" {
		return "", ErrExpressionEmpty
	}

	ex := Expression{
		//!!! Если придётся реализовывать удаление выражения, то нужно изменить систему выдачи индекса !!!
		//!!! При удалении элемента, длина уменьшается, следовательно следующее добавленное выражение, будет иметь такой же индекс, что и предпоследний !!!
		Id:     fmt.Sprint(len(expressionsQueue)),
		Status: calc.StatusAnalyze,
		String: expression,
	}

	if err := ex.setTasksQueue(); err != nil {
		return "", err
	}

	expressionsQueue = append(expressionsQueue, ex)
	return ex.Id, nil
}

func GetExpression(id string) (Expression, error) {
	if id == "" {
		for _, ex := range expressionsQueue {
			if ex.Status == calc.StatusInProgress || ex.Status == calc.StatusBacklog {
				return ex, nil
			}
		}
		return Expression{}, DHT
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

func GetTask(id string) (calc.Example, error) {
	if id == "" {
		expression_local, err := GetExpression("")
		if err != nil {
			return calc.Example{}, DHT
		}

		expId, err := strconv.Atoi(expression_local.Id)
		if err != nil {
			return calc.Example{}, err
		}

		p_expression := &expressionsQueue[expId]
		for i := range p_expression.TasksQueue {
			p_example := &p_expression.TasksQueue[i]
			if p_example.Status == calc.StatusBacklog {
				p_example.Status = calc.StatusInProgress
				p_expression.Status = calc.StatusInProgress
				return *p_example, nil
			}
		}
		return calc.Example{}, DHT
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

func SetExampleResult(id string, result float64) error {
	example_local, err := GetTask(id)
	if err != nil {
		return err
	}

	if example_local.Status == calc.StatusComplete {
		fmt.Println("example", id, "already complete")
		return nil
	}

	low_line_idx := strings.IndexRune(example_local.Id, '_')
	exp_local, err := GetExpression(example_local.Id[:low_line_idx])
	if err != nil {
		return err
	}

	expressionId_int, err := strconv.Atoi(exp_local.Id)
	if err != nil {
		return err
	}

	exampleId_int, err := strconv.Atoi(example_local.Id[low_line_idx+1:])
	if err != nil {
		return err
	}

	p_expression := &expressionsQueue[expressionId_int]
	p_example := &p_expression.TasksQueue[exampleId_int]

	p_example.Answer = result
	p_example.Status = calc.StatusComplete

	if exampleId_int == len(p_expression.TasksQueue)-1 {
		p_expression.Result = result
		p_expression.Status = calc.StatusComplete
		return nil
	}

	// Return true, if example result expected
	delExpectation := func(arg *calc.Argument) {
		if arg.Expected == p_example.Id {
			arg.Value = result
			arg.Expected = ""
		}
	}

	for i := range p_expression.TasksQueue {
		p_local_example := &p_expression.TasksQueue[i]
		if p_example.Id == p_local_example.Id {
			continue
		}

		delExpectation(&p_local_example.FirstArgument)
		delExpectation(&p_local_example.SecondArgument)

		if p_local_example.FirstArgument.Expected == "" && p_local_example.SecondArgument.Expected == "" {
			p_local_example.Status = calc.StatusBacklog
		}
	}

	return nil
}
