package config

import "time"

type Todo struct {
	ID        int
	ProjectID int
	File      string
	Line      int
	Text      string
	Hash      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type WebTodo struct {
	ID   int    `json:"id"`
	File string `json:"file"`
	Line int    `json:"line"`
	Text string `json:"text"`
}

type Project struct {
	ID        int
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
	Active    bool
}
type TodoStatus int

const (
	TodoUnchanged TodoStatus = iota
	TodoUpdated
	TodoInserted
)
