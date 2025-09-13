package data

import (
	"log"

	"github.com/K-Road/extract_todos/internal/logging"
)

func getLog() *log.Logger {
	return logging.DATA()
}
