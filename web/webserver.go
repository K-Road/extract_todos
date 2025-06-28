package main

import (
	"encoding/json"
	"log"
	"net/http"
	"text/template"
	"time"

	"github.com/K-Road/extract_todos/config"
	"github.com/K-Road/extract_todos/internal/db"
	_ "github.com/K-Road/extract_todos/web/docs"
	httpSwagger "github.com/swaggo/http-swagger"
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
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

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

	mux.HandleFunc("/swagger/", httpSwagger.WrapHandler)

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

// @Summary Get all project names
// @Description Returns a list of project buckets
// @Tags Projects
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {string} string
// @Router /api/projects [get]
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
	todos, err := db.FetchProjectTodos(cfg.DB, name)
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

// @Summary Get all todos
// @Description Returns a list of project todos
// @Tags Todos
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {string} string
// @Router /api/projects/{name}/todos [get]
func apiTodosHandler(w http.ResponseWriter, r *http.Request, cfg *config.Config) {
	name := r.PathValue("name")
	todos, err := db.FetchProjectTodos(cfg.DB, name)
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
