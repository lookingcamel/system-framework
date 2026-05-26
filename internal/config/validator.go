package config

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

type ValidationError struct {
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

type Validator struct {
	errors []ValidationError
}

func NewValidator() *Validator {
	return &Validator{
		errors: []ValidationError{},
	}
}

func (v *Validator) AddError(field, message string) {
	v.errors = append(v.errors, ValidationError{
		Field:   field,
		Message: message,
	})
}

func (v *Validator) Validate(cfg *Config) error {
	v.errors = []ValidationError{}

	v.validateApp(&cfg.App)
	v.validateServer(&cfg.Server)
	v.validateLogging(&cfg.Logging)
	v.validatePrometheus(&cfg.Prometheus)
	v.validatePProf(&cfg.PProf)
	v.validateTracing(&cfg.Tracing)
	v.validateHealth(&cfg.Health)
	v.validateGracefulShutdown(&cfg.GracefulShutdown)
	v.validateDatabase(&cfg.Database)
	v.validateRedis(&cfg.Redis)
	v.validateNacos(&cfg.Nacos)
	v.validateCircuitBreaker(&cfg.CircuitBreaker)
	v.validateAuth(&cfg.Auth)

	if len(v.errors) > 0 {
		return v
	}
	return nil
}

func (v *Validator) Error() string {
	if len(v.errors) == 0 {
		return ""
	}
	var sb strings.Builder
	sb.WriteString("Configuration validation failed:\n")
	for _, err := range v.errors {
		sb.WriteString(fmt.Sprintf("  - %s: %s\n", err.Field, err.Message))
	}
	return sb.String()
}

func (v *Validator) Field() string {
	return "validation_errors"
}

func (v *Validator) Unwrap() []ValidationError {
	return v.errors
}

func (v *Validator) validateApp(cfg *AppConfig) {
	if cfg.Name == "" {
		v.AddError("app.name", "app name is required")
		cfg.Name = "system-framework"
	}

	if cfg.Host == "" {
		v.AddError("app.host", "app host is required")
		cfg.Host = "0.0.0.0"
	}

	if cfg.Port <= 0 || cfg.Port > 65535 {
		v.AddError("app.port", fmt.Sprintf("invalid port: %d, using default 8080", cfg.Port))
		cfg.Port = 8080
	}

	validModes := map[string]bool{"debug": true, "release": true, "test": true}
	if !validModes[cfg.Mode] {
		v.AddError("app.mode", fmt.Sprintf("invalid mode: %s, using default 'release'", cfg.Mode))
		cfg.Mode = "release"
	}
}

func (v *Validator) validateServer(cfg *ServerConfig) {
	if cfg.ReadTimeout <= 0 {
		v.AddError("server.read_timeout", fmt.Sprintf("invalid read_timeout: %d, using default 30", cfg.ReadTimeout))
		cfg.ReadTimeout = 30
	}

	if cfg.WriteTimeout <= 0 {
		v.AddError("server.write_timeout", fmt.Sprintf("invalid write_timeout: %d, using default 30", cfg.WriteTimeout))
		cfg.WriteTimeout = 30
	}

	if cfg.IdleTimeout <= 0 {
		v.AddError("server.idle_timeout", fmt.Sprintf("invalid idle_timeout: %d, using default 120", cfg.IdleTimeout))
		cfg.IdleTimeout = 120
	}

	if cfg.MaxBodySize <= 0 {
		v.AddError("server.max_body_size", fmt.Sprintf("invalid max_body_size: %d, using default 8388608", cfg.MaxBodySize))
		cfg.MaxBodySize = 8 << 20
	}

	if cfg.MaxBodySize > 100<<20 {
		v.AddError("server.max_body_size", fmt.Sprintf("max_body_size too large: %d, max 100MB, using default", cfg.MaxBodySize))
		cfg.MaxBodySize = 100 << 20
	}
}

func (v *Validator) validateLogging(cfg *LoggingConfig) {
	validLevels := map[string]bool{"debug": true, "info": true, "warn": true, "error": true, "fatal": true}
	if !validLevels[cfg.Level] {
		v.AddError("logging.level", fmt.Sprintf("invalid level: %s, using default 'info'", cfg.Level))
		cfg.Level = "info"
	}

	validFormats := map[string]bool{"json": true, "console": true}
	if !validFormats[cfg.Format] {
		v.AddError("logging.format", fmt.Sprintf("invalid format: %s, using default 'json'", cfg.Format))
		cfg.Format = "json"
	}

	validOutputs := map[string]bool{"console": true, "file": true, "both": true}
	if !validOutputs[cfg.Output] {
		v.AddError("logging.output", fmt.Sprintf("invalid output: %s, using default 'console'", cfg.Output))
		cfg.Output = "console"
	}

	if cfg.Output == "file" || cfg.Output == "both" {
		if cfg.OutputPath == "" {
			v.AddError("logging.output_path", "output_path is required when output is file or both")
			cfg.OutputPath = "logs/app.log"
		}
	}

	if cfg.MaxSize <= 0 {
		cfg.MaxSize = 100
	}

	if cfg.MaxBackups <= 0 {
		cfg.MaxBackups = 7
	}

	if cfg.MaxAge <= 0 {
		cfg.MaxAge = 30
	}
}

func (v *Validator) validatePrometheus(cfg *PrometheusConfig) {
	if cfg.Path == "" {
		cfg.Path = "/metrics"
	}
}

func (v *Validator) validatePProf(cfg *PProfConfig) {
	if cfg.Path == "" {
		cfg.Path = "/debug/pprof"
	}
}

func (v *Validator) validateTracing(cfg *TracingConfig) {
	if cfg.Exporter == "" {
		cfg.Exporter = "stdout"
	}

	validExporters := map[string]bool{"stdout": true, "jaeger": true, "otlp": true}
	if !validExporters[cfg.Exporter] {
		v.AddError("tracing.exporter", fmt.Sprintf("invalid exporter: %s, using default 'stdout'", cfg.Exporter))
		cfg.Exporter = "stdout"
	}

	if cfg.SampleRate < 0 || cfg.SampleRate > 1 {
		v.AddError("tracing.sample_rate", fmt.Sprintf("invalid sample_rate: %f, using default 1.0", cfg.SampleRate))
		cfg.SampleRate = 1.0
	}

	if cfg.ServiceName == "" {
		cfg.ServiceName = "system-framework"
	}
}

func (v *Validator) validateHealth(cfg *HealthConfig) {
	if cfg.Path == "" {
		cfg.Path = "/health"
	}
}

func (v *Validator) validateGracefulShutdown(cfg *GracefulShutdownConfig) {
	if cfg.Timeout <= 0 {
		v.AddError("graceful_shutdown.timeout", fmt.Sprintf("invalid timeout: %d, using default 30", cfg.Timeout))
		cfg.Timeout = 30
	}
}

func (v *Validator) validateDatabase(cfg *DatabaseConfig) {
	validTypes := map[string]bool{"mysql": true, "postgres": true, "postgresql": true, "sqlite": true, "sqlite3": true}
	if !validTypes[cfg.Type] {
		v.AddError("database.type", fmt.Sprintf("invalid type: %s, using default 'sqlite'", cfg.Type))
		cfg.Type = "sqlite"
	}

	if cfg.Type == "sqlite" || cfg.Type == "sqlite3" {
		if cfg.SQLitePath == "" {
			cfg.SQLitePath = "./data/app.db"
		}
		dir := getDirFromPath(cfg.SQLitePath)
		if dir != "" {
			if err := os.MkdirAll(dir, 0755); err != nil {
				v.AddError("database.sqlite_path", fmt.Sprintf("failed to create directory: %s", err))
			}
		}
	} else {
		if cfg.Host == "" {
			v.AddError("database.host", "database host is required for non-sqlite databases")
			cfg.Host = "localhost"
		}

		if cfg.Name == "" {
			v.AddError("database.name", "database name is required for non-sqlite databases")
			cfg.Name = "app"
		}

		if cfg.Port <= 0 || cfg.Port > 65535 {
			v.AddError("database.port", fmt.Sprintf("invalid port: %d, using default 3306", cfg.Port))
			cfg.Port = 3306
		}
	}

	if cfg.MaxOpenConnections <= 0 {
		cfg.MaxOpenConnections = 10
	}

	if cfg.MaxIdleConnections <= 0 {
		cfg.MaxIdleConnections = 5
	}

	if cfg.ConnectionMaxLifetime <= 0 {
		cfg.ConnectionMaxLifetime = 300
	}

	if cfg.MaxOpenConnections < cfg.MaxIdleConnections {
		v.AddError("database.max_open_connections", "max_open_connections should be >= max_idle_connections")
		cfg.MaxIdleConnections = cfg.MaxOpenConnections
	}
}

func (v *Validator) validateRedis(cfg *RedisConfig) {
	if !cfg.Enabled {
		return
	}

	if cfg.Host == "" {
		v.AddError("redis.host", "redis host is required when redis is enabled")
		cfg.Host = "localhost"
	}

	if cfg.Port <= 0 || cfg.Port > 65535 {
		v.AddError("redis.port", fmt.Sprintf("invalid port: %d, using default 6379", cfg.Port))
		cfg.Port = 6379
	}

	if cfg.PoolSize <= 0 {
		cfg.PoolSize = 10
	}

	if cfg.MinIdleConns <= 0 {
		cfg.MinIdleConns = 5
	}
}

func (v *Validator) validateNacos(cfg *NacosConfig) {
	if !cfg.Enabled {
		return
	}

	if cfg.ServerAddr == "" {
		v.AddError("nacos.server_addr", "nacos server_addr is required when nacos is enabled")
		cfg.ServerAddr = "localhost"
	}

	if cfg.ServerPort <= 0 || cfg.ServerPort > 65535 {
		v.AddError("nacos.server_port", fmt.Sprintf("invalid port: %d, using default 8848", cfg.ServerPort))
		cfg.ServerPort = 8848
	}

	if cfg.DataId == "" {
		cfg.DataId = "application.yaml"
	}

	if cfg.Group == "" {
		cfg.Group = "DEFAULT_GROUP"
	}

	if cfg.ConfigType == "" {
		cfg.ConfigType = "yaml"
	}

	validConfigTypes := map[string]bool{"yaml": true, "yml": true, "json": true, "properties": true, "xml": true}
	if !validConfigTypes[cfg.ConfigType] {
		v.AddError("nacos.config_type", fmt.Sprintf("invalid config_type: %s, using default 'yaml'", cfg.ConfigType))
		cfg.ConfigType = "yaml"
	}
}

func (v *Validator) validateCircuitBreaker(cfg *CircuitBreakerConfig) {
	if !cfg.Enabled {
		return
	}

	if cfg.DefaultTimeout <= 0 {
		v.AddError("circuit_breaker.default_timeout", fmt.Sprintf("invalid timeout: %d, using default 3000", cfg.DefaultTimeout))
		cfg.DefaultTimeout = 3000
	}

	if cfg.DefaultMaxConcurrent <= 0 {
		v.AddError("circuit_breaker.default_max_concurrent", fmt.Sprintf("invalid max_concurrent: %d, using default 100", cfg.DefaultMaxConcurrent))
		cfg.DefaultMaxConcurrent = 100
	}

	if cfg.DefaultErrorPercentage < 0 || cfg.DefaultErrorPercentage > 100 {
		v.AddError("circuit_breaker.default_error_percentage", fmt.Sprintf("invalid error_percentage: %d, using default 50", cfg.DefaultErrorPercentage))
		cfg.DefaultErrorPercentage = 50
	}

	if cfg.DefaultRequestVolume <= 0 {
		v.AddError("circuit_breaker.default_request_volume", fmt.Sprintf("invalid request_volume: %d, using default 20", cfg.DefaultRequestVolume))
		cfg.DefaultRequestVolume = 20
	}
}

func (v *Validator) validateAuth(cfg *AuthConfig) {
	if !cfg.Enabled {
		return
	}

	if cfg.JWTEnabled {
		if cfg.JWTSecret == "" {
			v.AddError("auth.jwt_secret", "jwt_secret is required when JWT is enabled")
		} else if len(cfg.JWTSecret) < 32 {
			v.AddError("auth.jwt_secret", "jwt_secret should be at least 32 characters for security")
		}

		if cfg.JWTIssuer == "" {
			cfg.JWTIssuer = "system-framework"
		}

		if cfg.JWTAudience == "" {
			cfg.JWTAudience = "system-framework-api"
		}

		if cfg.JWTExpireSeconds <= 0 {
			v.AddError("auth.jwt_expire_seconds", fmt.Sprintf("invalid expire_seconds: %d, using default 86400", cfg.JWTExpireSeconds))
			cfg.JWTExpireSeconds = 86400
		}
	}

	if cfg.SignatureEnabled {
		if cfg.SignatureSecret == "" {
			v.AddError("auth.signature_secret", "signature_secret is required when signature is enabled")
		} else if len(cfg.SignatureSecret) < 32 {
			v.AddError("auth.signature_secret", "signature_secret should be at least 32 characters for security")
		}

		if cfg.SignatureHeader == "" {
			cfg.SignatureHeader = "X-Signature"
		}

		if cfg.TimestampHeader == "" {
			cfg.TimestampHeader = "X-Timestamp"
		}

		if cfg.TimestampMaxDiff <= 0 {
			v.AddError("auth.timestamp_max_diff", fmt.Sprintf("invalid max_diff: %d, using default 300", cfg.TimestampMaxDiff))
			cfg.TimestampMaxDiff = 300
		}
	}
}

func getDirFromPath(path string) string {
	for i := len(path) - 1; i >= 0; i-- {
		if path[i] == '/' || path[i] == '\\' {
			return path[:i]
		}
	}
	return ""
}

func ValidatePattern(value, pattern, fieldName string) error {
	if pattern == "" {
		return nil
	}
	matched, err := regexp.MatchString(pattern, value)
	if err != nil {
		return fmt.Errorf("invalid pattern for %s: %w", fieldName, err)
	}
	if !matched {
		return fmt.Errorf("%s does not match required pattern: %s", fieldName, pattern)
	}
	return nil
}

func ValidateRange(value, min, max int, fieldName string) error {
	if value < min {
		return fmt.Errorf("%s must be at least %d", fieldName, min)
	}
	if value > max {
		return fmt.Errorf("%s must be at most %d", fieldName, max)
	}
	return nil
}

func ValidateNotEmpty(value, fieldName string) error {
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("%s cannot be empty", fieldName)
	}
	return nil
}

func ValidateEmail(email string) error {
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, err := regexp.MatchString(emailRegex, email)
	if err != nil {
		return err
	}
	if !matched {
		return fmt.Errorf("invalid email format: %s", email)
	}
	return nil
}

func ValidateURL(url string) error {
	urlRegex := `^https?://[^\s/$.?#].[^\s]*$`
	matched, err := regexp.MatchString(urlRegex, url)
	if err != nil {
		return err
	}
	if !matched {
		return fmt.Errorf("invalid URL format: %s", url)
	}
	return nil
}