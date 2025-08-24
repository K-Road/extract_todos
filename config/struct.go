package config

import "time"

type Todo struct {
	ID        int64
	ProjectID int64
	File      string
	Line      int
	Text      string
	Hash      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Project struct {
	ID        int64
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
	Active    bool
}
