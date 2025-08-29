package ui

import (
	"log"
	"time"

	"github.com/K-Road/extract_todos/config"
	"github.com/K-Road/extract_todos/internal/extract"
	"github.com/K-Road/extract_todos/web"
	tea "github.com/charmbracelet/bubbletea"
)

func StartWebServerCmd(log *log.Logger, sendMsg func(tea.Msg)) tea.Cmd {
	return func() tea.Msg {
		if web.IsWebServerRunning() {
			log.Println("Web server is already running.")
			return statusMsg("Webserver already running!✅")
		}

		go func() {

			err := web.StartWebServerDetached()
			if err != nil {
				log.Println("Failed to start webserver:", err)
				// ch := make(chan tea.Msg)
				// ch <- statusMsg("❌ Failed to start webserver")
				sendMsg(statusMsg("Failed to start webserver"))
				return
			}

			sendMsg(statusMsg("Webserver started!✅"))
		}()
		return statusMsg("Webserver starting...")
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
	if m.activeProject == "" {
		m.statusMessage = "⚠️ No active project selected"
		return m, nil
	}
	m.progressVisible = true
	m.progressPercent = 0
	m.spinnerRunning = true
	m.statusMessage = "Starting extraction..."
	msgCh := make(chan tea.Msg, 10)
	m.progressChan = msgCh
	log.Println("Running extraction command...")

	go func(project string, dp config.DataProvider) {
		err := extract.RunWithProgress(project, dp, func(p float64) {
			msgCh <- progressMsg(p)
		})
		if err != nil {
			msgCh <- statusMsg("❌ Extraction failed")
		} else {
			msgCh <- doneExtractingMsg{}
		}
		close(msgCh)
	}(m.activeProject, m.dataProvider)

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

func projectSettingsMenuChoices() []string {
	return []string{
		"List Projects",
		"Add Project",
		"Back",
	}
}

func CheckWebServerStatusCmd() tea.Cmd {
	return func() tea.Msg {
		running := web.IsWebServerRunning()
		return WebServerStatusMsg(running)
	}
}
