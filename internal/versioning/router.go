package versioning

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/lookingcamel/system-framework/internal/logger"
)

type VersionedRouter struct {
	engine              *gin.Engine
	registeredVersions  map[string]*VersionGroup
	registeredRoutes    map[string]bool
}

type VersionGroup struct {
	group  *gin.RouterGroup
	prefix string
}

var globalRouter *VersionedRouter

func NewVersionedRouter(engine *gin.Engine) *VersionedRouter {
	if globalRouter == nil {
		globalRouter = &VersionedRouter{
			engine:             engine,
			registeredVersions: make(map[string]*VersionGroup),
			registeredRoutes:   make(map[string]bool),
		}
	}
	return globalRouter
}

func (vr *VersionedRouter) RegisterVersion(version string, middlewares ...gin.HandlerFunc) *VersionGroup {
	if _, exists := vr.registeredVersions[version]; exists {
		logger.Log.Warn("Version already registered, skipping",
			zap.String("version", version),
		)
		return vr.registeredVersions[version]
	}

	group := vr.engine.Group("/api/" + version, middlewares...)

	vr.registeredVersions[version] = &VersionGroup{
		group:  group,
		prefix: "/api/" + version,
	}

	logger.Log.Info("API version registered",
		zap.String("version", version),
		zap.String("prefix", "/api/"+version),
	)

	return vr.registeredVersions[version]
}

func (vr *VersionedRouter) GetGroup(version string) *gin.RouterGroup {
	if group, exists := vr.registeredVersions[version]; exists {
		return group.group
	}
	return nil
}

func (vr *VersionedRouter) RegisterRoute(version, method, path string, handlers ...gin.HandlerFunc) {
	routeKey := version + ":" + method + ":" + path

	if vr.registeredRoutes[routeKey] {
		logger.Log.Warn("Route already registered, skipping",
			zap.String("version", version),
			zap.String("method", method),
			zap.String("path", path),
		)
		return
	}

	group := vr.GetGroup(version)
	if group == nil {
		logger.Log.Error("Version group not found, cannot register route",
			zap.String("version", version),
			zap.String("method", method),
			zap.String("path", path),
		)
		return
	}

	fullPath := "/api/" + version + path

	switch method {
	case "GET":
		group.GET(path, handlers...)
	case "POST":
		group.POST(path, handlers...)
	case "PUT":
		group.PUT(path, handlers...)
	case "PATCH":
		group.PATCH(path, handlers...)
	case "DELETE":
		group.DELETE(path, handlers...)
	case "OPTIONS":
		group.OPTIONS(path, handlers...)
	default:
		group.Any(path, handlers...)
	}

	vr.registeredRoutes[routeKey] = true

	logger.Log.Debug("Route registered",
		zap.String("version", version),
		zap.String("method", method),
		zap.String("path", fullPath),
	)
}

func (vr *VersionedRouter) GetRegisteredVersions() []string {
	versions := make([]string, 0, len(vr.registeredVersions))
	for version := range vr.registeredVersions {
		versions = append(versions, version)
	}
	return versions
}

func GetGlobalRouter() *VersionedRouter {
	return globalRouter
}
