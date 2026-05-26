package circuitbreaker

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
	cfg := config.CircuitBreakerConfig{
		Enabled: false,
	}
	Init(cfg)
}

func TestExecute(t *testing.T) {
	cfg := config.CircuitBreakerConfig{
		Enabled:                true,
		DefaultTimeout:         1000,
		DefaultMaxConcurrent:   10,
		DefaultErrorPercentage: 50,
	}
	Init(cfg)

	err := Execute("test-command", func() error {
		return nil
	}, func(err error) error {
		return err
	})

	if err != nil {
		t.Errorf("Execute failed: %v", err)
	}
}

func TestServerError_Error(t *testing.T) {
	err := &serverError{}
	if err.Error() != "server error" {
		t.Errorf("Expected 'server error', got '%s'", err.Error())
	}
}
