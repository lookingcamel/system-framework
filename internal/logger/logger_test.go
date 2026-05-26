package logger

import (
	"testing"

	"github.com/lookingcamel/system-framework/internal/config"
)

func TestInit_ConsoleOutput(t *testing.T) {
	cfg := config.LoggingConfig{
		Level:  "info",
		Format: "console",
		Output: "console",
	}

	err := Init(cfg)
	if err != nil {
		t.Errorf("Init failed: %v", err)
	}
	if Log == nil {
		t.Error("Log should not be nil")
	}
}

func TestInit_JsonFormat(t *testing.T) {
	cfg := config.LoggingConfig{
		Level:  "info",
		Format: "json",
		Output: "console",
	}

	err := Init(cfg)
	if err != nil {
		t.Errorf("Init failed: %v", err)
	}
	if Log == nil {
		t.Error("Log should not be nil")
	}
}

func TestInit_DebugLevel(t *testing.T) {
	cfg := config.LoggingConfig{
		Level:  "debug",
		Format: "json",
		Output: "console",
	}

	err := Init(cfg)
	if err != nil {
		t.Errorf("Init failed: %v", err)
	}
}

func TestInit_WarnLevel(t *testing.T) {
	cfg := config.LoggingConfig{
		Level:  "warn",
		Format: "console",
		Output: "console",
	}

	err := Init(cfg)
	if err != nil {
		t.Errorf("Init failed: %v", err)
	}
}

func TestInit_ErrorLevel(t *testing.T) {
	cfg := config.LoggingConfig{
		Level:  "error",
		Format: "console",
		Output: "console",
	}

	err := Init(cfg)
	if err != nil {
		t.Errorf("Init failed: %v", err)
	}
}

func TestInit_InvalidLevel(t *testing.T) {
	cfg := config.LoggingConfig{
		Level:  "invalid",
		Format: "console",
		Output: "console",
	}

	err := Init(cfg)
	if err != nil {
		t.Errorf("Init failed: %v", err)
	}
}

func TestInit_DefaultOutput(t *testing.T) {
	cfg := config.LoggingConfig{
		Level:  "info",
		Format: "json",
		Output: "unknown",
	}

	err := Init(cfg)
	if err != nil {
		t.Errorf("Init failed: %v", err)
	}
}

func TestInit_StdoutOutput(t *testing.T) {
	cfg := config.LoggingConfig{
		Level:  "info",
		Format: "json",
		Output: "stdout",
	}

	err := Init(cfg)
	if err != nil {
		t.Errorf("Init failed: %v", err)
	}
}

func TestInit_TerminalOutput(t *testing.T) {
	cfg := config.LoggingConfig{
		Level:  "info",
		Format: "json",
		Output: "terminal",
	}

	err := Init(cfg)
	if err != nil {
		t.Errorf("Init failed: %v", err)
	}
}

func TestSync(t *testing.T) {
	cfg := config.LoggingConfig{
		Level:  "info",
		Format: "console",
		Output: "console",
	}

	Init(cfg)
	Sync()
}

func TestSync_NilLogger(t *testing.T) {
	Log = nil
	Sync()
}
