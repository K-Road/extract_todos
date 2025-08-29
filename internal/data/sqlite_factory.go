package data

import (
	"github.com/K-Road/extract_todos/config"
	"github.com/K-Road/extract_todos/internal/db"
	_ "modernc.org/sqlite"
)

func SQLiteFactory(dbfile string) func() (config.DataProvider, error) {
	return func() (config.DataProvider, error) {
		dbconn, err := db.OpenDB(dbfile)
		if err != nil {
			return nil, err
		}
		sp := &SQLiteProvider{
			DB: dbconn,
		}
		return sp, nil
	}
}
