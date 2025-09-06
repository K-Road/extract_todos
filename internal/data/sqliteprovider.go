package data

import (
	"fmt"

	"github.com/K-Road/extract_todos/config"
	"github.com/K-Road/extract_todos/internal/db"
	"github.com/K-Road/extract_todos/internal/helper"
)

type SQLiteProvider struct {
	DB *db.DB
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
	projects, err := sp.DB.ListProjects()
	if err != nil {
		return nil, err
	}
	names := make([]string, len(projects))
	for i, p := range projects {
		names[i] = p.Name
	}
	return names, nil
}

func (sp *SQLiteProvider) ListProjectTodos(name string) ([]config.WebTodo, error) {
	todos, err := sp.DB.FetchProjectTodos(name)
	if err != nil {
		return nil, err
	}

	webTodos := make([]config.WebTodo, len(todos))
	for i, t := range todos {
		webTodos[i] = config.WebTodo{
			ID:   t.ID,
			File: t.File,
			Line: t.Line,
			Text: t.Text,
		}
	}
	return webTodos, nil
}

func (sp *SQLiteProvider) DeleteTodoById(id int) error {
	return sp.DB.DeleteTodoById(id)
}

func (sp *SQLiteProvider) SaveTodo(projectName string, todo config.Todo) (config.TodoStatus, error) {
	hash := helper.HashTodo(todo.File, todo.Text)

	projectID, err := sp.DB.GetProjectID(projectName)
	if err != nil {
		return config.TodoUnchanged, fmt.Errorf("project %q not found: %w", projectName, err)
	}

	existing, err := sp.DB.GetTodoByHash(hash)
	if err == nil {
		//hash exists -> check line
		if existing.Line != todo.Line {
			if err := sp.DB.UpdateTodoLine(existing.ID, todo.Line); err != nil {
				return config.TodoUnchanged, err
			}
			return config.TodoUpdated, nil
		}
		return config.TodoUnchanged, nil

	}

	//New todo
	if err := sp.DB.InsertTodo(projectID, todo, hash); err != nil {
		return config.TodoUnchanged, err
	}
	return config.TodoInserted, nil
	//return sp.DB.SaveTodo(projectName, todo)
}

// TODO add logging back
func (sp *SQLiteProvider) RemoveTodos(projectName string, scannedTodos []config.Todo) error {
	storedTodos, err := sp.DB.FetchProjectTodos(projectName)
	if err != nil {
		return fmt.Errorf("failed to fetch todos for project %s: %w", projectName, err)
	}
	scannedHashes := make(map[string]struct{})
	for _, todo := range scannedTodos {
		hash := hashTodo(todo)
		scannedHashes[hash] = struct{}{}
	}

	for _, todo := range storedTodos {
		//id := hashTodo(todo) //dont need now have hash in db
		if _, exists := scannedHashes[todo.Hash]; !exists {
			//getLog().Printf("Detected deleted TODO: %s:%s", todo.File, todo.Text)

			// //Delete from db
			if err := sp.DB.DeleteTodoById(todo.ID); err != nil {
				//getLog().Printf("Failed to delete from DB: %v", err)
			}
		}
	}
	return nil
}

func (sp *SQLiteProvider) OpenDB(path string) error {
	return nil
}

func (sp *SQLiteProvider) Close() error {
	if sp.DB != nil {
		return sp.DB.Close()
	}
	return nil
}

func (sp *SQLiteProvider) SetActiveProject(name string) error {
	if err := sp.DB.UnSetActiveProjects(); err != nil {
		return err
	}
	return sp.DB.SetActiveProject(name)
}

func (sp *SQLiteProvider) GetActiveProject() (string, int, error) {
	return sp.DB.GetActiveProject()
}
