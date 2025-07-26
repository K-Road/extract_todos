package ui

import (
	"fmt"
	"time"

	"github.com/K-Road/extract_todos/internal/extract"
	"github.com/K-Road/extract_todos/web"
	tea "github.com/charmbracelet/bubbletea"
)

func StartWebServerCmd() tea.Cmd {
	return func() tea.Msg {
		time.Sleep(2 * time.Second)
		err := web.StartWebServerDetached()
		if err != nil {
			return statusMsg("Failed to start webserver")
		}
		return statusMsg("Webserver started")
	}
}

func StopWebServerCmd() tea.Cmd {
	return func() tea.Msg {
		time.Sleep(2 * time.Second)
		web.StopWebServer()
		return statusMsg("Webserver stopped")
	}
}

func StartExtractTodos() tea.Cmd {
	return func() tea.Msg {
		err := extract.Run()
		if err != nil {
			return statusMsg(fmt.Sprintf("%s", err.Error()))
		}
		return statusMsg("Extraction complete")
	}
}
