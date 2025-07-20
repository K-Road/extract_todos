package ui

import (
	"github.com/K-Road/extract_todos/web"
	tea "github.com/charmbracelet/bubbletea"
)

func StartWebServerCmd() tea.Cmd {
	return func() tea.Msg {
		err := web.StartWebServerDetached()
		if err != nil {
			return statusMsg("Failed to start webserver")
		}
		return statusMsg("Webserver started")
	}
}

func StopWebServerCmd() tea.Cmd {
	return func() tea.Msg {
		web.StopWebServer()
		return statusMsg("Webserver stopped")
	}
}
