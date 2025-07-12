package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"

	"github.com/K-Road/extract_todos/config"
	"github.com/K-Road/extract_todos/internal/db"
	bolt "go.etcd.io/bbolt"
)

// Stops webserver to avoid conflicts
// Returns true if the server was already running
func stopWebServer() bool {
	data, err := os.ReadFile("webserver.pid")
	if err != nil {
		log.Printf("Failed to read webserver PID file: %v", err)
		return false
	}

	pidStr := strings.TrimSpace(string(data))
	pid, err := strconv.Atoi(pidStr)
	if err != nil {
		log.Printf("Invalid PID in webserver PID file: %v", err)
		return false
	}

	process, err := os.FindProcess(pid)
	if err != nil {
		log.Printf("Failed to find webserver process with PID %d: %v", pid, err)
		return false
	}

	//check if the process is running
	if err := process.Signal(syscall.Signal(0)); err != nil {
		log.Printf("Webserver with PID %d is not running: %v", pid, err)
		_ = os.Remove("webserver.pid")
		return true
	}
	if err = process.Signal(os.Interrupt); err != nil {
		log.Printf("Failed to send interrupt signal to webserver with PID %d: %v", pid, err)
		return false
	}
	log.Printf("Sent interrupt signal to webserver with PID %d", pid)
	return true
}

func startWebServer() error {
	binaryPath := "./web/webserver.go"
	cmd := exec.Command("bash", "-c", fmt.Sprintf("nohup go run %s > webserver.log 2>&1 &", binaryPath))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Start()
}

func hashTodo(todo config.Todo) string {
	s := fmt.Sprintf("%s:%s", todo.File, todo.Text)
	h := sha256.Sum256([]byte(s))
	return hex.EncodeToString(h[:])
}

func saveTodo(bdb *bolt.DB, todo config.Todo, projectName string) (bool, error) {
	id := hashTodo(todo)
	var saved bool

	err := bdb.Update(func(tx *bolt.Tx) error {
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
	var scannedTodos []config.Todo
	//TODO Implement this is a flag
	root := "/home/chrode/workspace/github.com/K-Road/discord-moodbot/"
	projectName := filepath.Base(strings.TrimRight(root, string(os.PathSeparator)))
	fmt.Println(root)
	fmt.Println(projectName)

	wasServerRunning := stopWebServer()

	//Open bolt db
	bdb, err := bolt.Open("todos.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer bdb.Close()

	db.CheckDBVersionOrExit(bdb)

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
				todo := config.Todo{
					File: relPath,
					Line: lineNum,
					Text: strings.TrimSpace(trimmed[len("//TODO"):]),
				}
				scannedTodos = append(scannedTodos, todo) // Collect all todos for delete sync
				//Check if duplicate
				if saved, err := saveTodo(bdb, todo, projectName); err != nil {
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
	viewTodos(bdb)

	err = removeTodos(bdb, projectName, scannedTodos)

	viewTodos(bdb)

	//Restart webserver
	if wasServerRunning {
		log.Println("Restarting webserver...")
		if err := startWebServer(); err != nil {
			log.Printf("Failed to start webserver: %v", err)
		}
	}
}

func removeTodos(bdb *bolt.DB, projectName string, scannedTodos []config.Todo) error {
	storedTodos, err := db.FetchProjectTodos(bdb, projectName)
	if err != nil {
		return fmt.Errorf("failed to fetch todos for project %s: %w", projectName, err)
	}
	scannedIDs := make(map[string]struct{})
	for _, todo := range scannedTodos {
		id := hashTodo(todo)
		scannedIDs[id] = struct{}{}
	}

	for _, todo := range storedTodos {
		id := hashTodo(todo)
		if _, exists := scannedIDs[id]; !exists {
			log.Printf("Detected deleted TODO: %s:%s", todo.File, todo.Text)

			// //Delete from bolt db
			if err := db.DeleteTodoById(bdb, projectName, id); err != nil {
				log.Printf("Failed to delete from DB: %v", err)
			}
		}
	}
	return nil
}

func viewTodos(bdb *bolt.DB) {
	err := bdb.View(func(tx *bolt.Tx) error {
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
