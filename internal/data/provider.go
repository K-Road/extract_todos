package data

import "github.com/K-Road/extract_todos/config"

type DataProvider interface {
	ListProjects() ([]string, error)
	ListProjectTodos(name string) ([]config.Todo, error)
	DeleteTodoById(name, id string) error
	GetActiveProject() (string, error)
	SetActiveProject(name string) error
}
