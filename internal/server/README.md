# HTTP 服务器模块

基于 Gin 框架的 HTTP 服务器封装，提供优雅启动和关闭。

## 功能特性

- ⚡ **高性能**: 基于 Gin 的高性能服务器
- 🔄 **优雅关闭**: 支持零停机部署
- 🏥 **健康检查**: 存活和就绪探针
- 🌐 **中间件支持**: 完整的中间件链
- 🔧 **配置灵活**: 支持多种配置选项
- 📊 **性能监控**: 内置 pprof 支持
- 🔐 **安全特性**: CORS、Security Headers
- 🔄 **热重载**: 支持配置热更新

## 快速开始

### 创建服务器

```go
import "github.com/lookingcamel/system-framework/internal/server"

func main() {
    srv, err := server.New(cfg)
    if err != nil {
        log.Fatal(err)
    }

    if err := srv.Start(); err != nil {
        log.Fatal(err)
    }
}
```

### 服务器配置

```go
type Config struct {
    App                AppConfig
    Server             ServerConfig
    Logging            LoggingConfig
    Prometheus         PrometheusConfig
    PProf              PProfConfig
    Tracing            TracingConfig
    Health             HealthConfig
    Database           DatabaseConfig
    Redis              RedisConfig
    Auth               AuthConfig
    CircuitBreaker     CircuitBreakerConfig
}
```

## 服务器启动

```go
func main() {
    cfg, err := config.Load("config.yaml")
    if err != nil {
        log.Fatal(err)
    }

    // 初始化日志
    logger.Init(cfg.Logging)
    defer logger.Sync()

    // 初始化追踪
    shutdown, err := tracing.Init(cfg.Tracing)
    if err != nil {
        logger.Log.Warn("Failed to init tracing", zap.Error(err))
    } else {
        defer shutdown(context.Background())
    }

    // 初始化数据库
    if err := database.Init(cfg.Database); err != nil {
        logger.Log.Warn("Failed to init database", zap.Error(err))
    } else {
        defer database.Close()
    }

    // 创建服务器
    srv, err := server.New(cfg)
    if err != nil {
        log.Fatal(err)
    }

    // 启动服务器
    if err := srv.Start(); err != nil {
        log.Fatal(err)
    }
}
```

## 内置路由

### 健康检查

```
GET /health      # 存活探针
GET /ready       # 就绪探针
```

### 指标端点

```
GET /metrics     # Prometheus 指标
```

### 性能分析

```
GET /debug/pprof/           # pprof 首页
GET /debug/pprof/profile    # CPU 性能分析
GET /debug/pprof/heap       # 内存分析
GET /debug/pprof/goroutine  # Goroutine 分析
GET /debug/pprof/threadcreate # 线程分析
```

### 熔断器状态

```
GET /circuit/status  # 熔断器状态
```

## 中间件配置

### 自动注册

服务器会自动注册以下中间件：

```go
engine.Use(middleware.RequestID())     // 请求 ID
engine.Use(middleware.Recovery())     // 异常恢复
engine.Use(middleware.CORS())         // 跨域

if cfg.Auth.Enabled {
    if cfg.Auth.SignatureEnabled {
        engine.Use(auth.Signature())   // 签名验证
    }
    if cfg.Auth.JWTEnabled {
        engine.Use(auth.JWT())        // JWT 认证
    }
}

if cfg.Logging.Level != "" {
    engine.Use(middleware.Logger())   // 请求日志
}
```

### 手动注册

```go
srv, _ := server.New(cfg)

// 添加自定义中间件
srv.Engine().Use(customMiddleware)

// 注册路由
srv.Engine().GET("/custom", customHandler)

// 启动服务器
srv.Start()
```

## 优雅关闭

### 信号处理

```go
func (s *Server) waitForShutdown() {
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit

    logger.Log.Info("Shutting down server...")

    // 创建超时上下文
    ctx, cancel := context.WithTimeout(
        context.Background(),
        time.Duration(s.cfg.GracefulShutdown.Timeout)*time.Second,
    )
    defer cancel()

    // 优雅关闭
    if err := s.server.Shutdown(ctx); err != nil {
        logger.Log.Error("Server forced to shutdown", zap.Error(err))
    }

    // 刷新日志
    logger.Sync()

    fmt.Println("Server exited")
}
```

### 关闭流程

```
1. 接收 SIGINT/SIGTERM
2. 停止接收新连接
3. 等待正在处理的请求完成
4. 关闭数据库连接
5. 关闭 Redis 连接
6. 刷新日志
7. 退出进程
```

### 超时配置

```yaml
graceful_shutdown:
  timeout: 30  # 最多等待 30 秒
```

## 超时配置

```yaml
server:
  read_timeout: 30     # 读取超时（秒）
  write_timeout: 30    # 写入超时（秒）
  idle_timeout: 120    # 空闲连接超时（秒）
```

