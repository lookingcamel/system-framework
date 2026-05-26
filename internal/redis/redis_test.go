package redis

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

func TestInit_Disabled(t *testing.T) {
	cfg := config.RedisConfig{
		Enabled: false,
	}

	err := Init(cfg)
	if err != nil {
		t.Errorf("Init should not fail when disabled: %v", err)
	}
	if Client != nil {
		t.Error("Client should be nil when disabled")
	}
}

func TestClose_NilClient(t *testing.T) {
	Client = nil
	Close()
}

func TestClose(t *testing.T) {
	cfg := config.RedisConfig{
		Enabled: true,
		Host:    "localhost",
		Port:    6379,
		Password: "",
		DB:      0,
		PoolSize: 10,
		MinIdleConns: 5,
	}

	err := Init(cfg)
	if err != nil {
		t.Skip("Redis not available, skipping test")
	}

	Close()
}
