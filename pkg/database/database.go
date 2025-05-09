package database

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"sync"

	dbreq "github.com/braginantonev/gcalc-server/pkg/database/requests"
	"github.com/braginantonev/gcalc-server/pkg/orchestrator"
	orch_pb "github.com/braginantonev/gcalc-server/proto/orchestrator"
	_ "github.com/mattn/go-sqlite3"
)

//! Если ничего не будет работать - чекать кэш, т.к. я хз как делать тесты для него, так что надеюсь, что он работает нормально
//Todo: Добавить использование базы данных в орхестраторе

type CacheConstrains interface {
	*orch_pb.Expression
}

type Cache[T CacheConstrains] struct {
	Query map[string]T
}

func (c *Cache[T]) Init() {
	c.Query = make(map[string]T, 0)
}

// if cache keep *orch_pb.Expression - id in the form: username+expression_id
func (c *Cache[T]) Set(mux *sync.RWMutex, id string, value T) {
	mux.Lock()
	defer mux.Unlock()
	c.Query[id] = value
}

// Чертов колхоз! Я хотел сделать все системно, чтобы не было методов только для одного типа данных, но я так заебался, так что пусть будет так

// Id in the form: username+expression_id
func (c *Cache[T]) UpdateExpression(mux *sync.RWMutex, id string, params ...any) bool {
	mux.Lock()
	defer mux.Unlock()

	element, ok := c.Query[id]
	if !ok {
		return false
	}

	expression := (*orch_pb.Expression)(element)
	expression.Status = params[0].(orch_pb.ETStatus)
	expression.Result = params[1].(float64)

	return true
}

// if cache keep *orch_pb.Expression - id in the form: username+expression_id
func (c *Cache[T]) Get(mux *sync.RWMutex, id string) (T, bool) {
	mux.Lock()
	defer mux.Unlock()
	res, ok := c.Query[id]
	return res, ok
}

func (c *Cache[T]) GetUserValues(mux *sync.RWMutex, username string) ([]T, bool) {
	res := make([]T, 0)
	mux.Lock()
	defer mux.Unlock()
	for key, value := range c.Query {
		if strings.Contains(key, username) {
			res = append(res, value)
		}
	}
	return res, len(res) > 0
}

// if cache keep *orch_pb.Expression - id in the form: username+expression_id
func (c *Cache[T]) Pop(mux *sync.RWMutex, id string) {
	mux.Lock()
	defer mux.Unlock()
	delete(c.Query, id)
}

type DataBase struct {
	*sql.DB
	Path string

	mux               *sync.RWMutex
	expressions_cache Cache[*orch_pb.Expression]
}

func (db *DataBase) Init(ctx context.Context) (err error) {
	db.DB, err = sql.Open("sqlite3", db.Path)
	if err != nil {
		return
	}

	if err = db.PingContext(ctx); err != nil {
		return err
	}

	db.mux = &sync.RWMutex{}

	//Orchestrator table
	_, err = db.ExecContext(ctx, orchestrator.DBRequest_CREATE_Table.ToString())
	db.expressions_cache.Init()
	return
}

