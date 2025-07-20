package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	cursor  int
	choices []string
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

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit // Exit the TUI
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "enter":
			//trigger commands
		}

	}
	return m, nil
}

func (m model) View() string {
	s := "Select an option:\n\n"

	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}
		s += fmt.Sprintf("%s %s\n", cursor, choice)
	}
	s += "\nPress q to quit.\n"
	return s

}
