package ui

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	SelectedItemStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#2fc711ff"))

	CursorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00")).Bold(true)
)
