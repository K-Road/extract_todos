package ui

import tea "github.com/charmbracelet/bubbletea"

type statusMsg string

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
			return m.handleSelection()
		}
	case statusMsg:
		m.statusMessage = string(msg)
		return m, nil
	}
	return m, nil
}

func (m model) handleSelection() (tea.Model, tea.Cmd) {
	switch m.cursor {
	case 0:
		//Extract TODOs
	case 1:
		//Start web server
		return m, tea.Batch(setStatus("Webserver started!✅"), StartWebServerCmd())
	case 2:
		//Stop web server
		return m, tea.Batch(setStatus("Webserver stopped!"), StopWebServerCmd())
	case 3:
		//Exit TUI
		return m, tea.Quit
	case 4:
		//Exit & Shutdown web server
		return m, tea.Batch(StopWebServerCmd(), tea.Quit)
	}
	return m, nil
}

func setStatus(msg string) tea.Cmd {
	return func() tea.Msg {
		return statusMsg(msg)
	}
}
