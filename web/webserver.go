package web

import (
	"context"
	"encoding/json"
	"fmt"

	"net/http"
	"os"
	"os/signal"
	"syscall"
	"text/template"
	"time"

	"github.com/K-Road/extract_todos/config"
	_ "github.com/K-Road/extract_todos/web/docs"
	"github.com/joho/godotenv"
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
	_ = godotenv.Load(".env")
	config.LoadUsersFromEnv()
}

func StartServer(factory config.ProviderFactory) {
	var err error

	dp, err := factory()
	if err != nil {
		getLog().Fatalf("Failed to create data provider: %v", err)
	}
	defer dp.Close()
	getLog().Println("Data provider initialized")
	//TODO load users for webserver authentication
	//config.LoadUsersFromEnv()
	// if err != nil {
	// 	getLog().Fatalf("Failed to load users from env: %v", err)
	// }

	mux := http.NewServeMux()
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	mux.HandleFunc("GET /projects", func(w http.ResponseWriter, r *http.Request) {
		projectsHandler(w, r, dp)
	})
	mux.HandleFunc("/projects/{name}/todos", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		todosHandler(w, r, dp)
	})

	apiMux := http.NewServeMux()
	apiMux.HandleFunc("/api/projects", func(w http.ResponseWriter, r *http.Request) {
		apiProjectHandler(w, r, dp)
	})
	apiMux.HandleFunc("/api/projects/{name}/todos", func(w http.ResponseWriter, r *http.Request) {
		apiTodosHandler(w, r, dp)
	})
	mux.Handle("/api/", AuthenticateAPIKey(apiMux))

	mux.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	server = &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	getLog().Println("server setup for 8080")
	//shutdown handling
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		//run webserver
		getLog().Println("Starting Server on :8080")
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			getLog().Fatalf("HTTP server error: %v", err)
		}
	}()

	<-stop
	ShutdownServer()
}

func ShutdownServer() {
	getLog().Println("Shutting down webserver...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		getLog().Printf("Forced shutdown of webserver: %v", err)
	}

	if dbfile != nil {
		getLog().Println("Closing database connection...")
		dbfile.Close()
	}
	_ = os.Remove(pidFile)
	//server.Shutdown(ctx)
	getLog().Println("Webserver stopped successfully")
}

func projectsHandler(w http.ResponseWriter, r *http.Request, dp config.DataProvider) {
	buckets, err := dp.ListProjects()
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
func apiProjectHandler(w http.ResponseWriter, r *http.Request, dp config.DataProvider) {
	buckets, err := dp.ListProjects()
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

func todosHandler(w http.ResponseWriter, r *http.Request, dp config.DataProvider) {
	name := r.PathValue("name")
	todos, err := dp.ListProjectTodos(name)
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
func apiTodosHandler(w http.ResponseWriter, r *http.Request, dp config.DataProvider) {
	name := r.PathValue("name")
	fmt.Println("Fetching todos for project:", name)
	todos, err := dp.ListProjectTodos(name)
	if err != nil {
		http.Error(w, fmt.Sprintf("DB Error %v", err), http.StatusInternalServerError)
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
