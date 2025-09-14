package ui

import (
	"log"

	"github.com/K-Road/extract_todos/config"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	log              *log.Logger
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
	state            string // "list", "settings", "set", "add"
	dataProvider     config.DataProvider
	activeProject    string // currently active project
	extractionLogs   []string
	showExtraction   bool
}
type tickMsg struct{}
type progressMsg float64
type doneExtractingMsg struct{}
type tickContinueMsg struct{}
type WebServerStatusMsg bool
type extractionLogMsg struct {
	lines []string
}
type extractionDoneMsg struct{}

func InitialModel(logger *log.Logger, dp config.DataProvider) model {
	s := spinner.New(spinner.WithSpinner(spinner.Dot))
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	p := progress.New(
		progress.WithDefaultGradient(),
		progress.WithScaledGradient("10", "200"),
	)
	activeProject, _, err := dp.GetActiveProject()
	if err != nil {
		logger.Printf("Error getting active project: %v", err)
		activeProject = ""
	}
	// activeProject, err := dp.GetActiveProject()
	// if err != nil {
	// 	logger.Printf("Error getting active project: %v", err)
	// }

	// webServerRunning := false
	// if web.IsWebServerRunning() {
	// 	webServerRunning = true
	// 	logger.Println("Web server is currently running.")
	// } else {
	// 	logger.Println("Web server is not running.")
	// }

	m := model{
		log: logger,
		choices: []string{
			"Extract TODOs",
			"Project Settings",
			"Start Web Server",
			"Stop Web Server",
			"Exit TUI",
			"Exit & Shutdown Web Server",
		},
		spinner:          s,
		webServerRunning: false,
		progress:         p,
		state:            "main",
		dataProvider:     dp,
		activeProject:    activeProject,
	}
	return m
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		CheckWebServerStatusCmd(),
	)
}

func (m *model) resetToMainMenu() {
	m.state = "main"
	m.choices = mainMenuChoices()
	m.cursor = 0
}

func (m *model) resetToProjectSettings() {
	m.state = "settings"
	m.choices = projectSettingsMenuChoices()
	m.cursor = 0
}
