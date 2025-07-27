package ui

import (
	"github.com/K-Road/extract_todos/web"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	cursor           int
	choices          []string
	statusMessage    string
	spinner          spinner.Model
	spinnerRunning   bool
	webServerRunning bool
	progress         progress.Model
	progressVisible  bool
	progressPercent  float64
	progressChan     chan tea.Msg
}
type tickMsg struct{}
type progressMsg float64
type doneExtractingMsg struct{}
type tickContinueMsg struct{}

func InitialModel() model {
	s := spinner.New(spinner.WithSpinner(spinner.Dot))
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	p := progress.New(
		progress.WithDefaultGradient(),
		progress.WithScaledGradient("10", "200"),
	)

	m := model{
		choices: []string{
			"Extract TODOs",
			"Start Web Server",
			"Stop Web Server",
			"Exit TUI",
			"Exit & Shutdown Web Server",
		},
		spinner:          s,
		webServerRunning: web.IsWebServerRunning(),
		progress:         p,
	}
	return m
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
	)
}
