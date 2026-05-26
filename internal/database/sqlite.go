package database

import (
	"database/sql"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"

	"github.com/lookingcamel/system-framework/internal/config"
)

func initSQLite(cfg config.DatabaseConfig) (*sql.DB, error) {
	path := cfg.SQLitePath
	if path == "" {
		path = "./data/example_db.sqlite"
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	return db, nil
}
