package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/K-Road/extract_todos/internal/db"
	bolt "go.etcd.io/bbolt"
)

func main() {
	dryRun := flag.Bool("dry-run", false, "Show what would change without modifying db")
	force := flag.Bool("force", false, "Skip confirmation prompt and proceed with mgiration")
	flag.Parse()

	//TODO Add helper to locate projectroot
	absPath, _ := filepath.Abs("todos.db")
	fmt.Println("Opening DB at:", absPath)

	//TODO handle lockouts
	boltdb, err := bolt.Open("todos.db", 0666, &bolt.Options{
		Timeout: 2 + time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer boltdb.Close()

	version, err := db.GetDBVersion(boltdb)
	if err != nil {
		log.Fatalf("Failed to read DB version: %v", err)
	}

	if version != db.CurrentVersion {
		log.Printf("DB Version is %q, current is %q. Running migration.", version, db.CurrentVersion)

		if !*dryRun && !*force {
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("This will modify the db. Continue? (yes/no): ")
			fmt.Print("> ")
			os.Stdout.Sync()
			confirm, _ := reader.ReadString('\n')
			confirm = strings.TrimSpace(strings.ToLower(confirm))
			if confirm != "yes" {
				fmt.Println("Aborted")
				return
			}
		}

		inserted, deleted, err := db.MigrateOldKeys(boltdb, *dryRun)
		if err != nil {
			log.Fatalf("Migration failed: %v", err)
		}

		if *dryRun {
			log.Printf("Dry-run completed. Inserted: %d, Deleted: %d", inserted, deleted)
		} else {
			log.Printf("Migration done. Inserted: %d, Deleted: %d", inserted, deleted)
			err = db.SetDBVersion(boltdb, db.CurrentVersion)
			if err != nil {
				log.Fatalf("Failed to update DB version: %v", err)
			}
		}
		if inserted != deleted {
			log.Println("Mismatch between inserted and deleted keys - check logs.")
		}
	} else {
		log.Println("DB version is current. Skipping migration.")
	}
}
