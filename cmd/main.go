package main

import (
	"fmt"

	"github.com/braginantonev/gcalc-server/internal/application"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

//Todo: Почистить код
//Todo: Добавить Task уникальный id, без привязки к Expression id

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file", err)
	}

	grpcServer := grpc.NewServer()

	app := application.NewApplication()
	err = app.Run(grpcServer)
	if err != nil {
		fmt.Println(err)
	}
}
