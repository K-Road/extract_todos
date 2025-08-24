package db

import (
	"database/sql"

	"github.com/K-Road/extract_todos/config"
	_ "modernc.org/sqlite"
)

// Open connection to sqlite DB
func OpenDB(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`PRAGMA journal_mode=WAL;`)
	if err != nil {
		return nil, err
	}
	return db, nil
}

// Create schema if not exists
func InitSchema(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS projects (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT UNIQUE NOT NULL,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL,
		active BOOLEAN NOT NULL DEFAULT 0
	);

	CREATE TABLE IF NOT EXISTS todos (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		project_id INTEGER NOT NULL,
		hash TEXT NOT NULL UNIQUE,
		file TEXT NOT NULL,
		line INTEGER NOT NULL,
		text TEXT NOT NULL,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL,
		FOREIGN KEY(project_id) REFERENCES projects(id)
	);
	PRAGMA user_version = 1;
	`)
	return err
}

// List all projects
func ListProjects(db *sql.DB) ([]config.Project, error) {
	rows, err := db.Query(`SELECT id, name, created_at, updated_at, active FROM projects`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []config.Project
	for rows.Next() {
		var p config.Project
		if err := rows.Scan(&p.ID, &p.Name, &p.CreatedAt, &p.UpdatedAt, &p.Active); err != nil {
			return nil, err
		}
		projects = append(projects, p)
	}
	return projects, nil
}

func FetchProjectTodos(db *sql.DB, name string) ([]config.Todo, error) {
	return nil, nil
}

func DeleteTodoById(db *sql.DB, projectName, id string) error {
	return nil
}

func SaveTodo(bdb *sql.DB, projectName string, todo config.Todo) (bool, error) {
	return false, nil
}

func SetActiveProject(db *sql.DB, projectName string) error {
	return nil
}

func GetActiveproject(db *sql.DB) (string, error) {
	return "", nil
}
