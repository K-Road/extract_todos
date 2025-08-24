package data

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/K-Road/extract_todos/config"
	"github.com/K-Road/extract_todos/internal/db"
	bolt "go.etcd.io/bbolt"
)

type BoltProvider struct {
	DB *bolt.DB
}

// func (bp *BoltProvider) ListProjects() ([]string, error) {
// 	return db.ListBuckets(bp.DB)
// }

// func (bp *BoltProvider) ListProjectTodos(name string) ([]config.Todo, error) {
// 	return db.FetchProjectTodos(bp.DB, name)
// }

// func (bp *BoltProvider) DeleteTodoById(projectName, id string) error {
// 	return db.DeleteTodoById(bp.DB, projectName, id)
// }

// func (bp *BoltProvider) SaveTodo(projectName string, todo config.Todo) (bool, error) {
// 	return db.SaveTodo(bp.DB, projectName, todo)
// }

// TODO add logging back
// func (bp *BoltProvider) RemoveTodos(projectName string, scannedTodos []config.Todo) error {
// 	storedTodos, err := db.FetchProjectTodos(bp.DB, projectName)
// 	if err != nil {
// 		return fmt.Errorf("failed to fetch todos for project %s: %w", projectName, err)
// 	}
// 	scannedIDs := make(map[string]struct{})
// 	for _, todo := range scannedTodos {
// 		id := hashTodo(todo)
// 		scannedIDs[id] = struct{}{}
// 	}

// 	for _, todo := range storedTodos {
// 		id := hashTodo(todo)
// 		if _, exists := scannedIDs[id]; !exists {
// 			//getLog().Printf("Detected deleted TODO: %s:%s", todo.File, todo.Text)

// 			// //Delete from bolt db
// 			if err := db.DeleteTodoById(bp.DB, projectName, id); err != nil {
// 				//getLog().Printf("Failed to delete from DB: %v", err)
// 			}
// 		}
// 	}
// 	return nil
// }

func (bp *BoltProvider) OpenDBWithOptions(path string, opts *bolt.Options) error {
	var err error
	bp.DB, err = bolt.Open(path, 0600, opts)
	if err != nil {
		return err
	}
	return db.CheckDBVersionOrExit(bp.DB)
}

func (bp *BoltProvider) OpenDB(path string) error {
	return bp.OpenDBWithOptions(path, nil)
}

func (bp *BoltProvider) Close() error {
	if bp.DB != nil {
		return bp.DB.Close()
	}
	return nil
}

// func (bp *BoltProvider) SetActiveProject(name string) error {
// 	return db.SetActiveProject(bp.DB, name)
// }

// func (bp *BoltProvider) GetActiveProject() (string, error) {
// 	return db.GetActiveproject(bp.DB)
// }

func hashTodo(todo config.Todo) string {
	s := fmt.Sprintf("%s:%s", todo.File, todo.Text)
	h := sha256.Sum256([]byte(s))
	return hex.EncodeToString(h[:])
}
