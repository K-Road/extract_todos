package ui

import (
	"time"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

type statusMsg string

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	//if m.state == "settings" || m.state == "list" || m.state == "set" || m.state == "add" {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch m.state {
		case "settings":
			return m.updateProjectSettings(keyMsg)
		case "list":
			return m.updateListProjects(keyMsg)
		case "add":
			//TODO add add route
		}
	}

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
		m.spinnerRunning = false
		m.statusMessage = string(msg)

		switch m.statusMessage {
		case "Webserver started!✅":
			m.webServerRunning = true
		case "Webserver stopped 🛑":
			m.webServerRunning = false
		}
		if m.statusMessage != "" {
			cmds = append(cmds, clearStatus())
		}
	case WebServerStatusMsg:
		m.webServerRunning = bool(msg)

	case progressMsg:
		m.progressPercent = float64(msg)
	case extractionLogMsg:
		m.extractionLogs = append(m.extractionLogs, msg.lines...)
		return m, nil

	case doneExtractingMsg:
		m.progressPercent = 1.0
		m.statusMessage = "✅ Extraction complete"
		m.spinnerRunning = false
		m.progressVisible = false
		m.progressChan = nil
		m.showExtraction = false
		m.extractionLogs = append(m.extractionLogs, "Extraction complete")
		return m, nil

	case spinner.TickMsg:
		if m.spinnerRunning {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			cmds = append(cmds, cmd)
		}
	case matrixTickMsg:
		return m, matrixTicker()

	}

	if m.progressVisible && m.progressChan != nil {
		var progressCmd tea.Cmd
		updatedProgress, progressCmd := m.progress.Update(msg)
		m.progress = updatedProgress.(progress.Model)
		cmds = append(cmds, progressCmd)

		//if m.progressChan != nil {
		cmds = append(cmds, readProgressChan(m.progressChan))
		//}
	}

	return m, tea.Batch(cmds...)
}

func (m model) handleSelection() (tea.Model, tea.Cmd) {
	selected := m.choices[m.cursor]

	switch selected.Action {
	case "extract":
		//Extract TODOs
		if m.activeProject == "" {
			return m, func() tea.Msg {
				return statusMsg("❌ No active project set. Please set an active project in Project Settings.")
			}
		}
		return m.RunExtractionCmd(m.log)
	case "project_settings":
		m.state = "settings"
		m.cursor = 0
		m.choices = projectSettingsMenuChoices()
		return m, nil
	case "start_web_server":
		//Start web server
		return m.withSpinner("Starting webserver...", StartWebServerCmd(m.log, func(msg tea.Msg) {
			m.Update(msg)
		}))
	case "stop_web_server":
		//Stop web server
		return m.withSpinner("Stopping webserver...", StopWebServerCmd(m.log))
	case "exit_tui":
		//Exit TUI
		return m, tea.Quit
	case "exit_shutdown_web_server":
		//Exit & Shutdown web server
		//TODO this isn't shutting down the webserver
		return m, tea.Batch(StopWebServerCmd(m.log), tea.Quit)

	default:
		return m, nil
	}
}

func (m model) withSpinner(status string, cmd tea.Cmd) (tea.Model, tea.Cmd) {
	m.statusMessage = status
	m.spinnerRunning = true
	m.progressVisible = false
	return m, tea.Batch(
		m.spinner.Tick,
		cmd,
	)
}

func clearStatus() tea.Cmd {
	return tea.Tick(time.Second*3, func(time.Time) tea.Msg {
		return statusMsg("")
	})
}

func readProgressChan(ch <-chan tea.Msg) tea.Cmd {
	return tea.Tick(time.Millisecond*60, func(t time.Time) tea.Msg {
		select {
		case msg, ok := <-ch:
			if !ok {
				return nil
			}
			return msg
		default:
			return nil
		}
	})
}

func tick() tea.Cmd {
	return tea.Tick(100*time.Millisecond, func(time.Time) tea.Msg {
		return tickMsg{}
	})
}
