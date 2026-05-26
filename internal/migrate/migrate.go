package migrate

import (
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/lookingcamel/system-framework/internal/config"
)

func RunMigrations(cfg config.DatabaseConfig) error {
	var dsn string
	switch cfg.Type {
	case "mysql":
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
			cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Name)
	case "postgres", "postgresql":
		dsn = fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
			cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Name)
	case "sqlite":
		dsn = cfg.SQLitePath
	default:
		dsn = cfg.SQLitePath
	}

	m, err := migrate.New(
		"file://"+cfg.MigrationPath,
		fmt.Sprintf("%s://%s", cfg.Type, dsn),
	)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}

func RollbackMigrations(cfg config.DatabaseConfig, steps int) error {
	var dsn string
	switch cfg.Type {
	case "mysql":
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
			cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Name)
	case "postgres", "postgresql":
		dsn = fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
			cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Name)
	case "sqlite":
		dsn = cfg.SQLitePath
	default:
		dsn = cfg.SQLitePath
	}

	m, err := migrate.New(
		"file://"+cfg.MigrationPath,
		fmt.Sprintf("%s://%s", cfg.Type, dsn),
	)
	if err != nil {
		return err
	}

	if err := m.Steps(-steps); err != nil {
		return err
	}

	return nil
}
