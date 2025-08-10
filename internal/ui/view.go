package ui

import (
	"fmt"
	"strings"

	"github.com/K-Road/extract_todos/web"
)

func (m model) View() string {
	m.webServerRunning = web.IsWebServerRunning()
	var s strings.Builder

	//header
	s.WriteString("\n")
	if m.activeProject != "" {
		s.WriteString(fmt.Sprintf("Active Project: %s\n", m.activeProject))
	} else {
		s.WriteString("No active project set.\n")
	}

	s.WriteString("\nSelect an option:\n\n")

	for i, choice := range m.choices {
		cursor := " "
		label := choice
		suffix := ""

		if choice == "Start Web Server" && m.webServerRunning {
			suffix = RunningItemStyle.Render(" (running)")
		}

		if m.cursor == i {
			cursor = ">"
			line := SelectedItemStyle.Render(fmt.Sprintf("%s %s", cursor, label)) + suffix
			s.WriteString(line + "\n")
		} else {
			line := fmt.Sprintf("%s %s%s", cursor, label, suffix)
			s.WriteString(line + "\n")
		}
	}

	//footer
	switch m.state {
	case "main":
		s.WriteString("\nPress 'q' to quit.\n")
	default:
		s.WriteString("\nPress 'esc' to return to the main menu.\n")
	}

	//status area
	if m.statusMessage != "" {
		if m.spinnerRunning {
			//s += fmt.Sprintf("%s %s\n", m.spinner.View(), m.statusMessage)
			s.WriteString(fmt.Sprintf("%s %s\n", m.spinner.View(), m.statusMessage))
		} else //if m.statusMessage != "" {
		//s += fmt.Sprintf("💬 %s\n", m.statusMessage)
		{
			s.WriteString(fmt.Sprintf("💬 %s\n", m.statusMessage))
		}
	} else {
		s.WriteString("\n")
	}

	if m.progressVisible {
		//s += fmt.Sprintf("\n%s", m.progress.ViewAs(m.progressPercent))
		s.WriteString(fmt.Sprintf("%s\n", m.progress.ViewAs(m.progressPercent)))
	} else {
		s.WriteString("\n")
	} //else if m.statusMessage != "" && m.progressVisible {
	//	s += fmt.Sprintf("\n💬 %s", m.statusMessage)
	//}
	return s.String()

}
