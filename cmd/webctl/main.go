package main

import (
	"fmt"
	"log"
	"os"

	"github.com/K-Road/extract_todos/web"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: webctl <start|stop|status>")
		os.Exit(1)
	}
	switch os.Args[1] {
	case "start":
		log.Println("Starting web server...")
		if err := web.StartWebServerDetached(); err != nil {
			log.Fatalf("Failed to start web server: %v", err)
		}
	case "stop":
		log.Println("Stopping web server...")
		if err := web.StopWebServer(); err != nil {
			log.Println("Web server was not running or failed to stop:", err)
		} else {
			log.Println("Web server stopped successfully.")
		}
	case "status":
		if web.IsWebServerRunning() {
			fmt.Println("Web server is running.")
		} else {
			fmt.Println("Web server is not running.")
		}
	default:
		fmt.Printf("Unkown command: %s\n", os.Args[1])
		fmt.Println("Usage: webctl <start|stop|status>")
		os.Exit(1)
	}
}
