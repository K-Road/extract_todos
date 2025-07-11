package db

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/K-Road/extract_todos/config"
	bolt "go.etcd.io/bbolt"
)

const (
	MetaBucket     = "meta"
	VersionKey     = "version"
	CurrentVersion = "2" // increment this each time you change DB format/keys
)

func CheckDBVersionOrExit(dbfile *bolt.DB) {
	version, err := GetDBVersion(dbfile)
	if err != nil {
		log.Fatalf("Failed to read DB version: %v", err)
	}
	if version != CurrentVersion {
		log.Fatalf("DB is out of date (got %q, expected %q). Please run the migration tool.", version, CurrentVersion)
	}
}

func GetDBVersion(db *bolt.DB) (string, error) {
	var version string
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(MetaBucket))
		if b == nil {
			version = ""
			return nil
		}
		v := b.Get([]byte(VersionKey))
		if v == nil {
			version = ""
			return nil
		}
		version = string(v)
		return nil
	})
	return version, err
}

func SetDBVersion(db *bolt.DB, version string) error {
	return db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(MetaBucket))
		if err != nil {
			return err
		}
		return b.Put([]byte(VersionKey), []byte(version))
	})
}

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

func FetchProjectTodos(db *bolt.DB, name string) ([]config.Todo, error) {
	var todos []config.Todo
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(name))
		if b == nil {
			return fmt.Errorf("project bucket %q not found", name)
		}
		return b.ForEach(func(k, v []byte) error {
			parts := strings.SplitN(string(v), ":", 3)
			if len(parts) != 3 {
				return fmt.Errorf("invalid todo format in bucket %q: %s", name, v)
			}
			line, err := strconv.Atoi(parts[1])
			if err != nil {
				return fmt.Errorf("invalid line number in todo %q: %v", string(k), err)
			}
			todos = append(todos, config.Todo{
				File: parts[0],
				Line: line,
				Text: parts[2],
			})
			return nil
		})
	})
	if err != nil {
		fmt.Println("Erroring reading from DB:", err)
		return nil, err
	}

	return todos, err
}

func DeleteTodoById(bdb *bolt.DB, projectName, id string) error {
	return bdb.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(projectName))
		if b == nil {
			return fmt.Errorf("project bucket %q not found", projectName)
		}
		if err := b.Delete([]byte(id)); err != nil {
			return fmt.Errorf("failed to delete todo %q from project %q: %w", id, projectName, err)
		}
		return nil
	})
}
