package ui

import (
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

type statusMsg string

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit // Exit the TUI
		case "up", "u":
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
		m.spinnerRunning = false

		switch string(msg) {
		case "Webserver started!✅":
			m.webServerRunning = true
		case "Webserver stopped 🛑":
			m.webServerRunning = false
		}
		if msg != "" {
			return m, clearStatus()
		}

	case spinner.TickMsg:
		if m.spinnerRunning {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			cmds = append(cmds, cmd)
		}

	}
	return m, tea.Batch(cmds...)
}

func (m model) handleSelection() (tea.Model, tea.Cmd) {
	switch m.cursor {
	case 0:
		//Extract TODOs
		m.statusMessage = "Starting extraction..."
		m.spinnerRunning = true
		return m, tea.Batch(
			m.spinner.Tick,
			StartExtractTodos(),
		)
	case 1:
		//Start web server
		m.statusMessage = "Starting webserver..."
		m.spinnerRunning = true
		return m, tea.Batch(
			m.spinner.Tick,
			StartWebServerCmd(),
		)
	case 2:
		//Stop web server
		m.statusMessage = "Stopping webserver..."
		return m, tea.Batch(
			m.spinner.Tick,
			StopWebServerCmd(),
		)
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

func clearStatus() tea.Cmd {
	return tea.Tick(time.Second*3, func(time.Time) tea.Msg {
		return statusMsg("")
	})
}
