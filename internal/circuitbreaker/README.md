# 熔断降级模块

基于 Hystrix-go 的熔断器实现，保护系统免受级联故障影响。

## 功能特性

- 🔥 **熔断保护**: 快速失败，防止雪崩
- ⏱️ **超时控制**: 防止请求无限等待
- 📊 **并发限制**: 控制最大并发数
- 📈 **指标监控**: 实时熔断器状态
- 🔄 **自动恢复**: 熔断后自动尝试恢复
- 🎯 **降级处理**: 提供备选响应

## 快速开始

### 初始化

```go
import "github.com/lookingcamel/system-framework/internal/circuitbreaker"

func main() {
    circuitbreaker.Init(cfg.CircuitBreaker)
}
```

### 基本使用

```go
import (
    "github.com/lookingcamel/system-framework/internal/circuitbreaker"
)

// 执行带熔断保护的调用
err := circuitbreaker.Execute("GetUser",
    func() error {
        // 业务逻辑
        return getUserFromDB(123)
    },
    func(err error) error {
        // 降级逻辑
        return getUserFromCache(123)
    },
)
```

### Gin 中间件

```go
import "github.com/lookingcamel/system-framework/internal/circuitbreaker"

// 使用默认降级
router.GET("/api/users/:id",
    circuitbreaker.GinMiddleware("GetUser"),
    getUserHandler,
)

// 自定义降级
router.GET("/api/users/:id",
    circuitbreaker.GinMiddlewareWithFallback("GetUser",
        func(c *gin.Context) {
            c.JSON(200, gin.H{
                "id":      c.Param("id"),
                "name":    "Fallback User",
                "message": "Service temporarily unavailable",
            })
        },
    ),
    getUserHandler,
)
```

## 配置说明

```yaml
circuit_breaker:
  enabled: true
  default_timeout: 3000              # 默认超时时间（毫秒）
  default_max_concurrent: 100        # 最大并发数
  default_error_percentage: 50        # 错误百分比阈值
  default_request_volume: 20           # 请求量阈值
```

### 配置参数

| 参数 | 说明 | 默认值 |
|------|------|--------|
| `enabled` | 是否启用熔断器 | true |
| `default_timeout` | 默认超时时间（ms） | 3000 |
| `default_max_concurrent` | 最大并发数 | 100 |
| `default_error_percentage` | 错误百分比阈值 | 50 |
| `default_request_volume` | 滚动窗口最小请求数 | 20 |

## 工作原理

### 熔断器状态

```
            ┌───────────────┐
            │    CLOSED     │
            │  正常运行     │
            └───────┬───────┘
                    │
                    │ 失败率 > 阈值
                    ▼
            ┌───────────────┐
            │     OPEN      │
            │  熔断拒绝     │
            └───────┬───────┘
                    │
                    │ 休眠期结束
                    ▼
            ┌───────────────┐
            │ HALF-OPEN     │
            │  尝试恢复     │
            └───────┬───────┘
                    │
        ┌───────────┴───────────┐
        │                       │
        │  成功                  │  失败
        ▼                       ▼
┌───────────────┐       ┌───────────────┐
│    CLOSED     │       │     OPEN      │
│  恢复正常     │       │  继续熔断     │
└───────────────┘       └───────────────┘
```

### 配置说明

- **CLOSE**: 正常状态，所有请求正常执行
- **OPEN**: 熔断状态，所有请求快速失败
- **HALF-OPEN**: 半开状态，允许部分请求尝试恢复

## API 文档

### 主要函数

| 函数 | 说明 | 参数 | 返回值 |
|------|------|------|--------|
| `Init(cfg)` | 初始化熔断器 | CircuitBreakerConfig | void |
| `Execute(cmd, run, fallback)` | 执行熔断保护 | command, runFunc, fallbackFunc | error |
| `ExecuteWithResult(cmd, run, fallback)` | 执行并返回结果 | command, runFunc, fallbackFunc | (result, error) |
| `ExecuteWithContext(ctx, cmd, run, fallback)` | 带上下文的执行 | ctx, command, runFunc, fallbackFunc | error |
| `GinMiddleware(command)` | Gin 中间件 | command name | gin.HandlerFunc |
| `GinMiddlewareWithFallback(command, fallback)` | 带降级的中间件 | command, fallbackHandler | gin.HandlerFunc |
| `GetCircuitStatus()` | 获取熔断状态 | - | gin.HandlerFunc |

## 使用场景

### 1. 外部 API 调用

```go
err := circuitbreaker.Execute("ExternalAPI",
    func() error {
        resp, err := http.Get("https://api.example.com/data")
        if err != nil {
            return err
        }
        return processResponse(resp)
    },
    func(err error) error {
        // 返回缓存数据
        return getCachedData()
    },
)
```

### 2. 数据库查询

```go
err := circuitbreaker.Execute("DatabaseQuery",
    func() error {
        return db.QueryUser(123)
    },
    func(err error) error {
        // 返回默认数据
        return defaultUser()
    },
)
```

### 3. 缓存降级

```go
result, err := circuitbreaker.ExecuteWithResult("CacheOrDB",
    func() error {
        data, err := getFromDB()
        if err == nil {
            cache.Set("key", data)
        }
        return nil
    },
    func(err error) error {
        data := cache.Get("key")
        if data == nil {
            return errors.New("no data available")
        }
        return nil
    },
)
```

## 监控指标

### Prometheus 指标

```
# 熔断器调用指标
hystrix_commands{command="GetUser", event="success"} 100
hystrix_commands{command="GetUser", event="failure"} 5
hystrix_commands{command="GetUser", event="timeout"} 2
hystrix_commands{command="GetUser", event="fallback"} 3

# 熔断器状态
hystrix_circuit_breaker{command="GetUser", name="OPEN"} 1
hystrix_circuit_breaker{command="GetUser", name="CLOSED"} 0

# 并发执行数
hystrix_current_execution{command="GetUser"} 5
```

### 查询熔断器状态

```bash
curl http://localhost:8080/circuit/status
```

## 最佳实践

1. **合理设置超时**
   - 根据服务响应时间设置
   - 避免过长或过短
   - 考虑网络延迟

2. **降级策略**
   - 返回默认值
   - 返回缓存数据
   - 返回友好错误信息

3. **监控告警**
   - 监控熔断器触发次数
   - 监控降级调用次数
   - 设置告警阈值

4. **配置隔离**
   - 不同服务使用不同配置
   - 核心服务优先保护
   - 合理设置错误阈值

## 故障排查

### 常见问题

**Q: 所有请求都失败**
```
A: 检查熔断器是否处于 OPEN 状态，等待恢复或手动关闭
```

**Q: 降级逻辑未执行**
```
A: 确认降级函数返回值，确保 run 函数返回错误
```

**Q: 超时配置无效**
```
A: 检查是否使用 ExecuteWithContext，使用 context 设置超时
```

## 相关文档

- [认证模块](../auth/README.md)
- [Redis 模块](../redis/README.md)
- [数据库模块](../database/README.md)
- [项目总览](../../README.md)
