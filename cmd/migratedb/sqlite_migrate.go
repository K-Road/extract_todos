package main

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"log"

	bolt "go.etcd.io/bbolt"
	_ "modernc.org/sqlite"
)

type Project struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type Todo struct {
	ID          string `json:"id"`
	ProjectID   string `json:"project_id"`
	Description string `json:"description"`
	Status      string `json:"status"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

func main() {
	bdb, err := bolt.Open("todos.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer bdb.Close()

	//Open SQLite
	sdb, err := sql.Open("sqlite", "todos.sqlite")
	if err != nil {
		log.Fatal(err)
	}
	defer sdb.Close()

	//Create schema in SQLite
	_, err = sdb.Exec(`
	CREATE TABLE IF NOT EXISTS projects (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT UNIQUE NOT NULL,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		active BOOLEAN NOT NULL DEFAULT 0
	);
	CREATE TABLE IF NOT EXISTS todos (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		project_id INTEGER NOT NULL,
		hash TEXT NOT NULL UNIQUE,
		file TEXT NOT NULL,
		line INTEGER NOT NULL,
		text TEXT NOT NULL,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (project_id) REFERENCES projects(id)
	);
	`)
	if err != nil {
		log.Fatal(err)
	}

	//migrate data from BoltDB to SQLite
	err = bdb.View(func(tx *bolt.Tx) error {
		return tx.ForEach(func(bucketName []byte, b *bolt.Bucket) error {
			projectName := string(bucketName)
			if projectName == "meta" {
				return nil
			}

			//Insert projects
			res, err := sdb.Exec("INSERT INTO projects(name) VALUES (?)", projectName)
			if err != nil {
				return err
			}

			var projectID int64
			if id, _ := res.LastInsertId(); id > 0 {
				projectID = id
			} else {
				err = sdb.QueryRow("SELECT id FROM projects WHERE name=?", projectName).Scan(&projectID)
				if err != nil {
					return err
				}
			}

			//migrate todos
			return b.ForEach(func(k, v []byte) error {
				hash := string(k)
				parts := strings.SplitN(string(v), ":", 3)
				if len(parts) != 3 {
					return fmt.Errorf("invalid todo value: %s", v)
				}
				file := parts[0]
				line, _ := strconv.Atoi(parts[1])
				text := parts[2]

				_, err := sdb.Exec(`
				INSERT INTO todos (project_id, hash, file, line, text)
				VALUES (?,?,?,?,?)`,
					projectID, hash, file, line, text,
				)
				return err
			})
		})
	})
	if err != nil {
		log.Fatal(err)
	}

	_, err = sdb.Exec("PRAGMA user_version = 1;")
	if err != nil {
		log.Fatal(err)
	}
}
