package ui

import "fmt"

func (m model) View() string {
	s := "Select an option:\n\n"

	for i, choice := range m.choices {
		cursor := " "
		line := fmt.Sprintf(" %s", choice)

		if m.cursor == i {
			cursor = ">"
			line = SelectedItemStyle.Render(fmt.Sprintf("%s %s", cursor, choice))
		} else {
			line = fmt.Sprintf("%s %s", cursor, choice)
		}
		s += line + "\n"
	}
	s += "\nPress q to quit.\n"

	if m.statusMessage != "" {
		s += m.statusMessage
	}
	return s

}