## CORS 配置

```go
// 允许所有来源（开发环境）
engine.Use(middleware.CORS())

// 生产环境建议配置
```

### CORS Headers

```
Access-Control-Allow-Origin: *
Access-Control-Allow-Credentials: true
Access-Control-Allow-Headers: Content-Type, Authorization, X-Request-ID
Access-Control-Allow-Methods: GET, POST, PUT, DELETE, OPTIONS
Access-Control-Expose-Headers: X-Request-ID
```

## 性能优化

### 1. 连接池配置

```yaml
server:
  max_connections: 10000
  read_timeout: 30
  write_timeout: 30
  idle_timeout: 120
```

### 2. GOMAXPROCS

```go
import "runtime"

func init() {
    // 设置 CPU 核心数
    numCPU := runtime.NumCPU()
    runtime.GOMAXPROCS(numCPU)
}
```

### 3. 禁用访问日志

```yaml
logging:
  level: "warn"  # 生产环境使用 warn
```

## 监控集成

### 健康检查端点

```go
func healthHandler(c *gin.Context) {
    checks := map[string]string{
        "database": "ok",
        "redis": "ok",
    }

    for name, status := range checks {
        if status != "ok" {
            c.JSON(503, gin.H{
                "status": "unhealthy",
                "checks": checks,
            })
            return
        }
    }

    c.JSON(200, gin.H{
        "status": "healthy",
        "checks": checks,
    })
}
```

### 自定义路由组

```go
func (s *Server) setupRoutes() {
    // API v1
    v1 := s.engine.Group("/api/v1")
    {
        v1.GET("/users", userHandler.List)
        v1.POST("/users", userHandler.Create)
    }

    // API v2
    v2 := s.engine.Group("/api/v2")
    {
        v2.GET("/users", userHandler.ListV2)
        v2.POST("/users", userHandler.CreateV2)
    }

    // 管理接口
    admin := s.engine.Group("/admin")
    admin.Use(auth.RequireRole("admin"))
    {
        admin.GET("/stats", adminHandler.Stats)
    }
}
```

## 错误处理

### 全局错误处理

```go
func errorHandler() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Next()

        if len(c.Errors) > 0 {
            c.JSON(-1, gin.H{
                "code":    -1,
                "message": c.Errors.String(),
            })
        }
    }
}
```

### 404 处理

```go
s.engine.NoRoute(func(c *gin.Context) {
    c.JSON(404, gin.H{
        "code":    -1,
        "message": "route not found",
    })
})
```

### 500 处理

```go
s.engine.Use(gin.Recovery())
```

## 安全配置

### Security Headers

```go
func SecurityHeaders() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Header("X-Frame-Options", "DENY")
        c.Header("X-Content-Type-Options", "nosniff")
        c.Header("X-XSS-Protection", "1; mode=block")
        c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
        c.Next()
    }
}
```

### 请求大小限制

```go
engine.MaxMultipartMemory = 8 << 20  // 8 MB
```

## 测试

### 单元测试

```go
func TestServer(t *testing.T) {
    cfg := &config.Config{
        App: config.AppConfig{
            Port: 8080,
            Mode: "test",
        },
    }

    srv, err := server.New(cfg)
    if err != nil {
        t.Fatal(err)
    }

    req, _ := http.NewRequest("GET", "/health", nil)
    w := httptest.NewRecorder()

    srv.Engine().ServeHTTP(w, req)

    assert.Equal(t, 200, w.Code)
}
```

### 集成测试

```go
func TestServerIntegration(t *testing.T) {
    srv, _ := server.New(cfg)
    go srv.Start()

    time.Sleep(2 * time.Second)

    req, _ := http.NewRequest("GET", "/health", nil)
    resp, err := http.Get("http://localhost:8080/health")
    if err != nil {
        t.Fatal(err)
    }
    defer resp.Body.Close()

    assert.Equal(t, 200, resp.StatusCode)
}
```

## 最佳实践

1. **启动顺序**
   - 日志 → 追踪 → 数据库 → Redis → 服务器

2. **关闭顺序**
   - 服务器 → Redis → 数据库 → 日志

3. **错误处理**
   - 不panic，使用日志记录
   - 返回友好错误信息
   - 监控错误率

4. **性能监控**
   - 监控请求延迟
   - 监控并发数
   - 监控资源使用

## 故障排查

### 常见问题

**Q: 端口被占用**
```
A: 检查端口占用进程: lsof -i:8080
```

**Q: 连接超时**
```
A: 增加 server.read_timeout 和 server.write_timeout
```

**Q: 优雅关闭失败**
```
A: 增加 graceful_shutdown.timeout
```

## 相关文档

- [中间件模块](../middleware/README.md)
- [认证模块](../auth/README.md)
- [日志模块](../logger/README.md)
- [链路追踪模块](../tracing/README.md)
- [项目总览](../../README.md)
