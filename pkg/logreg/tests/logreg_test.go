package logreg_test

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/braginantonev/gcalc-server/pkg/database"
	dbreq "github.com/braginantonev/gcalc-server/pkg/database/requests-types"
	"github.com/braginantonev/gcalc-server/pkg/logreg/lrreq"
	"google.golang.org/grpc/status"

	"github.com/braginantonev/gcalc-server/pkg/logreg"
	pb "github.com/braginantonev/gcalc-server/proto/logreg"

	"github.com/golang-jwt/jwt/v5"
)

const TEST_DB_NAME string = "test.db"
const SIGNATURE string = "super_secret_signature"

func tokenIsValid(token string) (string, bool) {
	tokenFromString, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(SIGNATURE), nil
	})

	if err != nil {
		fmt.Println(err)
		return "", false
	}

	if claims, ok := tokenFromString.Claims.(jwt.MapClaims); ok {
		return claims["name"].(string), true
	} else {
		return "", false
	}
}

func TestRegisterLogin(t *testing.T) {
	ctx := context.Background()
	db, err := database.NewDataBase(ctx, TEST_DB_NAME)
	if err != nil {
		t.Error(err)
	}

	s := logreg.NewServer(db, SIGNATURE)
	if err := db.Create(ctx, dbreq.NewDBRequest(lrreq.CREATE_LogReg_Table)); err != nil {
		t.Error(err)
	}

	users := []*pb.User{
		{
			Name:     "Anton",
			Password: "12345",
		},
		{
			Name:     "San4hi",
			Password: "54321",
		},
	}

	cases := []struct {
		name          string
		user          *pb.User
		is_registered bool
		expected_err  error
	}{
		{
			name:          "Anton register",
			user:          users[0],
			is_registered: false,
		},
		{
			name:          "San4hi register",
			user:          users[1],
			is_registered: false,
		},
		{
			name:          "Anton login with correct password",
			user:          users[0],
			is_registered: true,
		},
		{
			name:          "San4hi login with correct password",
			user:          users[1],
			is_registered: true,
		},
		{
			name:          "Anton login with incorrect password",
			user:          &pb.User{Name: users[0].Name, Password: "motivation"},
			is_registered: true,
			expected_err:  logreg.ErrPasswordIncorrect,
		},
		{
			name:          "San4hi login with incorrect password",
			user:          &pb.User{Name: users[1].Name, Password: "power"},
			is_registered: true,
			expected_err:  logreg.ErrPasswordIncorrect,
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			if !test.is_registered {
				token, err := s.Register(ctx, test.user)
				if !errors.Is(err, test.expected_err) {
					t.Error("expected:", test.expected_err, "but got:", err)
				}

				name, is_valid := tokenIsValid(token.Token)
				if !is_valid {
					t.Error("token not valid")
				}

				if name != test.user.Name {
					t.Error("expected:", test.user.Name, "but got:", name)
				}
			}

			//Todo: Исправить проверку ошибок
			token, err := s.Login(ctx, test.user)
			if st, ok := status.FromError(err); ok && err != nil {
				if st.Message() != test.expected_err.Error() {
					t.Error("expected:", test.expected_err, "but got:", st.Message())
				}
			}

			name, is_valid := tokenIsValid(token.Token)
			if !is_valid {
				t.Error("token not valid")
			}

			if name != test.user.Name {
				t.Error("expected:", test.user.Name, "but got:", name)
			}
		})
	}

	db.Close()
	os.Remove(TEST_DB_NAME)
}
