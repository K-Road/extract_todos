package ui

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	SelectedItemStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#2fc711ff"))
	RunningItemStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
	ActiveItemStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00")).Bold(true)

	CursorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00")).Bold(true)

	MatrixStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00")).Background(lipgloss.Color("#000000"))
)

func renderChar(r rune, brightness float64) string {
	switch {
	case brightness > 0.7:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00")).Render(string(r))
	case brightness > 0.4:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#33FF33")).Render(string(r))
	default:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#007700")).Render(string(r))
	}
}
