package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Todo struct {
	File string
	Line int
	Text string
}

func main() {
	var todos []Todo
	root := "/home/chrode/workspace/github.com/K-Road/discord-moodbot/"
	fmt.Println(root)
	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			//fmt.Println("Error:", err)
			return nil
		}

		if info.IsDir() {
			//fmt.Println("Skipping directory:", path)
			return nil
		}
		if !strings.HasSuffix(path, ".go") {
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
			if strings.HasPrefix(trimmed, "//TODO") {
				todos = append(todos, Todo{
					File: path,
					Line: lineNum,
					Text: strings.TrimSpace(trimmed[len("//TODO"):]),
				})
			}
			lineNum++
		}
		return nil
	})

	for _, todo := range todos {
		fmt.Printf("gh issue create --title \"%s\" --body \"Found in %s:%d\"\n", todo.Text, todo.File, todo.Line)
	}
}
