package agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Antibrag/gcalc-server/pkg/calc"
)

const (
	TIME_ADDITION_MS       = 100
	TIME_SUBTRACTION_MS    = 150
	TIME_MULTIPLICATION_MS = 300
	TIME_DIVISIONS_MS      = 350
	COMPUTING_POWER        = 1
	INTERNAL_TASK_URL      = "localhost/internal/task"
)

type Request struct {
	Id     string  `json:"id"`
	Result float64 `json:"result,omitempty"`
	Error  string  `json:"error,omitempty"`
}

func EnableAgents() {
	for range COMPUTING_POWER {
		go func() {
			ex, err := GetExample()
			if err != nil {
				SendRequest(ex, err)
			}

			if err = Solve(&ex); err != nil {
				SendRequest(ex, err)
			}

			SendRequest(ex, nil)
		}()
	}
}

func GetExample() (calc.Example, error) {
	resp, err := http.Get(INTERNAL_TASK_URL)
	if err != nil {
		return calc.Example{}, err
	}

	resp_json := make([]byte, 1024)
	n, err := resp.Body.Read(resp_json)
	if err != nil {
		return calc.Example{}, err
	}
	resp_json = resp_json[n:]

	var example calc.Example
	if err = json.Unmarshal(resp_json, &example); err != nil {
		return calc.Example{}, err
	}

	return example, nil
}

func SendRequest(example calc.Example, err error) {
	req := Request{Id: example.Id, Result: example.Answer, Error: err.Error()}
	req_json, err := json.Marshal(req)
	if err != nil {
		fmt.Printf("Error: SendRequest() - %s", err.Error())
	}

	_, err = http.Post(INTERNAL_TASK_URL, "application/json", bytes.NewReader(req_json))
	if err != nil {
		fmt.Printf("Error: SendRequest(): %s", err.Error())
	}
}

func Solve(ex *calc.Example) error {
	if ex.SecondArgument.Value == 0 && ex.Operation == calc.Division {
		return calc.ErrDivideByZero
	}

	switch ex.Operation {
	case calc.Plus:
		<-time.NewTimer(TIME_ADDITION_MS * time.Millisecond).C
		ex.Answer = ex.FirstArgument.Value + ex.SecondArgument.Value
		return nil

	case calc.Minus:
		<-time.NewTimer(TIME_SUBTRACTION_MS * time.Millisecond).C
		ex.Answer = ex.FirstArgument.Value - ex.SecondArgument.Value
		return nil

	case calc.Multiply:
		<-time.NewTimer(TIME_MULTIPLICATION_MS * time.Millisecond).C
		ex.Answer = ex.FirstArgument.Value * ex.SecondArgument.Value
		return nil

	case calc.Division:
		<-time.NewTimer(TIME_DIVISIONS_MS * time.Millisecond).C
		ex.Answer = ex.FirstArgument.Value / ex.SecondArgument.Value
		return nil

		// case calc.Equals:
		// 	return ex.FirstArgument.Value, nil
	}
	return calc.ErrExpressionIncorrect
}
