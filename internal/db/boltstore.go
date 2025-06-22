package db

import (
	"fmt"
	"log"

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
