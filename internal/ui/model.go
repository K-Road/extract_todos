package ui

import (
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	cursor        int
	choices       []string
	statusMessage string
}

func InitialModel() model {
	return model{
		choices: []string{
			"Extract TODOs",
			"Start Web Server",
			"Stop Web Server",
			"Exit TUI",
			"Exit & Shutdown Web Server",
		},
	}
}

func (m model) Init() tea.Cmd {
	return nil
}
