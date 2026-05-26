package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
	"go.uber.org/zap"

	"github.com/lookingcamel/system-framework/internal/config"
	"github.com/lookingcamel/system-framework/internal/logger"
	"github.com/lookingcamel/system-framework/pkg/utils"
)

type RateLimiter struct {
	limiters map[string]*rate.Limiter
	mu       sync.RWMutex
	rps      float64
	burst    int
}

var rateLimiter *RateLimiter

func InitRateLimiter(cfg config.RateLimitConfig) {
	rateLimiter = &RateLimiter{
		limiters: make(map[string]*rate.Limiter),
		rps:      cfg.RequestsPerSecond,
		burst:    cfg.Burst,
	}

	if cfg.Enabled {
		logger.Log.Info("Rate limiter initialized",
			zap.Float64("requests_per_second", cfg.RequestsPerSecond),
			zap.Int("burst", cfg.Burst),
		)
	}
}

func (rl *RateLimiter) getLimiter(key string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	limiter, exists := rl.limiters[key]
	if !exists {
		limiter = rate.NewLimiter(rate.Limit(rl.rps), rl.burst)
		rl.limiters[key] = limiter
	}

	return limiter
}

func RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		if rateLimiter == nil {
			c.Next()
			return
		}

		if isExcludedPath(c.Request.URL.Path, rateLimiter.rps, rateLimiter.burst) {
			c.Next()
			return
		}

		key := getClientKey(c)
		limiter := rateLimiter.getLimiter(key)

		if !limiter.Allow() {
			requestID := utils.GinGetRequestID(c)
			logger.Log.Warn("Rate limit exceeded",
				zap.String("request_id", requestID),
				zap.String("client_ip", c.ClientIP()),
				zap.String("path", c.Request.URL.Path),
			)

			c.JSON(http.StatusTooManyRequests, gin.H{
				"code":    -1,
				"message": "Too many requests, please try again later",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

func getClientKey(c *gin.Context) string {
	if apiKey := c.GetHeader("X-API-Key"); apiKey != "" {
		return "api_key:" + apiKey
	}

	userID := utils.GinGetRequestID(c)
	if userID != "" {
		return "user:" + userID
	}

	return "ip:" + c.ClientIP()
}

func isExcludedPath(path string, rps float64, burst int) bool {
	excludePaths := []string{
		"/health",
		"/ready",
		"/metrics",
	}

	for _, excluded := range excludePaths {
		if path == excluded || path == excluded+"/" {
			return true
		}
	}

	return false
}

func CleanupRateLimiters() {
	if rateLimiter == nil {
		return
	}

	rateLimiter.mu.Lock()
	defer rateLimiter.mu.Unlock()

	ticker := time.NewTicker(10 * time.Minute)
	go func() {
		for range ticker.C {
			rateLimiter.mu.Lock()
			now := time.Now()
			for key, limiter := range rateLimiter.limiters {
				if limiter == nil {
					delete(rateLimiter.limiters, key)
					continue
				}

				_ = limiter.Allow()
			}
			_ = now
		}
	}()
}
