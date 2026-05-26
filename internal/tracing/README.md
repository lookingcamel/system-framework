# 链路追踪模块

基于 OpenTelemetry 的分布式追踪系统，支持 OTLP、Jaeger、Zipkin 等多种后端。

## 功能特性

- 🔍 **分布式追踪**: 完整请求链路追踪
- 🌐 **多 Exporter**: OTLP、Stdout 等
- 📊 **Span 管理**: 自动创建和管理 Span
- 🔗 **Context 传递**: Trace Context 跨服务传递
- 🎯 **采样控制**: 灵活的采样策略
- 📈 **性能监控**: 追踪性能指标

## 快速开始

### 初始化

```go
import "github.com/lookingcamel/system-framework/internal/tracing"

func main() {
    shutdown, err := tracing.Init(cfg.Tracing)
    if err != nil {
        log.Fatal(err)
    }
    defer shutdown(context.Background())
}
```

### 创建 Span

```go
import "github.com/lookingcamel/system-framework/internal/tracing"

// 创建新的 Span
ctx, span := tracing.StartSpan(c.Request.Context(), "OperationName")
defer span.End()

// 添加属性
span.SetAttributes(
    attribute.String("user_id", "123"),
    attribute.Int("order_id", 456),
)

// 添加事件
span.AddEvent("Processing started")

// 记录错误
span.RecordError(err)

// 完成子操作
ctx2, childSpan := tracing.StartSpan(ctx, "ChildOperation")
defer childSpan.End()
```

## 配置说明

### OTLP Exporter

```yaml
tracing:
  enabled: true
  service_name: "system-framework"
  exporter: "otlp"               # otlp, jaeger, stdout
  endpoint: "localhost:4317"    # OTLP gRPC endpoint
  sample_rate: 1.0              # 采样率 (0-1)
```

### Jaeger Exporter

```yaml
tracing:
  enabled: true
  service_name: "system-framework"
  exporter: "jaeger"
  endpoint: "http://localhost:14268/api/traces"
  sample_rate: 1.0
```

### Stdout Exporter

```yaml
tracing:
  enabled: true
  service_name: "system-framework"
  exporter: "stdout"
  sample_rate: 1.0
```

## 工作原理

### 追踪模型

```
Trace
├── Span 1 (Root)
│   ├── Span 2 (Child)
│   │   ├── Span 4 (Child)
│   │   └── Span 5 (Child)
│   └── Span 3 (Child)
│       └── Span 6 (Child)
```

### Context 传递

```
HTTP 请求 (Trace ID: abc123)
    ↓
中间件创建 Root Span
    ↓
Handler 处理
    ↓
Service 调用 (注入 Context)
    ↓
Database 查询
    ↓
返回响应
```

## 使用示例

### 1. HTTP 请求追踪

```go
import "github.com/lookingcamel/system-framework/internal/tracing"

func MyHandler(c *gin.Context) {
    ctx := c.Request.Context()

    ctx, span := tracing.StartSpan(ctx, "MyHandler")
    defer span.End()

    span.SetAttributes(
        attribute.String("http.method", c.Request.Method),
        attribute.String("http.url", c.Request.URL.Path),
    )

    // 处理请求
    result, err := processRequest(ctx)

    if err != nil {
        span.RecordError(err)
        span.SetAttributes(attribute.Int("http.status_code", 500))
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }

    span.SetAttributes(attribute.Int("http.status_code", 200))
    c.JSON(200, result)
}
```

### 2. 数据库查询追踪

```go
func queryUser(ctx context.Context, userID int) (*User, error) {
    ctx, span := tracing.StartSpan(ctx, "Database.QueryUser")
    defer span.End()

    span.SetAttributes(
        attribute.String("db.system", "mysql"),
        attribute.String("db.operation", "SELECT"),
        attribute.String("db.table", "users"),
    )

    user, err := db.QueryRowContext(ctx, "SELECT * FROM users WHERE id = ?", userID)
    if err != nil {
        span.RecordError(err)
        return nil, err
    }

    return user, nil
}
```

### 3. 外部 API 调用追踪

