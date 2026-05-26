# System Framework 安全评审报告

**项目**: github.com/lookingcamel/system-framework  
**评审日期**: 2026-05-26  
**评审范围**: 认证、授权、配置、中间件、数据库  
**技术栈**: Go 1.25 + Gin + JWT + MySQL/PostgreSQL/SQLite + Redis

---

## 执行摘要

本报告对 System Framework 项目进行了全面的安全评审，识别出 **3 个高危漏洞**、**5 个中危漏洞** 和 **4 个低危漏洞**。项目在认证机制、请求签名和配置验证方面表现良好，但在 CORS 配置、敏感信息管理和生产环境安全配置方面存在需要立即修复的问题。

---

## 关键发现

### 🔴 高危漏洞

#### **[HIGH-1] CORS 配置允许所有来源访问**

**位置**: [internal/middleware/middleware.go#L127-142](file:///d:/code/system-framework/internal/middleware/middleware.go#L127-142)

**问题描述**:
```go
c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
```
生产环境中使用通配符 `*` 允许所有来源访问，存在跨域请求伪造风险。

**影响**: 攻击者可以从未授权的网站发起跨域请求，可能导致数据泄露和会话劫持。

**建议修复**:
```go
// 根据环境配置 CORS
allowedOrigins := viper.GetStringSlice("cors.allowed_origins")
if len(allowedOrigins) > 0 && allowedOrigins[0] != "*" {
    c.Writer.Header().Set("Access-Control-Allow-Origin", strings.Join(allowedOrigins, ","))
}
```

**配置文件**:
```yaml
cors:
  allowed_origins:
    - "https://example.com"
    - "https://app.example.com"
  allow_credentials: true
```

**优先级**: 🔴 高  
**CVSS 评分**: 8.1

---

#### **[HIGH-2] 配置文件包含明文敏感信息**

**位置**: [config.yaml#L89-102](file:///d:/code/system-framework/config.yaml#L89-102)

**问题描述**:
```yaml
auth:
  jwt_secret: "your-256-bit-secret-key-here-change-in-production"
  signature_secret: "your-signature-secret-key-here-change-in-production"

database:
  password: "password"

nacos:
  password: "nacos"
```

敏感信息（密钥、密码）以明文形式存储在配置文件中，可能被泄露。

**影响**: 
- JWT secret 泄露可导致攻击者伪造任意用户的 Token
- 数据库密码泄露可导致未授权数据库访问
- 签名密钥泄露可导致请求签名验证失效

**建议修复**:
1. 使用环境变量存储敏感信息
```yaml
auth:
  jwt_secret: ${JWT_SECRET}
  signature_secret: ${SIGNATURE_SECRET}

database:
  password: ${DB_PASSWORD}

nacos:
  password: ${NACOS_PASSWORD}
```

2. 使用密钥管理服务（KMS）
```go
secret, err := kmsClient.GetSecretValue("jwt_secret")
```

3. 添加配置加密支持
```yaml
auth:
  jwt_secret_encrypted: "encrypted_value_here"
```

**优先级**: 🔴 高  
**CVSS 评分**: 9.1

---

#### **[HIGH-3] 生产环境暴露 pprof 端点**

**位置**: [config.yaml#L27-29](file:///d:/code/system-framework/config.yaml#L27-29)  
**位置**: [internal/server/server.go](file:///d:/code/system-framework/internal/server/server.go)

**问题描述**:
```yaml
pprof:
  enabled: true
  path: "/debug/pprof"
```

pprof 端点暴露了运行时性能分析数据，包括内存分配、CPU 使用和 Goroutine 状态。

**影响**: 
- 可能泄露敏感内存数据
- 可能泄露代码结构和内部实现细节
- 可能被用于 DoS 攻击

**建议修复**:
```yaml
pprof:
  enabled: ${PPROF_ENABLED:false}  # 生产环境默认禁用
  path: "/debug/pprof"
```

或通过环境变量控制：
```go
if viper.GetBool("pprof.enabled") && viper.GetString("app.mode") != "release" {
    // 注册 pprof 路由
}
```

**优先级**: 🔴 高  
**CVSS 评分**: 7.5

---

### 🟡 中危漏洞

#### **[MED-1] 缺少速率限制中间件**

**位置**: [internal/middleware/middleware.go](file:///d:/code/system-framework/internal/middleware/middleware.go)

**问题描述**: 项目未实现请求速率限制，可能受到暴力破解和 DoS 攻击。

**影响**: 攻击者可以无限制地尝试认证或消耗服务器资源。

**建议实现**:
```go
// 使用令牌桶算法
import "golang.org/x/time/rate"

func RateLimitMiddleware(rps float64, burst int) gin.HandlerFunc {
    limiter := rate.NewLimiter(rate.Limit(rps), burst)
    return func(c *gin.Context) {
        if !limiter.Allow() {
            c.JSON(http.StatusTooManyRequests, gin.H{
                "code":    -1,
                "message": "Too many requests",
            })
            c.Abort()
            return
        }
        c.Next()
    }
}
```

**配置文件**:
```yaml
rate_limit:
  enabled: true
  requests_per_second: 100
  burst: 200
```

**优先级**: 🟡 中  
**CVSS 评分**: 5.3

---

#### **[MED-2] 缺少安全响应头配置**

**位置**: [internal/middleware/middleware.go#L127-142](file:///d:/code/system-framework/internal/middleware/middleware.go#L127-142)

**问题描述**: 缺少关键的安全响应头，包括：
- `X-Frame-Options`
- `X-Content-Type-Options`
- `X-XSS-Protection`
- `Strict-Transport-Security`

**影响**: 应用容易受到点击劫持、内容类型嗅探和 XSS 攻击。

**建议实现**:
```go
func SecurityHeaders() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Header("X-Frame-Options", "DENY")
        c.Header("X-Content-Type-Options", "nosniff")
        c.Header("X-XSS-Protection", "1; mode=block")
        c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
        
        // 仅在 HTTPS 环境下启用
        if c.Request.TLS != nil {
            c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
        }
        
        c.Next()
    }
}
```

**优先级**: 🟡 中  
**CVSS 评分**: 6.1

---

#### **[MED-3] 缺少请求体大小限制**

**位置**: [internal/server/server.go](file:///d:/code/system-framework/internal/server/server.go)

**问题描述**: 未限制请求体大小，可能导致内存耗尽。

**影响**: 攻击者可以发送超大请求体导致服务拒绝。

**建议实现**:
```go
import "github.com/gin-gonic/gin"

func (s *Server) New() {
    engine := gin.New()
    
    // 设置最大请求体大小为 8MB
    engine.MaxMultipartMemory = 8 << 20
    
    // 添加请求大小限制中间件
    engine.Use(RequestBodySizeLimit(8 << 20))
}

func RequestBodySizeLimit(maxSize int64) gin.HandlerFunc {
    return func(c *gin.Context) {
        if c.Request.ContentLength > maxSize {
            c.JSON(http.StatusRequestEntityTooLarge, gin.H{
                "code":    -1,
                "message": fmt.Sprintf("Request body too large, max size: %d bytes", maxSize),
            })
            c.Abort()
            return
        }
        c.Next()
    }
}
```

**优先级**: 🟡 中  
**CVSS 评分**: 4.3

---

#### **[MED-4] 数据库连接未使用 TLS**

**位置**: [internal/database/db.go](file:///d:/code/system-framework/internal/database/db.go)

**问题描述**: 数据库连接未配置 TLS/SSL，可能被中间人攻击。

**影响**: 敏感数据在传输过程中可能被窃取。

**建议实现**:
```go
// MySQL TLS 配置
dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local&tls=true",
    cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Name)

// PostgreSQL TLS 配置
pgxConfig, err := pgxpool.ParseConfig(connStr)
pgxConfig.ConnConfig.TLSConfig = &tls.Config{
    MinVersion: tls.VersionTLS12,
    InsecureSkipVerify: false,
}
```

**配置文件**:
```yaml
database:
  tls_enabled: true
  tls_min_version: "TLSv1.2"
  ca_cert_path: "/path/to/ca.crt"
```

**优先级**: 🟡 中  
**CVSS 评分**: 6.5

---

#### **[MED-5] 缺少 API 版本控制**

**位置**: 项目全局

**问题描述**: API 未实现版本控制，可能导致安全更新困难。

**影响**: 无法优雅地废弃旧版本 API，难以实施安全补丁。

**建议实现**:
```go
// 路由分组
v1 := engine.Group("/api/v1")
v2 := engine.Group("/api/v2")

// 显式版本控制头
v1.Use(func(c *gin.Context) {
    c.Header("API-Version", "v1")
    c.Next()
})
```

**优先级**: 🟡 中  
**CVSS 评分**: 3.1

---

### 🟢 低危漏洞

#### **[LOW-1] JWT 错误消息可能泄露信息**

**位置**: [internal/auth/auth.go#L119-127](file:///d:/code/system-framework/internal/auth/auth.go#L119-127)

**问题描述**:
```go
c.JSON(http.StatusUnauthorized, gin.H{
    "code":    -1,
    "message": "Invalid token: " + err.Error(),
})
```

将错误详情返回给客户端可能泄露内部信息。

**建议修复**:
```go
// 生产环境不返回详细错误
if viper.GetString("app.mode") == "production" {
    c.JSON(http.StatusUnauthorized, gin.H{
        "code":    -1,
        "message": "Invalid token",
    })
} else {
    c.JSON(http.StatusUnauthorized, gin.H{
        "code":    -1,
        "message": "Invalid token: " + err.Error(),
    })
}
```

**优先级**: 🟢 低  
**CVSS 评分**: 2.1

---

#### **[LOW-2] 缺少审计日志**

**位置**: [internal/auth/auth.go](file:///d:/code/system-framework/internal/auth/auth.go)

**问题描述**: 未记录关键安全事件的审计日志，如登录失败、权限拒绝等。

**建议实现**:
```go
type AuditLog struct {
    Timestamp     time.Time `json:"timestamp"`
    EventType     string    `json:"event_type"`
    UserID        string    `json:"user_id,omitempty"`
    IPAddress     string    `json:"ip_address"`
    UserAgent     string    `json:"user_agent"`
    RequestPath   string    `json:"request_path"`
    Success       bool      `json:"success"`
    FailureReason string    `json:"failure_reason,omitempty"`
}

func LogSecurityEvent(event AuditLog) {
    logger.Log.Info("Security event",
        zap.Any("audit", event),
    )
}
```

**优先级**: 🟢 低  
**CVSS 评分**: 1.2

---

#### **[LOW-3] Nonce 验证过于简单**

**位置**: [internal/auth/auth.go#L237-247](file:///d:/code/system-framework/internal/auth/auth.go#L237-247)

**问题描述**:
```go
func validateNonce(nonce string) error {
    if nonce == "" {
        return nil
    }
    if len(nonce) < 16 {
        return fmt.Errorf("nonce too short")
    }
    return nil
}
```

仅检查长度，未验证唯一性和存在时间。

**建议实现**:
```go
import "github.com/google/uuid"

var nonceCache = cache.New(5*time.Minute, 10*time.Minute)

func validateNonce(nonce string) error {
    if nonce == "" {
        return nil
    }
    if len(nonce) < 16 {
        return fmt.Errorf("nonce too short")
    }
    
    // 检查 nonce 是否已使用
    if nonceCache.Get(nonce) != nil {
        return fmt.Errorf("nonce already used")
    }
    
    // 存储 nonce（5分钟过期）
    nonceCache.Set(nonce, true, 5*time.Minute)
    return nil
}
```

**优先级**: 🟢 低  
**CVSS 评分**: 3.1

---

#### **[LOW-4] 缺少 SQL 注入防护最佳实践**

**位置**: [internal/database/db.go](file:///d:/code/system-framework/internal/database/db.go)

**问题描述**: 虽然使用参数化查询，但缺少以下最佳实践：
- 查询超时设置
- 危险 SQL 命令黑名单
- 审计日志

**建议实现**:
```go
// 设置查询超时
ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
defer cancel()

// 危险操作二次确认
var dangerousCommands = []string{"DROP", "TRUNCATE", "DELETE"}
for _, cmd := range dangerousCommands {
    if strings.Contains(query, cmd) {
        // 记录审计日志或拒绝执行
        return fmt.Errorf("dangerous command detected: %s", cmd)
    }
}
```

**优先级**: 🟢 低  
**CVSS 评分**: 2.5

---

## 安全亮点 ✅

项目在以下方面表现良好：

1. **JWT 实现**: 使用 HS256 算法，正确设置过期时间和签名验证
2. **请求签名**: 实现了 HMAC-SHA256 签名验证和重放保护
3. **配置验证**: 实现了完整的配置校验机制
4. **熔断器**: 集成了 Hystrix-go 防止级联故障
5. **链路追踪**: 集成 OpenTelemetry 支持分布式追踪
6. **优雅关闭**: 支持信号处理和超时等待
7. **日志结构化**: 使用 Zap 进行结构化日志记录

---

## 总体安全评估

| 类别 | 评分 | 说明 |
|------|------|------|
| **整体评分** | 6.5/10 | 需要改进 |
| **认证授权** | 8.0/10 | JWT 实现良好 |
| **数据传输** | 5.0/10 | 缺少 TLS 配置 |
| **API 安全** | 5.5/10 | 缺少限流和版本控制 |
| **配置安全** | 3.0/10 | 敏感信息明文存储 |
| **基础设施** | 6.0/10 | 生产环境暴露调试端点 |

---

## 建议修复优先级

### 🔴 立即修复（24小时内）

1. **[HIGH-2]** 配置敏感信息加密存储
2. **[HIGH-3]** 生产环境禁用 pprof

### 🟡 本周内修复

1. **[HIGH-1]** 配置 CORS 白名单
2. **[MED-1]** 实现速率限制
3. **[MED-2]** 添加安全响应头

### 🟢 后续迭代

1. **[MED-3]** 添加请求体大小限制
2. **[MED-4]** 数据库连接启用 TLS
3. **[LOW-1-4]** 完善审计日志和防护机制

---

## 附录

### A. 依赖安全检查

建议定期运行以下命令检查依赖漏洞：

```bash
# Go 依赖漏洞扫描
go install github.com/securego/gosec/v2/cmd/gosec@latest
gosec ./...

# 依赖更新检查
go outdated -d
```

### B. 安全配置检查清单

- [ ] 生产环境 CORS 配置白名单
- [ ] 所有密钥使用环境变量或 KMS
- [ ] pprof 仅在非生产环境启用
- [ ] 启用请求速率限制
- [ ] 添加安全响应头
- [ ] 数据库连接启用 TLS
- [ ] 实现审计日志
- [ ] 启用 API 版本控制

---

**报告生成时间**: 2026-05-26  
**评审工具**: Manual Code Review + Security Best Practices  
**下一步行动**: 请查看上述建议修复优先级，开始修复工作
