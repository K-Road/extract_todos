package data

import (
	"database/sql"
	"fmt"

	"github.com/K-Road/extract_todos/config"
	"github.com/K-Road/extract_todos/internal/db"
)

type SQLiteProvider struct {
	DB *sql.DB
}

func NewSQLiteProvider(path string) (*SQLiteProvider, error) {
	dbConn, err := db.OpenDB(path)
	if err != nil {
		return nil, err
	}
	if err := db.InitSchema(dbConn); err != nil {
		return nil, err
	}
	return &SQLiteProvider{DB: dbConn}, nil
}

func (sp *SQLiteProvider) ListProjects() ([]string, error) {
	projects, err := db.ListProjects(sp.DB)
	if err != nil {
		return nil, err
	}
	names := make([]string, len(projects))
	for i, p := range projects {
		names[i] = p.Name
	}
	return names, nil
}

func (sp *SQLiteProvider) ListProjectTodos(name string) ([]config.Todo, error) {
	return db.FetchProjectTodos(sp.DB, name)
}

func (sp *SQLiteProvider) DeleteTodoById(projectName, id string) error {
	return db.DeleteTodoById(sp.DB, projectName, id)
}

func (sp *SQLiteProvider) SaveTodo(projectName string, todo config.Todo) (bool, error) {
	return db.SaveTodo(sp.DB, projectName, todo)
}

// TODO add logging back
func (sp *SQLiteProvider) RemoveTodos(projectName string, scannedTodos []config.Todo) error {
	storedTodos, err := db.FetchProjectTodos(sp.DB, projectName)
	if err != nil {
		return fmt.Errorf("failed to fetch todos for project %s: %w", projectName, err)
	}
	scannedIDs := make(map[string]struct{})
	for _, todo := range scannedTodos {
		id := hashTodo(todo)
		scannedIDs[id] = struct{}{}
	}

	for _, todo := range storedTodos {
		id := hashTodo(todo)
		if _, exists := scannedIDs[id]; !exists {
			//getLog().Printf("Detected deleted TODO: %s:%s", todo.File, todo.Text)

			// //Delete from bolt db
			if err := db.DeleteTodoById(sp.DB, projectName, id); err != nil {
				//getLog().Printf("Failed to delete from DB: %v", err)
			}
		}
	}
	return nil
}

func (sp *SQLiteProvider) OpenDB(path string) error {
	return sp.OpenDB(path)
}

func (sp *SQLiteProvider) Close() error {
	if sp.DB != nil {
		return sp.DB.Close()
	}
	return nil
}

func (sp *SQLiteProvider) SetActiveProject(name string) error {
	return db.SetActiveProject(sp.DB, name)
}

func (sp *SQLiteProvider) GetActiveProject() (string, error) {
	return db.GetActiveproject(sp.DB)
}
