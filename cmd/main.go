package main

import (
	"fmt"

	"github.com/Antibrag/gcalc-server/internal/application"
)

func main() {
	app := application.NewApplication()
	err := app.Run()
	if err != nil {
		fmt.Println(err)
	}
}
