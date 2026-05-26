package handler

import (
	"testing"

	"go.uber.org/zap"

	"github.com/lookingcamel/system-framework/internal/logger"
)

func TestMain(m *testing.M) {
	// Initialize logger for tests
	logger.Log = zap.NewNop()
	m.Run()
}
