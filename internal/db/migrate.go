package db

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

// Used to update bbolt db keys
// old %s:%d:%s
// new %s:%s  --remove line number from key
func MigrateOldKeys(db *bolt.DB, dryRun bool) (int, int, error) {
	insertedCount := 0
	deletedCount := 0

	err := db.Update(func(tx *bolt.Tx) error {
		return tx.ForEach(func(bucketName []byte, b *bolt.Bucket) error {
			toMigrate := make(map[string]string)
			oldToNew := make(map[string]string)

			//collect keys
			err := b.ForEach(func(k, v []byte) error {
				parts := strings.SplitN(string(v), ":", 3)
				if len(parts) != 3 {
					return nil
				}
				file := parts[0]
				text := parts[2]
				newID := hashFileAndText(file, text)

				if newID != string(k) {
					toMigrate[newID] = string(v)
					oldToNew[string(k)] = newID
				}
				return nil
			})
			if err != nil {
				return err
			}

			//Insert new keys
			for newID, val := range toMigrate {
				if !dryRun {
					if err := b.Put([]byte(newID), []byte(val)); err != nil {
						return err
					}
					insertedCount++
				} else {
					log.Printf("[dry-run] Would insert: key=%s val=%s", newID, val)
				}
			}

			//Delete old keys
			for oldID := range oldToNew {
				if !dryRun {
					if err := b.Delete([]byte(oldID)); err != nil {
						return err
					}
					deletedCount++
				} else {
					log.Printf("[dry-run] Would delete old key: %s", oldID)
				}
			}

			return nil
		})
	})
	return insertedCount, deletedCount, err
}
