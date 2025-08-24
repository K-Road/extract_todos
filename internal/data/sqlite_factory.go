package data

import (
	"database/sql"

	"github.com/K-Road/extract_todos/config"
	_ "modernc.org/sqlite"
)

func SQLiteFactory(dbfile string) func() (config.DataProvider, error) {
	return func() (config.DataProvider, error) {
		db, err := sql.Open("sqlite", dbfile)
		if err != nil {
			return nil, err
		}
		sp := &SQLiteProvider{
			DB: db,
		}
		return sp, nil
	}
}
