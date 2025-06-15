package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	bolt "go.etcd.io/bbolt"
)

type Todo struct {
	File string
	Line int
	Text string
}

func hashTodo(todo Todo) string {
	s := fmt.Sprintf("%s:%s", todo.File, todo.Text)
	h := sha256.Sum256([]byte(s))
	return hex.EncodeToString(h[:])
}

func saveTodo(db *bolt.DB, todo Todo, projectName string) (bool, error) {
	id := hashTodo(todo)
	var saved bool

	err := db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(projectName))
		if err != nil {
			return err
		}
		if b.Get([]byte(id)) != nil {
			saved = false // already exists
			return nil
		}
		val := fmt.Sprintf("%s:%d:%s", todo.File, todo.Line, todo.Text)
		err = b.Put([]byte(id), []byte(val))
		if err == nil {
			saved = true
		}
		return err
	})
	return saved, err
}

func main() {
	//var todos []Todo
	//TODO Implement this is a flag
	root := "/home/chrode/workspace/github.com/K-Road/discord-moodbot/"
	projectName := filepath.Base(strings.TrimRight(root, string(os.PathSeparator)))
	fmt.Println(root)
	fmt.Println(projectName)

	//Open bolt db
	db, err := bolt.Open("todos.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || !strings.HasSuffix(path, ".go") {
			//fmt.Println("Skipping non-Go file:", path)
			return nil
		}
		//fmt.Println("Processing:", path)

		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		scanner := bufio.NewScanner(f)
		lineNum := 1
		for scanner.Scan() {
			line := scanner.Text()
			trimmed := strings.TrimSpace(line)
			relPath, err := filepath.Rel(root, path)
			if err != nil {
				relPath = path
			}
			if strings.HasPrefix(trimmed, "//TODO") {
				todo := Todo{
					File: relPath,
					Line: lineNum,
					Text: strings.TrimSpace(trimmed[len("//TODO"):]),
				}
				//Check if duplicate
				if saved, err := saveTodo(db, todo, projectName); err != nil {
					log.Println("Failed to save TODO:", err)
				} else if saved {
					fmt.Printf("New TODO saved: %s:%d %s\n", todo.File, todo.Line, todo.Text)
				}
			}
			lineNum++
		}
		return nil
	})
	//DEBUG to list all entries
	viewTodos(db)
}

func viewTodos(db *bolt.DB) {
	err := db.View(func(tx *bolt.Tx) error {
		return tx.ForEach(func(name []byte, b *bolt.Bucket) error {
			fmt.Printf("Project: %s\n", name)
			return b.ForEach(func(k, v []byte) error {
				fmt.Printf(" - %s\n", v)
				return nil
			})
		})
	})
	if err != nil {
		fmt.Println("Erroring reading from DB:", err)
	}
}
