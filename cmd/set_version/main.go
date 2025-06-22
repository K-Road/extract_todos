package main

import (
	"log"

	bolt "go.etcd.io/bbolt"
)

const (
	dbPath         = "todos.db"
	metaBucket     = "meta"
	currentVersion = "2"
)

func main() {
	db, err := bolt.Open(dbPath, 0666, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(metaBucket))
		if err != nil {
			return err
		}
		return b.Put([]byte("version"), []byte(currentVersion))
	})
	if err != nil {
		log.Fatal("Failed to set version:", err)
	}

	log.Println("✅ Set DB version to", currentVersion)
}
