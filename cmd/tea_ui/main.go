package main

import (
	"log"

	"github.com/K-Road/extract_todos/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	p := tea.NewProgram(ui.InitialModel())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
