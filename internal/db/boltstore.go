package db

import (
	"fmt"

	bolt "go.etcd.io/bbolt"
)

func ListBuckets(db *bolt.DB) ([]string, error) {
	var buckets []string
	err := db.View(func(tx *bolt.Tx) error {
		return tx.ForEach(func(name []byte, _ *bolt.Bucket) error {
			buckets = append(buckets, string(name))
			return nil
		})
	})
	if err != nil {
		fmt.Println("Erroring reading from DB:", err)
		return nil, err
	}
	return buckets, nil
}

func ListProjectTodos(db *bolt.DB, name string) ([]string, error) {
	var todos []string
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(name))
		if b == nil {
			return fmt.Errorf("project bucket %q not found", name)
		}
		return b.ForEach(func(k, v []byte) error {
			todos = append(todos, string(v))
			return nil
		})
	})
	if err != nil {
		fmt.Println("Erroring reading from DB:", err)
		return nil, err
	}

	return todos, err
}
