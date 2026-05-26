# 中间件模块

提供 Gin 框架的中间件集合，包括日志、认证、追踪、限流等。

## 功能特性

- 📝 **请求日志**: 结构化请求日志记录
- 🔍 **链路追踪**: 自动创建追踪 Span
- 🔐 **认证授权**: JWT 和签名验证
- 🔄 **限流控制**: 请求频率限制
- 🔒 **安全 Headers**: 安全响应头设置
- 🔧 **请求 ID**: 请求唯一标识追踪
- 📊 **指标采集**: 请求指标自动采集
- ⏰ **超时控制**: 请求超时处理
- 🔄 **CORS**: 跨域资源共享
- 💊 **异常恢复**: Panic 恢复

## 快速开始

### 注册中间件

```go
import "github.com/lookingcamel/system-framework/internal/server"

func main() {
    srv, _ := server.New(cfg)

    // 中间件已在 server.New 中自动注册
    // RequestID, Recovery, CORS, Logger, Metrics, Tracing
}
```

### 自定义中间件

```go
func CustomMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()

        // 处理请求
        c.Next()

        // 记录日志
        duration := time.Since(start)
        fmt.Printf("Request: %s %s - %d (%s)\n",
            c.Request.Method,
            c.Request.URL.Path,
            c.Writer.Status(),
            duration,
        )
    }
}

// 使用自定义中间件
engine.Use(CustomMiddleware())
```

## 中间件列表

### 1. RequestID 中间件

自动为每个请求生成唯一 ID，支持请求间传递。

```go
// 生成的 Request ID 会添加以下 Header
X-Request-ID: uuid
```

### 2. Logger 中间件

记录请求日志，包括方法、路径、状态码、延迟等。

```json
{
  "level": "info",
  "msg": "HTTP request",
  "request_id": "abc123",
  "method": "GET",
  "path": "/api/users",
  "status": 200,
  "latency": 0.023,
  "client_ip": "192.168.1.1",
  "user_agent": "Mozilla/5.0"
}
```

### 3. Recovery 中间件

捕获 Panic 异常，返回 500 错误并记录堆栈。

```go
func(c *gin.Context) {
    defer func() {
        if err := recover(); err != nil {
            log.Printf("Panic recovered: %v\n%s", err, debug.Stack())
            c.AbortWithStatus(500)
        }
    }()
    c.Next()
}
```

### 4. CORS 中间件

处理跨域请求。

```go
// 默认配置
// Access-Control-Allow-Origin: *
// Access-Control-Allow-Methods: GET, POST, PUT, DELETE, OPTIONS
// Access-Control-Allow-Headers: *
```

### 5. Metrics 中间件

自动采集 HTTP 请求指标。

### 6. Tracing 中间件

自动创建追踪 Span。

### 7. Auth 中间件

JWT 认证和签名验证。

## 链式调用

```go
engine := gin.New()

// 按顺序注册中间件
engine.Use(middleware.RequestID())      // 1. 生成 Request ID
engine.Use(middleware.Recovery())       // 2. 异常恢复
engine.Use(middleware.CORS())           // 3. CORS 处理
engine.Use(auth.JWT())                  // 4. JWT 认证
engine.Use(middleware.Logger())         // 5. 日志记录
engine.Use(middleware.Metrics())        // 6. 指标采集
engine.Use(middleware.Tracing())        // 7. 链路追踪
```

## 自定义中间件

### 1. 限流中间件

```go
import "golang.org/x/time/rate"

func RateLimitMiddleware(requestsPerSecond float64, burst int) gin.HandlerFunc {
    limiter := rate.NewLimiter(rate.Limit(requestsPerSecond), burst)

    return func(c *gin.Context) {
        if !limiter.Allow() {
            c.JSON(429, gin.H{
                "code":    -1,
                "message": "Too many requests",
            })
            c.Abort()
            return
        }
        c.Next()
    }
}

// 使用
engine.Use(RateLimitMiddleware(100, 200))
```

### 2. 权限检查中间件

```go
func RequirePermission(permission string) gin.HandlerFunc {
    return func(c *gin.Context) {
        userPerms := GetUserPermissions(c)

        hasPermission := false
        for _, perm := range userPerms {
            if perm == permission {
                hasPermission = true
                break
            }
        }

        if !hasPermission {
            c.JSON(403, gin.H{
                "code":    -1,
                "message": "Insufficient permissions",
            })
            c.Abort()
            return
        }

        c.Next()
    }
}

// 使用
engine.GET("/admin", RequirePermission("admin:access"), adminHandler)
```

### 3. 请求验证中间件

```go
func ValidateRequest(schema string) gin.HandlerFunc {
    return func(c *gin.Context) {
        if err := validateJSONSchema(c.Request.Body, schema); err != nil {
            c.JSON(400, gin.H{
                "code":    -1,
                "message": fmt.Sprintf("Invalid request: %s", err.Error()),
            })
            c.Abort()
            return
        }
        c.Next()
    }
}

// 使用
engine.POST("/api/users",
    ValidateRequest("user_create_schema"),
    createUserHandler,
)
```

