package main

import (
	"log"
	todo "todo-app"
	"todo-app/pkg/handler"
)

func main() {
	handlers := new(handler.Handler)

	srv := new(todo.Server)
	if err := srv.Run("8080", handlers.InitRouter()); err != nil {
		log.Fatalf("error occured while running server: %s", err.Error())
	}
}
