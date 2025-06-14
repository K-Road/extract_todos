package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/K-Road/extract_todos/internal/githubsync"
	"github.com/joho/godotenv"
)

type TodosResponse struct {
	Projects  string   `json:"projects"`
	TimeStamp string   `json:"timestamp"`
	Todos     []string `json:"todos"`
}

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	baseURL := "http://localhost:8080"
	project := "discord-moodbot"

	todos, err := getTodos(baseURL, project)
	if err != nil {
		log.Fatal("Error fetching todos:", err)
	}

	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		log.Fatal("Missing GITHUB_TOKEN env var")
	}

	owner := "K-Road"
	repo := project

	for _, todo := range todos {
		err := githubsync.CreateIssue(githubsync.NewGitHubClient(), owner, repo, todo, "Created from local todos")
		if err != nil {
			log.Printf("Failed to create issue for: %s. Error: %v\n", todo, err)
			continue
		}
	}
}

func getTodos(baseURL, project string) ([]string, error) {
	resp, err := http.Get(fmt.Sprintf("%s/api/projects/%s/todos", baseURL, project))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data TodosResponse
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		log.Fatal("Error decoding JSON:", err)
	}

	return data.Todos, nil
}
