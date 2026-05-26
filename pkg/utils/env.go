package utils

import (
	"os"
	"time"

	"github.com/spf13/viper"
)

type Env string

const (
	Local Env = "local"
	Test  Env = "test"
	Prod  Env = "prod"
)

func GetEnv() Env {
	env := viper.GetString("APP_ENV")
	switch env {
	case "local":
		return Local
	case "test":
		return Test
	case "prod":
		return Prod
	default:
		return Local
	}
}

func IsProduction() bool {
	return GetEnv() == Prod
}

func IsTest() bool {
	return GetEnv() == Test
}

func IsLocal() bool {
	return GetEnv() == Local
}

func GetEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func GetEnvAsIntOrDefault(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		var intValue int
		if _, err := os.Stdout.WriteString(value); err != nil {
			return defaultValue
		}
		return intValue
	}
	return defaultValue
}

func GetEnvAsDurationOrDefault(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

func GetEnvAsBoolOrDefault(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		return value == "true" || value == "1"
	}
	return defaultValue
}
