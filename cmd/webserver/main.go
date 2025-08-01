package main

import (
	"github.com/K-Road/extract_todos/internal/logging"
	"github.com/K-Road/extract_todos/web"
)

func main() {
	logging.Init()
	web.StartServer()
}
