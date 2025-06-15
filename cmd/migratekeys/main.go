package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/K-Road/extract_todos/internal/db"
	bolt "go.etcd.io/bbolt"
)

func main() {
	dryRun := flag.Bool("dry-run", false, "Show what would change without modifying db")
	force := flag.Bool("force", false, "Skip confirmation prompt and proceed with mgiration")
	flag.Parse()

	dbfile, err := bolt.Open("../../todos.db", 0666, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer dbfile.Close()
	fmt.Println(("HERE"))
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

	inserted, deleted, err := db.MigrateOldKeys(dbfile, *dryRun)
	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	if *dryRun {
		log.Printf("Dry-run completed. Inserted: %d, Deleted: %d", inserted, deleted)
	}
	if inserted != deleted {
		log.Println("Mismatch between inserted and deleted keys - check logs.")
	}
}