### 4. 缓存中间件

```go
func CacheMiddleware(expiration time.Duration) gin.HandlerFunc {
    return func(c *gin.Context) {
        key := fmt.Sprintf("cache:%s:%s", c.Request.Method, c.Request.URL.Path)
        
        cached, err := redis.Get(key)
        if err == nil && cached != "" {
            c.Header("X-Cache", "HIT")
            c.Data(200, "application/json", []byte(cached))
            c.Abort()
            return
        }

        c.Header("X-Cache", "MISS")
        c.Next()

        if c.Writer.Status() == 200 {
            body := c.Writer.Body.String()
            redis.Set(key, body, expiration)
        }
    }
}

// 使用
engine.GET("/api/users", CacheMiddleware(5*time.Minute), getUsersHandler)
```

### 5. 日志增强中间件

```go
func EnhancedLogger() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 请求开始时间
        start := time.Now()

        // 获取 Request ID
        requestID := utils.GinGetRequestID(c)

        // 处理请求
        c.Next()

        // 计算延迟
        latency := time.Since(start)

        // 获取用户信息
        userID := auth.GetUserID(c)

        // 结构化日志
        logger.Log.Info("HTTP request",
            zap.String("request_id", requestID),
            zap.String("user_id", userID),
            zap.String("method", c.Request.Method),
            zap.String("path", c.Request.URL.Path),
            zap.Int("status", c.Writer.Status()),
            zap.Duration("latency", latency),
            zap.String("client_ip", c.ClientIP()),
            zap.Int("response_size", c.Writer.Size()),
        )
    }
}
```

## 中间件顺序

中间件执行顺序很重要，推荐顺序：

```
1. RequestID    - 生成请求 ID
2. Recovery     - 异常恢复
3. CORS         - 跨域处理
4. Auth         - 认证授权
5. Logger       - 日志记录
6. Metrics      - 指标采集
7. Tracing      - 链路追踪
8. Custom       - 自定义中间件
9. Handler      - 业务处理
```

## 中间件参数

### 从 Context 获取数据

```go
// 获取 Request ID
requestID := utils.GinGetRequestID(c)

// 获取 Claims
claims := auth.GetClaims(c)

// 获取自定义值
value, exists := c.Get("custom_key")
```

### 设置 Context 值

```go
// 设置自定义值
c.Set("custom_key", "custom_value")

// 修改请求上下文
ctx := context.WithValue(c.Request.Context(), "key", "value")
c.Request = c.Request.WithContext(ctx)
```

## 性能优化

### 1. 异步中间件

```go
func AsyncMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        cCopy := c.Copy()
        go func() {
            // 异步日志
            logRequest(cCopy)
        }()
        c.Next()
    }
}
```

### 2. 条件中间件

```go
func ConditionalMiddleware(skipPaths []string) gin.HandlerFunc {
    skipMap := make(map[string]bool)
    for _, path := range skipPaths {
        skipMap[path] = true
    }

    return func(c *gin.Context) {
        if skipMap[c.Request.URL.Path] {
            c.Next()
            return
        }
        // 执行中间件逻辑
        doSomething(c)
        c.Next()
    }
}
```

## 错误处理

### 中间件错误返回

```go
func ErrorHandler() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Next()

        // 处理错误
        if len(c.Errors) > 0 {
            c.JSON(-1, gin.H{
                "code":    -1,
                "message": c.Errors.String(),
            })
        }
    }
}
```

## 测试中间件

```go
func TestMiddleware(t *testing.T) {
    engine := gin.New()
    engine.Use(Logger())
    engine.GET("/test", func(c *gin.Context) {
        c.JSON(200, gin.H{"message": "ok"})
    })

    req, _ := http.NewRequest("GET", "/test", nil)
    w := httptest.NewRecorder()
    engine.ServeHTTP(w, req)

    assert.Equal(t, 200, w.Code)
    assert.NotEmpty(t, w.Header().Get("X-Request-ID"))
}
```

## 最佳实践

1. **中间件顺序**
   - 基础中间件在前（RequestID, Recovery）
   - 认证授权次之
   - 日志和监控居中
   - 业务中间件在后

2. **错误处理**
   - 统一错误格式
   - 记录详细错误信息
   - 优雅降级

3. **性能优化**
   - 避免不必要的计算
   - 使用异步处理
   - 条件执行

4. **安全性**
   - 验证所有输入
   - 过滤敏感信息
   - 记录安全事件

## 相关文档

- [认证模块](../auth/README.md)
- [链路追踪模块](../tracing/README.md)
- [日志模块](../logger/README.md)
- [指标模块](../metrics/README.md)
- [项目总览](../../README.md)
