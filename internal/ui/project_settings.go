package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

func (m model) updateProjectSettings(msg tea.KeyMsg) (model, tea.Cmd) {
	if m.state != "settings" && m.state != "list" && m.state != "set" && m.state != "add" {
		return m, nil
	}

	switch msg.String() {
	case "esc":
		m.resetToMainMenu()

	case "up":
		if m.cursor > 0 {
			m.cursor--
		}
	case "down":
		if m.cursor < len(m.choices)-1 {
			m.cursor++
		}
	case "enter":
		switch m.state {
		case "settings":
			switch m.cursor {
			case 0: //List Projects
				projects, err := m.dataProvider.ListProjects()
				if err != nil {
					m.statusMessage = fmt.Sprintf("❌ Error: %v", err)
				} else {
					m.choices = projects
					m.state = "list"
					if len(projects) > 0 {
						m.cursor = 0
					} else {
						m.cursor = -1 // No projects available
					}
				}
			case 3: //Back
				m.resetToMainMenu()
			}
		case "list":
			//handle list of projects
		}
	}
	return m, nil
}
