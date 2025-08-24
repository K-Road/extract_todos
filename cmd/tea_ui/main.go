package main

import (
	"github.com/K-Road/extract_todos/internal/data"
	"github.com/K-Road/extract_todos/internal/logging"
	"github.com/K-Road/extract_todos/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	logging.Init()
	log := logging.TUI()
	log.Println("Starting TUI...")

	//Initialize the database connection
	dbfile := "todos.sqlite"
	// bp := &data.BoltProvider{}
	// if err := bp.OpenDB(dbfile); err != nil {
	// 	log.Fatalf("Failed to open database: %v", err)
	// }

	factory := data.SQLiteFactory(dbfile)
	dp, err := factory()
	if err != nil {
		log.Fatalf("Failed to create data provider: %v", err)
	}
	defer dp.Close()

	//var dp config.DataProvider = bp

	p := tea.NewProgram(ui.InitialModel(log, dp))
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
