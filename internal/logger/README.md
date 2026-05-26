# 日志系统模块

基于 Uber Zap 的高性能结构化日志系统，支持多种输出格式和日志轮转。

## 功能特性

- ⚡ **高性能**: Uber Zap 高性能日志库
- 📋 **结构化日志**: JSON 格式，易于解析和分析
- 🔄 **多输出**: 控制台、文件、双输出
- 🔄 **日志轮转**: 自动轮转，避免日志文件过大
- 📦 **字段丰富**: 支持添加自定义字段
- 🎯 **日志级别**: debug, info, warn, error, fatal
- 📝 **Caller 信息**: 显示调用位置信息
- 🚀 **异步写入**: 高性能异步日志写入

## 快速开始

### 初始化

```go
import "github.com/lookingcamel/system-framework/internal/logger"

func main() {
    // 从配置初始化
    err := logger.Init(cfg.Logging)
    if err != nil {
        log.Fatal(err)
    }
    defer logger.Sync()
}
```

### 基本使用

```go
import "go.uber.org/zap"

// Info 日志
logger.Log.Info("Server started",
    zap.String("host", "localhost"),
    zap.Int("port", 8080),
)

// 带字段的日志
logger.Log.Info("User login",
    zap.String("user_id", "123"),
    zap.String("username", "john"),
    zap.String("ip", "192.168.1.1"),
)

// 错误日志
logger.Log.Error("Database connection failed",
    zap.Error(err),
)

// 警告日志
logger.Log.Warn("Rate limit approaching",
    zap.Int("requests", 999),
    zap.Int("limit", 1000),
)

// Debug 日志
logger.Log.Debug("Processing request",
    zap.String("request_id", "abc123"),
)
```

## 配置说明

```yaml
logging:
  level: "info"              # 日志级别: debug, info, warn, error, fatal
  format: "json"             # 输出格式: json, console
  output: "console"          # 输出方式: console, file, both
  output_path: "logs/app.log"  # 日志文件路径
  max_size: 100              # 单个日志文件大小（MB）
  max_backups: 7             # 保留的日志文件数量
  max_age: 30                # 日志文件保留天数
  compress: true             # 是否压缩历史日志
  enable_caller: true        # 是否显示调用者信息
  enable_stacktrace: true    # 是否显示堆栈跟踪
```

## 日志输出格式

### JSON 格式（生产环境）

```json
{
  "level": "info",
  "ts": 1704096000.123,
  "caller": "main.go:45",
  "msg": "Server started",
  "host": "localhost",
  "port": 8080,
  "request_id": "abc123",
  "trace_id": "xyz789"
}
```

### Console 格式（开发环境）

```
2024-01-01 12:00:00.000  INFO   main.go:45   Server started  host=localhost port=8080
2024-01-01 12:00:01.234  ERROR  db.go:78    Database error  error="connection refused"
```

## 高级用法

### 记录请求 ID

```go
// 在中间件中设置请求 ID
logger.Log.Info("HTTP request",
    zap.String("request_id", requestID),
    zap.String("method", c.Request.Method),
    zap.String("path", c.Request.URL.Path),
    zap.Int("status", status),
    zap.Duration("latency", latency),
)
```

### 结构化错误日志

```go
// 使用字段记录错误
logger.Log.Error("Operation failed",
    zap.String("operation", "database_query"),
    zap.String("table", "users"),
    zap.Int("user_id", 123),
    zap.Error(err),
    zap.Stack(),
)
```

### 上下文日志

```go
// 在请求上下文中添加日志字段
ctx := logger.WithContext(context.Background())
ctx = logger.WithField(ctx, "request_id", requestID)
ctx = logger.WithFields(ctx, zap.String("user_id", userID))

logger.LogFromContext(ctx).Info("Processing request")
```

### 性能日志

```go
// 记录性能指标
logger.Log.Info("Query completed",
    zap.String("query", "SELECT * FROM users"),
    zap.Duration("duration", elapsed),
    zap.Int64("rows_affected", rows),
)
```

## 日志文件管理

### 文件轮转策略

```
logs/
├── app.log          # 当前日志文件
├── app-20240101.log.gz  # 压缩的历史日志
├── app-20240102.log.gz
└── app-20240107.log.gz  # 超过 7 天自动删除
```

### 自定义文件名

```go
logger.Init(LoggingConfig{
    OutputPath: "logs/service-%d.log",  // 支持日期格式化
})
```

## 日志级别

| 级别 | 使用场景 | 输出 |
|------|----------|------|
| Debug | 开发调试信息 | 详细调试信息 |
| Info | 正常运行信息 | 重要业务事件 |
| Warn | 警告信息 | 潜在问题 |
| Error | 错误信息 | 错误事件 |
| Fatal | 致命错误 | 致命错误并退出 |

## 日志收集

### Filebeat 配置

```yaml
filebeat.inputs:
  - type: log
    paths:
      - /var/log/app/*.log
    json.keys_under_root: true
    json.add_error_key: true

output.elasticsearch:
  hosts: ["elasticsearch:9200"]
```

### Fluentd 配置

```xml
<source>
  @type tail
  path /var/log/app/app.log
  pos_file /var/log/app/app.log.pos
  <parse>
    @type json
  </parse>
</source>

<match **>
  @type elasticsearch
  hosts elasticsearch:9200
</match>
```

## 最佳实践

1. **日志级别使用**
   - Debug: 仅开发环境使用
   - Info: 记录正常业务事件
   - Warn: 记录可恢复的错误
   - Error: 记录需要关注的错误
   - Fatal: 仅记录致命错误

2. **字段设计**
   - 使用统一字段命名
   - 避免记录敏感信息
   - 保持字段一致性

3. **性能优化**
   - 生产环境使用异步模式
   - 合理设置日志级别
   - 定期清理历史日志

4. **安全考虑**
   - 不记录密码
   - 不记录完整信用卡号
   - 脱敏处理敏感数据

## 故障排查

### 常见问题

**Q: 日志文件无法写入**
```
A: 检查目录权限，确认磁盘空间充足
```

**Q: 日志丢失**
```
A: 确保调用 logger.Sync() 刷新缓冲区
```

**Q: 日志格式错误**
```
A: 检查 JSON 格式配置，确认无特殊字符
```

## 相关文档

- [中间件模块](../middleware/README.md)
- [链路追踪模块](../tracing/README.md)
- [项目总览](../../README.md)
