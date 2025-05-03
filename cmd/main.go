package main

import (
	"fmt"

	"github.com/braginantonev/gcalc-server/internal/application"
	"google.golang.org/grpc"
)

//Todo: Почистить код
//Todo: Добавить Task уникальный id, без привязки к Expression id

func main() {
	grpcServer := grpc.NewServer()

	app := application.NewApplication()
	err := app.Run(grpcServer)
	if err != nil {
		fmt.Println(err)
	}
}
