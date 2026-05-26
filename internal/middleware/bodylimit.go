package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/lookingcamel/system-framework/internal/config"
	"github.com/lookingcamel/system-framework/internal/logger"
	"github.com/lookingcamel/system-framework/pkg/utils"
)

var maxBodySize int64

func InitBodyLimit(cfg config.ServerConfig) {
	maxBodySize = cfg.MaxBodySize
	logger.Log.Info("Body limit middleware initialized",
		zap.Int64("max_body_size", maxBodySize),
		zap.String("max_body_size_formatted", formatBytes(maxBodySize)),
	)
}

func BodyLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		if maxBodySize <= 0 {
			c.Next()
			return
		}

		if isExcludedPathForBodyLimit(c.Request.URL.Path) {
			c.Next()
			return
		}

		requestID := utils.GinGetRequestID(c)

		if c.Request.ContentLength > maxBodySize {
			logger.Log.Warn("Request body too large",
				zap.String("request_id", requestID),
				zap.String("client_ip", c.ClientIP()),
				zap.String("path", c.Request.URL.Path),
				zap.Int64("content_length", c.Request.ContentLength),
				zap.Int64("max_body_size", maxBodySize),
			)

			c.JSON(http.StatusRequestEntityTooLarge, gin.H{
				"code":    -1,
				"message": "Request body too large",
			})
			c.Abort()
			return
		}

		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxBodySize)

		c.Next()
	}
}

func isExcludedPathForBodyLimit(path string) bool {
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

func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return "0 B"
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return "0 B"
}
