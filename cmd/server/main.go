package main

import (
	"context"
	"flag"
	"os"

	"go.uber.org/zap"

	"github.com/lookingcamel/system-framework/internal/auth"
	"github.com/lookingcamel/system-framework/internal/circuitbreaker"
	"github.com/lookingcamel/system-framework/internal/config"
	"github.com/lookingcamel/system-framework/internal/database"
	"github.com/lookingcamel/system-framework/internal/logger"
	"github.com/lookingcamel/system-framework/internal/metrics"
	"github.com/lookingcamel/system-framework/internal/redis"
	"github.com/lookingcamel/system-framework/internal/server"
	"github.com/lookingcamel/system-framework/internal/tracing"
)

func main() {
	configPath := flag.String("config", "config.yaml", "Path to configuration file")
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		panic(err)
	}

	if err := logger.Init(cfg.Logging); err != nil {
		panic(err)
	}
	defer logger.Sync()

	circuitbreaker.Init(cfg.CircuitBreaker)

	auth.Init(cfg.Auth)

	shutdown, err := tracing.Init(cfg.Tracing)
	if err != nil {
		logger.Log.Warn("Failed to initialize tracing",
			zap.Error(err),
		)
	} else {
		defer shutdown(context.Background())
	}

	if err := database.Init(cfg.Database); err != nil {
		logger.Log.Warn("Failed to initialize database, continuing without DB",
			zap.Error(err),
		)
	} else {
		defer database.Close()
	}

	if err := redis.Init(cfg.Redis); err != nil {
		logger.Log.Warn("Failed to initialize redis, continuing without Redis",
			zap.Error(err),
		)
	} else {
		defer redis.Close()
	}

	metrics.SetAppInfo(cfg.App.Name, cfg.App.Version)

	srv, err := server.New(cfg)
	if err != nil {
		logger.Log.Error("Failed to create server",
			zap.Error(err),
		)
		os.Exit(1)
	}

	logger.Log.Info("Application started",
		zap.String("app", cfg.App.Name),
		zap.String("version", cfg.App.Version),
	)

	if err := srv.Start(); err != nil {
		logger.Log.Error("Server error",
			zap.Error(err),
		)
		os.Exit(1)
	}
}
