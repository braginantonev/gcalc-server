package database_test

/*
import (
	"context"
	"os"
	"reflect"
	"testing"

	"github.com/braginantonev/gcalc-server/pkg/database"
	dbreq "github.com/braginantonev/gcalc-server/pkg/database/requests"
	"github.com/braginantonev/gcalc-server/pkg/orchestrator/orchreq"
	pb "github.com/braginantonev/gcalc-server/proto/orchestrator"
)

const TEST_DB_PATH string = "test.db"

func TestDBExpression(t *testing.T) {
	cases := []struct {
		name       string
		expression *pb.Expression
	}{
		{
			name: "simple expression",
			expression: &pb.Expression{
				User:   "Anton",
				Id:     0,
				Str:    "1+1",
				Status: pb.ETStatus_Backlog,
			},
		},
		{
			name: "simple expression 2",
			expression: &pb.Expression{
				User:   "Maksim",
				Id:     1,
				Str:    "1-1+1-1",
				Status: pb.ETStatus_Backlog,
			},
		},
	}

	ctx := context.Background()
	db, err := database.NewDataBase(ctx, TEST_DB_PATH)
	if err != nil {
		t.Error(err)
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			if err := db.Add(ctx, dbreq.NewDBRequest(orchreq.DBRequest_INSERT_Expression,
				test.expression.User, test.expression.Id, test.expression.Str, int32(test.expression.Status), test.expression.Result)); err != nil {
				t.Error(err)
			}

			res, err := db.Get(ctx, dbreq.NewDBRequest(orchreq.DBRequest_SELECT_Expression, test.expression.User, test.expression.Id))
			if err != nil {
				t.Error(err)
			}

			got_exp := res.(*pb.Expression)
			if !reflect.DeepEqual(test.expression, got_exp) {
				t.Error("expected: ", test.expression, "\nbut got:", got_exp)
			}
		})
	}

	db.Close()
	os.Remove(TEST_DB_PATH)
}

func TestDBExpressions(t *testing.T) {
	expression1 := &pb.Expression{
		User:   "Crystal Maiden",
		Id:     0,
		Str:    "1+1",
		Status: pb.ETStatus_Complete,
		Result: 2,
	}

	expression2 := &pb.Expression{
		User:   "Lina",
		Id:     6,
		Str:    "1-1+(13*9)-8",
		Status: pb.ETStatus_InProgress,
	}

	expression3 := &pb.Expression{
		User:   "Lina",
		Id:     7,
		Str:    "1/1/1/1/1/1/1/1/1/1/1/1/1/1/1/1/1/1/1/1/1/1/1/1/1/1/1",
		Status: pb.ETStatus_InProgress,
	}

	expressions := []*pb.Expression{expression1, expression2, expression3}

	ctx := context.Background()
	db, err := database.NewDataBase(ctx, TEST_DB_PATH)
	if err != nil {
		t.Error(err)
	}

	for _, exp := range expressions {
		if err := db.Add(ctx, dbreq.NewDBRequest(orchreq.DBRequest_INSERT_Expression, exp.User, exp.Id, exp.Str, int32(exp.Status), exp.Result)); err != nil {
			t.Error(err)
		}
	}

	res, err := db.Get(ctx, dbreq.NewDBRequest(orchreq.DBRequest_SELECT_Expressions, "Lina"))
	if err != nil {
		t.Error(err)
	}

	got_exps := res.(*pb.Expressions)
	if !reflect.DeepEqual(expressions[1:], got_exps.GetQueue()) {
		t.Error("expected: ", expressions[1:], "\nbut got:", got_exps.GetQueue())
	}

	res, err = db.Get(ctx, dbreq.NewDBRequest(orchreq.DBRequest_SELECT_Expressions, "Crystal Maiden"))
	if err != nil {
		t.Error(err)
	}

	got_exps = res.(*pb.Expressions)
	if !reflect.DeepEqual(expressions[:1], got_exps.GetQueue()) {
		t.Error("expected: ", expressions[:1], "\nbut got:", got_exps.GetQueue())
	}

	db.Close()
	os.Remove(TEST_DB_PATH)
}

func TestDBUpdateExpression(t *testing.T) {
	old_expression := &pb.Expression{
		User:   "Crystal Maiden",
		Id:     0,
		Str:    "1+1",
		Status: pb.ETStatus_InProgress,
	}

	new_expression := &pb.Expression{
		User:   "Crystal Maiden",
		Id:     0,
		Str:    "1+1",
		Status: pb.ETStatus_Complete,
		Result: 2,
	}

	ctx := context.Background()
	db, err := database.NewDataBase(ctx, TEST_DB_PATH)
	if err != nil {
		t.Error(err)
	}

	err = db.Add(ctx, dbreq.NewDBRequest(orchreq.DBRequest_INSERT_Expression,
		old_expression.User, old_expression.Id, old_expression.Str, int32(old_expression.Status), old_expression.Result))
	if err != nil {
		t.Error(err)
	}

	err = db.Update(ctx, dbreq.NewDBRequest(orchreq.DBRequest_UPDATE_Expression,
		int32(new_expression.Status), new_expression.Result, new_expression.User, new_expression.Id))
	if err != nil {
		t.Error(err)
	}

	res, err := db.Get(ctx, dbreq.NewDBRequest(orchreq.DBRequest_SELECT_Expression, "Crystal Maiden", 0))
	if err != nil {
		t.Error(err)
	}

	got_exp := res.(*pb.Expression)
	if !reflect.DeepEqual(new_expression, got_exp) {
		t.Error("expected:", new_expression, "\nbut got:", got_exp)
	}

	db.Close()
	os.Remove(TEST_DB_PATH)
}
*/
