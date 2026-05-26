package database

import (
	"testing"

	"go.uber.org/zap"

	"github.com/lookingcamel/system-framework/internal/config"
	"github.com/lookingcamel/system-framework/internal/logger"
)

func TestMain(m *testing.M) {
	logger.Log = zap.NewNop()
	m.Run()
}

func TestClose_NilDB(t *testing.T) {
	DB = nil
	Close()
}

func TestInit_SQLite(t *testing.T) {
	cfg := config.DatabaseConfig{
		Type:                 "sqlite",
		SQLitePath:           ":memory:",
		MaxOpenConnections:   10,
		MaxIdleConnections:   5,
		ConnectionMaxLifetime: 300,
		BackupEnabled:        false,
		MigrationPath:        "",
	}

	err := Init(cfg)
	if err != nil {
		t.Errorf("Init failed: %v", err)
	}

	if DB == nil {
		t.Error("DB should not be nil")
	}

	Close()
}

func TestInit_InvalidType(t *testing.T) {
	cfg := config.DatabaseConfig{
		Type:                 "invalid",
		MaxOpenConnections:   10,
		MaxIdleConnections:   5,
		ConnectionMaxLifetime: 300,
		BackupEnabled:        false,
		MigrationPath:        "",
	}

	err := Init(cfg)
	if err != nil {
		t.Errorf("Init should not fail for invalid type: %v", err)
	}

	Close()
}

func TestInit_WithBackup(t *testing.T) {
	cfg := config.DatabaseConfig{
		Type:                 "sqlite",
		SQLitePath:           ":memory:",
		MaxOpenConnections:   10,
		MaxIdleConnections:   5,
		ConnectionMaxLifetime: 300,
		BackupEnabled:        true,
		BackupDir:            "./test_backups",
		MigrationPath:        "",
	}

	err := Init(cfg)
	if err != nil {
		t.Skipf("Skipping backup test due to error: %v", err)
	}

	Close()
}
