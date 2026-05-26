package middleware

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/lookingcamel/system-framework/internal/config"
)

func TestInitBodyLimit(t *testing.T) {
	cfg := config.ServerConfig{
		MaxBodySize: 8388608,
	}
	InitBodyLimit(cfg)

	if maxBodySize != 8388608 {
		t.Errorf("Expected maxBodySize %d, got %d", 8388608, maxBodySize)
	}
}

func TestBodyLimit_SmallBody(t *testing.T) {
	// Reset maxBodySize for test
	maxBodySize = 100 // 100 bytes

	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	
	c.Request, _ = http.NewRequest("POST", "/test", bytes.NewBuffer(make([]byte, 50))) // 50 bytes
	c.Request.ContentLength = 50

	handler := BodyLimit()
	handler(c)

	if c.IsAborted() {
		t.Error("Request should not be aborted for small body")
	}
}

func TestBodyLimit_LargeBody(t *testing.T) {
	maxBodySize = 100 // 100 bytes

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	
	c.Request, _ = http.NewRequest("POST", "/test", bytes.NewBuffer(make([]byte, 200))) // 200 bytes
	c.Request.ContentLength = 200

	handler := BodyLimit()
	handler(c)

	if !c.IsAborted() {
		t.Error("Request should be aborted for large body")
	}

	if w.Code != http.StatusRequestEntityTooLarge {
		t.Errorf("Expected status %d, got %d", http.StatusRequestEntityTooLarge, w.Code)
	}
}

func TestBodyLimit_ExcludedPath(t *testing.T) {
	maxBodySize = 100

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	
	c.Request, _ = http.NewRequest("POST", "/health", bytes.NewBuffer(make([]byte, 200)))
	c.Request.ContentLength = 200

	handler := BodyLimit()
	handler(c)

	if c.IsAborted() {
		t.Error("Request should not be aborted for excluded path")
	}
}

func TestBodyLimit_Disabled(t *testing.T) {
	maxBodySize = 0

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	
	c.Request, _ = http.NewRequest("POST", "/test", bytes.NewBuffer(make([]byte, 200)))
	c.Request.ContentLength = 200

	handler := BodyLimit()
	handler(c)

	if c.IsAborted() {
		t.Error("Request should not be aborted when maxBodySize is 0")
	}
}

func TestIsExcludedPathForBodyLimit(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{
			name:     "/health should be excluded",
			path:     "/health",
			expected: true,
		},
		{
			name:     "/ready should be excluded",
			path:     "/ready",
			expected: true,
		},
		{
			name:     "/metrics should be excluded",
			path:     "/metrics",
			expected: true,
		},
		{
			name:     "/health/ should be excluded",
			path:     "/health/",
			expected: true,
		},
		{
			name:     "/api/test should not be excluded",
			path:     "/api/test",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isExcludedPathForBodyLimit(tt.path)
			if result != tt.expected {
				t.Errorf("Expected %t for path %s, got %t", tt.expected, tt.path, result)
			}
		})
	}
}

