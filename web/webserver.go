package main

import (
	"encoding/json"
	"log"
	"net/http"
	"text/template"
	"time"

	"github.com/K-Road/extract_todos/config"
	"github.com/K-Road/extract_todos/db"
	bolt "go.etcd.io/bbolt"
)

type ResponseWithTime struct {
	TimeStamp string      `json:"timestamp"`
	Data      interface{} `json:"data"`
}

var templates *template.Template

func init() {
	templates = template.Must(template.ParseGlob("templates/*.html"))
}

func main() {
	//TODO handle flag for db name
	dbfile, err := bolt.Open("todos.db", 0666, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer dbfile.Close()

	cfg := &config.Config{DB: dbfile}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /projects", func(w http.ResponseWriter, r *http.Request) {
		projectsHandler(w, r, cfg)
	})
	mux.HandleFunc("GET /projects/{name}/todos", func(w http.ResponseWriter, r *http.Request) {
		todosHandler(w, r, cfg)
	})

	mux.HandleFunc("GET /api/projects", func(w http.ResponseWriter, r *http.Request) {
		apiProjectHandler(w, r, cfg)
	})
	mux.HandleFunc("GET /api/projects/{name}/todos", func(w http.ResponseWriter, r *http.Request) {
		apiTodosHandler(w, r, cfg)
	})

	log.Println("Starting Server on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

func projectsHandler(w http.ResponseWriter, r *http.Request, cfg *config.Config) {
	buckets, err := db.ListBuckets(cfg.DB)
	if err != nil {
		http.Error(w, "DB Error", http.StatusInternalServerError)
		return
	}

	err = templates.ExecuteTemplate(w, "projects.html", map[string]interface{}{
		"Projects": buckets,
	})
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
	}
}

func apiProjectHandler(w http.ResponseWriter, r *http.Request, cfg *config.Config) {
	buckets, err := db.ListBuckets(cfg.DB)
	if err != nil {
		http.Error(w, "DB Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	//TODO remove timestamps
	json.NewEncoder(w).Encode(map[string]interface{}{
		"timestamp": time.Now(),
		"projects":  buckets,
	})

}

func todosHandler(w http.ResponseWriter, r *http.Request, cfg *config.Config) {
	name := r.PathValue("name")
	todos, err := db.ListProjectTodos(cfg.DB, name)
	if err != nil {
		http.Error(w, "DB Error", http.StatusInternalServerError)
		return
	}

	err = templates.ExecuteTemplate(w, "todos.html", map[string]interface{}{
		"Project": name,
		"Todos":   todos,
	})
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
	}
}

func apiTodosHandler(w http.ResponseWriter, r *http.Request, cfg *config.Config) {
	name := r.PathValue("name")
	todos, err := db.ListProjectTodos(cfg.DB, name)
	if err != nil {
		http.Error(w, "DB Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	//TODO remove timestamp and project name
	json.NewEncoder(w).Encode(map[string]interface{}{
		"timestamp": time.Now(),
		"projects":  name,
		"todos":     todos,
	})
}
