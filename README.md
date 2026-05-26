# System Framework

基于 Go + Gin 的现代化云原生微服务框架。

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
- [性能优化](#-性能优化)
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
- 🔐 **安全认证**: JWT + 请求签名验证
- 🗄️ **数据库迁移**: golang-migrate 支持版本化迁移
- 💾 **自动备份**: 数据库定时备份和恢复
- 🐳 **容器化部署**: Docker 和 Docker Compose 支持
- ☸️ **Kubernetes 部署**: Helm Chart 完整支持
- 🔄 **灰度发布**: 金丝雀发布和流量策略
- 📝 **结构化日志**: Zap 日志系统，支持 JSON 格式
- 🏥 **健康检查**: 就绪探针和存活探针
- 🎨 **优雅关闭**: 支持零停机部署
- 🔧 **性能分析**: pprof 性能分析工具

## 🚀 快速开始

### 前置要求

- Go 1.21+
- MySQL 8.0+ / PostgreSQL 14+ / SQLite 3
- Redis 6.0+ (可选)

### 安装

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

### Docker 部署

```bash
# 使用 Docker Compose 启动
docker-compose up -d

# 查看日志
docker-compose logs -f
```

### Kubernetes 部署

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
│   ├── config/                  # 配置管理
│   ├── database/                # 数据库连接
│   ├── handler/                 # HTTP 处理器
│   ├── logger/                  # 日志系统
│   ├── metrics/                 # 指标采集
│   ├── middleware/              # 中间件
│   ├── migrate/                 # 数据库迁移
│   ├── redis/                   # Redis 缓存
│   ├── server/                  # HTTP 服务器
│   └── tracing/                 # 链路追踪
├── pkg/                         # 公共工具包
│   └── utils/                   # 工具函数
├── migrations/                  # 数据库迁移文件
├── deployments/                 # 部署配置
│   ├── docker/                  # Docker 配置
│   └── helm/                    # Helm Chart
├── config.yaml                 # 配置文件
├── docker-compose.yml           # Docker Compose 配置
├── Dockerfile                  # Docker 镜像构建
└── README.md                    # 项目说明
```

## 🔧 核心模块

### 认证模块 (auth)

JWT 认证和请求签名验证，支持 RBAC 权限控制。

详细文档: [internal/auth/README.md](internal/auth/README.md)

### 配置模块 (config)

支持 YAML 配置、环境变量和 Nacos 配置中心，运行时热更新。

详细文档: [internal/config/README.md](internal/config/README.md)

### 日志模块 (logger)

基于 Zap 的结构化日志，支持多种输出格式和日志轮转。

详细文档: [internal/logger/README.md](internal/logger/README.md)

### 数据库模块 (database)

支持 MySQL、PostgreSQL、SQLite，自动迁移和备份。

详细文档: [internal/database/README.md](internal/database/README.md)

### Redis 模块 (redis)

连接池管理，自动重连，分布式缓存支持。

详细文档: [internal/redis/README.md](internal/redis/README.md)

### 熔断器模块 (circuitbreaker)

基于 Hystrix-go 的熔断降级，防止级联故障。

详细文档: [internal/circuitbreaker/README.md](internal/circuitbreaker/README.md)

### 链路追踪模块 (tracing)

OpenTelemetry + OTLP 分布式追踪。

详细文档: [internal/tracing/README.md](internal/tracing/README.md)

### 指标模块 (metrics)

Prometheus 指标采集和暴露。

详细文档: [internal/metrics/README.md](internal/metrics/README.md)

### 中间件模块 (middleware)

请求日志、追踪、认证、重试等中间件。

详细文档: [internal/middleware/README.md](internal/middleware/README.md)

### 数据库迁移 (migrate)

版本化管理数据库 schema 变更。

详细文档: [internal/migrate/README.md](internal/migrate/README.md)

### 数据库备份 (backup)

自动备份和手动恢复功能。

详细文档: [internal/backup/README.md](internal/backup/README.md)

## 🚢 部署指南

### Docker 部署

```yaml
# docker-compose.yml
services:
  app:
    build: .
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

  mysql:
    image: mysql:8.0
    environment:
      MYSQL_ROOT_PASSWORD: password
      MYSQL_DATABASE: app

  redis:
    image: redis:7-alpine
```

### Kubernetes 部署

#### 安装 Helm Chart

```bash
# 安装
helm install system-framework ./deployments/helm \
  --set image.tag=v1.0.0

# 查看部署状态
helm status system-framework

# 升级
helm upgrade system-framework ./deployments/helm \
  --set image.tag=v1.1.0

# 回滚
helm rollback system-framework
```

#### 金丝雀发布

```bash
# 部署金丝雀版本（10% 流量）
helm upgrade system-framework ./deployments/helm \
  --set traffic.canary.enabled=true \
  --set traffic.canary.weight=10

# 全量切换
helm upgrade system-framework ./deployments/helm \
  --set traffic.canary.enabled=false
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
  mode: "release"

logging:
  level: "info"
  format: "json"
  output: "console"
  output_path: "logs/app.log"
  max_size: 100
  max_backups: 7
  max_age: 30
  compress: true

database:
  type: "sqlite"
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

auth:
  enabled: true
  jwt_enabled: true
  jwt_secret: "your-256-bit-secret"
  jwt_expire_seconds: 86400
  signature_enabled: false

circuit_breaker:
  enabled: true
  default_timeout: 3000
  default_max_concurrent: 100
  default_error_percentage: 50

tracing:
  enabled: true
  service_name: "system-framework"
  exporter: "otlp"
  endpoint: "localhost:4317"
  sample_rate: 1.0
```

详细配置文档: [internal/config/README.md](internal/config/README.md)

## 📚 API 文档

### 健康检查

```
GET /health              # 存活探针
GET /ready               # 就绪探针
```

### Prometheus 指标

```
GET /metrics             # Prometheus 指标端点
```

### 熔断器状态

```
GET /circuit/status      # 熔断器状态查询
```

### 示例 API

```
GET /api/v1/example/:id      # 获取单个示例
GET /api/v1/examples         # 列出所有示例
```

详细 API 文档: [internal/handler/README.md](internal/handler/README.md)

## 📊 监控与可观测性

### Prometheus 监控

项目集成完整的 Prometheus 指标：

- **HTTP 请求指标**: 请求数、延迟、状态码分布
- **数据库指标**: 连接池状态、查询延迟
- **Redis 指标**: 连接数、缓存命中率
- **熔断器指标**: 断路器状态、调用次数、降级次数
- **自定义业务指标**: 可扩展的指标采集

### 链路追踪

支持 OpenTelemetry 分布式追踪：

- **Trace Context**: 自动在请求间传递
- **Span 关联**: HTTP 请求、数据库查询、Redis 操作
- **Exporter**: OTLP (Jaeger, Zipkin, Tempo 等)

### 日志聚合

结构化 JSON 日志输出，支持：

- **请求 ID 追踪**: 每个请求关联唯一 ID
- **Trace ID 关联**: 日志与链路追踪关联
- **多级别日志**: debug, info, warn, error, fatal

## 🔐 安全

### 认证授权

- **JWT 认证**: HS256/HS512 签名支持
- **请求签名**: HMAC-SHA256 签名验证
- **RBAC 权限控制**: 角色基础访问控制

### 数据安全

- **敏感信息加密**: 配置文件加密存储
- **安全 Headers**: CORS、Security Headers
- **输入验证**: 请求参数自动校验

## ⚡ 性能优化

### pprof 性能分析

内置 pprof 支持：

```bash
# CPU 分析
go tool pprof http://localhost:8080/debug/pprof/profile

# 内存分析
go tool pprof http://localhost:8080/debug/pprof/heap

# Goroutine 分析
go tool pprof http://localhost:8080/debug/pprof/goroutine
```

### 性能优化建议

- 合理配置数据库连接池
- 使用 Redis 缓存热点数据
- 启用熔断器保护下游服务
- 使用链路追踪定位性能瓶颈
- 监控 Goroutine 数量和内存使用

## 🤝 贡献指南

欢迎提交 Issue 和 Pull Request！

1. Fork 本仓库
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 创建 Pull Request

## 📄 许可证

本项目基于 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情

## 📞 联系方式

- 项目主页: https://github.com/lookingcamel/system-framework
- 问题反馈: https://github.com/lookingcamel/system-framework/issues

## 🙏 致谢

- [Gin](https://github.com/gin-gonic/gin) - Web 框架
- [Zap](https://github.com/uber-go/zap) - 日志库
- [Prometheus](https://github.com/prometheus/client_golang) - 监控
- [OpenTelemetry](https://opentelemetry.io/) - 可观测性
- [Hystrix-go](https://github.com/afex/hystrix-go) - 熔断器
