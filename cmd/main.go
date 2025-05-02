package main

import (
	"fmt"

	"github.com/braginantonev/gcalc-server/internal/application"
	"google.golang.org/grpc"
)

//Todo: Пофиксить отображение статусов при вызове списка выражений
//Todo: Заменить все вызовы полей gRPC на гетеры
//Todo: Все &wrapperspb.StringValue{} заменить на сеттеры - wrapperspb.String()
//Todo: Почистить код

func main() {
	grpcServer := grpc.NewServer()

	app := application.NewApplication()
	err := app.Run(grpcServer)
	if err != nil {
		fmt.Println(err)
	}
}
