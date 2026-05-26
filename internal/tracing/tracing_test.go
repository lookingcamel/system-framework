package tracing

import (
	"context"
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
	cfg := config.TracingConfig{
		Enabled:     false,
		ServiceName: "test-service",
	}

	shutdown, err := Init(cfg)
	if err != nil {
		t.Errorf("Init should not fail when disabled: %v", err)
	}
	if shutdown == nil {
		t.Error("Shutdown function should not be nil")
	}
}

func TestTracer(t *testing.T) {
	cfg := config.TracingConfig{
		Enabled:     false,
		ServiceName: "test-service",
	}

	Init(cfg)

	tracer := Tracer()
	if tracer == nil {
		t.Error("Tracer should not be nil")
	}
}

func TestStartSpan(t *testing.T) {
	cfg := config.TracingConfig{
		Enabled:     false,
		ServiceName: "test-service",
	}

	Init(cfg)

	ctx := context.Background()
	_, span := StartSpan(ctx, "test-span")

	if span == nil {
		t.Error("Span should not be nil")
	}
}
