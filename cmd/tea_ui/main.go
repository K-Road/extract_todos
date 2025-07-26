package main

import (
	"log"
	"os"

	"github.com/K-Road/extract_todos/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	logFile, err := os.OpenFile("debug.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(logFile)
	log.SetFlags((log.Ldate | log.Ltime | log.Lshortfile))

	p := tea.NewProgram(ui.InitialModel())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
