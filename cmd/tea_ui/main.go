package main

import (
	"github.com/K-Road/extract_todos/internal/logging"
	"github.com/K-Road/extract_todos/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	logging.Init()
	log := logging.TUI()
	log.Println("Starting TUI...")

	p := tea.NewProgram(ui.InitialModel(log))
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