/*
Available requests:
  - DBRequest_SELECT_Expressions (required 1 arg - user (string)). Return *pb.Expressions array
  - DBRequest_SELECT_Expression (required 2 args - user (string), internal_id (int32)). Return *pb.Expression
*/
func (db *DataBase) Get(ctx context.Context, request dbreq.DBRequest) (any, error) {
	if db.DB == nil {
		return nil, ErrDBNotInit
	}

	switch request.Type {
	case orchestrator.DBRequest_SELECT_Expressions:
		if !request.ArgsIsValid(1) {
			return nil, ErrBadArguments
		}

		cache_res, ok := db.expressions_cache.GetUserValues(db.mux, request.Args[0].(string))
		if ok {
			return &orch_pb.Expressions{Queue: cache_res}, nil
		}

		db.mux.Lock()
		rows, err := db.QueryContext(ctx, request.Type.ToString(), request.Args...)
		if err != nil {
			db.mux.Unlock()
			return nil, err
		}
		db.mux.Unlock()
		defer rows.Close()

		expressions := orch_pb.Expressions{}
		for rows.Next() {
			exp := orch_pb.Expression{}
			var status int32

			err := rows.Scan(&exp.User, &exp.Id, &exp.Str, &status, &exp.Result)
			if err != nil {
				return nil, err
			}

			exp.Status = orch_pb.ETStatus(status)
			expressions.Queue = append(expressions.Queue, &exp)
		}
		return &expressions, nil

	case orchestrator.DBRequest_SELECT_Expression:
		if !request.ArgsIsValid(2) {
			return nil, ErrBadArguments
		}

		cache_res, ok := db.expressions_cache.Get(db.mux, fmt.Sprintf("%s_%d", request.Args[0], request.Args[1]))
		if ok {
			return cache_res, nil
		}

		exp := orch_pb.Expression{}
		var status int32

		db.mux.Lock()
		defer db.mux.Unlock()
		if err := db.QueryRowContext(ctx, request.Type.ToString(), request.Args...).Scan(&exp.User, &exp.Id, &exp.Str, &status, &exp.Result); err != nil {
			return nil, err
		}

		exp.Status = orch_pb.ETStatus(status)
		return &exp, nil
	}

	return nil, ErrUnexpectedRequestType
}

/*
Available requests:
  - DBRequest_INSERT_Expression (required 5 args - user (string), internal_id (int32), str (string), status (int32), result (float64))
*/
func (db *DataBase) Add(ctx context.Context, request dbreq.DBRequest) error {
	if db.DB == nil {
		return ErrDBNotInit
	}

	switch request.Type {
	case orchestrator.DBRequest_INSERT_Expression:
		if !request.ArgsIsValid(5) {
			return ErrBadArguments
		}

		user, id := request.Args[0].(string), request.Args[1].(int32)
		db.expressions_cache.Set(db.mux, fmt.Sprintf("%s_%d", user, id), &orch_pb.Expression{
			User:   user,
			Id:     id,
			Str:    request.Args[2].(string),
			Status: orch_pb.ETStatus(request.Args[3].(int32)),
			Result: request.Args[4].(float64),
		})

		db.mux.Lock()
		defer db.mux.Unlock()
		if _, err := db.ExecContext(ctx, request.Type.ToString(), request.Args...); err != nil {
			return err
		}
		return nil
	}

	return ErrUnexpectedRequestType
}

/*
Available requests:
  - DBRequest_UPDATE_Expression (required 4 args - status (int32), result (float64), user (string), internal_id (int32))
*/
func (db *DataBase) Update(ctx context.Context, request dbreq.DBRequest) error {
	if db.DB == nil {
		return ErrDBNotInit
	}

	switch request.Type {
	case orchestrator.DBRequest_UPDATE_Expression:
		if !request.ArgsIsValid(4) {
			return ErrBadArguments
		}

		cache_id := fmt.Sprintf("%s_%d", request.Args[2], request.Args[3])
		if orch_pb.ETStatus(request.Args[0].(int32)) == orch_pb.ETStatus_Complete {
			db.expressions_cache.Pop(db.mux, cache_id)
		} else {
			if !db.expressions_cache.UpdateExpression(db.mux, cache_id, request.Args[:2]...) {
				return ErrCacheUpdateFailed
			}
		}

		db.mux.Lock()
		defer db.mux.Unlock()
		if _, err := db.ExecContext(ctx, request.Type.ToString(), request.Args...); err != nil {
			return err
		}
		return nil
	}

	return ErrUnexpectedRequestType
}

func NewDataBase(ctx context.Context, path string) (*DataBase, error) {
	if path == "" {
		return nil, ErrDBPathIsEmpty
	}

	db := &DataBase{Path: path}
	if err := db.Init(ctx); err != nil {
		return nil, err
	}

	return db, nil
}
