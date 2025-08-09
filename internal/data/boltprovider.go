package data

import (
	"github.com/K-Road/extract_todos/config"
	"github.com/K-Road/extract_todos/internal/db"
	bolt "go.etcd.io/bbolt"
)

type BoltProvider struct {
	DB *bolt.DB
}

func (bp *BoltProvider) ListProjects() ([]string, error) {
	return db.ListBuckets(bp.DB)
}

func (bp *BoltProvider) ListProjectTodos(name string) ([]config.Todo, error) {
	return db.FetchProjectTodos(bp.DB, name)
}

func (bp *BoltProvider) DeleteTodoById(projectName, id string) error {
	return db.DeleteTodoById(bp.DB, projectName, id)
}

func (bp *BoltProvider) OpenDB(path string) error {
	var err error
	bp.DB, err = bolt.Open(path, 0600, nil)
	if err != nil {
		return err
	}
	return db.CheckDBVersionOrExit(bp.DB)
}

func (bp *BoltProvider) Close() error {
	if bp.DB != nil {
		return bp.DB.Close()
	}
	return nil
}
