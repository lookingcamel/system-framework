package database

import (
	"database/sql"
	"time"

	"github.com/lookingcamel/system-framework/internal/backup"
	"github.com/lookingcamel/system-framework/internal/config"
	"github.com/lookingcamel/system-framework/internal/logger"
	"github.com/lookingcamel/system-framework/internal/migrate"
	"go.uber.org/zap"
)

var DB *sql.DB

func Init(cfg config.DatabaseConfig) error {
	var err error

	switch cfg.Type {
	case "mysql":
		DB, err = initMySQL(cfg)
	case "postgres", "postgresql":
		DB, err = initPostgreSQL(cfg)
	case "sqlite":
		DB, err = initSQLite(cfg)
	default:
		logger.Log.Warn("Unsupported database type, using SQLite as default",
			zap.String("type", cfg.Type),
		)
		DB, err = initSQLite(cfg)
	}

	if err != nil {
		return err
	}

	DB.SetMaxOpenConns(cfg.MaxOpenConnections)
	DB.SetMaxIdleConns(cfg.MaxIdleConnections)
	DB.SetConnMaxLifetime(time.Duration(cfg.ConnectionMaxLifetime) * time.Second)

	if err := DB.Ping(); err != nil {
		return err
	}

	logger.Log.Info("Database connection established",
		zap.String("type", cfg.Type),
	)

	if cfg.BackupEnabled {
		if backupFile, err := backup.BackupDatabase(DB, cfg.Type, cfg.BackupDir); err != nil {
			logger.Log.Warn("Failed to backup database",
				zap.Error(err),
			)
		} else {
			logger.Log.Info("Database backup created",
				zap.String("file", backupFile),
			)
		}
	}

	if cfg.MigrationPath != "" {
		if err := migrate.RunMigrations(cfg); err != nil {
			logger.Log.Warn("Failed to run database migrations",
				zap.Error(err),
			)
		} else {
			logger.Log.Info("Database migrations completed")
		}
	}

	return nil
}

func Close() {
	if DB != nil {
		if err := DB.Close(); err != nil {
			logger.Log.Warn("Failed to close database connection",
				zap.Error(err),
			)
		}
	}
}
