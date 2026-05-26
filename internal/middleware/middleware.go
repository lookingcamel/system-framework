package middleware

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"

	"github.com/lookingcamel/system-framework/internal/logger"
	"github.com/lookingcamel/system-framework/internal/metrics"
	"github.com/lookingcamel/system-framework/internal/tracing"
	"github.com/lookingcamel/system-framework/pkg/utils"
)

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = utils.GenerateRequestID()
		}

		utils.GinSetRequestID(c, requestID)

		c.Writer.Header().Set("X-Request-ID", requestID)

		ctx := utils.SetRequestID(c.Request.Context(), requestID)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		requestID := utils.GinGetRequestID(c)

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()

		if raw != "" {
			path = path + "?" + raw
		}

		logger.Log.Info("HTTP request",
			zap.String("request_id", requestID),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.Int("status", status),
			zap.Duration("latency", latency),
			zap.String("client_ip", c.ClientIP()),
			zap.String("user_agent", c.Request.UserAgent()),
		)
	}
}

func Metrics() gin.HandlerFunc {
	return func(c *gin.Context) {
		metrics.IncrementInFlight()
		start := time.Now()

		c.Next()

		duration := time.Since(start).Seconds()
		status := strconv.Itoa(c.Writer.Status())
		path := c.FullPath()
		if path == "" {
			path = "unknown"
		}

		metrics.RecordRequest(
			c.Request.Method,
			path,
			status,
			duration,
		)
		metrics.DecrementInFlight()
	}
}

func Tracing() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := utils.GinGetRequestID(c)

		ctx, span := tracing.StartSpan(c.Request.Context(), c.Request.Method+" "+c.FullPath())
		defer span.End()

		if requestID != "" {
			span.SetAttributes(
				attribute.String("request_id", requestID),
			)
		}

		c.Request = c.Request.WithContext(ctx)

		c.Next()

		span.SetAttributes(
			attribute.Int("http.status_code", c.Writer.Status()),
		)
	}
}

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				requestID := utils.GinGetRequestID(c)
				logger.Log.Error("Panic recovered",
					zap.String("request_id", requestID),
					zap.Any("error", err),
				)
				c.AbortWithStatus(500)
			}
		}()
		c.Next()
	}
}

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, X-Request-ID")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "X-Request-ID")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}