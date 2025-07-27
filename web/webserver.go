package web

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
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
var server *http.Server
var dbfile *bolt.DB

func init() {
	templates = template.Must(template.ParseGlob("templates/*.html"))
}

func StartServer() {
	var err error
	//TODO handle flag for db name
	dbfile, err = bolt.Open("todos.db", 0666, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Write PID file here
	pid := os.Getpid()
	pidStr := strconv.Itoa(pid)
	if err := os.WriteFile(pidFile, []byte(pidStr), 0644); err != nil {
		log.Fatalf("Failed to write PID file: %v", err)
	}
	log.Printf("Webserver started with PID %s", pidStr)

	cfg := &config.Config{DB: dbfile}

	mux := http.NewServeMux()
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	mux.HandleFunc("GET /projects", func(w http.ResponseWriter, r *http.Request) {
		projectsHandler(w, r, cfg)
	})
	mux.HandleFunc("GET /projects/{name}/todos", func(w http.ResponseWriter, r *http.Request) {
		todosHandler(w, r, cfg)
	})

	apiMux := http.NewServeMux()
	apiMux.HandleFunc("/api/projects", func(w http.ResponseWriter, r *http.Request) {
		apiProjectHandler(w, r, cfg)
	})
	apiMux.HandleFunc("/api/projects/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/todos") {
			apiTodosHandler(w, r, cfg)
			return
		}
		http.NotFound(w, r)
	})
	mux.Handle("/api/", AuthenticateAPIKey(apiMux))

	mux.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	server = &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	//shutdown handling
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		//run webserver
		log.Println("Starting Server on :8080")
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	<-stop
	ShutdownServer()
}

func ShutdownServer() {
	log.Println("Shutting down webserver...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Forced shutdown of webserver: %v", err)
	}

	if dbfile != nil {
		log.Println("Closing database connection...")
		dbfile.Close()
	}
	_ = os.Remove(pidFile)
	//server.Shutdown(ctx)
	log.Println("Webserver stopped successfully")
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
