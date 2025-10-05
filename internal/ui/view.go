package ui

import (
	"fmt"
	"strings"

	"github.com/K-Road/extract_todos/web"
	"github.com/charmbracelet/lipgloss"
)

func (m model) View() string {
	switch {
	case m.showExtraction:
		return m.extractionModalView()
	default:
		return m.mainView()
	}
}

func (m *model) extractionModalView() string {
	termWidth, termHeight := getTerminalSize()
	m.modalWidth = termWidth - 10
	m.modalHeight = termHeight - 6

	if len(m.streams) == 0 {
		m.initStreams(m.modalWidth-4, m.modalHeight)
	}

	// advance streams
	for i := range m.streams {
		m.streams[i].step(m.modalHeight)
	}

	// generate full frame
	m.displayLines = m.generateMatrixFrame()

	// join lines and render modal
	logsStr := strings.Join(m.displayLines, "\n")
	logsStyle := lipgloss.NewStyle().
		Width(m.modalWidth).
		Border(lipgloss.RoundedBorder()).
		Render(logsStr)

	progressBar := m.progress.ViewAs(m.progressPercent)
	barStyle := lipgloss.NewStyle().Width(m.modalWidth).Align(lipgloss.Center).Render(progressBar)

	combined := logsStyle + "\n" + barStyle
	return lipgloss.Place(termWidth, termHeight, lipgloss.Center, lipgloss.Center, combined)
}

func (m model) mainView() string {
	m.webServerRunning = web.IsWebServerRunning()
	var s strings.Builder

	//header
	s.WriteString("\n")
	if m.activeProject != "" {
		activeProject := ActiveItemStyle.Render(m.activeProject)
		s.WriteString(fmt.Sprintf("Active Project: %s\n", activeProject))
	} else {
		s.WriteString("No active project set.\n")
	}

	s.WriteString("\nSelect an option:\n\n")

	for i, choice := range m.choices {
		cursor := " "
		//label := choice
		suffix := ""

		if m.state == "list" && choice.Label == m.activeProject {
			suffix = ActiveItemStyle.Render(" (active)")
		} else if choice.Label == "Start Web Server" && m.webServerRunning {
			suffix = RunningItemStyle.Render(" (running)")
		}

		if m.cursor == i {
			cursor = ">"
			line := SelectedItemStyle.Render(fmt.Sprintf("%s %s", cursor, choice.Label)) + suffix
			s.WriteString(line + "\n")
		} else {
			line := fmt.Sprintf("%s %s%s", cursor, choice.Label, suffix)
			s.WriteString(line + "\n")
		}
	}

	//footer
	switch m.state {
	case "main":
		s.WriteString("\nPress 'q' to quit.\n")
	case "list_projects":
		s.WriteString("\nEnter: Set Active | d: Delete | esc: Back\n")
	default:
		s.WriteString("\nPress 'esc' to return to the main menu.\n")
	}

	//status area
	if m.statusMessage != "" {
		if m.spinnerRunning {
			s.WriteString(fmt.Sprintf("%s %s\n", m.spinner.View(), m.statusMessage))
		} else {
			s.WriteString(fmt.Sprintf("💬 %s\n", m.statusMessage))
		}
	} else {
		s.WriteString("\n")
	}

	if m.progressVisible {
		s.WriteString(fmt.Sprintf("%s\n", m.progress.ViewAs(m.progressPercent)))
	} else {
		s.WriteString("\n")
	}
	return s.String()
}
