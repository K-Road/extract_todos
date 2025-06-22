package main

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"log"
	"strings"

	bolt "go.etcd.io/bbolt"
)

func hashFileAndText(file, text string) string {
	data := fmt.Sprintf("%s:%s", file, text)
	h := sha1.Sum([]byte(data))
	return hex.EncodeToString(h[:])
}

func main() {
	db, err := bolt.Open("todos.db", 0666, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		return tx.ForEach(func(bucketName []byte, b *bolt.Bucket) error {
			fmt.Printf("Scanning bucket: %s\n", bucketName)
			var deleted int
			err := b.ForEach(func(k, v []byte) error {
				keyStr := string(k)
				valStr := string(v)

				fmt.Printf("Key: %s\n", keyStr)
				fmt.Printf("Val: %s\n", valStr)

				parts := strings.SplitN(valStr, ":", 3)
				if len(parts) != 3 {
					fmt.Println("  → Skipped: value not in file:line:text format")
					return nil
				}

				file := parts[0]
				text := parts[2]
				newID := hashFileAndText(file, text)

				if newID != keyStr {
					fmt.Printf("  → Deleting legacy key: %s\n", keyStr)
					err := b.Delete(k)
					if err != nil {
						return err
					}
					deleted++
				} else {
					fmt.Println("  → Already uses new ID")
				}
				return nil
			})
			fmt.Printf("Deleted %d keys from bucket %s\n\n", deleted, bucketName)
			return err
		})
	})

	if err != nil {
		log.Fatalf("cleanup failed: %v", err)
	}
}
