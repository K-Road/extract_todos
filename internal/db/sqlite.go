package db

import (
	"database/sql"
	"fmt"

	"github.com/K-Road/extract_todos/config"
	"github.com/K-Road/extract_todos/internal/helper"
	_ "modernc.org/sqlite"
)

type DB struct {
	Conn *sql.DB
}

// Open connection to sqlite DB
func OpenDB(path string) (*DB, error) {
	dsn := fmt.Sprintf("file:%s?_pragma=busy_timeout(5000)", path)
	sqlDB, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}

	_, err = sqlDB.Exec(`PRAGMA journal_mode=WAL;`)
	if err != nil {
		return nil, err
	}

	dbConn := &DB{sqlDB}

	if err := InitSchema(dbConn); err != nil {
		return nil, err
	}

	return dbConn, nil
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

	rows, err := db.Conn.Query(`SELECT t.id, t.project_id, t.file, t.line, t.text, t.hash,t.created_at, t.updated_at
	FROM todos t
	JOIN projects p ON t.project_id = p.id
	WHERE p.name = ?`, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []config.Todo
	for rows.Next() {
		var t config.Todo
		if err := rows.Scan(&t.ID, &t.ProjectID, &t.File, &t.Line, &t.Text, &t.Hash, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		todos = append(todos, t)
	}
	return todos, nil
}

func (db *DB) DeleteTodoById(id int) error {
	tx, err := db.Conn.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	_, err = tx.Exec(`DELETE FROM todos WHERE id = ?`, id)
	if err != nil {
		return err
	}
	return tx.Commit()

}

// SaveTodo writes a new todo to the database, updates the line number if it already but the line number has changed
// Pull this back to data package
func (db *DB) SaveTodo(projectName string, todo config.Todo) (config.TodoStatus, error) {
	//Uses a hash of the filename and the text of the todo to determine if it altready exists
	hash := helper.HashTodo(todo.File, todo.Text)
	var existingID, existingLine int

	err := db.Conn.QueryRow(`SELECT id, line FROM todos WHERE hash = ?`, hash).Scan(&existingID, &existingLine)
	if err != nil {
		if err == sql.ErrNoRows {
			//Not found - insert new
			_, err := db.Conn.Exec(`INSERT INTO todos (project_id, file, line, text, hash, created_at, updated_at)
			VALUES(
			(SELECT id FROM projects WHERE name = ?),
			?,?,?,?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
			)
			`, projectName, todo.File, todo.Line, todo.Text, hash)
			if err != nil {
				return config.TodoUnchanged, err
			}
			return config.TodoInserted, nil
		}
		return config.TodoUnchanged, err
	}

	//Found, update line number is changed.
	if existingLine != todo.Line {
		_, err := db.Conn.Exec(`
		UPDATE todos
		SET line = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
		`, todo.Line, existingID)
		if err != nil {
			return config.TodoUnchanged, err
		}
		return config.TodoUpdated, nil
	}
	//Found but no change
	return config.TodoUnchanged, nil
}

// Get todo by hash
func (db *DB) GetTodoByHash(hash string) (*config.Todo, error) {
	var t config.Todo
	err := db.Conn.QueryRow(`
	SELECT id, project_id, file, line, text, created_at, updated_at
	FROM todos
	WHERE hash = ?
	`, hash).Scan(&t.ID, &t.ProjectID, &t.File, &t.Line, &t.Text, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

// Insert new Todo
func (db *DB) InsertTodo(projectID int, todo config.Todo, hash string) error {
	_, err := db.Conn.Exec(`
	INSERT INTO todos (project_id, file, line, text, hash, created_at, updated_at)
	VALUES (?,?,?,?,?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`, projectID, todo.File, todo.Line, todo.Text, hash)
	return err
}

// Update existing todo line number
func (db *DB) UpdateTodoLine(id, line int) error {
	_, err := db.Conn.Exec(`
	UPDATE todos SET line = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?
	`, line, id)
	return err
}

func (db *DB) SetActiveProject(projectName string) error {
	tx, err := db.Conn.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	//set selected project to active
	if _, err := tx.Exec(`UPDATE projects SET active = 1 WHERE name = ?`, projectName); err != nil {
		return err
	}
	return tx.Commit()
}

func (db *DB) UnSetActiveProjects() error {
	tx, err := db.Conn.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.Exec(`UPDATE projects SET active = 0 WHERE active = 1`); err != nil {
		return err
	}
	return tx.Commit()
}

func (db *DB) GetActiveProject() (string, int, error) {
	var active string
	var activeID int
	err := db.Conn.QueryRow(`SELECT name, ID From projects WHERE active = 1 LIMIT 1`).Scan(&active, &activeID)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", 0, nil // No active project set
		}
		return "", 0, err
	}
	return active, activeID, nil
}

// Get project ID by name
func (db *DB) GetProjectID(projectName string) (int, error) {
	var id int
	err := db.Conn.QueryRow(`SELECT id from projects WHERE name = ?`, projectName).Scan(&id)
	return id, err
}
