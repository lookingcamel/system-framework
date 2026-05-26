# 配置管理模块

支持 YAML 配置文件、环境变量和 Nacos 配置中心的配置管理解决方案。

## 功能特性

- 📄 **多格式支持**: YAML、JSON、ENV
- 🔄 **热更新**: 运行时配置动态更新
- ☁️ **Nacos 集成**: 云原生配置中心支持
- ✅ **配置校验**: 自动校验配置合法性
- 🔧 **默认值**: 智能设置合理默认值
- 🌐 **环境变量**: 支持环境变量覆盖
- 🔍 **配置变更通知**: 配置变更事件回调

## 快速开始

### 加载配置

```go
import "github.com/lookingcamel/system-framework/internal/config"

func main() {
    cfg, err := config.Load("config.yaml")
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("App: %s:%d\n", cfg.App.Host, cfg.App.Port)
}
```

### 获取全局配置

```go
// 在任意位置获取全局配置
cfg := config.GetGlobalConfig()
```

### 重载配置

```go
// 重新加载配置文件
err := config.Reload()
if err != nil {
    log.Fatal(err)
}
```

## 配置结构

### 完整配置示例

```yaml
app:
  name: "system-framework"
  version: "1.0.0"
  host: "0.0.0.0"
  port: 8080
  mode: "release"

server:
  read_timeout: 30
  write_timeout: 30
  idle_timeout: 120

logging:
  level: "info"
  format: "json"
  output: "console"
  output_path: "logs/app.log"
  max_size: 100
  max_backups: 7
  max_age: 30
  compress: true

prometheus:
  enabled: true
  path: "/metrics"

pprof:
  enabled: true
  path: "/debug/pprof"

tracing:
  enabled: true
  service_name: "system-framework"
  exporter: "otlp"
  endpoint: "localhost:4317"
  sample_rate: 1.0

health:
  enabled: true
  path: "/health"

graceful_shutdown:
  timeout: 30

database:
  type: "sqlite"
  host: "localhost"
  port: 3306
  name: "app"
  username: "admin"
  password: "password"
  sqlite_path: "./data/app.db"
  max_open_connections: 10
  max_idle_connections: 5
  connection_max_lifetime: 300
  migration_path: "./migrations"
  backup_enabled: true
  backup_dir: "./backups"

redis:
  enabled: true
  host: "localhost"
  port: 6379
  password: ""
  db: 0
  pool_size: 10
  min_idle_conns: 5

nacos:
  enabled: false
  server_addr: "localhost"
  server_port: 8848
  namespace_id: ""
  group: "DEFAULT_GROUP"
  data_id: "system-framework.yaml"
  username: "nacos"
  password: "nacos"
  config_type: "yaml"
  refresh_enable: true
  refresh_delay: 30

circuit_breaker:
  enabled: true
  default_timeout: 3000
  default_max_concurrent: 100
  default_error_percentage: 50
  default_request_volume: 20

auth:
  enabled: false
  jwt_enabled: true
  jwt_secret: "your-secret"
  jwt_issuer: "system-framework"
  jwt_audience: "system-framework-api"
  jwt_expire_seconds: 86400
```

## 配置校验

### 自动校验规则

| 配置项 | 校验规则 | 默认值 |
|--------|----------|--------|
| `app.port` | 1-65535 | 8080 |
| `app.mode` | debug/release/test | release |
| `server.*_timeout` | > 0 | 30 |
| `logging.level` | debug/info/warn/error | info |
| `logging.output` | console/file/both | console |
| `database.port` | 1-65535 | 3306 |
| `redis.port` | 1-65535 | 6379 |
| `auth.jwt_secret` | 至少 32 字符 | - |

### 校验示例

```go
validator := config.NewValidator()
if err := validator.Validate(cfg); err != nil {
    fmt.Println(err.Error())
}
```

### 自定义验证函数

```go
// 范围验证
err := config.ValidateRange(value, min, max, "field")

// 模式验证
err := config.ValidatePattern(value, `^[a-z]+$`, "field")

// 非空验证
err := config.ValidateNotEmpty(value, "field")

// 邮箱验证
err := config.ValidateEmail("user@example.com")

// URL 验证
err := config.ValidateURL("https://example.com")
```

## Nacos 配置中心

### 启用 Nacos

```yaml
nacos:
  enabled: true
  server_addr: "nacos-server"
  server_port: 8848
  namespace_id: "public"
  group: "DEFAULT_GROUP"
  data_id: "application.yaml"
  username: "nacos"
  password: "nacos"
  config_type: "yaml"
  refresh_enable: true
  refresh_delay: 30
```

### 配置监听

```go
// 配置变更时自动回调
// 通过 config.GetGlobalConfig() 获取最新配置
```

### 环境变量覆盖

配置支持环境变量覆盖，格式：`NACOS_ENABLED=true`

```yaml
nacos:
  enabled: ${NACOS_ENABLED:false}
  server_addr: ${NACOS_SERVER_ADDR:localhost}
```

## 默认值设置

### 自动设置的默认值

```go
// App 配置
if cfg.App.Name == "" {
    cfg.App.Name = "system-framework"
}
if cfg.App.Port <= 0 {
    cfg.App.Port = 8080
}

// 数据库配置
if cfg.Database.MaxOpenConnections <= 0 {
    cfg.Database.MaxOpenConnections = 10
}

// Redis 配置
if cfg.Redis.PoolSize <= 0 {
    cfg.Redis.PoolSize = 10
}

// 认证配置
if cfg.Auth.JWTIssuer == "" {
    cfg.Auth.JWTIssuer = "system-framework"
}
```

## 配置热更新

### Nacos 配置热更新

```go
// 启用热更新
nacos:
  enabled: true
  refresh_enable: true
  refresh_delay: 30  # 检查间隔（秒）
```

### 手动重载

```go
// 监听配置变更
go func() {
    ticker := time.NewTicker(30 * time.Second)
    for range ticker.C {
        err := config.Reload()
        if err != nil {
            log.Printf("Reload config failed: %v", err)
        }
    }
}()
```

## 最佳实践

1. **配置分离**
   - 开发环境：config-dev.yaml
   - 测试环境：config-test.yaml
   - 生产环境：config-prod.yaml

2. **敏感信息**
   - 使用环境变量存储敏感信息
   - 不将密码提交到版本控制
   - 使用密钥管理服务

3. **配置文档化**
   - 为每个配置项添加注释
   - 维护配置变更日志
   - 文档示例配置

## 故障排查

### 常见问题

**Q: 配置文件读取失败**
```
A: 检查文件路径是否正确，确认文件格式是否为有效的 YAML
```

**Q: Nacos 连接失败**
```
A: 检查 Nacos 服务是否可用，确认网络连通性
```

**Q: 配置校验失败**
```
A: 查看日志中的校验错误信息，修正配置或使用默认值
```

## 相关文档

- [认证模块](../auth/README.md)
- [数据库模块](../database/README.md)
- [Redis 模块](../redis/README.md)
- [项目总览](../../README.md)
