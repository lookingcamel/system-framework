package versioning

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/lookingcamel/system-framework/internal/config"
)

func TestInitAPIVersion(t *testing.T) {
	cfg := config.APIVersionConfig{
		Enabled:          true,
		DefaultVersion:   "v1",
		SupportedVersions: []string{"v1", "v2"},
		DeprecationNotice: "API v1 will be deprecated",
	}
	InitAPIVersion(cfg)

	if apiVersionMiddleware == nil {
		t.Fatal("apiVersionMiddleware should not be nil")
	}
	if apiVersionMiddleware.defaultVersion != "v1" {
		t.Errorf("Expected defaultVersion 'v1', got '%s'", apiVersionMiddleware.defaultVersion)
	}
	if len(apiVersionMiddleware.supportedVersions) != 2 {
		t.Errorf("Expected 2 supported versions, got %d", len(apiVersionMiddleware.supportedVersions))
	}
}

func TestInitAPIVersion_Disabled(t *testing.T) {
	cfg := config.APIVersionConfig{
		Enabled: false,
	}
	apiVersionMiddleware = nil // Reset
	InitAPIVersion(cfg)

	if apiVersionMiddleware != nil {
		t.Error("apiVersionMiddleware should still be nil when disabled")
	}
}

func TestAPIVersion_SupportedVersion(t *testing.T) {
	cfg := config.APIVersionConfig{
		Enabled:          true,
		DefaultVersion:   "v1",
		SupportedVersions: []string{"v1", "v2"},
	}
	InitAPIVersion(cfg)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	
	c.Request, _ = http.NewRequest("GET", "/api/v1/test", nil)

	handler := APIVersion()
	handler(c)

	if c.IsAborted() {
		t.Error("Request should not be aborted for supported version")
	}

	version, exists := c.Get("API-Version")
	if !exists {
		t.Error("API-Version should be set in context")
	}
	if version != "v1" {
		t.Errorf("Expected version 'v1', got '%s'", version)
	}
}

