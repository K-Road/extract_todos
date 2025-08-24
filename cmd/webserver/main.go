package main

import (
	"fmt"

	"github.com/K-Road/extract_todos/internal/data"
	"github.com/K-Road/extract_todos/internal/logging"
	"github.com/K-Road/extract_todos/web"
)

func main() {
	fmt.Println("=== Webserver starting ===")
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Webserver panicked:", r)
		}
	}()
	logging.Init()
	fmt.Println("Logging initialized")
	factory := data.SQLiteFactory("todos.sqlite")
	fmt.Println("Data provider factory created")
	web.StartServer(factory)
	fmt.Println("Webserver exited normally")
}
