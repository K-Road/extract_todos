package ui

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/K-Road/extract_todos/web"
	"github.com/charmbracelet/lipgloss"
)

func (m model) View() string {
	switch {
	case m.showExtraction:
		return m.extractionModalView()
	default:
		return m.mainView()
	}
}

func (m model) extractionModalView() string {
	termWidth, termHeight := getTerminalSize()
	m.modalWidth = termWidth - 10
	m.modalHeight = termHeight - 6

	//init streams
	if len(m.streams) == 0 {
		m.initStreams(m.modalWidth-4, m.modalHeight)
	}

	//falling effect
	for i := range m.streams {
		m.streams[i].step(m.modalHeight)
	}

	displayLines := make([]string, m.modalHeight)

	for row := 0; row < m.modalHeight; row++ {
		var b strings.Builder
		for col := 0; col < m.modalWidth-4; col++ {
			s := &m.streams[col]
			//var char rune = ' '
			//var brightness float64 = 0.3

			// if row < len(s.column) {
			char := s.column[row]
			brightness := s.bright[row]
			// }

			logIndex := len(m.extractionLogs) - m.modalHeight + row
			if logIndex >= 0 && logIndex < len(m.extractionLogs) {
				line := m.extractionLogs[logIndex]

				if col < len(line) {
					orig := rune(line[col])

					r := rand.Float64()
					switch {
					case r < 0.6:
						char = randomMatrixChar()
						brightness = 0.1 + rand.Float64()*0.2
					case r < 0.8:
						char = unicodeCorrupt(orig)
						brightness = 0.3 + rand.Float64()*0.3
					default:
						char = orig
						brightness = 0.4 + rand.Float64()*0.4
					}
				}
			}
			b.WriteString(renderChar(char, brightness))
		}

		displayLines[row] = b.String()
	}

	logsStr := strings.Join(displayLines, "\n")

	logsStyle := lipgloss.NewStyle().
		Width(m.modalWidth).
		Height(m.modalHeight).
		Border(lipgloss.RoundedBorder()).
		Render(logsStr)

	progressWidth := m.modalWidth
	bar := m.progress.ViewAs(m.progressPercent)
	barStyle := lipgloss.NewStyle().Width(progressWidth).Align(lipgloss.Center).Render(bar)

	combined := logsStyle + "\n" + barStyle
	return lipgloss.Place(termWidth, termHeight, lipgloss.Center, lipgloss.Center, combined)
}

func truncateLine(line string, width int) string {
	if len(line) > width {
		return line[:width]
	}
	return line
}

func obfuscateLine(line string, width int) string {
	var b strings.Builder
	if line == "" {
		return randomNoiseLine(width)
	}
	for _, r := range line {
		if rand.Float64() < 0.3 {
			b.WriteRune(randomMatrixChar())
		} else {
			b.WriteRune(r)
		}
	}
	for b.Len() < width {
		b.WriteRune(randomMatrixChar())
	}
	return b.String()
}

func randomNoiseLine(width int) string {
	var b strings.Builder
	for i := 0; i < width; i++ {
		b.WriteRune(randomMatrixChar())
	}
	return b.String()
}

func (m model) mainView() string {
	m.webServerRunning = web.IsWebServerRunning()
	var s strings.Builder

	//header
	s.WriteString("\n")
	if m.activeProject != "" {
		activeProject := ActiveItemStyle.Render(m.activeProject)
		s.WriteString(fmt.Sprintf("Active Project: %s\n", activeProject))
	} else {
		s.WriteString("No active project set.\n")
	}

	s.WriteString("\nSelect an option:\n\n")

	if m.state == "list" {
		s.WriteString("Enter: Set Active | d: Delete | esc: Back\n\n")
	}

	for i, choice := range m.choices {
		cursor := " "
		label := choice
		suffix := ""

		if m.state == "list" && choice == m.activeProject {
			suffix = ActiveItemStyle.Render(" (active)")
		} else if choice == "Start Web Server" && m.webServerRunning {
			suffix = RunningItemStyle.Render(" (running)")
		}

		if m.cursor == i {
			cursor = ">"
			line := SelectedItemStyle.Render(fmt.Sprintf("%s %s", cursor, label)) + suffix
			s.WriteString(line + "\n")
		} else {
			line := fmt.Sprintf("%s %s%s", cursor, label, suffix)
			s.WriteString(line + "\n")
		}
	}

	//footer
	switch m.state {
	case "main":
		s.WriteString("\nPress 'q' to quit.\n")
	case "list":
		s.WriteString("\nEnter: Set Active | d: Delete | esc: Back\n")
	default:
		s.WriteString("\nPress 'esc' to return to the main menu.\n")
	}

	//status area
	if m.statusMessage != "" {
		if m.spinnerRunning {
			s.WriteString(fmt.Sprintf("%s %s\n", m.spinner.View(), m.statusMessage))
		} else {
			s.WriteString(fmt.Sprintf("💬 %s\n", m.statusMessage))
		}
	} else {
		s.WriteString("\n")
	}

	if m.progressVisible {
		s.WriteString(fmt.Sprintf("%s\n", m.progress.ViewAs(m.progressPercent)))
	} else {
		s.WriteString("\n")
	}
	return s.String()
}
