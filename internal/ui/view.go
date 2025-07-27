package ui

import (
	"fmt"

	"github.com/K-Road/extract_todos/web"
)

func (m model) View() string {
	m.webServerRunning = web.IsWebServerRunning()
	s := "Select an option:\n\n"

	for i, choice := range m.choices {
		cursor := " "
		label := choice
		suffix := ""

		if choice == "Start Web Server" && m.webServerRunning {
			suffix = RunningItemStyle.Render(" (running)")
		}

		//line := fmt.Sprintf(" %s", choice)

		if m.cursor == i {
			cursor = ">"
			line := SelectedItemStyle.Render(fmt.Sprintf("%s %s", cursor, label)) + suffix
			s += line + "\n"
		} else {
			line := fmt.Sprintf("%s %s%s", cursor, label, suffix)
			s += line + "\n"
		}
		//s += line + "\n"
	}

	if m.spinnerRunning {
		s += fmt.Sprintf("%s %s\n", m.spinner.View(), m.statusMessage)
	} else if m.statusMessage != "" {
		s += fmt.Sprintf("💬 %s\n", m.statusMessage)
	}

	if m.progressVisible {
		s += fmt.Sprintf("\n%s", m.progress.ViewAs(m.progressPercent))
	} //else if m.statusMessage != "" && m.progressVisible {
	//	s += fmt.Sprintf("\n💬 %s", m.statusMessage)
	//}

	s += "\nPress q to quit.\n"

	// if m.statusMessage != "" {
	// 	s += m.statusMessage
	// }
	return s

}
