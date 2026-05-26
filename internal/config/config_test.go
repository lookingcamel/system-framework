package config

import (
	"testing"
)

func TestNewValidator(t *testing.T) {
	validator := NewValidator()
	if validator == nil {
		t.Fatal("NewValidator() returned nil")
	}
	if len(validator.errors) != 0 {
		t.Errorf("Expected empty errors, got %d errors", len(validator.errors))
	}
}

func TestValidator_Validate_EmptyConfig(t *testing.T) {
	validator := NewValidator()
	cfg := &Config{}

	err := validator.Validate(cfg)
	if err == nil {
		t.Error("Expected validation errors for empty config, got nil")
	}

	if cfg.App.Name != "system-framework" {
		t.Errorf("Expected default app name 'system-framework', got %s", cfg.App.Name)
	}
}

func TestValidator_Validate_ValidConfig(t *testing.T) {
	validator := NewValidator()
	cfg := &Config{
		App: AppConfig{
			Name: "test-app",
			Host: "0.0.0.0",
			Port: 8080,
			Mode: "release",
		},
		Server: ServerConfig{
			ReadTimeout:  30,
			WriteTimeout: 30,
			IdleTimeout:  120,
			MaxBodySize:  8388608,
		},
		Logging: LoggingConfig{
			Level:  "info",
			Format: "json",
			Output: "console",
		},
		APIVersion: APIVersionConfig{
			Enabled:          true,
			DefaultVersion:   "v1",
			SupportedVersions: []string{"v1", "v2"},
		},
		GracefulShutdown: GracefulShutdownConfig{
			Timeout: 30,
		},
		Database: DatabaseConfig{
			Type: "sqlite",
		},
	}

	err := validator.Validate(cfg)
	if err != nil {
		t.Errorf("Unexpected validation errors: %v", err)
	}
}

func TestValidator_Validate_InvalidPort(t *testing.T) {
	validator := NewValidator()
	cfg := &Config{
		App: AppConfig{
			Name: "test-app",
			Host: "0.0.0.0",
			Port: 99999, // Invalid port
			Mode: "release",
		},
	}

	validator.Validate(cfg)

	if cfg.App.Port != 8080 {
		t.Errorf("Expected default port 8080, got %d", cfg.App.Port)
	}
}

func TestValidator_Validate_InvalidMode(t *testing.T) {
	validator := NewValidator()
	cfg := &Config{
		App: AppConfig{
			Name: "test-app",
			Host: "0.0.0.0",
			Port: 8080,
			Mode: "invalid-mode",
		},
	}

	validator.Validate(cfg)

	if cfg.App.Mode != "release" {
		t.Errorf("Expected default mode 'release', got %s", cfg.App.Mode)
	}
}

func TestValidator_Validate_InvalidMaxBodySize(t *testing.T) {
	validator := NewValidator()
	cfg := &Config{
		App: AppConfig{
			Name: "test-app",
			Host: "0.0.0.0",
			Port: 8080,
			Mode: "release",
		},
		Server: ServerConfig{
			MaxBodySize: 200 * 1024 * 1024, // 200MB, exceeds max 100MB
		},
	}

	validator.Validate(cfg)

	if cfg.Server.MaxBodySize != 100<<20 {
		t.Errorf("Expected max body size capped at 100MB, got %d", cfg.Server.MaxBodySize)
	}
}

func TestValidator_Validate_NegativeMaxBodySize(t *testing.T) {
	validator := NewValidator()
	cfg := &Config{
		App: AppConfig{
			Name: "test-app",
			Host: "0.0.0.0",
			Port: 8080,
			Mode: "release",
		},
		Server: ServerConfig{
			MaxBodySize: -1,
		},
	}

	validator.Validate(cfg)

	if cfg.Server.MaxBodySize != 8<<20 {
		t.Errorf("Expected default max body size 8MB, got %d", cfg.Server.MaxBodySize)
	}
}

func TestValidator_Validate_InvalidAPIVersion(t *testing.T) {
	validator := NewValidator()
	cfg := &Config{
		App: AppConfig{
			Name: "test-app",
			Host: "0.0.0.0",
			Port: 8080,
			Mode: "release",
		},
		APIVersion: APIVersionConfig{
			Enabled:          true,
			DefaultVersion:   "version1", // Invalid format
			SupportedVersions: []string{"version1", "v2"},
		},
	}

	validator.Validate(cfg)

	if cfg.APIVersion.DefaultVersion != "v1" {
		t.Errorf("Expected default version 'v1', got %s", cfg.APIVersion.DefaultVersion)
	}
}

func TestValidator_Validate_DefaultVersionNotInSupported(t *testing.T) {
	validator := NewValidator()
	cfg := &Config{
		App: AppConfig{
			Name: "test-app",
			Host: "0.0.0.0",
			Port: 8080,
			Mode: "release",
		},
		APIVersion: APIVersionConfig{
			Enabled:          true,
			DefaultVersion:   "v3",
			SupportedVersions: []string{"v1", "v2"},
		},
	}

	validator.Validate(cfg)

	found := false
	for _, v := range cfg.APIVersion.SupportedVersions {
		if v == "v3" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected default version 'v3' to be added to supported versions")
	}
}

