package calc

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

//TODO: Изменить вывод ошибок strconv.Parse... в вывод ошибок через переменную

type Operator rune

const (
	Plus     Operator = '+'
	Minus    Operator = '-'
	Multiply Operator = '*'
	Division Operator = '/'
	Equals   Operator = '='
)

var (
	DivideByZero          error = errors.New("divide by zero")
	UnkownOperator        error = errors.New("unkown operator")
	ExpressionEmpty       error = errors.New("expression empty")
	OperationWithoutValue error = errors.New("operation dont have a value")
	BracketsNotFound      error = errors.New("not found opened or closed bracket")
)

type Example struct {
	First_value  float64
	Second_value float64
	Operation    Operator
}

func SolveExample(ex Example) (float64, error) {
	if ex.Second_value == 0 {
		return 0, DivideByZero
	}

	switch ex.Operation {
	case Plus:
		return ex.First_value + ex.Second_value, nil
	case Minus:
		return ex.First_value - ex.Second_value, nil
	case Multiply:
		return ex.First_value * ex.Second_value, nil
	case Division:
		return ex.First_value / ex.Second_value, nil
	case Equals:
		return ex.First_value, nil
	}
	return 0, UnkownOperator
}

// * Заменяет выражение на его ответ
// ! Реализация данной функции довольно таки херовая
// ! Для оптимизации можно будет попробовать превратить в такую строку где не требуется постоянная замена
func GetExample(example string) (string, int, Example, error) {
	var err error
	var ex Example
	var local_ex string = example

	var begin, end int = -1, -1
	if strings.ContainsRune(local_ex, '(') {
		for i, rn := range local_ex {
			if rn == '(' {
				begin = i
				continue
			} else if rn == ')' {
				end = i
				break
			}
		}

		if (begin == -1 && end != -1) || (begin != -1 && end == -1) {
			return "", 0, Example{}, BracketsNotFound
		}

		local_ex = local_ex[begin : end+1]
	}

	var actionIdx int
	if op := "*/"; strings.ContainsAny(local_ex, op) {
		actionIdx = strings.IndexAny(local_ex, op)
	} else if op := "+-"; strings.ContainsAny(local_ex, op) {
		actionIdx = strings.IndexAny(local_ex, op)
	} else if strings.ContainsAny(local_ex, "()") {
		var value float64
		value, err = strconv.ParseFloat(local_ex[1:len(local_ex)-1], 64)
		return local_ex[:], strings.IndexRune(example, rune(local_ex[0])), Example{First_value: value, Second_value: 52, Operation: Equals}, err //52 - по рофлу, чтобы при калькулировании не возникала ошибка. Крч костыль
	} else {
		var value float64
		value, err = strconv.ParseFloat(local_ex, 64)
		return "end", 0, Example{First_value: value, Second_value: 52, Operation: Equals}, err
	}

	if actionIdx == 0 || actionIdx == len(local_ex)-1 {
		return "", 0, Example{}, OperationWithoutValue
	}

	ex.Operation = Operator(local_ex[actionIdx])

	//Нахождение концов двух чисел
	var exampleLen = len(local_ex)
	if actionIdx == 0 || actionIdx == exampleLen-1 {

		//TODO: Изменить вывод ошибки на вывод с использованием переменной

		return "", 0, Example{}, errors.New("action in first or lst place")
	}

	for i := actionIdx - 1; i >= 0; i-- {
		if strings.ContainsRune("+-/*()", rune(local_ex[i])) {
			ex.First_value, err = strconv.ParseFloat(local_ex[i+1:actionIdx], 64)
			begin = i + 1
			break
		} else if i == 0 {
			ex.First_value, err = strconv.ParseFloat(local_ex[i:actionIdx], 64)
			begin = i
			break
		}
	}

	for i := actionIdx + 1; i < exampleLen; i++ {
		if strings.ContainsRune("+-/*()", rune(local_ex[i])) {
			ex.Second_value, err = strconv.ParseFloat(local_ex[actionIdx+1:i], 64)
			end = i
			break
		} else if i+1 == exampleLen {
			ex.Second_value, err = strconv.ParseFloat(local_ex[actionIdx+1:i+1], 64)
			end = exampleLen
			break
		}
	}

	if err != nil {
		return "", 0, Example{}, err
	}

	return local_ex[begin:end], strings.IndexRune(example, rune(local_ex[0])), ex, nil
}

func EraseExample(example, erase_ex string, pri_idx int, answ float64) string {
	return example[:pri_idx] + strings.Replace(example[pri_idx:], erase_ex, fmt.Sprintf("%f", answ), 1)
}

func Calc(expression string) (result float64, err error) {
	if expression == "" {
		return 0, ExpressionEmpty
	}

	for {
		ex_str, pri_idx, example, err := GetExample(expression)
		if err != nil {
			return 0, err
		}

		result, _ = SolveExample(example)

		if ex_str == "end" {
			break
		}

		expression = EraseExample(expression, ex_str, pri_idx, result)
	}
	return
}

func main() {}