```go
import (
    "go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func callExternalAPI(ctx context.Context) error {
    ctx, span := tracing.StartSpan(ctx, "ExternalAPI.Call")
    defer span.End()

    client := &http.Client{
        Transport: otelhttp.NewTransport(http.DefaultTransport),
    }

    req, _ := http.NewRequestWithContext(ctx, "GET", "https://api.example.com", nil)
    resp, err := client.Do(req)

    if err != nil {
        span.RecordError(err)
        return err
    }
    defer resp.Body.Close()

    span.SetAttributes(
        attribute.Int("http.status_code", resp.StatusCode),
    )

    return nil
}
```

### 4. Context 注入和提取

```go
import "go.opentelemetry.io/otel/propagation"

// 注入 Context 到 HTTP Header
func InjectTraceID(c *gin.Context) {
    ctx := c.Request.Context()
    carrier := propagation.HeaderCarrier(c.Request.Header)
    tracing.Inject(ctx, carrier)
}

// 提取 Context 从 HTTP Header
func ExtractTraceID(c *gin.Context) context.Context {
    carrier := propagation.HeaderCarrier(c.Request.Header)
    return tracing.Extract(c.Request.Context(), carrier)
}
```

## 采样策略

### 1. 全量采样

```yaml
tracing:
  sample_rate: 1.0  # 100% 采样
```

### 2. 概率采样

```yaml
tracing:
  sample_rate: 0.1  # 10% 采样
```

### 3. 父级采样

自动继承父 Trace 的采样决策。

## 与日志集成

### 结构化日志

```go
import "go.opentelemetry.io/otel/trace"

func logWithTrace(ctx context.Context, msg string) {
    span := trace.SpanFromContext(ctx)
    if span.SpanContext().HasTraceID() {
        logger.Info(msg,
            zap.String("trace_id", span.SpanContext().TraceID().String()),
            zap.String("span_id", span.SpanContext().SpanID().String()),
        )
    } else {
        logger.Info(msg)
    }
}
```

## 追踪后端配置

### Jaeger

```bash
# 启动 Jaeger
docker run -d --name jaeger \
  -e COLLECTOR_ZIPKIN_HOST_PORT=:9411 \
  -p 5775:5775/udp \
  -p 6831:6831/udp \
  -p 6832:6832/udp \
  -p 5778:5778 \
  -p 16686:16686 \
  -p 14268:14268 \
  -p 14250:14250 \
  -p 9411:9411 \
  jaegertracing/all-in-one:1.48
```

### OpenTelemetry Collector

```yaml
# otel-collector-config.yaml
receivers:
  otlp:
    protocols:
      grpc:
        endpoint: 0.0.0.0:4317
      http:
        endpoint: 0.0.0.0:4318

exporters:
  jaeger:
    endpoint: jaeger:14250
    tls:
      insecure: true

  prometheus:
    endpoint: "0.0.0.0:8889"

processors:
  batch:
    timeout: 10s

service:
  pipelines:
    traces:
      receivers: [otlp]
      exporters: [jaeger]
      processors: [batch]
    metrics:
      receivers: [otlp]
      exporters: [prometheus]
      processors: [batch]
```

## 最佳实践

1. **命名规范**
   - 操作名称清晰
   - 使用小写字母
   - 使用点号分隔层级

2. **属性设计**
   - 添加关键业务属性
   - 避免添加过多属性
   - 使用标准属性名称

3. **错误处理**
   - 记录所有错误
   - 使用 RecordError
   - 添加错误类型

4. **性能优化**
   - 合理设置采样率
   - 使用异步导出
   - 避免追踪高频操作

## 故障排查

### 常见问题

**Q: 追踪数据未发送**
```
A: 检查 Exporter 配置，确认网络连通性
```

**Q: Trace ID 不连续**
```
A: 确认所有服务使用相同 TracerProvider
```

**Q: Span 数据丢失**
```
A: 检查 BatchSpanProcessor 缓冲区，使用 otelcol 缓冲
```

## 相关文档

- [日志模块](../logger/README.md)
- [中间件模块](../middleware/README.md)
- [指标模块](../metrics/README.md)
- [项目总览](../../README.md)
