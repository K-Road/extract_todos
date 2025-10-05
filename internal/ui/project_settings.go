package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func (m *model) updateProjectSettings(msg tea.KeyMsg) (model, tea.Cmd) {
	if m.state != "settings_project" && m.state != "list_projects" && m.state != "set" && m.state != "add" {
		return *m, nil
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
		switch m.cursor {
		case 0: //List Projects
			projects, err := m.dataProvider.ListProjects()
			if err != nil {
				m.statusMessage = fmt.Sprintf("❌ Error: %v", err)
			} else {
				projectItems := make([]MenuItem, len(projects))
				for i, name := range projects {
					projectItems[i] = MenuItem{Label: name, Action: "select_project"}
				}
				m.choices = projectItems
				m.state = "list_projects"
				if len(projects) > 0 {
					m.cursor = 0
				} else {
					m.cursor = -1 // No projects available
				}
			}
		case 1: //Add Project
		//TODO add project logic
		case 2: //Back
			m.resetToMainMenu()
		}
	}

	return *m, nil
}

func (m *model) updateListProjects(msg tea.KeyMsg) (model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.resetToProjectSettings()
	case "up":
		if m.cursor > 0 {
			m.cursor--
		}
	case "down":
		if m.cursor < len(m.choices)-1 {
			m.cursor++
		}
	case "enter":
		if m.cursor >= 0 && m.cursor < len(m.choices) {
			selected := m.choices[m.cursor]
			selectedProject := stripActiveSuffix(selected.Label)
			err := m.dataProvider.SetActiveProject(selectedProject)
			if err != nil {
				return *m, func() tea.Msg {
					return statusMsg(fmt.Sprintf("❌ Error setting active project: %v", err))
				}
			}
			m.activeProject = selectedProject
			//rebuild list to mark active
			projects, _ := m.dataProvider.ListProjects()
			projectItems := markActiveProject(projects, *m)
			m.choices = projectItems

			return *m, func() tea.Msg {
				return statusMsg(fmt.Sprintf("✅ Active project set to: %s", selectedProject))
			}

		}
	case "d":
		//TODO implement delete project logic
	}
	return *m, nil
}

func markActiveProject(projects []string, m model) []MenuItem {
	active, _, _ := m.dataProvider.GetActiveProject()
	items := make([]MenuItem, len(projects))
	for i, p := range projects {
		label := p
		if p == active {
			label = fmt.Sprintf("%s (active)", p)
		}
		items[i] = MenuItem{Label: label, Action: "active_project"}
	}
	return items
}

// stripActiveSuffix removes "(active)" so we store the clean name
func stripActiveSuffix(name string) string {
	return strings.TrimSpace(strings.Replace(name, "(active)", "", 1))
}
