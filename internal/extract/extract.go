package extract

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/K-Road/extract_todos/config"
	"github.com/K-Road/extract_todos/internal/db"
	"github.com/K-Road/extract_todos/web"
	bolt "go.etcd.io/bbolt"
)

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

func internalRun(updateProgress func(p float64)) error {
	time.Sleep(2 * time.Second)
	//var todos []Todo
	var scannedTodos []config.Todo
	//TODO Implement this is a flag
	root := "/home/chrode/workspace/github.com/K-Road/discord-moodbot/"
	projectName := filepath.Base(strings.TrimRight(root, string(os.PathSeparator)))
	log.Println(root)
	log.Println(projectName)

	wasServerRunning := web.IsWebServerRunning()
	if wasServerRunning {
		log.Println("Webserver is running, stopping it to avoid conflicts...")
		if err := web.StopWebServer(); err != nil {
			return fmt.Errorf("Failed to stop webserver: %v", err)
		}
	}
	// if web.StopWebServer() != nil {
	// 	log.Println("Webserver not running or already stopped.")
	// }
	log.Println("Waiting DB lock...")
	for i := 0; i < 10; i++ {
		if !isDBLocked("todos.db") {
			log.Println("DB is not locked, proceeding...")
			break
		}
		log.Println("DB is locked, waiting...")
		time.Sleep(500 * time.Millisecond)
	}
	//time.Sleep(500 * time.Millisecond)
	//wasServerRunning := stopWebServer()

	log.Println("Opening DB... ")
	//Open bolt db
	bdb, err := bolt.Open("todos.db", 0600, nil)
	if err != nil {
		return err
	}
	defer bdb.Close()

	db.CheckDBVersionOrExit(bdb)

	var goFiles []string
	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err == nil || !info.IsDir() || strings.HasSuffix(path, ".go") {
			goFiles = append(goFiles, path)
		}
		return nil

	})
	if err != nil {
		return err
	}
	total := len(goFiles)
	current := 0

	for _, path := range goFiles {
		current++
		if updateProgress != nil {
			//log.Println("1")
			updateProgress(float64(current) / float64(total))
			//log.Println("2")
		}

		f, err := os.Open(path)
		if err != nil {
			return err
		}

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
					log.Printf("New TODO saved: %s:%d %s\n", todo.File, todo.Line, todo.Text)
				}
			}
			lineNum++
		}
		f.Close()

	}
	log.Println("Finished scan")
	//DEBUG to list all entries
	viewTodos(bdb)

	err = removeTodos(bdb, projectName, scannedTodos)

	viewTodos(bdb)

	//Restart webserver
	if wasServerRunning {
		log.Println("Restarting webserver...")
		if err := web.StartWebServerDetached(); err != nil {
			log.Printf("Failed to restart webserver: %v", err)
		}
	}
	updateProgress(1)
	return nil
}

func Run() error {
	return internalRun(nil)
}

func RunWithProgress(updateProgress func(p float64)) error {
	return internalRun(updateProgress)
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
			log.Printf("Project: %s\n", name)
			return b.ForEach(func(k, v []byte) error {
				log.Printf(" - %s\n", v)
				return nil
			})
		})
	})
	if err != nil {
		log.Println("Erroring reading from DB:", err)
	}
}

func isDBLocked(path string) bool {
	cmd := exec.Command("lsof", path)
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	lines := strings.Split(string(output), "\n")
	return len(lines) > 2 // more than header and one line means still open
}
