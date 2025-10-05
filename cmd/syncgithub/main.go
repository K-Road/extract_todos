package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/K-Road/extract_todos/config"
	"github.com/K-Road/extract_todos/internal/githubsync"
	"github.com/joho/godotenv"
)

type TodosResponse struct {
	Project   string        `json:"project"`
	TimeStamp string        `json:"timestamp"`
	Todos     []config.Todo `json:"todos"`
}

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	//TODO specify user.
	config.LoadUsersFromEnv()
}

func main() {
	//TODO parse url and projects
	retag := flag.String("retag", "", "Apply updates to existing issues...")
	flag.Parse()
	baseURL := "http://localhost:8080"
	project := "extract_todos"
	apikey := os.Getenv("API_KEY")

	todos, err := getTodos(baseURL, project, apikey)
	if err != nil {
		log.Fatal("Error fetching todos:", err)
	}

	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		log.Fatal("Missing GITHUB_TOKEN env var")
	}

	//TODO create hithub config struct for ctx, client, onwer, repo
	owner := "K-Road"
	repo := project
	ctx := context.Background()

	client := githubsync.NewGitHubClient()
	if client == nil {
		log.Fatal("Failed to create GitHub client")
	}

	//TODO failsafe check for token validation

	//Fetch existing issues to avoid duplicates
	existingIssues, err := githubsync.FetchAllOpenIssues(ctx, client, owner, repo)
	if err != nil {
		log.Fatalf("Failed to fetch existing issues: %v", err)
	}

	//TODO refactor into own function
	if *retag != "" {
		fmt.Println("Retagging existing issues...")

		//uncomment the line below to update existing issues with new labels
		updateValue := []string{}
		//updateValue := []string{"todo", "sync-generated"} //labels to add
		//updateValue := []string{"Created from local todos - Line#"} //body to add
		if len(updateValue) == 0 {
			log.Print("Update value is nil, skipping update...\n")
			return
		}
		for _, issue := range existingIssues {
			err := githubsync.UpdateIssueIfNeeded(ctx, client, owner, repo, issue, *retag, updateValue)
			if err != nil {
				log.Printf("Failed to update issue: %d. Error: %v\n", issue.GetNumber(), err)
			} else {
				log.Printf("Updated issue #%d with %s", issue.GetNumber(), *retag)
			}
		}
		return
	}
	//Create github issues
	for _, todo := range todos {
		title := fmt.Sprintf("%s:%s", todo.File, todo.Text)
		body := fmt.Sprintf("Created from local todos - Line# %d:\n", todo.Line)
		err := githubsync.CreateIssueIfNotExists(ctx, client, owner, repo, title, body, existingIssues)
		if err != nil {
			log.Printf("Failed to create issue for: %s. Error: %v\n", title, err)
			continue
		}
	}

	//Close deleted todos
	//todos
	//existingIssues
	todoSet := make(map[string]struct{})
	for _, t := range todos {
		title := fmt.Sprintf("%s:%s", t.File, t.Text)
		todoSet[title] = struct{}{}
	}

	githubsync.CloseDeletedTodos(ctx, client, owner, repo, project, todoSet)

}

func getTodos(baseURL, project, apikey string) ([]config.Todo, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/projects/%s/todos", baseURL, project), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-API-Key", apikey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("error: %s\n%s", resp.Status, body)
	}

	var data TodosResponse
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, fmt.Errorf("error decoding JSON: %v", err)
	}

	return data.Todos, nil

}
