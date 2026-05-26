package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/lookingcamel/system-framework/internal/config"
	"github.com/lookingcamel/system-framework/internal/logger"
)

func TestMain(m *testing.M) {
	logger.Log = zap.NewNop()
	gin.SetMode(gin.TestMode)
	m.Run()
}

func TestInit_Enabled(t *testing.T) {
	cfg := config.AuthConfig{
		Enabled:          true,
		JWTEnabled:       true,
		SignatureEnabled: false,
	}
	Init(cfg)
}

func TestInit_Disabled(t *testing.T) {
	cfg := config.AuthConfig{
		Enabled:          false,
		JWTEnabled:       false,
		SignatureEnabled: false,
	}
	Init(cfg)
}

func TestGenerateToken(t *testing.T) {
	cfg := config.AuthConfig{
		JWTSecret:       "test-secret-key-1234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890",
		JWTIssuer:       "test-issuer",
		JWTAudience:     "test-audience",
		JWTExpireSeconds: 3600,
	}
	Init(cfg)

	token, err := GenerateToken("user123", "testuser", "admin")
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}
	if token == "" {
		t.Error("Token should not be empty")
	}
}

func TestParseToken_ValidToken(t *testing.T) {
	cfg := config.AuthConfig{
		JWTSecret:       "test-secret-key-1234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890",
		JWTIssuer:       "test-issuer",
		JWTAudience:     "test-audience",
		JWTExpireSeconds: 3600,
	}
	Init(cfg)

	token, err := GenerateToken("user123", "testuser", "admin")
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	claims, err := ParseToken(token)
	if err != nil {
		t.Fatalf("Failed to parse token: %v", err)
	}
	if claims.UserID != "user123" {
		t.Errorf("Expected UserID 'user123', got '%s'", claims.UserID)
	}
	if claims.Username != "testuser" {
		t.Errorf("Expected Username 'testuser', got '%s'", claims.Username)
	}
	if claims.Role != "admin" {
		t.Errorf("Expected Role 'admin', got '%s'", claims.Role)
	}
}

func TestParseToken_InvalidToken(t *testing.T) {
	cfg := config.AuthConfig{
		JWTSecret:       "test-secret-key-1234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890",
		JWTIssuer:       "test-issuer",
		JWTAudience:     "test-audience",
		JWTExpireSeconds: 3600,
	}
	Init(cfg)

	_, err := ParseToken("invalid-token")
	if err == nil {
		t.Error("Expected error for invalid token")
	}
}

func TestJWT_Middleware_Disabled(t *testing.T) {
	cfg := config.AuthConfig{
		Enabled:    false,
		JWTEnabled: false,
	}
	Init(cfg)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/test", nil)

	handler := JWT()
	handler(c)

	if c.IsAborted() {
		t.Error("Request should not be aborted when JWT is disabled")
	}
}

