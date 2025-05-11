package logreg

import (
	"context"
	"time"

	"github.com/braginantonev/gcalc-server/pkg/database"
	dbreq "github.com/braginantonev/gcalc-server/pkg/database/requests-types"
	"github.com/golang-jwt/jwt/v5"

	"github.com/braginantonev/gcalc-server/pkg/logreg/lrreq"
	pb "github.com/braginantonev/gcalc-server/proto/logreg"

	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
)

type Server struct {
	pb.LogRegServiceServer
	db            *database.DataBase
	jwt_signature string
}

func NewServer(server_db *database.DataBase, jwt_secret_signature string) *Server {
	return &Server{
		db:            server_db,
		jwt_signature: jwt_secret_signature,
	}
}

func (s *Server) Login(ctx context.Context, user *pb.User) (*pb.JWT, error) {
	if user.Name == "" {
		return nil, ErrUsernameIsEmpty
	}

	got_user, err := s.db.Get(ctx, dbreq.NewDBRequest(lrreq.SELECT_UserPass, user.Name))
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(got_user.(*pb.User).Password), []byte(user.Password))
	if err != nil {
		return nil, ErrPasswordIncorrect
	}

	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"name": got_user.(*pb.User).Name,
		"nbf":  now.Unix(),
		"exp":  now.Add(24 * time.Hour).Unix(),
		"iat":  now.Unix(),
	})

	token_str, err := token.SignedString([]byte(s.jwt_signature))
	if err != nil {
		return nil, err
	}

	return &pb.JWT{Token: token_str}, nil
}

func (s *Server) Register(ctx context.Context, user *pb.User) (*pb.JWT, error) {
	if user.Name == "" {
		return nil, ErrUsernameIsEmpty
	}

	_, err := s.db.Get(ctx, dbreq.NewDBRequest(lrreq.SELECT_UserPass, user.Name))
	if err == nil {
		return nil, ErrAlreadyRegistered
	}

	salt := []byte(user.Password)
	hash, err := bcrypt.GenerateFromPassword(salt, bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	err = s.db.Add(ctx, dbreq.NewDBRequest(lrreq.INSERT_UserPass, user.Name, string(hash)))
	if err != nil {
		return nil, err
	}

	jwt, err := s.Login(ctx, user)
	return jwt, err
}

func RegisterServer(ctx context.Context, grpcServer *grpc.Server, server_db *database.DataBase, jwt_secret_signature string) error {
	if server_db == nil {
		return database.ErrDBNotInit
	}

	if err := server_db.Create(ctx, dbreq.NewDBRequest(lrreq.CREATE_LogReg_Table)); err != nil {
		return err
	}

	pb.RegisterLogRegServiceServer(grpcServer, NewServer(server_db, jwt_secret_signature))
	return nil
}
