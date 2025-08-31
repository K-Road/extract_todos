package ui

import (
	"fmt"
	"strings"

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
		//switch m.state {
		//case "settings":
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
		case 1: //Add Project
		//TODO add project logic
		case 2: //Back
			m.resetToMainMenu()
		}
	}

	return m, nil
}

func (m model) updateListProjects(msg tea.KeyMsg) (model, tea.Cmd) {
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
			selectedProject := stripActiveSuffix(m.choices[m.cursor])
			err := m.dataProvider.SetActiveProject(selectedProject)
			if err != nil {
				return m, func() tea.Msg {
					return statusMsg(fmt.Sprintf("❌ Error setting active project: %v", err))
				}
			} else {
				m.activeProject = selectedProject
				projects, _ := m.dataProvider.ListProjects()
				m.choices = markActiveProject(projects, m)

				return m, func() tea.Msg {
					return statusMsg(fmt.Sprintf("✅ Active project set to: %s", selectedProject))
				}
			}
		}
	case "d":
		//TODO implement delete project logic
	}
	return m, nil
}

// Do I need this anymore?
func markActiveProject(projects []string, m model) []string {
	active, _, _ := m.dataProvider.GetActiveProject()
	out := make([]string, len(projects))
	for i, p := range projects {
		if p == active {
			out[i] = fmt.Sprintf("%s (active)", p)
		} else {
			out[i] = p
		}
	}
	return out
}

// stripActiveSuffix removes "(active)" so we store the clean name
func stripActiveSuffix(name string) string {
	return strings.TrimSpace(strings.Replace(name, "(active)", "", 1))
}
