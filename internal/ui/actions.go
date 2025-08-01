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

func StopWebServerCmd() tea.Cmd {
	return func() tea.Msg {
		time.Sleep(1 * time.Second)
		err := web.StopWebServer()
		if err != nil {
			return statusMsg("❌ Failed to stop webserver")
		}
		return statusMsg("Webserver stopped 🛑")
	}
}

func (m model) RunExtractionCmd() (tea.Model, tea.Cmd) {
	m.progressVisible = true
	m.progressPercent = 0
	m.spinnerRunning = true
	m.statusMessage = "Starting extraction..."
	msgCh := make(chan tea.Msg, 10)
	m.progressChan = msgCh

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

// func runExtractionCmd() tea.Cmd {
// 	return func() tea.Msg {
// 		ch := make(chan tea.Msg)
// 		go func() {
// 			_ = extract.RunWithProgress(func(p float64) {
// 				ch <- progressMsg(p)
// 			})
// 			ch <- doneExtractingMsg{}
// 			close(ch)
// 		}()

// 		return func() tea.Msg {
// 			msg, ok := <-ch
// 			if !ok {
// 				return nil
// 			}
// 			return msg
// 		}
// 	}
// }

// func runExtractionCmd() tea.Cmd {
// 	return func() tea.Msg {
// 		msgCh := make(chan tea.Msg, 100) // buffered to avoid deadlock

// 		go func() {
// 			err := extract.RunWithProgress(func(p float64) {
// 				select {
// 				case msgCh <- progressMsg(p):
// 				default: // avoid blocking if UI is not ready
// 				}
// 			})

// 			if err != nil {
// 				msgCh <- statusMsg("❌ Extraction failed")
// 			} else {
// 				msgCh <- doneExtractingMsg{}
// 			}
// 			close(msgCh)
// 		}()

// 		return func() tea.Msg {
// 			msg, ok := <-msgCh
// 			if !ok {
// 				return nil
// 			}
// 			return msg
// 		}
// 	}
// }