func TestValidator_AddError(t *testing.T) {
	validator := NewValidator()
	validator.AddError("test.field", "test error message")

	if len(validator.errors) != 1 {
		t.Errorf("Expected 1 error, got %d", len(validator.errors))
	}
	if validator.errors[0].Field != "test.field" {
		t.Errorf("Expected field 'test.field', got %s", validator.errors[0].Field)
	}
}

func TestValidator_Error(t *testing.T) {
	validator := NewValidator()
	validator.AddError("field1", "message1")
	validator.AddError("field2", "message2")

	errMsg := validator.Error()
	if len(errMsg) == 0 {
		t.Error("Expected error message, got empty string")
	}
	if !contains(errMsg, "field1") {
		t.Error("Expected error message to contain 'field1'")
	}
	if !contains(errMsg, "message2") {
		t.Error("Expected error message to contain 'message2'")
	}
}

func TestValidator_Unwrap(t *testing.T) {
	validator := NewValidator()
	validator.AddError("field1", "message1")
	validator.AddError("field2", "message2")

	errors := validator.Unwrap()
	if len(errors) != 2 {
		t.Errorf("Expected 2 errors, got %d", len(errors))
	}
}

func TestValidationError_Error(t *testing.T) {
	err := ValidationError{
		Field:   "test.field",
		Message: "test message",
	}

	errMsg := err.Error()
	if errMsg != "test.field: test message" {
		t.Errorf("Expected 'test.field: test message', got %s", errMsg)
	}
}

func TestValidatePattern(t *testing.T) {
	tests := []struct {
		name      string
		value     string
		pattern   string
		shouldErr bool
	}{
		{
			name:      "valid pattern match",
			value:     "v1",
			pattern:   `^v\d+$`,
			shouldErr: false,
		},
		{
			name:      "invalid pattern match",
			value:     "version1",
			pattern:   `^v\d+$`,
			shouldErr: true,
		},
		{
			name:      "empty pattern",
			value:     "any",
			pattern:   "",
			shouldErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePattern(tt.value, tt.pattern, "test")
			if tt.shouldErr && err == nil {
				t.Error("Expected error, got nil")
			}
			if !tt.shouldErr && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestValidateRange(t *testing.T) {
	tests := []struct {
		name      string
		value     int
		min       int
		max       int
		shouldErr bool
	}{
		{
			name:      "value in range",
			value:     50,
			min:       0,
			max:       100,
			shouldErr: false,
		},
		{
			name:      "value below min",
			value:     -1,
			min:       0,
			max:       100,
			shouldErr: true,
		},
		{
			name:      "value above max",
			value:     101,
			min:       0,
			max:       100,
			shouldErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateRange(tt.value, tt.min, tt.max, "test")
			if tt.shouldErr && err == nil {
				t.Error("Expected error, got nil")
			}
			if !tt.shouldErr && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestValidateNotEmpty(t *testing.T) {
	tests := []struct {
		name      string
		value     string
		shouldErr bool
	}{
		{
			name:      "non-empty string",
			value:     "hello",
			shouldErr: false,
		},
		{
			name:      "empty string",
			value:     "",
			shouldErr: true,
		},
		{
			name:      "only whitespace",
			value:     "   ",
			shouldErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateNotEmpty(tt.value, "test")
			if tt.shouldErr && err == nil {
				t.Error("Expected error, got nil")
			}
			if !tt.shouldErr && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		name      string
		email     string
		shouldErr bool
	}{
		{
			name:      "valid email",
			email:     "test@example.com",
			shouldErr: false,
		},
		{
			name:      "invalid email - no @",
			email:     "testexample.com",
			shouldErr: true,
		},
		{
			name:      "invalid email - no domain",
			email:     "test@",
			shouldErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateEmail(tt.email)
			if tt.shouldErr && err == nil {
				t.Error("Expected error, got nil")
			}
			if !tt.shouldErr && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestValidateURL(t *testing.T) {
	tests := []struct {
		name      string
		url       string
		shouldErr bool
	}{
		{
			name:      "valid http url",
			url:       "http://example.com",
			shouldErr: false,
		},
		{
			name:      "valid https url",
			url:       "https://example.com/path",
			shouldErr: false,
		},
		{
			name:      "invalid url - no protocol",
			url:       "example.com",
			shouldErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateURL(tt.url)
			if tt.shouldErr && err == nil {
				t.Error("Expected error, got nil")
			}
			if !tt.shouldErr && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestGetDirFromPath(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected string
	}{
		{
			name:     "linux path with directory",
			path:     "/var/log/app.log",
			expected: "/var/log",
		},
		{
			name:     "windows path",
			path:     "C:\\logs\\app.log",
			expected: "C:\\logs",
		},
		{
			name:     "no directory",
			path:     "app.log",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getDirFromPath(tt.path)
			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || contains(s[1:], substr)))
}
