package versioning

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/lookingcamel/system-framework/internal/config"
	"github.com/lookingcamel/system-framework/internal/logger"
)

type APIVersionMiddleware struct {
	defaultVersion    string
	supportedVersions []string
	deprecationNotice string
}

var apiVersionMiddleware *APIVersionMiddleware

func InitAPIVersion(cfg config.APIVersionConfig) {
	if !cfg.Enabled {
		logger.Log.Info("API versioning is disabled")
		return
	}

	apiVersionMiddleware = &APIVersionMiddleware{
		defaultVersion:    cfg.DefaultVersion,
		supportedVersions: cfg.SupportedVersions,
		deprecationNotice: cfg.DeprecationNotice,
	}

	logger.Log.Info("API versioning initialized",
		zap.String("default_version", cfg.DefaultVersion),
		zap.Strings("supported_versions", cfg.SupportedVersions),
	)
}

func APIVersion() gin.HandlerFunc {
	return func(c *gin.Context) {
		if apiVersionMiddleware == nil {
			c.Next()
			return
		}

		path := c.Request.URL.Path
		version := extractVersionFromPath(path)

		if version == "" {
			version = apiVersionMiddleware.defaultVersion
			c.Set("API-Version", version)
			c.Set("API-Version-Deprecated", "false")
			c.Next()
			return
		}

		if !isVersionSupported(version) {
			logger.Log.Warn("Unsupported API version requested",
				zap.String("version", version),
				zap.String("path", path),
				zap.Strings("supported_versions", apiVersionMiddleware.supportedVersions),
			)

			c.JSON(http.StatusBadRequest, gin.H{
				"code":    -1,
				"message": "Unsupported API version",
				"error": gin.H{
					"requested_version": version,
					"supported_versions": apiVersionMiddleware.supportedVersions,
					"default_version": apiVersionMiddleware.defaultVersion,
				},
			})
			c.Abort()
			return
		}

		isDeprecated := isVersionDeprecated(version)
		c.Set("API-Version", version)
		c.Set("API-Version-Deprecated", isDeprecated)

		if isDeprecated && apiVersionMiddleware.deprecationNotice != "" {
			c.Header("X-API-Deprecation-Notice", apiVersionMiddleware.deprecationNotice)
		}

		c.Next()
	}
}

func extractVersionFromPath(path string) string {
	parts := strings.Split(strings.TrimPrefix(path, "/"), "/")

	if len(parts) >= 2 && strings.HasPrefix(parts[1], "v") {
		return parts[1]
	}

	return ""
}

func isVersionSupported(version string) bool {
	for _, v := range apiVersionMiddleware.supportedVersions {
		if v == version {
			return true
		}
	}
	return false
}

func isVersionDeprecated(version string) bool {
	if len(apiVersionMiddleware.supportedVersions) == 0 {
		return false
	}

	lastIndex := len(apiVersionMiddleware.supportedVersions) - 1
	latestVersion := apiVersionMiddleware.supportedVersions[lastIndex]

	return version != latestVersion
}

func GetCurrentVersion(c *gin.Context) string {
	if version, exists := c.Get("API-Version"); exists {
		return version.(string)
	}
	return ""
}

func IsDeprecated(c *gin.Context) bool {
	if deprecated, exists := c.Get("API-Version-Deprecated"); exists {
		return deprecated.(bool)
	}
	return false
}
