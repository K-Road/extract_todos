package ui

import (
	tea "github.com/charmbracelet/bubbletea"
)

func (m *model) updateGitHubIntegration(msg tea.KeyMsg) (model, tea.Cmd) {

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
		selected := m.choices[m.cursor]
		switch selected.Action {
		case "sync": //Sync
			//TODO add sync logic
		case "back":
			m.resetToMainMenu()
		}
	}
	return *m, nil
}
