# 认证授权模块

基于 JWT 和请求签名的双因素认证解决方案，支持 RBAC 权限控制。

## 功能特性

- 🔐 **JWT 认证**: HS256/HS512 签名算法，支持自定义 Claims
- 📝 **请求签名验证**: HMAC-SHA256 签名，防止请求篡改
- ⏰ **时间戳校验**: 自动校验请求时间，防止重放攻击
- 🎲 **Nonce 防重放**: 随机数机制，确保请求唯一性
- 👥 **RBAC 权限控制**: 基于角色的访问控制
- 🚫 **路径排除**: 白名单机制，无需认证的路径

## 快速开始

### 初始化

```go
import "github.com/lookingcamel/system-framework/internal/auth"

// 初始化认证模块
auth.Init(cfg.Auth)
```

### 生成 Token

```go
// 生成用户 Token
token, err := auth.GenerateToken("user123", "john", "admin")
if err != nil {
    log.Fatal(err)
}
fmt.Println(token)
```

### 解析 Token

```go
// 解析并验证 Token
claims, err := auth.ParseToken(tokenString)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("UserID: %s\n", claims.UserID)
fmt.Printf("Username: %s\n", claims.Username)
fmt.Printf("Role: %s\n", claims.Role)
```

### Gin 中间件

```go
import (
    "github.com/gin-gonic/gin"
    "github.com/lookingcamel/system-framework/internal/auth"
)

// 所有需要认证的路由
api := engine.Group("/api/v1", auth.JWT())

// 获取用户信息
api.GET("/profile", func(c *gin.Context) {
    userID := auth.GetUserID(c)
    username := auth.GetUsername(c)
    role := auth.GetRole(c)

    c.JSON(200, gin.H{
        "user_id":  userID,
        "username": username,
        "role":     role,
    })
})

// 管理员专属接口
api.GET("/admin", auth.RequireRole("admin"), adminHandler)

// 经理或管理员可访问
api.GET("/manager", auth.RequireAnyRole("admin", "manager"), managerHandler)
```

## 配置说明

```yaml
auth:
  enabled: true                      # 认证总开关
  jwt_enabled: true                  # JWT 认证开关
  jwt_secret: "your-256-bit-secret-key-here"  # JWT 密钥
  jwt_issuer: "system-framework"     # 签发者
  jwt_audience: "system-framework-api"  # 受众
  jwt_expire_seconds: 86400         # Token 过期时间（秒）

  signature_enabled: true            # 请求签名开关
  signature_secret: "your-signature-secret"  # 签名密钥
  signature_header: "X-Signature"   # 签名 Header
  timestamp_header: "X-Timestamp"    # 时间戳 Header
  nonce_header: "X-Nonce"           # Nonce Header
  timestamp_max_diff: 300           # 时间差容忍度（秒）

  exclude_paths:                    # 排除认证的路径
    - "/health"
    - "/ready"
    - "/metrics"
    - "/debug/pprof"
```

## 使用示例

### 请求签名算法

签名格式：
```
signature = HMAC-SHA256(secret, method + path + timestamp + body + nonce)
```

### 签名请求示例

```bash
# 生成签名
TIMESTAMP=$(date +%s)
NONCE=$(openssl rand -hex 16)
BODY='{"key": "value"}'
SIGNATURE=$(echo -n "POST/api/v1/data${TIMESTAMP}${BODY}${NONCE}" | openssl dgst -sha256 -hmac "your-signature-secret" | cut -d' ' -f2)

# 发送请求
curl -X POST http://localhost:8080/api/v1/data \
  -H "Content-Type: application/json" \
  -H "X-Timestamp: ${TIMESTAMP}" \
  -H "X-Nonce: ${NONCE}" \
  -H "X-Signature: ${SIGNATURE}" \
  -d "${BODY}"
```

### JWT Token 示例

```json
{
  "user_id": "user123",
  "username": "john",
  "role": "admin",
  "iss": "system-framework",
  "sub": "user123",
  "aud": ["system-framework-api"],
  "exp": 1704168000,
  "nbf": 1704081600,
  "iat": 1704081600,
  "jti": "550e8400-e29b-41d4-a716-446655440000"
}
```

## API 文档

### 主要函数

| 函数 | 说明 | 参数 | 返回值 |
|------|------|------|--------|
| `Init(config)` | 初始化认证模块 | AuthConfig | void |
| `GenerateToken(userID, username, role)` | 生成 JWT | userID, username, role | (token, error) |
| `ParseToken(token)` | 解析 JWT | token string | (*Claims, error) |
| `JWT()` | JWT 中间件 | - | gin.HandlerFunc |
| `Signature()` | 签名验证中间件 | - | gin.HandlerFunc |
| `RequireRole(role)` | 角色要求中间件 | role string | gin.HandlerFunc |
| `RequireAnyRole(roles...)` | 多角色要求中间件 | roles ...string | gin.HandlerFunc |
| `GetUserID(c)` | 获取用户 ID | *gin.Context | string |
| `GetUsername(c)` | 获取用户名 | *gin.Context | string |
| `GetRole(c)` | 获取用户角色 | *gin.Context | string |
| `GetClaims(c)` | 获取 Claims | *gin.Context | *Claims |

### Claims 结构

```go
type Claims struct {
    UserID   string `json:"user_id"`
    Username string `json:"username"`
    Role     string `json:"role"`
    jwt.RegisteredClaims
}
```

## 安全最佳实践

1. **密钥管理**
   - 生产环境使用强密钥（至少 32 字符）
   - 定期轮换密钥
   - 使用密钥管理服务（KMS）存储

2. **Token 安全**
   - 设置合理的过期时间
   - 使用 HTTPS 传输
   - 不在 URL 中传递 Token

3. **签名验证**
   - 启用时间戳校验
   - 使用 Nonce 防止重放
   - 设置合理的时间差容忍度

4. **日志审计**
   - 记录认证失败事件
   - 记录权限不足访问
   - 定期审计日志

## 故障排查

### 常见问题

**Q: Token 验证失败**
```
A: 检查 JWT 密钥是否一致，确认 Token 未过期
```

**Q: 签名验证失败**
```
A: 确认签名算法正确，检查时间戳是否过期
```

**Q: 权限不足**
```
A: 检查用户角色是否符合要求，使用 RequireAnyRole 允许多角色
```

## 相关文档

- [配置模块](../config/README.md)
- [中间件模块](../middleware/README.md)
- [项目总览](../../README.md)
