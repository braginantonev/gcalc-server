package orchestrator

import (
	"strconv"
	"strings"

	"github.com/Antibrag/gcalc-server/pkg/calc"
)

// Получает строку с выражением
func GetExample(example string) (calc.Example, int, error) {
	var ex calc.Example
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
			return calc.Example{}, 0, ErrBracketsNotFound
		}

		local_ex = local_ex[begin : end+1]
	}

	// ---- Нахождение оператора ----
	var actionIdx int
	if op := "*/"; strings.ContainsAny(local_ex, op) {
		actionIdx = strings.IndexAny(local_ex, op)
	} else if op := "+-"; strings.ContainsAny(local_ex, op) {
		actionIdx = strings.IndexAny(local_ex, op)
	} else if strings.ContainsAny(local_ex, "()") {
		if strings.Contains(local_ex, "id:") {
			return calc.Example{Operation: calc.Equals, String: local_ex}, strings.IndexRune(example, rune(local_ex[0])), nil
		}

		value, err := strconv.ParseFloat(local_ex[1:len(local_ex)-1], 64)

		if err != nil {
			return calc.Example{}, 0, calc.ErrExpressionIncorrect
		}
		return calc.Example{FirstArgument: calc.Argument{Value: value}, Operation: calc.Equals, String: local_ex[:]}, strings.IndexRune(example, rune(local_ex[0])), nil
	} else {
		if strings.Contains(local_ex, "id:") {
			return calc.Example{String: END_STR}, 0, nil
		}

		value, err := strconv.ParseFloat(local_ex, 64)
		if err != nil {
			return calc.Example{}, 0, calc.ErrExpressionIncorrect
		}
		return calc.Example{FirstArgument: calc.Argument{Value: value}, Operation: calc.Equals, String: END_STR}, 0, nil
	}

	if actionIdx == 0 || actionIdx == len(local_ex)-1 {
		return calc.Example{}, 0, ErrOperationWithoutValue
	}

	ex.Operation = calc.Operator(local_ex[actionIdx])

	// ---- Нахождение концов двух чисел ----
	var exampleLen = len(local_ex)
	if actionIdx == 0 || actionIdx == exampleLen-1 {
		return calc.Example{}, 0, ErrOperationWithoutValue
	}

	convertArgument := func(arg *calc.Argument, str string) (err error) {
		if strings.Contains(str, "id:") {
			arg.Expected = str[3:]
			ex.Status = calc.StatusIsWaitingValues
		} else {
			arg.Value, err = strconv.ParseFloat(str, 64)
		}
		return
	}

	var str_firstValue, str_secondValue string
	for i := actionIdx - 1; i >= 0; i-- {
		if strings.ContainsRune("+-/*()", rune(local_ex[i])) {
			str_firstValue = local_ex[i+1 : actionIdx]
			begin = i + 1
			break
		} else if i == 0 {
			str_firstValue = local_ex[i:actionIdx]
			begin = i
			break
		}
	}

	if err := convertArgument(&ex.FirstArgument, str_firstValue); err != nil {
		return calc.Example{}, 0, err
	}

	for i := actionIdx + 1; i < exampleLen; i++ {
		if strings.ContainsRune("+-/*()", rune(local_ex[i])) {
			str_secondValue = local_ex[actionIdx+1 : i]
			end = i
			break
		} else if i+1 == exampleLen {
			str_secondValue = local_ex[actionIdx+1 : i+1]
			end = exampleLen
			break
		}
	}

	if err := convertArgument(&ex.SecondArgument, str_secondValue); err != nil {
		return calc.Example{}, 0, err
	}

	ex.String = local_ex[begin:end]
	return ex, strings.IndexRune(example, rune(local_ex[0])), nil
}

// Заменяет выражение на его ответ
func EraseExample(example, erase_ex string, pri_idx int, id string) string {
	return example[:pri_idx] + strings.Replace(example[pri_idx:], erase_ex, "id:"+id, 1)
}
