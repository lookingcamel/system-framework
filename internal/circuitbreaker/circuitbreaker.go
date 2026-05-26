package circuitbreaker

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/afex/hystrix-go/hystrix"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/lookingcamel/system-framework/internal/config"
	"github.com/lookingcamel/system-framework/internal/logger"
)

type CircuitBreakerConfig struct {
	Enabled                bool    `mapstructure:"enabled"`
	DefaultTimeout         int     `mapstructure:"default_timeout"`
	DefaultMaxConcurrent   int     `mapstructure:"default_max_concurrent"`
	DefaultErrorPercentage int     `mapstructure:"default_error_percentage"`
	DefaultRequestVolume   int     `mapstructure:"default_request_volume"`
}

func Init(cfg config.CircuitBreakerConfig) {
	if !cfg.Enabled {
		logger.Log.Info("Circuit breaker is disabled")
		return
	}

	hystrix.DefaultTimeout = cfg.DefaultTimeout
	hystrix.DefaultMaxConcurrent = cfg.DefaultMaxConcurrent
	hystrix.DefaultErrorPercentThreshold = cfg.DefaultErrorPercentage

	logger.Log.Info("Circuit breaker initialized",
		zap.Int("timeout", cfg.DefaultTimeout),
		zap.Int("max_concurrent", cfg.DefaultMaxConcurrent),
		zap.Int("error_percentage", cfg.DefaultErrorPercentage),
		zap.Int("request_volume", cfg.DefaultRequestVolume),
	)

	hystrixStreamHandler := hystrix.NewStreamHandler()
	hystrixStreamHandler.Start()

	go func() {
		http.Handle("/hystrix.stream", hystrixStreamHandler)
		logger.Log.Info("Hystrix stream handler started on /hystrix.stream")
	}()
}

func Execute(commandName string, run func() error, fallback func(error) error) error {
	return hystrix.Do(commandName, run, fallback)
}

func ExecuteWithResult[T any](commandName string, run func() (T, error), fallback func(error) (T, error)) (T, error) {
	var result T
	err := hystrix.Do(commandName, func() error {
		var err error
		result, err = run()
		return err
	}, func(err error) error {
		var fallbackErr error
		result, fallbackErr = fallback(err)
		return fallbackErr
	})
	return result, err
}

func GinMiddleware(commandName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		done := make(chan bool, 1)
		errChan := make(chan error, 1)

		go func() {
			err := hystrix.Do(commandName, func() error {
				c.Next()
				if c.Writer.Status() >= 500 {
					return &serverError{}
				}
				return nil
			}, func(err error) error {
				return err
			})
			if err != nil {
				errChan <- err
				return
			}
			done <- true
		}()

		select {
		case <-done:
			return
		case err := <-errChan:
			logger.Log.Warn("Circuit breaker tripped",
				zap.String("command", commandName),
				zap.Error(err),
			)
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"code":    -1,
				"message": "Service unavailable, circuit breaker is open",
				"data":    nil,
			})
			c.Abort()
		case <-time.After(time.Duration(hystrix.DefaultTimeout) * time.Millisecond):
			c.JSON(http.StatusGatewayTimeout, gin.H{
				"code":    -1,
				"message": "Request timeout",
				"data":    nil,
			})
			c.Abort()
		}
	}
}

func GinMiddlewareWithFallback(commandName string, fallback gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		done := make(chan bool, 1)
		errChan := make(chan error, 1)

		go func() {
			err := hystrix.Do(commandName, func() error {
				c.Next()
				if c.Writer.Status() >= 500 {
					return &serverError{}
				}
				return nil
			}, func(err error) error {
				return err
			})
			if err != nil {
				errChan <- err
				return
			}
			done <- true
		}()

		select {
		case <-done:
			return
		case <-errChan:
			fallback(c)
			c.Abort()
		case <-time.After(time.Duration(hystrix.DefaultTimeout) * time.Millisecond):
			fallback(c)
			c.Abort()
		}
	}
}

type serverError struct{}

func (e *serverError) Error() string {
	return "server error"
}

func GetCircuitStatus() gin.HandlerFunc {
	return func(c *gin.Context) {
		status := gin.H{
			"message": "Circuit breaker status endpoint",
			"note":    "Detailed metrics require Hystrix Dashboard",
		}
		c.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "success",
			"data":    status,
		})
	}
}

func GetCircuitStatusJSON() ([]byte, error) {
	status := map[string]interface{}{
		"message": "Circuit breaker status",
		"note":    "Detailed metrics require Hystrix Dashboard",
	}
	return json.Marshal(status)
}

func ExecuteWithContext(ctx context.Context, commandName string, run func(context.Context) error, fallback func(context.Context, error) error) error {
	return hystrix.Do(commandName, func() error {
		return run(ctx)
	}, func(err error) error {
		return fallback(ctx, err)
	})
}