func TestJWT_Middleware_MissingAuthHeader(t *testing.T) {
	cfg := config.AuthConfig{
		Enabled:    true,
		JWTEnabled: true,
		JWTSecret: "test-secret-key-1234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890",
		JWTIssuer:  "test-issuer",
		JWTAudience: "test-audience",
		JWTExpireSeconds: 3600,
	}
	Init(cfg)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/test", nil)

	handler := JWT()
	handler(c)

	if !c.IsAborted() {
		t.Error("Request should be aborted when Authorization header is missing")
	}
	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestJWT_Middleware_InvalidAuthFormat(t *testing.T) {
	cfg := config.AuthConfig{
		Enabled:    true,
		JWTEnabled: true,
		JWTSecret: "test-secret-key-1234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890",
		JWTIssuer:  "test-issuer",
		JWTAudience: "test-audience",
		JWTExpireSeconds: 3600,
	}
	Init(cfg)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/test", nil)
	c.Request.Header.Set("Authorization", "InvalidFormat")

	handler := JWT()
	handler(c)

	if !c.IsAborted() {
		t.Error("Request should be aborted for invalid auth format")
	}
}

func TestJWT_Middleware_ValidToken(t *testing.T) {
	cfg := config.AuthConfig{
		Enabled:    true,
		JWTEnabled: true,
		JWTSecret: "test-secret-key-1234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890",
		JWTIssuer:  "test-issuer",
		JWTAudience: "test-audience",
		JWTExpireSeconds: 3600,
	}
	Init(cfg)

	token, _ := GenerateToken("user123", "testuser", "admin")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/test", nil)
	c.Request.Header.Set("Authorization", "Bearer "+token)

	handler := JWT()
	handler(c)

	if c.IsAborted() {
		t.Error("Request should not be aborted for valid token")
	}

	userID := GetUserID(c)
	if userID != "user123" {
		t.Errorf("Expected userID 'user123', got '%s'", userID)
	}
}

func TestJWT_Middleware_ExcludedPath(t *testing.T) {
	cfg := config.AuthConfig{
		Enabled:    true,
		JWTEnabled: true,
		JWTSecret: "test-secret-key-1234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890",
		JWTIssuer:  "test-issuer",
		JWTAudience: "test-audience",
		JWTExpireSeconds: 3600,
		ExcludePaths: []string{"/health", "/ready"},
	}
	Init(cfg)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/health", nil)

	handler := JWT()
	handler(c)

	if c.IsAborted() {
		t.Error("Request should not be aborted for excluded path")
	}
}

func TestSignature_Middleware_Disabled(t *testing.T) {
	cfg := config.AuthConfig{
		Enabled:          false,
		SignatureEnabled: false,
	}
	Init(cfg)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/test", nil)

	handler := Signature()
	handler(c)

	if c.IsAborted() {
		t.Error("Request should not be aborted when signature is disabled")
	}
}

func TestSignature_Middleware_MissingSignature(t *testing.T) {
	cfg := config.AuthConfig{
		Enabled:          true,
		SignatureEnabled: true,
		SignatureSecret:  "test-secret",
		SignatureHeader:  "X-Signature",
		TimestampHeader:  "X-Timestamp",
		TimestampMaxDiff: 300,
	}
	Init(cfg)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/test", nil)

	handler := Signature()
	handler(c)

	if !c.IsAborted() {
		t.Error("Request should be aborted when signature is missing")
	}
}

func TestIsExcludedPath(t *testing.T) {
	cfg := config.AuthConfig{
		ExcludePaths: []string{"/health", "/ready"},
	}
	Init(cfg)

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
			name:     "/health/extra should be excluded",
			path:     "/health/extra",
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
			result := isExcludedPath(tt.path)
			if result != tt.expected {
				t.Errorf("Expected %t, got %t", tt.expected, result)
			}
		})
	}
}

func TestAbs(t *testing.T) {
	tests := []struct {
		name     string
		input    int64
		expected int64
	}{
		{
			name:     "positive number",
			input:    10,
			expected: 10,
		},
		{
			name:     "negative number",
			input:    -10,
			expected: 10,
		},
		{
			name:     "zero",
			input:    0,
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := abs(tt.input)
			if result != tt.expected {
				t.Errorf("Expected %d, got %d", tt.expected, result)
			}
		})
	}
}

func TestGetUserID(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user_id", "user123")

	userID := GetUserID(c)
	if userID != "user123" {
		t.Errorf("Expected 'user123', got '%s'", userID)
	}
}

func TestGetUserID_NotSet(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	userID := GetUserID(c)
	if userID != "" {
		t.Errorf("Expected empty string, got '%s'", userID)
	}
}

func TestGetUsername(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("username", "testuser")

	username := GetUsername(c)
	if username != "testuser" {
		t.Errorf("Expected 'testuser', got '%s'", username)
	}
}

func TestGetRole(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("role", "admin")

	role := GetRole(c)
	if role != "admin" {
		t.Errorf("Expected 'admin', got '%s'", role)
	}
}

func TestGetClaims(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("claims", &Claims{UserID: "user123"})

	claims := GetClaims(c)
	if claims == nil {
		t.Error("Claims should not be nil")
	}
	if claims.UserID != "user123" {
		t.Errorf("Expected 'user123', got '%s'", claims.UserID)
	}
}

func TestRequireRole_Success(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("role", "admin")

	handler := RequireRole("admin")
	handler(c)

	if c.IsAborted() {
		t.Error("Request should not be aborted when role matches")
	}
}

func TestRequireRole_Failure(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("role", "user")

	handler := RequireRole("admin")
	handler(c)

	if !c.IsAborted() {
		t.Error("Request should be aborted when role does not match")
	}
	if w.Code != http.StatusForbidden {
		t.Errorf("Expected status %d, got %d", http.StatusForbidden, w.Code)
	}
}

func TestRequireAnyRole_Success(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("role", "editor")

	handler := RequireAnyRole("admin", "editor")
	handler(c)

	if c.IsAborted() {
		t.Error("Request should not be aborted when role matches any")
	}
}

func TestRequireAnyRole_Failure(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("role", "guest")

	handler := RequireAnyRole("admin", "editor")
	handler(c)

	if !c.IsAborted() {
		t.Error("Request should be aborted when role does not match any")
	}
}
