package calc

import (
	"fmt"
	"strconv"
	"strings"
)

var ExamplesQueue []Example

// Получает строку с выражением
func GetExample(example string) (int, Example, error) {
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
			return 0, Example{}, ErrBracketsNotFound
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
		// ex_without_brackets := local_ex[1:len(local_ex)-1]
		// var value float64
		// var err error

		// if strings.Contains(ex_without_brackets, "id") {

		// }
		value, err := strconv.ParseFloat(local_ex[1:len(local_ex)-1], 64)
		if err != nil {
			return 0, Example{}, ErrExpressionIncorrect
		}
		return strings.IndexRune(example, rune(local_ex[0])), Example{FirstArgument: Argument{Value: value}, Operation: Equals, String: local_ex[:]}, nil
	} else {
		value, err := strconv.ParseFloat(local_ex, 64)
		if err != nil {
			return 0, Example{}, ErrExpressionIncorrect
		}
		return 0, Example{FirstArgument: Argument{Value: value}, Operation: Equals, String: "end"}, nil
	}

	if actionIdx == 0 || actionIdx == len(local_ex)-1 {
		return 0, Example{}, ErrOperationWithoutValue
	}

	ex.Operation = Operator(local_ex[actionIdx])

	// Нахождение концов двух чисел
	var exampleLen = len(local_ex)
	if actionIdx == 0 || actionIdx == exampleLen-1 {
		return 0, Example{}, ErrOperationWithoutValue
	}

	var err error
	for i := actionIdx - 1; i >= 0; i-- {
		if strings.ContainsRune("+-/*()", rune(local_ex[i])) {
			ex.FirstArgument.Value, err = strconv.ParseFloat(local_ex[i+1:actionIdx], 64)
			if err != nil {
				return 0, Example{}, ErrExpressionIncorrect
			}
			begin = i + 1
			break
		} else if i == 0 {
			ex.FirstArgument.Value, err = strconv.ParseFloat(local_ex[i:actionIdx], 64)
			if err != nil {
				return 0, Example{}, ErrExpressionIncorrect
			}
			begin = i
			break
		}
	}

	for i := actionIdx + 1; i < exampleLen; i++ {
		if strings.ContainsRune("+-/*()", rune(local_ex[i])) {
			ex.SecondArgument.Value, err = strconv.ParseFloat(local_ex[actionIdx+1:i], 64)
			if err != nil {
				return 0, Example{}, ErrExpressionIncorrect
			}
			end = i
			break
		} else if i+1 == exampleLen {
			ex.SecondArgument.Value, err = strconv.ParseFloat(local_ex[actionIdx+1:i+1], 64)
			if err != nil {
				return 0, Example{}, ErrExpressionIncorrect
			}
			end = exampleLen
			break
		}
	}

	ex.String = local_ex[begin:end]
	return strings.IndexRune(example, rune(local_ex[0])), ex, nil
}

// Заменяет выражение на его ответ
func EraseExample(example, erase_ex string, pri_idx int, id int) string {
	return example[:pri_idx] + strings.Replace(example[pri_idx:], erase_ex, fmt.Sprintf("id:%d", id), 1)
}
