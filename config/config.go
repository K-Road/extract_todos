package config

type Config struct {
	//DB *bolt.DB
	DataProvider DataProvider
}

// // TODO Move this to extract package?
// type Todo struct {
// 	File string
// 	Line int
// 	Text string
// }

type DataProvider interface {
	ListProjects() ([]string, error)
	ListProjectTodos(name string) ([]WebTodo, error)
	DeleteTodoById(name, id string) error
	GetActiveProject() (string, int, error)
	SetActiveProject(name string) error
	SaveTodo(projectName string, todo Todo) (TodoStatus, error)
	OpenDB(path string) error
	Close() error
	RemoveTodos(projectName string, scannedTodos []Todo) error
	//OpenDBWithOptions(path string, opts *bolt.Options) error
}

type ProviderFactory func() (DataProvider, error)
