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
	"github.com/K-Road/extract_todos/internal/logging"
	"github.com/K-Road/extract_todos/web"
)

func getLog() *log.Logger {
	return logging.Extract()
}

func hashTodo(todo config.Todo) string {
	s := fmt.Sprintf("%s:%s", todo.File, todo.Text)
	h := sha256.Sum256([]byte(s))
	return hex.EncodeToString(h[:])
}

func internalRun(project string, dp config.DataProvider, updateProgress func(p float64)) error {
	time.Sleep(2 * time.Second)
	//var todos []Todo
	var scannedTodos []config.Todo
	//TODO Implement this is a flag
	//root := "/home/chrode/workspace/github.com/K-Road/extract_todos/"
	root := "/home/chrode/workspace/github.com/K-Road/"
	//projectName := filepath.Base(strings.TrimRight(root, string(os.PathSeparator)))

	projectRoot := root + project

	getLog().Println(root)
	getLog().Println(project)
	getLog().Println(projectRoot)

	wasServerRunning := web.IsWebServerRunning()
	//Dont need to force stop webserver anymore

	// if wasServerRunning {
	// 	getLog().Println("Webserver is running, stopping it to avoid conflicts...")
	// 	if err := web.StopWebServer(); err != nil {
	// 		return fmt.Errorf("Failed to stop webserver: %v", err)
	// 	}
	// }
	// if web.StopWebServer() != nil {
	// 	getLog().Println("Webserver not running or already stopped.")
	// }

	//Dont need to checm on db lock anymore
	// getLog().Println("Waiting DB lock...")
	// for i := 0; i < 10; i++ {
	// 	if !isDBLocked("todos.sqlite") {
	// 		getLog().Println("DB is not locked, proceeding...")
	// 		break
	// 	}
	// 	getLog().Println("DB is locked, waiting...")
	// 	time.Sleep(500 * time.Millisecond)
	// }

	//time.Sleep(500 * time.Millisecond)
	//wasServerRunning := stopWebServer()

	// getLog().Println("Opening DB... ")
	// //Open bolt db
	// bdb, err := bolt.Open("todos.db", 0600, nil)
	// if err != nil {
	// 	logging.ExitWithError(getLog(), "Failed to open database:", err)
	// }
	// defer bdb.Close()

	// getLog().Println("Checking DB version...")
	// if err = db.CheckDBVersionOrExit(bdb); err != nil {
	// 	logging.ExitWithError(getLog(), "DB version check failed:", err)
	// }

	var goFiles []string
	err := filepath.Walk(projectRoot, func(path string, info os.FileInfo, err error) error {
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
			//getLog().Println("1")
			updateProgress(float64(current) / float64(total))
			//getLog().Println("2")
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
			relPath, err := filepath.Rel(projectRoot, path)
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
				status, err := dp.SaveTodo(project, todo)
				if err != nil {
					getLog().Println("Failed to save TODO:", err)
				}

				switch status {
				case config.TodoInserted:
					getLog().Printf("New TODO saved: %s:%d %s\n", todo.File, todo.Line, todo.Text)
				case config.TodoUpdated:
					getLog().Printf("Updated TODO line number: %s:%d %s\n", todo.File, todo.Line, todo.Text)
				}
			}
			lineNum++
		}
		f.Close()

	}
	getLog().Println("Finished scan")
	//DEBUG to list all entries
	getLog().Println("DEBUG: Listing all entries in DB")
	//viewTodos(bdb)

	err = dp.RemoveTodos(project, scannedTodos)

	getLog().Println("DEBUG: Listing all entries in DB after removal")
	//viewTodos(bdb)

	//Restart webserver
	if wasServerRunning {
		getLog().Println("Webserver was running..Do some cache logic refresh??")
		// if err := web.StartWebServerDetached(); err != nil {
		// 	getLog().Printf("Failed to restart webserver: %v", err)
		// }
	}
	updateProgress(1)
	return nil
}

func Run(project string, dp config.DataProvider) error {
	return internalRun(project, dp, nil)
}

func RunWithProgress(project string, dp config.DataProvider, updateProgress func(p float64)) error {
	return internalRun(project, dp, updateProgress)
}

// func removeTodos(db *sql.DB, projectName string, scannedTodos []config.Todo) error {
// 	storedTodos, err := db.FetchProjectTodos(db, projectName)
// 	if err != nil {
// 		return fmt.Errorf("failed to fetch todos for project %s: %w", projectName, err)
// 	}
// 	scannedIDs := make(map[string]struct{})
// 	for _, todo := range scannedTodos {
// 		id := hashTodo(todo)
// 		scannedIDs[id] = struct{}{}
// 	}

// 	for _, todo := range storedTodos {
// 		id := hashTodo(todo)
// 		if _, exists := scannedIDs[id]; !exists {
// 			getLog().Printf("Detected deleted TODO: %s:%s", todo.File, todo.Text)

// 			// //Delete from bolt db
// 			if err := db.DeleteTodoById(db, projectName, id); err != nil {
// 				getLog().Printf("Failed to delete from DB: %v", err)
// 			}
// 		}
// 	}
// 	return nil
// }

// func viewTodos(dp data.DataProvider) {
// 	err := dp.View(func(tx *bolt.Tx) error {
// 		return tx.ForEach(func(name []byte, b *bolt.Bucket) error {
// 			getLog().Printf("Project: %s\n", name)
// 			return b.ForEach(func(k, v []byte) error {
// 				getLog().Printf(" - %s\n", v)
// 				return nil
// 			})
// 		})
// 	})
// 	if err != nil {
// 		getLog().Println("Erroring reading from DB:", err)
// 	}
// }

func isDBLocked(path string) bool {
	cmd := exec.Command("lsof", path)
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	lines := strings.Split(string(output), "\n")
	return len(lines) > 2 // more than header and one line means still open
}
