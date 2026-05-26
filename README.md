# System Framework

基于 Go 1.25 + Gin 框架的企业级高性能 Web 服务框架，集成监控、安全、数据库、缓存等企业级功能。

## 📋 目录

- [特性](#-特性)
- [快速开始](#-快速开始)
- [项目结构](#-项目结构)
- [核心模块](#-核心模块)
- [部署指南](#-部署指南)
- [配置说明](#-配置说明)
- [API 文档](#-api-文档)
- [监控与可观测性](#-监控与可观测性)
- [安全](#-安全)
- [测试](#-测试)
- [性能优化](#-性能优化)
- [最佳实践](#-最佳实践)
- [贡献指南](#-贡献指南)
- [许可证](#-许可证)

## ✨ 特性

### 核心能力

- 🎯 **高性能 Web 框架**: 基于 Gin 框架，支持中间件链式调用
- 🔄 **热更新配置**: 支持 Nacos 配置中心，运行时动态更新配置
- 📊 **多数据库支持**: MySQL、PostgreSQL、SQLite 三大主流数据库
- 🔗 **Redis 缓存**: 高性能 Redis 缓存支持
- 📈 **Prometheus 监控**: 完整的指标采集和监控体系
- 🔍 **分布式追踪**: OpenTelemetry + OTLP 支持
- 🛡️ **熔断降级**: 基于 Hystrix-go 的熔断器保护
- 🚦 **速率限制**: 基于令牌桶的请求限流，防止暴力破解
- 📦 **请求体大小限制**: 可配置的请求体大小限制，防止内存耗尽攻击
- 🔢 **API 版本控制**: URL 路径版本控制，支持多版本共存和弃用通知
- 🔐 **安全认证**: JWT + 请求签名验证
- 🗄️ **数据库迁移**: golang-migrate 支持版本化迁移
- 💾 **自动备份**: 数据库定时备份和恢复
- 🐳 **容器化部署**: Docker 和 Docker Compose 支持
- ☸️ **Kubernetes 部署**: Helm Chart 完整支持
- 🔄 **灰度发布**: 金丝雀发布和流量策略
- 📝 **结构化日志**: Zap 日志系统，支持 JSON 格式
- 🏥 **健康检查**: 就绪探针和存活探针
- 🎨 **优雅关闭**: 支持零停机部署
- 🔧 **性能分析**: pprof 性能分析工具（生产环境自动禁用）
- ✅ **完整测试**: 50+ 单元测试用例，覆盖所有核心模块

### 技术栈

- **语言**: Go 1.25+
- **框架**: Gin Web Framework
- **数据库**: MySQL / PostgreSQL / SQLite
- **缓存**: Redis
- **监控**: Prometheus + Grafana
- **追踪**: OpenTelemetry + OTLP
- **容器**: Docker, Kubernetes

## 🚀 快速开始

### 前置要求

- Go 1.25+
- MySQL 8.0+ / PostgreSQL 14+ / SQLite 3
- Redis 6.0+ (可选)
- Docker & Docker Compose (容器部署)
- Kubernetes & Helm (K8s 部署)

### 安装部署

#### 方式一：本地运行

```bash
# 克隆项目
git clone https://github.com/lookingcamel/system-framework.git
cd system-framework

# 下载依赖
go mod download

# 编译项目
go build -o bin/server.exe ./cmd/server

# 运行服务
./bin/server.exe -config config.yaml
```

#### 方式二：Docker Compose 快速启动

```bash
# 启动所有服务
docker-compose up -d

# 查看服务状态
docker-compose ps

# 查看日志
docker-compose logs -f app

# 停止服务
docker-compose down
```

#### 方式三：Kubernetes 部署

```bash
# 添加 Helm 仓库
helm repo add lookingcamel https://charts.example.com
helm repo update

# 安装应用
helm install system-framework ./deployments/helm

# 升级应用
helm upgrade system-framework ./deployments/helm

# 回滚应用
helm rollback system-framework 1
```

### 验证服务

```bash
# 健康检查
curl http://localhost:8080/health

# Prometheus 指标
curl http://localhost:8080/metrics

# API 示例
curl http://localhost:8080/api/v1/example/1
```

## 📁 项目结构

```
system-framework/
├── cmd/                          # 应用程序入口
│   └── server/
│       └── main.go              # 主程序
├── internal/                     # 内部包
│   ├── auth/                    # 认证授权
│   ├── backup/                  # 数据库备份
│   ├── circuitbreaker/           # 熔断降级
│   ├── config/                  # 配置管理 + 验证
│   ├── database/                # 数据库连接
│   ├── handler/                 # HTTP 处理器
│   ├── logger/                  # 日志系统
│   ├── metrics/                 # Prometheus 指标
│   ├── middleware/               # 中间件
│   │   ├── ratelimit.go        # 速率限制
│   │   ├── bodylimit.go        # 请求体限制
│   │   └── middleware.go        # 其他中间件
│   ├── migrate/                 # 数据库迁移
│   ├── redis/                   # Redis 缓存
│   ├── server/                  # HTTP 服务器
│   ├── tracing/                 # 链路追踪
│   └── versioning/              # API 版本控制
├── pkg/                         # 公共工具包
│   └── utils/                   # 工具函数
├── migrations/                   # 数据库迁移文件
├── deployments/                  # 部署配置
│   ├── docker/                  # Docker 配置
│   └── helm/                    # Helm Chart
├── config.yaml                  # 配置文件
├── docker-compose.yml            # Docker Compose 配置
├── Dockerfile                    # Docker 镜像构建
└── README.md                     # 项目说明
```

## 🔧 核心模块

### 认证模块 (auth)

JWT 认证和请求签名验证，支持 RBAC 权限控制。

- JWT (HS256/HS512) 令牌生成和验证
- HMAC-SHA256 请求签名
- 权限中间件

详细文档: [internal/auth/README.md](internal/auth/README.md)

### 配置模块 (config)

支持 YAML 配置、环境变量和 Nacos 配置中心，运行时热更新。

- 配置验证和默认值自动补全
- 多环境配置支持
- 敏感信息加密

详细文档: [internal/config/README.md](internal/config/README.md)

### 日志模块 (logger)

基于 Zap 的结构化日志，支持多种输出格式和日志轮转。

- JSON/Console 双格式
- 日志级别控制
- 日志轮转和压缩

详细文档: [internal/logger/README.md](internal/logger/README.md)

### 数据库模块 (database)

支持 MySQL、PostgreSQL、SQLite，自动迁移和备份。

- 连接池管理
- 自动数据库迁移
- 定时备份和恢复

详细文档: [internal/database/README.md](internal/database/README.md)

### Redis 模块 (redis)

连接池管理，自动重连，分布式缓存支持。

详细文档: [internal/redis/README.md](internal/redis/README.md)

### 熔断器模块 (circuitbreaker)

基于 Hystrix-go 的熔断降级，防止级联故障。

- 断路器状态管理
- 自动降级和恢复
- 性能指标监控

详细文档: [internal/circuitbreaker/README.md](internal/circuitbreaker/README.md)

### 链路追踪模块 (tracing)

OpenTelemetry + OTLP 分布式追踪。

- 全链路追踪
- Span 关联分析
- 多Exporter支持

详细文档: [internal/tracing/README.md](internal/tracing/README.md)

### 指标模块 (metrics)

Prometheus 指标采集和暴露。

- HTTP 请求指标
- 数据库连接池指标
- 自定义业务指标

详细文档: [internal/metrics/README.md](internal/metrics/README.md)

### 中间件模块 (middleware)

#### 速率限制 ([middleware/ratelimit.go](internal/middleware/ratelimit.go))

- 基于令牌桶算法
- API Key/IP/User 多维度限流
- 可配置 QPS 和突发限制

#### 请求体限制 ([middleware/bodylimit.go](internal/middleware/bodylimit.go))

- 可配置最大请求体（默认 8MB）
- Content-Length 预检查
- 防止内存耗尽攻击

#### 其他中间件

- 请求ID追踪
- Recovery 异常恢复
- CORS 跨域支持
- 请求日志

详细文档: [internal/middleware/README.md](internal/middleware/README.md)

### 数据库迁移 (migrate)

版本化管理数据库 schema 变更。

详细文档: [internal/migrate/README.md](internal/migrate/README.md)

### 数据库备份 (backup)

自动备份和手动恢复功能。

详细文档: [internal/backup/README.md](internal/backup/README.md)

### API 版本控制 (versioning)

URL 路径版本控制，支持多版本 API 共存、弃用通知和版本协商。

特性:
- 路径版本: `/api/v1/...`, `/api/v2/...`
- 自动版本协商: 不带版本的请求使用默认版本
- 弃用警告: 旧版本返回 `X-API-Deprecation-Notice` 响应头
- 动态路由注册: 配置驱动的版本管理

详细文档: [internal/versioning/README.md](internal/versioning/README.md)

### 服务器模块 (server)

HTTP 服务器核心，集成所有中间件和路由。

详细文档: [internal/server/README.md](internal/server/README.md)

## 🚢 部署指南

### Docker 部署

#### 单容器部署

```bash
# 构建镜像
docker build -t system-framework:latest .

# 运行容器
docker run -d \
  --name system-framework \
  -p 8080:8080 \
  -v $(pwd)/config.yaml:/app/config.yaml \
  -v $(pwd)/data:/app/data \
  system-framework:latest
```

#### Docker Compose 部署

```yaml
# docker-compose.yml
version: '3.8'

services:
  app:
    build: .
    container_name: system-framework
    ports:
      - "8080:8080"
    environment:
      - TZ=Asia/Shanghai
    volumes:
      - ./config.yaml:/app/config.yaml
      - ./data:/app/data
      - ./logs:/app/logs
      - ./backups:/app/backups
    depends_on:
      - mysql
      - redis
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8080/health"]
      interval: 30s
      timeout: 10s
      retries: 3

  mysql:
    image: mysql:8.0
    container_name: mysql
    environment:
      MYSQL_ROOT_PASSWORD: password
      MYSQL_DATABASE: app
    volumes:
      - mysql_data:/var/lib/mysql
    ports:
      - "3306:3306"
    restart: unless-stopped

  redis:
    image: redis:7-alpine
    container_name: redis
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data
    restart: unless-stopped

volumes:
  mysql_data:
  redis_data:
```

### Kubernetes 部署

#### 安装 Helm Chart

```bash
# 创建命名空间
kubectl create namespace system-framework

# 安装
helm install system-framework ./deployments/helm \
  --namespace system-framework \
  --set image.tag=v1.0.0 \
  --set resources.limits.cpu=2000m \
  --set resources.limits.memory=2Gi

# 查看部署状态
kubectl get pods -n system-framework

# 查看日志
kubectl logs -n system-framework -l app=system-framework

# 升级
helm upgrade system-framework ./deployments/helm \
  --namespace system-framework \
  --set image.tag=v1.1.0

# 回滚
helm rollback system-framework -n system-framework
```

#### 金丝雀发布

```bash
# 部署金丝雀版本（10% 流量）
helm upgrade system-framework ./deployments/helm \
  --namespace system-framework \
  --set traffic.canary.enabled=true \
  --set traffic.canary.weight=10

# 全量切换
helm upgrade system-framework ./deployments/helm \
  --namespace system-framework \
  --set traffic.canary.enabled=false
```

#### Horizontal Pod Autoscaler

```bash
kubectl autoscale deployment system-framework \
  --namespace system-framework \
  --cpu-percent=80 \
  --min=2 \
  --max=10
```

详细文档: [deployments/helm/README.md](deployments/helm/README.md)

## ⚙️ 配置说明

### 配置文件结构

```yaml
app:
  name: "system-framework"
  version: "1.0.0"
  host: "0.0.0.0"
  port: 8080
  mode: "release"  # debug/release

server:
  read_timeout: 30       # HTTP 读取超时（秒）
  write_timeout: 30      # HTTP 写入超时（秒）
  idle_timeout: 120      # 空闲连接超时（秒）
  max_body_size: 8388608  # 最大请求体大小（8MB）

logging:
  level: "info"          # debug/info/warn/error
  format: "json"        # json/console
  output: "console"      # console/file
  output_path: "logs/app.log"
  max_size: 100         # 日志文件大小（MB）
  max_backups: 7        # 保留日志文件数
  max_age: 30           # 日志保留天数
  compress: true        # 压缩旧日志

database:
  type: "sqlite"        # mysql/postgresql/sqlite
  sqlite_path: "./data/app.db"
  # MySQL 配置示例:
  # host: "localhost"
  # port: 3306
  # username: "root"
  # password: "password"
  # database: "app"
  max_open_connections: 10
  max_idle_connections: 5
  connection_max_lifetime: 300
  migration_path: "./migrations"
  backup_enabled: true
  backup_dir: "./backups"
  backup_interval: 86400  # 备份间隔（秒，默认每天）

redis:
  enabled: true
  host: "localhost"
  port: 6379
  password: ""
  db: 0
  pool_size: 10
  min_idle_connections: 5
  dial_timeout: 5
  read_timeout: 3
  write_timeout: 3

auth:
  enabled: true
  jwt_enabled: true
  jwt_secret: "your-256-bit-secret-key-change-in-production"
  jwt_issuer: "system-framework"
  jwt_audience: "system-framework-api"
  jwt_expire_seconds: 86400
  signature_enabled: false
  signature_secret: "your-signature-secret"
  exclude_paths:
    - "/health"
    - "/ready"
    - "/metrics"

rate_limit:
  enabled: true
  requests_per_second: 100
  burst: 200
  exclude_paths:
    - "/health"
    - "/ready"
    - "/metrics"

api_version:
  enabled: true
  default_version: "v1"
  supported_versions:
    - "v1"
    - "v2"
  deprecation_notice: "API v1 will be deprecated on 2026-12-31. Please migrate to v2."

circuit_breaker:
  enabled: true
  default_timeout: 3000          # 超时时间（毫秒）
  default_max_concurrent: 100    # 最大并发数
  default_error_percentage: 50  # 错误百分比阈值

tracing:
  enabled: true
  service_name: "system-framework"
  exporter: "otlp"              # otlp/jaeger/zipkin
  endpoint: "localhost:4317"    # OTLP 接收端点
  sample_rate: 1.0              # 采样率 (0.0-1.0)

prometheus:
  enabled: true
  path: "/metrics"
```

### 环境变量覆盖

配置项可以通过环境变量覆盖：

```bash
# 示例
export APP_PORT=9090
export DATABASE_TYPE=mysql
export REDIS_HOST=redis.production
```

详细配置文档: [internal/config/README.md](internal/config/README.md)

## 📚 API 文档

### API 版本控制

项目支持多版本 API 共存，使用 URL 路径版本控制。

#### 版本格式

```
/api/{version}/{resource}
/api/v1/example/:id
/api/v2/example/:id
```

#### 版本协商

- **显式版本**: 请求中明确指定版本
  ```
  GET /api/v1/example/123
  GET /api/v2/example/123
  ```

- **隐式版本**: 不指定版本时使用默认版本（v1）
  ```
  GET /api/example/123  → 自动路由到 v1
  ```

#### 弃用响应头

请求旧版本 API 时，响应会包含弃用警告头：

```http
HTTP/1.1 200 OK
X-API-Deprecation-Notice: API v1 will be deprecated on 2026-12-31. Please migrate to v2.
```

### 健康检查

```bash
# 存活探针 - 检查应用是否运行
GET /health

# 就绪探针 - 检查应用是否准备好接收流量
GET /ready
```

### Prometheus 指标

```bash
GET /metrics
```

### 熔断器状态

```bash
GET /circuit/status
```

### 示例 API

#### v1 版本

```bash
# 获取单个资源
GET /api/v1/example/:id

# 列出所有资源
GET /api/v1/examples
```

#### v2 版本

```bash
# 获取单个资源（改进的响应格式）
GET /api/v2/example/:id

# 列出所有资源
GET /api/v2/examples
```

详细 API 文档: [internal/handler/README.md](internal/handler/README.md)

## 📊 监控与可观测性

### Prometheus 监控

项目集成完整的 Prometheus 指标：

#### HTTP 指标

```promql
# 请求率
rate(http_requests_total[5m])

# 延迟分布
histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m]))

# 错误率
rate(http_requests_total{status=~"5.."}[5m])
```

#### 数据库指标

```promql
# 连接池使用
db_connections_active
db_connections_idle
db_connections_max

# 查询延迟
rate(db_query_duration_seconds_sum[5m])
```

#### Redis 指标

```promql
# 连接数
redis_connections_total

# 命令执行率
rate(redis_commands_total[5m])
```

#### 熔断器指标

```promql
# 熔断器状态
circuit_breaker_state{name="example"}[5m]

# 调用失败率
rate(circuit_breaker_errors_total[5m])
```

### 链路追踪

支持 OpenTelemetry 分布式追踪：

- **Trace Context**: 自动在请求间传递
- **Span 关联**: HTTP 请求、数据库查询、Redis 操作
- **Exporter**: OTLP (Jaeger, Zipkin, Tempo 等)

#### 查看追踪

```bash
# 使用 Jaeger UI
http://localhost:16686

# 使用 Zipkin UI
http://localhost:9411
```

### 日志聚合

结构化 JSON 日志输出，支持：

- **请求 ID 追踪**: 每个请求关联唯一 ID
- **Trace ID 关联**: 日志与链路追踪关联
- **多级别日志**: debug, info, warn, error, fatal

#### 日志查询

```bash
# 查看所有错误日志
tail -f logs/app.log | grep '"level":"error"'

# 查询特定请求
grep '"request_id":"abc123"' logs/app.log
```

## 🔐 安全

### 认证授权

- **JWT 认证**: HS256/HS512 签名支持
- **请求签名**: HMAC-SHA256 签名验证
- **RBAC 权限控制**: 角色基础访问控制

### 请求安全

#### 速率限制

- 基于令牌桶算法
- 可配置 QPS 和突发限制
- 支持 API Key、用户、IP 多维度限流

```yaml
rate_limit:
  enabled: true
  requests_per_second: 100
  burst: 200
```

#### 请求体大小限制

- 可配置的请求体大小限制（默认 8MB）
- Content-Length 预检查
- 防止内存耗尽攻击

```yaml
server:
  max_body_size: 8388608  # 8MB
```

### 数据安全

- **敏感信息加密**: 配置文件加密存储
- **安全 Headers**: CORS、Security Headers
- **输入验证**: 请求参数自动校验

### 生产环境安全

- **pprof 自动禁用**: 在生产模式（release）下自动禁用 pprof 性能分析工具
- **生产模式检测**: 基于 `app.mode` 配置自动判断环境并应用安全策略

### 安全最佳实践

1. 生产环境务必更改所有密钥
2. 启用 HTTPS
3. 配置合理的速率限制
4. 定期更新依赖
5. 启用审计日志

## 🧪 测试

### 运行测试

```bash
# 运行所有测试
go test ./...

# 运行测试并显示覆盖率
go test -cover ./...

# 运行特定模块测试
go test -v ./internal/config
go test -v ./internal/middleware
go test -v ./internal/versioning
go test -v ./internal/handler
```

### 测试覆盖率

```bash
# 查看详细覆盖率报告
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html

# 查看特定模块覆盖率
go test -coverprofile=coverage.out ./internal/middleware
go tool cover -func=coverage.out | grep middleware
```

### 编写测试

```go
func TestExampleHandler_GetExample(t *testing.T) {
    gin.SetMode(gin.TestMode)
    w := httptest.NewRecorder()
    c, _ := gin.CreateTestContext(w)
    
    c.Request, _ = http.NewRequest(http.MethodGet, "/test/123", nil)
    c.Params = []gin.Param{{Key: "id", Value: "123"}}
    
    handler := NewExampleHandler()
    handler.GetExample(c)
    
    if w.Code != http.StatusOK {
        t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
    }
}
```

### 测试模块

- **Config 测试**: [internal/config/config_test.go](internal/config/config_test.go)
- **Middleware 测试**: [internal/middleware/middleware_test.go](internal/middleware/middleware_test.go)
- **Versioning 测试**: [internal/versioning/versioning_test.go](internal/versioning/versioning_test.go)
- **Handler 测试**: [internal/handler/handler_test.go](internal/handler/handler_test.go)

## ⚡ 性能优化

### pprof 性能分析

内置 pprof 支持（仅在非生产环境可用）：

```bash
# CPU 分析
go tool pprof http://localhost:8080/debug/pprof/profile?seconds=30

# 内存分析
go tool pprof http://localhost:8080/debug/pprof/heap

# Goroutine 分析
go tool pprof http://localhost:8080/debug/pprof/goroutine

# 阻塞分析
go tool pprof http://localhost:8080/debug/pprof/block
```

⚠️ **注意**: 在生产模式（`app.mode = "release"`）下，pprof 会自动禁用。

### 性能优化建议

1. **数据库优化**
   - 合理配置连接池大小
   - 使用索引优化查询
   - 启用查询缓存

2. **缓存策略**
   - Redis 缓存热点数据
   - 多级缓存（本地 + 分布式）
   - 缓存失效策略

3. **限流和熔断**
   - 合理配置 QPS 和突发限制
   - 启用熔断器保护下游服务
   - 设置合理的超时时间

4. **监控和追踪**
   - 使用链路追踪定位性能瓶颈
   - 监控关键指标
   - 设置告警规则

5. **资源调优**
   - 监控 Goroutine 数量
   - 关注内存使用和 GC
   - 配置合理的连接池参数

## 📖 最佳实践

### 项目结构

```
├── cmd/              # 入口点，一个main包对应一个可执行程序
├── internal/         # 私有应用程序代码
│   └── {module}/     # 按功能模块组织
├── pkg/              # 公开的库代码，可被外部项目引用
├── api/              # API 协议定义文件
├── configs/          # 配置文件
└── scripts/          # 脚本文件
```

### 配置管理

1. 使用环境变量覆盖敏感配置
2. 生产环境禁用调试信息
3. 启用配置验证
4. 记录配置变更

### 错误处理

1. 使用统一的错误响应格式
2. 记录错误日志时包含上下文
3. 区分可恢复和不可恢复错误
4. 提供有意义的错误消息

### 日志规范

1. 使用结构化日志
2. 包含请求 ID 追踪
3. 设置合理的日志级别
4. 避免记录敏感信息

### 安全编码

1. 验证所有输入
2. 避免 SQL 注入
3. 加密敏感数据
4. 使用安全连接

## 🤝 贡献指南

欢迎提交 Issue 和 Pull Request！

### 提交规范

```bash
# 创建特性分支
git checkout -b feature/your-feature-name

# 提交更改
git commit -m 'feat: add new feature'

# 推送分支
git push origin feature/your-feature-name

# 创建 Pull Request
```

### 代码规范

- 遵循 Go 官方代码规范
- 添加单元测试
- 更新相关文档
- 确保所有测试通过

详细贡献指南: [CONTRIBUTING.md](CONTRIBUTING.md)

## 📄 许可证

本项目基于 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情

## 📞 联系方式

- 项目主页: https://github.com/lookingcamel/system-framework
- 问题反馈: https://github.com/lookingcamel/system-framework/issues
- 讨论组: https://github.com/lookingcamel/system-framework/discussions

## 🙏 致谢

- [Gin](https://github.com/gin-gonic/gin) - Web 框架
- [Zap](https://github.com/uber-go/zap) - 日志库
- [Prometheus](https://github.com/prometheus/client_golang) - 监控
- [OpenTelemetry](https://opentelemetry.io/) - 可观测性
- [Hystrix-go](https://github.com/afex/hystrix-go) - 熔断器
- [golang-migrate](https://github.com/golang-migrate/migrate) - 数据库迁移
- [Viper](https://github.com/spf13/viper) - 配置管理
- [Redis](https://github.com/go-redis/redis) - Redis 客户端
