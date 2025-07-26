package ui

import (
	"time"

	"github.com/charmbracelet/bubbles/progress"
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
			cmds = append(cmds, clearStatus())
		}

	case tickMsg:
		if m.progressPercent < 1.0 {
			m.progressPercent += 0.05
			cmds = append(cmds, tickProgressBar())
		} else {
			m.progressVisible = false
			m.statusMessage = "Extraction complete"
		}

	case spinner.TickMsg:
		if m.spinnerRunning {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			cmds = append(cmds, cmd)
		}

	}
	var progressCmd tea.Cmd
	updatedProgress, progressCmd := m.progress.Update(msg)
	m.progress = updatedProgress.(progress.Model)
	cmds = append(cmds, progressCmd)
	return m, tea.Batch(cmds...)
}

func (m model) handleSelection() (tea.Model, tea.Cmd) {
	switch m.cursor {
	case 0:
		//Extract TODOs
		m.progressVisible = true
		m.progressPercent = 0
		m.spinnerRunning = true
		m.statusMessage = "Starting extraction..."
		return m, tea.Batch(
			m.spinner.Tick,
			tickProgressBar(),
			StartExtractTodos(),
		)
	case 1:
		//Start web server
		m.statusMessage = "Starting webserver..."
		m.spinnerRunning = true
		m.progressVisible = false
		return m, tea.Batch(
			m.spinner.Tick,
			StartWebServerCmd(),
		)
	case 2:
		//Stop web server
		m.statusMessage = "Stopping webserver..."
		m.spinnerRunning = true
		m.progressVisible = false
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

func tickProgressBar() tea.Cmd {
	return tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg{}
	})
}
