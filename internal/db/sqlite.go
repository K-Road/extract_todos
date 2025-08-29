package db

import (
	"database/sql"

	"github.com/K-Road/extract_todos/config"
	_ "modernc.org/sqlite"
)

type DB struct {
	Conn *sql.DB
}

// Open connection to sqlite DB
func OpenDB(path string) (*DB, error) {
	sqlDB, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	_, err = sqlDB.Exec(`PRAGMA journal_mode=WAL;`)
	if err != nil {
		return nil, err
	}
	return &DB{sqlDB}, nil
}

func (db *DB) Close() error {
	if db.Conn != nil {
		return db.Conn.Close()
	}
	return nil
}

// Create schema if not exists
func InitSchema(db *DB) error {
	_, err := db.Conn.Exec(`CREATE TABLE IF NOT EXISTS projects (
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
func (db *DB) ListProjects() ([]config.Project, error) {
	rows, err := db.Conn.Query(`SELECT id, name, created_at, updated_at, active FROM projects`)
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

func (db *DB) FetchProjectTodos(name string) ([]config.Todo, error) {
	return nil, nil
}

func (db *DB) DeleteTodoById(projectName, id string) error {
	return nil
}

func (db *DB) SaveTodo(projectName string, todo config.Todo) (bool, error) {
	return false, nil
}

func (db *DB) SetActiveProject(projectName string) error {
	tx, err := db.Conn.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	//set all projects to inactive
	if _, err := tx.Exec(`UPDATE projects SET active = 0 WHERE active = 1`); err != nil {
		return err
	}

	//set selected project to active
	if _, err := tx.Exec(`UPDATE projects SET active = 1 WHERE name = ?`, projectName); err != nil {
		return err
	}
	return tx.Commit()

}

func (db *DB) GetActiveproject() (string, error) {
	var active string
	err := db.Conn.QueryRow(`SELECT name From projects WHERE active = 1 LIMIT 1`).Scan(&active)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil // No active project set
		}
		return "", err
	}
	return active, nil
}
