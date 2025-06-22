package main

import (
	"fmt"
	"log"

	bolt "go.etcd.io/bbolt"
)

func main() {
	db, err := bolt.Open("todos.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	err = db.View(func(tx *bolt.Tx) error {
		return tx.ForEach(func(bucketName []byte, b *bolt.Bucket) error {
			fmt.Printf("Bucket: %s\n", bucketName)
			return b.ForEach(func(k, v []byte) error {
				fmt.Printf("Key: %s\n", k)
				fmt.Printf("Value: %s\n", v)
				return nil
			})
		})
	})
	if err != nil {
		log.Fatal(err)
	}
}
