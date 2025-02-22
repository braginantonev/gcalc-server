package calc

import (
	"fmt"
	"strconv"
	"strings"
)

// Получает строку с выражением
func GetExample(example string) (string, int, Example, error) {
	var ex Example
	local_ex := example
	begin, end := -1, -1

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

	// Нахождение оператора
	var actionIdx int
	if op := "*/"; strings.ContainsAny(local_ex, op) {
		actionIdx = strings.IndexAny(local_ex, op)
	} else if op := "+-"; strings.ContainsAny(local_ex, op) {
		actionIdx = strings.IndexAny(local_ex, op)
	} else if strings.ContainsAny(local_ex, "()") {
		value, err := strconv.ParseFloat(local_ex[1:len(local_ex)-1], 64)
		if err != nil {
			return "", 0, Example{}, ExpressionIncorrect
		}
		return local_ex[:], strings.IndexRune(example, rune(local_ex[0])), Example{First_value: value, Second_value: 52, Operation: Equals}, nil //52 - по рофлу, чтобы при калькулировании не возникала ошибка. Крч костыль
	} else {
		value, err := strconv.ParseFloat(local_ex, 64)
		if err != nil {
			return "", 0, Example{}, ExpressionIncorrect
		}
		return "end", 0, Example{First_value: value, Second_value: 52, Operation: Equals}, nil
	}

	if actionIdx == 0 || actionIdx == len(local_ex)-1 {
		return "", 0, Example{}, OperationWithoutValue
	}

	ex.Operation = Operator(local_ex[actionIdx])

	// Нахождение концов двух чисел
	var exampleLen = len(local_ex)
	if actionIdx == 0 || actionIdx == exampleLen-1 {
		return "", 0, Example{}, OperationWithoutValue
	}

	var err error
	for i := actionIdx - 1; i >= 0; i-- {
		if strings.ContainsRune("+-/*()", rune(local_ex[i])) {
			ex.First_value, err = strconv.ParseFloat(local_ex[i+1:actionIdx], 64)
			if err != nil {
				return "", 0, Example{}, ExpressionIncorrect
			}
			begin = i + 1
			break
		} else if i == 0 {
			ex.First_value, err = strconv.ParseFloat(local_ex[i:actionIdx], 64)
			if err != nil {
				return "", 0, Example{}, ExpressionIncorrect
			}
			begin = i
			break
		}
	}

	for i := actionIdx + 1; i < exampleLen; i++ {
		if strings.ContainsRune("+-/*()", rune(local_ex[i])) {
			ex.Second_value, err = strconv.ParseFloat(local_ex[actionIdx+1:i], 64)
			if err != nil {
				return "", 0, Example{}, ExpressionIncorrect
			}
			end = i
			break
		} else if i+1 == exampleLen {
			ex.Second_value, err = strconv.ParseFloat(local_ex[actionIdx+1:i+1], 64)
			if err != nil {
				return "", 0, Example{}, ExpressionIncorrect
			}
			end = exampleLen
			break
		}
	}

	return local_ex[begin:end], strings.IndexRune(example, rune(local_ex[0])), ex, nil
}

// Заменяет выражение на его ответ
func EraseExample(example, erase_ex string, pri_idx int, answ float64) string {
	return example[:pri_idx] + strings.Replace(example[pri_idx:], erase_ex, fmt.Sprintf("%f", answ), 1)
}
