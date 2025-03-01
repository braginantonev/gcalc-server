package agent

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"sync"
	"time"

	"github.com/Antibrag/gcalc-server/pkg/calc"
)

const (
	TIME_ADDITION_MS       = 100
	TIME_SUBTRACTION_MS    = 150
	TIME_MULTIPLICATION_MS = 300
	TIME_DIVISIONS_MS      = 350
	COMPUTING_POWER        = 5
	INTERNAL_TASK_URL      = "http://localhost:8080/internal/task"
)

type Request struct {
	Id     string  `json:"id"`
	Result float64 `json:"result,omitempty"`
	Error  string  `json:"error,omitempty"`
}

var examplesQueue []calc.Example

func EnableAgents(ctx context.Context) {
	examplesQueue = make([]calc.Example, 0, 5)
	mux := sync.Mutex{}

	//! Для сервера нужно сделать очередь из запросов, для избежания получения повторных примеров

	for range COMPUTING_POWER {
		go func() {
			select {
			case <-ctx.Done():
				return

			default:
				ex, err := GetExample()
				if err != nil {
					SendRequest(ex, err)
				}

				var exampleIdx int

				mux.Lock()
				if !slices.Contains(examplesQueue, ex) {
					examplesQueue = append(examplesQueue, ex)
					exampleIdx = len(examplesQueue) - 1
				}
				mux.Unlock()

				if err = Solve(&ex); err != nil {
					SendRequest(ex, err)
				}

				mux.Lock()
				examplesQueue = append(examplesQueue[:exampleIdx], examplesQueue[exampleIdx+1:]...)
				mux.Unlock()

				SendRequest(ex, nil)
			}
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
		<-time.After(TIME_ADDITION_MS * time.Millisecond)
		ex.Answer = ex.FirstArgument.Value + ex.SecondArgument.Value
		return nil

	case calc.Minus:
		<-time.After(TIME_SUBTRACTION_MS * time.Millisecond)
		ex.Answer = ex.FirstArgument.Value - ex.SecondArgument.Value
		return nil

	case calc.Multiply:
		<-time.After(TIME_MULTIPLICATION_MS * time.Millisecond)
		ex.Answer = ex.FirstArgument.Value * ex.SecondArgument.Value
		return nil

	case calc.Division:
		<-time.After(TIME_DIVISIONS_MS * time.Millisecond)
		ex.Answer = ex.FirstArgument.Value / ex.SecondArgument.Value
		return nil

		// case calc.Equals:
		// 	return ex.FirstArgument.Value, nil
	}
	return calc.ErrExpressionIncorrect
}
