package database

import (
	"context"
	"database/sql"

	dbreq "github.com/braginantonev/gcalc-server/pkg/database/requests"
	"github.com/braginantonev/gcalc-server/pkg/orchestrator"
	orch_pb "github.com/braginantonev/gcalc-server/proto/orchestrator"
	_ "github.com/mattn/go-sqlite3"
)

//Todo: Доделать методы Add, Update
//Todo: Добавить тесты
//Todo: Реализовать кеш

// * -- Data base requests types -- * //

// * -- Data base -- * //

type DataBase struct {
	*sql.DB
	Path string
}

func (db *DataBase) Init(ctx context.Context) (err error) {
	db.DB, err = sql.Open("sqlite3", db.Path)
	if err != nil {
		return
	}

	if err = db.PingContext(ctx); err != nil {
		return err
	}

	//Orchestrator table
	_, err = db.ExecContext(ctx, orchestrator.DBRequest_CREATE_Table.ToString())
	return
}

/*
Available requests:
  - DBRequest_SELECT_Expressions (required 1 arg - user (string)). Return *pb.Expressions array
  - DBRequest_SELECT_Expression (required 2 args - user (string), internal_id (int32)). Return *pb.Expression
*/
func (db *DataBase) Get(ctx context.Context, request dbreq.DBRequest) (any, error) {
	switch request.Type {
	case orchestrator.DBRequest_SELECT_Expressions:
		if !request.ArgsIsValid(1) {
			return nil, ErrBadArguments
		}

		rows, err := db.QueryContext(ctx, request.Type.ToString(), request.Args...)
		if err != nil {
			return nil, err
		}
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

		exp := orch_pb.Expression{}
		var status int32

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
	switch request.Type {
	case orchestrator.DBRequest_INSERT_Expression:
		if !request.ArgsIsValid(5) {
			return ErrBadArguments
		}

		if _, err := db.ExecContext(ctx, request.Type.ToString(), request.Args...); err != nil {
			return err
		}
		return nil
	}

	return ErrUnexpectedRequestType
}

/*
Available requests:
  - DBRequest_UPDATE_Expression (required 4 args - user (string), internal_id (int32), status (int32), result (float64))
*/
func (db *DataBase) Update(ctx context.Context, request dbreq.DBRequest) error {
	switch request.Type {
	case orchestrator.DBRequest_UPDATE_Expression:
		if !request.ArgsIsValid(4) {
			return ErrBadArguments
		}

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
