package application

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/Antibrag/gcalc-server/pkg/calc"
	"github.com/Antibrag/gcalc-server/pkg/orchestrator"
)

const (
	TIME_ADDITION_MS       = 1000
	TIME_SUBTRACTION_MS    = 2000
	TIME_MULTIPLICATION_MS = 3000
	TIME_DIVISIONS_MS      = 4000
)

func ResultOrGet(result_fn http.HandlerFunc, get_fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Debug(r.Method)
		if r.Method == http.MethodGet {
			get_fn.ServeHTTP(w, r)
		} else if r.Method == http.MethodPost {
			result_fn.ServeHTTP(w, r)
		} else {
			slog.Error("ResultOrGet - wrong method")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

func ResultHandler(w http.ResponseWriter, r *http.Request) {
	req := make([]byte, 1024)
	n, err := r.Body.Read(req)
	if err != nil && err != io.EOF {
		slog.Error("ResultHandler(): Failed read request.", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	req = req[:n]

	var req_json struct {
		Id     string  `json:"id"`
		Result float64 `json:"result,omitempty"`
		Error  string  `json:"error,omitempty"`
	}

	err = json.Unmarshal(req, &req_json)
	if err != nil {
		fmt.Println("req - ", string(req))
		logFailedConvert("ResultHandler()", string(req), &w)
		return
	}

	err = orchestrator.SetExampleResult(req_json.Id, req_json.Result)
	if err != nil {
		if errors.Is(err, orchestrator.ErrTaskNotFound) {
			slog.Error("SetExampleResultHandler()", slog.String("error", err.Error()))
			w.WriteHeader(http.StatusNotFound)
			return
		} else {
			slog.Error("SetExampleResultHandler()", slog.String("error", err.Error()))
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}
	}

	slog.Info("Set result", slog.String("expression_id", req_json.Id), slog.Float64("result", req_json.Result))
	w.WriteHeader(http.StatusOK)
}

func GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	ex, err := orchestrator.GetTask("")
	if err == orchestrator.DHT {
		slog.Debug(err.Error())
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if err != nil {
		slog.Error("Internal error in getTask().", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var operation_time time.Duration
	switch ex.Operation {
	case calc.Plus:
		operation_time = TIME_ADDITION_MS * time.Millisecond
	case calc.Minus:
		operation_time = TIME_SUBTRACTION_MS * time.Millisecond
	case calc.Multiply:
		operation_time = TIME_MULTIPLICATION_MS * time.Millisecond
	case calc.Division:
		operation_time = TIME_DIVISIONS_MS * time.Millisecond
	}

	resp := calc.Example{
		Id:             ex.Id,
		FirstArgument:  ex.FirstArgument,
		SecondArgument: ex.SecondArgument,
		Operation:      ex.Operation,
		OperationTime:  operation_time,
	}

	resp_json, err := json.Marshal(resp)
	if err != nil {
		logFailedConvert("GetTaskHandler()", string(resp_json), &w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp_json)
}
