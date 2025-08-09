package ui

import (
	"log"
	"time"

	"github.com/K-Road/extract_todos/internal/extract"
	"github.com/K-Road/extract_todos/web"
	tea "github.com/charmbracelet/bubbletea"
)

func StartWebServerCmd(log *log.Logger) tea.Cmd {
	return func() tea.Msg {
		if web.IsWebServerRunning() {
			log.Println("Web server is already running.")
			return statusMsg("Webserver started!✅")
		}
		time.Sleep(2 * time.Second)
		err := web.StartWebServerDetached()
		if err != nil {
			return statusMsg("Failed to start webserver")
		}
		return statusMsg("Webserver started!✅")
	}
}

func StopWebServerCmd(log *log.Logger) tea.Cmd {
	return func() tea.Msg {
		log.Println("Stopping web server...")
		time.Sleep(1 * time.Second)
		err := web.StopWebServer()
		if err != nil {
			return statusMsg("❌ Failed to stop webserver")
		}
		return statusMsg("Webserver stopped 🛑")
	}
}

func (m model) RunExtractionCmd(log *log.Logger) (tea.Model, tea.Cmd) {
	m.progressVisible = true
	m.progressPercent = 0
	m.spinnerRunning = true
	m.statusMessage = "Starting extraction..."
	msgCh := make(chan tea.Msg, 10)
	m.progressChan = msgCh
	log.Println("Running extraction command...")

	go func() {
		err := extract.RunWithProgress(func(p float64) {
			msgCh <- progressMsg(p)
		})
		if err != nil {
			msgCh <- statusMsg("❌ Extraction failed")
		} else {
			msgCh <- doneExtractingMsg{}
		}
		close(msgCh)
	}()

	return m, tea.Batch(
		m.spinner.Tick,
		readProgressChan(msgCh),
	)
}

func mainMenuChoices() []string {
	return []string{
		"Extract TODOs",
		"Project Settings",
		"Start Web Server",
		"Stop Web Server",
		"Exit TUI",
		"Exit & Shutdown Web Server",
	}
}

func settingsMenuChoices() []string {
	return []string{
		"List Projects",
		"Add Project",
		"Set Active Project",
		"Back",
	}
}
