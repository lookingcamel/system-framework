package server

import (
	"context"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"

	"github.com/lookingcamel/system-framework/internal/auth"
	"github.com/lookingcamel/system-framework/internal/circuitbreaker"
	"github.com/lookingcamel/system-framework/internal/config"
	"github.com/lookingcamel/system-framework/internal/handler"
	"github.com/lookingcamel/system-framework/internal/logger"
	"github.com/lookingcamel/system-framework/internal/middleware"
	"github.com/lookingcamel/system-framework/internal/versioning"
)

type Server struct {
	engine *gin.Engine
	cfg    *config.Config
	server *http.Server
}

func New(cfg *config.Config) (*Server, error) {
	if cfg.App.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.New()

	engine.Use(middleware.RequestID())
	engine.Use(middleware.Recovery())
	engine.Use(middleware.CORS())
	engine.Use(middleware.RateLimit())
	middleware.InitBodyLimit(cfg.Server)
	engine.Use(middleware.BodyLimit())

	if cfg.Auth.Enabled {
		if cfg.Auth.SignatureEnabled {
			engine.Use(auth.Signature())
		}
		if cfg.Auth.JWTEnabled {
			engine.Use(auth.JWT())
		}
	}

	if cfg.Logging.Level != "" {
		engine.Use(middleware.Logger())
	}

	if cfg.Prometheus.Enabled {
		engine.GET(cfg.Prometheus.Path, gin.WrapH(promhttp.Handler()))
	}

	if cfg.PProf.Enabled && cfg.App.Mode != "release" {
		pprofPath := cfg.PProf.Path
		engine.Any(pprofPath+"/*any", gin.WrapH(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.DefaultServeMux.ServeHTTP(w, r)
		})))
		logger.Log.Info("PProf enabled (development mode only)",
			zap.String("path", pprofPath),
			zap.String("mode", cfg.App.Mode),
		)
	} else if cfg.App.Mode == "release" {
		logger.Log.Info("PProf disabled in production mode")
	}

	if cfg.Health.Enabled {
		healthHandler := handler.NewHealthHandler()
		engine.GET(cfg.Health.Path, healthHandler.Health)
		engine.GET("/ready", healthHandler.Ready)
	}

	if cfg.CircuitBreaker.Enabled {
		engine.GET("/circuit/status", circuitbreaker.GetCircuitStatus())
	}

	versioning.InitAPIVersion(cfg.APIVersion)
	engine.Use(versioning.APIVersion())

	vr := versioning.NewVersionedRouter(engine)
	exampleHandler := handler.NewExampleHandler()

	for _, version := range cfg.APIVersion.SupportedVersions {
		vr.RegisterVersion(version)

		switch version {
		case "v1":
			vr.RegisterRoute("v1", "GET", "/example/:id", exampleHandler.GetExample)
			vr.RegisterRoute("v1", "GET", "/examples", exampleHandler.ListExample)
		}
	}

	return &Server{
		engine: engine,
		cfg:    cfg,
	}, nil
}

func (s *Server) Start() error {
	logger.Log.Info("Starting server",
		zap.String("address", s.cfg.GetAddress()),
	)

	s.server = &http.Server{
		Addr:         s.cfg.GetAddress(),
		Handler:      s.engine,
		ReadTimeout:  time.Duration(s.cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(s.cfg.Server.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(s.cfg.Server.IdleTimeout) * time.Second,
	}

	go func() {
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Log.Error("Server failed to start",
				zap.Error(err),
			)
			os.Exit(1)
		}
	}()

	s.waitForShutdown()
	return nil
}

func (s *Server) waitForShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Log.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(s.cfg.GracefulShutdown.Timeout)*time.Second)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		logger.Log.Error("Server forced to shutdown",
			zap.Error(err),
		)
	}

	logger.Sync()

	fmt.Println("Server exited")
}

func (s *Server) Engine() *gin.Engine {
	return s.engine
}