func TestAPIVersion_UnsupportedVersion(t *testing.T) {
	cfg := config.APIVersionConfig{
		Enabled:          true,
		DefaultVersion:   "v1",
		SupportedVersions: []string{"v1", "v2"},
	}
	InitAPIVersion(cfg)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	
	c.Request, _ = http.NewRequest("GET", "/api/v3/test", nil)

	handler := APIVersion()
	handler(c)

	if !c.IsAborted() {
		t.Error("Request should be aborted for unsupported version")
	}

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestAPIVersion_NoVersionInPath(t *testing.T) {
	cfg := config.APIVersionConfig{
		Enabled:          true,
		DefaultVersion:   "v2",
		SupportedVersions: []string{"v1", "v2"},
	}
	InitAPIVersion(cfg)

	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	
	c.Request, _ = http.NewRequest("GET", "/api/test", nil)

	handler := APIVersion()
	handler(c)

	version, exists := c.Get("API-Version")
	if !exists {
		t.Error("API-Version should be set in context")
	}
	if version != "v2" {
		t.Errorf("Expected default version 'v2', got '%s'", version)
	}
}

func TestAPIVersion_DeprecatedVersion(t *testing.T) {
	cfg := config.APIVersionConfig{
		Enabled:          true,
		DefaultVersion:   "v1",
		SupportedVersions: []string{"v1", "v2"}, // v1 is deprecated (not latest)
		DeprecationNotice: "API v1 will be deprecated",
	}
	InitAPIVersion(cfg)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	
	c.Request, _ = http.NewRequest("GET", "/api/v1/test", nil)

	handler := APIVersion()
	handler(c)

	deprecated, exists := c.Get("API-Version-Deprecated")
	if !exists {
		t.Error("API-Version-Deprecated should be set in context")
	}
	if !deprecated.(bool) {
		t.Error("v1 should be marked as deprecated")
	}

	noticeHeader := w.Header().Get("X-API-Deprecation-Notice")
	if noticeHeader != "API v1 will be deprecated" {
		t.Errorf("Expected deprecation notice header, got '%s'", noticeHeader)
	}
}

func TestAPIVersion_LatestVersionNotDeprecated(t *testing.T) {
	cfg := config.APIVersionConfig{
		Enabled:          true,
		DefaultVersion:   "v2",
		SupportedVersions: []string{"v1", "v2"},
	}
	InitAPIVersion(cfg)

	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	
	c.Request, _ = http.NewRequest("GET", "/api/v2/test", nil)

	handler := APIVersion()
	handler(c)

	deprecated, exists := c.Get("API-Version-Deprecated")
	if !exists {
		t.Error("API-Version-Deprecated should be set in context")
	}
	if deprecated.(bool) {
		t.Error("v2 should not be marked as deprecated")
	}
}

func TestAPIVersion_MiddlewareDisabled(t *testing.T) {
	apiVersionMiddleware = nil // Reset

	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	
	c.Request, _ = http.NewRequest("GET", "/api/test", nil)

	handler := APIVersion()
	handler(c)

	if c.IsAborted() {
		t.Error("Request should not be aborted when middleware is disabled")
	}
}

func TestExtractVersionFromPath(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected string
	}{
		{
			name:     "/api/v1/test should extract v1",
			path:     "/api/v1/test",
			expected: "v1",
		},
		{
			name:     "/api/v2/items should extract v2",
			path:     "/api/v2/items",
			expected: "v2",
		},
		{
			name:     "/api/test should not extract",
			path:     "/api/test",
			expected: "",
		},
		{
			name:     "/test should not extract",
			path:     "/test",
			expected: "",
		},
		{
			name:     "/v1/test should not extract (needs api prefix)",
			path:     "/v1/test",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractVersionFromPath(tt.path)
			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

func TestIsVersionSupported(t *testing.T) {
	cfg := config.APIVersionConfig{
		Enabled:          true,
		DefaultVersion:   "v1",
		SupportedVersions: []string{"v1", "v2", "v3"},
	}
	InitAPIVersion(cfg)

	tests := []struct {
		name     string
		version  string
		expected bool
	}{
		{
			name:     "v1 should be supported",
			version:  "v1",
			expected: true,
		},
		{
			name:     "v2 should be supported",
			version:  "v2",
			expected: true,
		},
		{
			name:     "v4 should not be supported",
			version:  "v4",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isVersionSupported(tt.version)
			if result != tt.expected {
				t.Errorf("Expected %t for version %s, got %t", tt.expected, tt.version, result)
			}
		})
	}
}

func TestIsVersionDeprecated(t *testing.T) {
	cfg := config.APIVersionConfig{
		Enabled:          true,
		DefaultVersion:   "v1",
		SupportedVersions: []string{"v1", "v2", "v3"}, // v3 is latest
	}
	InitAPIVersion(cfg)

	tests := []struct {
		name     string
		version  string
		expected bool
	}{
		{
			name:     "v1 should be deprecated",
			version:  "v1",
			expected: true,
		},
		{
			name:     "v2 should be deprecated",
			version:  "v2",
			expected: true,
		},
		{
			name:     "v3 should not be deprecated",
			version:  "v3",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isVersionDeprecated(tt.version)
			if result != tt.expected {
				t.Errorf("Expected %t for version %s, got %t", tt.expected, tt.version, result)
			}
		})
	}
}

func TestGetCurrentVersion(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	// Test with version set
	c.Set("API-Version", "v2")
	version := GetCurrentVersion(c)
	if version != "v2" {
		t.Errorf("Expected 'v2', got '%s'", version)
	}

	// Test with no version set
	c2, _ := gin.CreateTestContext(httptest.NewRecorder())
	version2 := GetCurrentVersion(c2)
	if version2 != "" {
		t.Errorf("Expected empty string, got '%s'", version2)
	}
}

func TestIsDeprecated(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	// Test with true
	c.Set("API-Version-Deprecated", true)
	deprecated := IsDeprecated(c)
	if !deprecated {
		t.Error("Expected true, got false")
	}

	// Test with false
	c2, _ := gin.CreateTestContext(httptest.NewRecorder())
	c2.Set("API-Version-Deprecated", false)
	deprecated2 := IsDeprecated(c2)
	if deprecated2 {
		t.Error("Expected false, got true")
	}

	// Test with no value set
	c3, _ := gin.CreateTestContext(httptest.NewRecorder())
	deprecated3 := IsDeprecated(c3)
	if deprecated3 {
		t.Error("Expected false when no value set, got true")
	}
}

func TestNewVersionedRouter(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()

	router := NewVersionedRouter(engine)

	if router == nil {
		t.Fatal("NewVersionedRouter returned nil")
	}
	if router.engine != engine {
		t.Error("Engine not set correctly")
	}
	if router.registeredVersions == nil {
		t.Error("registeredVersions map not initialized")
	}
	if router.registeredRoutes == nil {
		t.Error("registeredRoutes map not initialized")
	}

	// Test singleton
	router2 := NewVersionedRouter(engine)
	if router != router2 {
		t.Error("NewVersionedRouter should return singleton")
	}
}

func TestVersionedRouter_RegisterVersion(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()

	// Reset globalRouter
	globalRouter = nil
	
	router := NewVersionedRouter(engine)
	group := router.RegisterVersion("v1")

	if group == nil {
		t.Fatal("RegisterVersion returned nil")
	}
	if _, exists := router.registeredVersions["v1"]; !exists {
		t.Error("v1 not in registeredVersions")
	}
}

func TestVersionedRouter_RegisterVersion_Duplicate(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()

	// Reset globalRouter
	globalRouter = nil
	
	router := NewVersionedRouter(engine)
	group1 := router.RegisterVersion("v1")
	group2 := router.RegisterVersion("v1") // Duplicate

	if group1 != group2 {
		t.Error("RegisterVersion should return same group for duplicate")
	}
}

func TestVersionedRouter_GetGroup(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()

	// Reset globalRouter
	globalRouter = nil
	
	router := NewVersionedRouter(engine)
	router.RegisterVersion("v1")

	group := router.GetGroup("v1")
	if group == nil {
		t.Error("GetGroup returned nil for existing version")
	}

	group2 := router.GetGroup("v2")
	if group2 != nil {
		t.Error("GetGroup should return nil for non-existent version")
	}
}

func TestVersionedRouter_GetRegisteredVersions(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()

	// Reset globalRouter
	globalRouter = nil
	
	router := NewVersionedRouter(engine)
	router.RegisterVersion("v1")
	router.RegisterVersion("v2")

	versions := router.GetRegisteredVersions()
	if len(versions) != 2 {
		t.Errorf("Expected 2 versions, got %d", len(versions))
	}
}

func TestGetGlobalRouter(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()

	// Reset globalRouter
	globalRouter = nil
	
	router := NewVersionedRouter(engine)
	global := GetGlobalRouter()

	if router != global {
		t.Error("GetGlobalRouter should return same instance as NewVersionedRouter")
	}
}
