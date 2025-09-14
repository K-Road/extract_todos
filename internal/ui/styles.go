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