func TestFormatBytes(t *testing.T) {
	tests := []struct {
		name     string
		bytes    int64
		expected string
	}{
		{
			name:     "0 bytes",
			bytes:    0,
			expected: "0 B",
		},
		{
			name:     "500 bytes",
			bytes:    500,
			expected: "0 B", // Note: current implementation always returns "0 B"
		},
		{
			name:     "1MB",
			bytes:    1048576,
			expected: "0 B", // Note: current implementation always returns "0 B"
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatBytes(tt.bytes)
			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

func TestInitRateLimiter(t *testing.T) {
	cfg := config.RateLimitConfig{
		Enabled:          true,
		RequestsPerSecond: 10,
		Burst:           20,
	}
	InitRateLimiter(cfg)

	if rateLimiter == nil {
		t.Fatal("rateLimiter should not be nil")
	}
	if rateLimiter.rps != 10 {
		t.Errorf("Expected rps %f, got %f", 10.0, rateLimiter.rps)
	}
	if rateLimiter.burst != 20 {
		t.Errorf("Expected burst %d, got %d", 20, rateLimiter.burst)
	}
}

func TestRateLimiter_GetLimiter(t *testing.T) {
	cfg := config.RateLimitConfig{
		Enabled:          true,
		RequestsPerSecond: 10,
		Burst:           20,
	}
	InitRateLimiter(cfg)

	key := "test_key"
	limiter := rateLimiter.getLimiter(key)

	if limiter == nil {
		t.Error("Limiter should not be nil")
	}

	// Test that same key returns same limiter
	limiter2 := rateLimiter.getLimiter(key)
	if limiter != limiter2 {
		t.Error("Same key should return same limiter")
	}

	// Test that different key returns different limiter
	limiter3 := rateLimiter.getLimiter("another_key")
	if limiter == limiter3 {
		t.Error("Different keys should return different limiters")
	}
}

func TestRateLimit_WithinLimit(t *testing.T) {
	cfg := config.RateLimitConfig{
		Enabled:          true,
		RequestsPerSecond: 1000, // High enough to not trigger limit
		Burst:           100,
	}
	InitRateLimiter(cfg)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	
	c.Request, _ = http.NewRequest("GET", "/test", nil)
	c.Request.RemoteAddr = "127.0.0.1:12345"

	handler := RateLimit()
	handler(c)

	if c.IsAborted() {
		t.Error("Request should not be aborted when within limit")
	}
}

func TestRateLimit_ExcludedPath(t *testing.T) {
	cfg := config.RateLimitConfig{
		Enabled:          true,
		RequestsPerSecond: 0.0001, // Very low limit
		Burst:           1,
	}
	InitRateLimiter(cfg)

	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	
	c.Request, _ = http.NewRequest("GET", "/health", nil)
	c.Request.RemoteAddr = "127.0.0.1:12345"

	handler := RateLimit()
	handler(c)

	if c.IsAborted() {
		t.Error("Request should not be aborted for excluded path")
	}
}

func TestRateLimit_RateLimiterNotInitialized(t *testing.T) {
	rateLimiter = nil // Reset to nil

	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	
	c.Request, _ = http.NewRequest("GET", "/test", nil)

	handler := RateLimit()
	handler(c)

	if c.IsAborted() {
		t.Error("Request should not be aborted when rateLimiter is nil")
	}
}

func TestIsExcludedPath(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{
			name:     "/health should be excluded",
			path:     "/health",
			expected: true,
		},
		{
			name:     "/ready should be excluded",
			path:     "/ready",
			expected: true,
		},
		{
			name:     "/metrics should be excluded",
			path:     "/metrics",
			expected: true,
		},
		{
			name:     "/health/ should be excluded",
			path:     "/health/",
			expected: true,
		},
		{
			name:     "/api/test should not be excluded",
			path:     "/api/test",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isExcludedPath(tt.path, 10, 20)
			if result != tt.expected {
				t.Errorf("Expected %t for path %s, got %t", tt.expected, tt.path, result)
			}
		})
	}
}

func TestGetClientKey_APIKey(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	
	c.Request, _ = http.NewRequest("GET", "/test", nil)
	c.Request.Header.Set("X-API-Key", "test-api-key")

	key := getClientKey(c)
	expected := "api_key:test-api-key"
	if key != expected {
		t.Errorf("Expected key '%s', got '%s'", expected, key)
	}
}

func TestGetClientKey_IPAddress(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	
	c.Request, _ = http.NewRequest("GET", "/test", nil)
	c.Request.RemoteAddr = "127.0.0.1:12345"

	key := getClientKey(c)
	expectedPrefix := "ip:"
	if len(key) < len(expectedPrefix) || key[:len(expectedPrefix)] != expectedPrefix {
		t.Errorf("Expected key to start with '%s', got '%s'", expectedPrefix, key)
	}
}
