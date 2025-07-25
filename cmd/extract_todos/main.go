package main

import (
	"log"

	"github.com/K-Road/extract_todos/internal/extract"
)

func main() {
	if err := extract.Run(); err != nil {
		log.Fatalf("Failed to run extract: %v", err)
	}
}
