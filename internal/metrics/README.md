# 指标监控模块

基于 Prometheus 的应用指标采集和监控，提供完整的可观测性。

## 功能特性

- 📊 **Prometheus 集成**: 标准 Prometheus 格式输出
- 🌐 **HTTP 指标**: 请求数、延迟、状态码分布
- 🔗 **gRPC 指标**: gRPC 调用监控
- 🗄️ **数据库指标**: 连接池、查询延迟
- 🔗 **Redis 指标**: 连接数、缓存命中率
- 🔥 **熔断器指标**: 断路器状态监控
- 🏗️ **自定义指标**: 支持自定义业务指标
- 📈 **Histogram/Summary**: 多维度指标统计

## 快速开始

### 初始化

```go
import "github.com/lookingcamel/system-framework/internal/metrics"

func main() {
    metrics.SetAppInfo("system-framework", "1.0.0")
}
```

### 注册指标端点

```go
import "github.com/prometheus/client_golang/prometheus/promhttp"

router.GET("/metrics", gin.WrapH(promhttp.Handler()))
```

## 内置指标

### HTTP 请求指标

| 指标名称 | 类型 | 说明 |
|----------|------|------|
| `http_requests_total` | Counter | HTTP 请求总数 |
| `http_request_duration_seconds` | Histogram | 请求延迟分布 |
| `http_requests_in_flight` | Gauge | 当前处理中的请求数 |
| `http_request_size_bytes` | Histogram | 请求体大小 |
| `http_response_size_bytes` | Histogram | 响应体大小 |

### 数据库指标

| 指标名称 | 类型 | 说明 |
|----------|------|------|
| `db_connections_open` | Gauge | 打开的连接数 |
| `db_connections_in_use` | Gauge | 使用中的连接数 |
| `db_connections_idle` | Gauge | 空闲连接数 |
| `db_query_duration_seconds` | Histogram | 查询延迟 |
| `db_query_total` | Counter | 查询总数 |
| `db_query_errors_total` | Counter | 查询错误数 |

### Redis 指标

| 指标名称 | 类型 | 说明 |
|----------|------|------|
| `redis_connections_open` | Gauge | 打开的连接数 |
| `redis_commands_total` | Counter | Redis 命令总数 |
| `redis_command_duration_seconds` | Histogram | 命令执行延迟 |
| `redis_errors_total` | Counter | Redis 错误总数 |

### 熔断器指标

| 指标名称 | 类型 | 说明 |
|----------|------|------|
| `hystrix_commands_total` | Counter | 熔断器命令调用总数 |
| `hystrix_circuit_breaker_state` | Gauge | 熔断器状态 |
| `hystrix_execution_duration_seconds` | Histogram | 命令执行延迟 |

## 自定义指标

### Counter 指标

```go
import "github.com/prometheus/client_golang/prometheus"

var (
    loginTotal = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "app_login_total",
            Help: "Total number of login attempts",
        },
        []string{"status", "method"},
    )

    orderTotal = prometheus.NewCounter(
        prometheus.CounterOpts{
            Name: "app_order_total",
            Help: "Total number of orders",
        },
    )
)

func init() {
    prometheus.MustRegister(loginTotal)
    prometheus.MustRegister(orderTotal)
}

// 使用指标
loginTotal.WithLabelValues("success", "password").Inc()
orderTotal.Inc()
```

### Gauge 指标

```go
var (
    activeUsers = prometheus.NewGauge(
        prometheus.GaugeOpts{
            Name: "app_active_users",
            Help: "Number of active users",
        },
    )

    queueLength = prometheus.NewGaugeVec(
        prometheus.GaugeOpts{
            Name: "app_queue_length",
            Help: "Length of processing queues",
        },
        []string{"queue_name"},
    )
)

func init() {
    prometheus.MustRegister(activeUsers)
    prometheus.MustRegister(queueLength)
}

// 使用指标
activeUsers.Set(123)
queueLength.WithLabelValues("email").Set(10)
```

### Histogram 指标

```go
var (
    requestDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "app_request_duration_seconds",
            Help:    "Request duration in seconds",
            Buckets: []float64{0.001, 0.005, 0.01, 0.05, 0.1, 0.5, 1, 5},
        },
        []string{"handler", "method"},
    )

    paymentAmount = prometheus.NewHistogram(
        prometheus.HistogramOpts{
            Name:    "app_payment_amount",
            Help:    "Payment amount distribution",
            Buckets: []float64{10, 50, 100, 500, 1000, 5000},
        },
    )
)

func init() {
    prometheus.MustRegister(requestDuration)
    prometheus.MustRegister(paymentAmount)
}

// 使用指标
requestDuration.WithLabelValues("/api/users", "GET").Observe(0.023)
paymentAmount.Observe(156.50)
```

### Summary 指标

```go
var (
    responseSize = prometheus.NewSummaryVec(
        prometheus.SummaryOpts{
            Name:       "app_response_size_bytes",
            Help:       "Response size in bytes",
            Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
        },
        []string{"handler"},
    )
)

func init() {
    prometheus.MustRegister(responseSize)
}

// 使用指标
responseSize.WithLabelValues("/api/users").Observe(1024)
```

## 业务指标示例

### 用户指标

```go
var (
    userLoginTotal = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "user_login_total",
            Help: "Total number of user logins",
        },
        []string{"status"},
    )

    userRegisterTotal = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "user_register_total",
            Help: "Total number of user registrations",
        },
        []string{"source"},
    )

    userActiveGauge = prometheus.NewGauge(
        prometheus.GaugeOpts{
            Name: "user_active_count",
            Help: "Number of active users",
        },
    )
)
```

### 订单指标

```go
var (
    orderCreatedTotal = prometheus.NewCounter(
        prometheus.CounterOpts{
            Name: "order_created_total",
            Help: "Total number of orders created",
        },
    )

    orderAmountHistogram = prometheus.NewHistogram(
        prometheus.HistogramOpts{
            Name:    "order_amount_histogram",
            Help:    "Order amount distribution",
            Buckets: []float64{10, 50, 100, 200, 500, 1000, 2000, 5000},
        },
    )

    orderProcessingGauge = prometheus.NewGauge(
        prometheus.GaugeOpts{
            Name: "order_processing_count",
            Help: "Number of orders being processed",
        },
    )
)
```

### 性能指标

```go
var (
    cacheHitTotal = prometheus.NewCounter(
        prometheus.CounterOpts{
            Name: "cache_hit_total",
            Help: "Total number of cache hits",
        },
    )

    cacheMissTotal = prometheus.NewCounter(
        prometheus.CounterOpts{
            Name: "cache_miss_total",
            Help: "Total number of cache misses",
        },
    )

    dbQueryDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "db_query_duration_seconds",
            Help:    "Database query duration",
            Buckets: []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1},
        },
        []string{"operation"},
    )
)
```

## Prometheus 配置

### Scrape 配置

```yaml
# prometheus.yml
scrape_configs:
  - job_name: 'system-framework'
    static_configs:
      - targets: ['localhost:8080']
    metrics_path: '/metrics'
    scrape_interval: 15s
```

### 指标查询

```promql
# 请求率
rate(http_requests_total[5m])

# 错误率
rate(http_requests_total{status="500"}[5m])

# P99 延迟
histogram_quantile(0.99, rate(http_request_duration_seconds_bucket[5m]))

# 活跃请求数
http_requests_in_flight
```

## Grafana 仪表盘

### JSON 仪表盘

```json
{
  "dashboard": {
    "title": "System Framework Dashboard",
    "panels": [
      {
        "title": "Request Rate",
        "targets": [
          {
            "expr": "rate(http_requests_total[5m])",
            "legendFormat": "{{handler}}"
          }
        ]
      },
      {
        "title": "Error Rate",
        "targets": [
          {
            "expr": "rate(http_requests_total{status=~\"5..\"}[5m])",
            "legendFormat": "{{handler}}"
          }
        ]
      },
      {
        "title": "P99 Latency",
        "targets": [
          {
            "expr": "histogram_quantile(0.99, rate(http_request_duration_seconds_bucket[5m]))",
            "legendFormat": "{{handler}}"
          }
        ]
      }
    ]
  }
}
```

## 告警规则

```yaml
# prometheus-rules.yml
groups:
  - name: system-framework
    rules:
      - alert: HighErrorRate
        expr: rate(http_requests_total{status=~"5.."}[5m]) > 0.01
        for: 5m
        labels:
          severity: critical
        annotations:
          summary: "High error rate detected"

      - alert: HighLatency
        expr: histogram_quantile(0.99, rate(http_request_duration_seconds_bucket[5m])) > 1
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "High latency detected"

      - alert: ServiceDown
        expr: up{job="system-framework"} == 0
        for: 1m
        labels:
          severity: critical
        annotations:
          summary: "Service is down"
```

## 最佳实践

1. **指标命名**
   - 使用有意义的名称
   - 统一命名规范
   - 添加单位后缀（如 _seconds, _bytes）

2. **标签设计**
   - 避免高基数标签
   - 使用标准标签名
   - 避免标签嵌套

3. **性能优化**
   - 使用异步指标收集
   - 合理设置 Bucket
   - 定期清理过期指标

4. **监控策略**
   - 设置合理的采集间隔
   - 配置告警阈值
   - 定期审查指标

## 相关文档

- [链路追踪模块](../tracing/README.md)
- [日志模块](../logger/README.md)
- [熔断器模块](../circuitbreaker/README.md)
- [项目总览](../../README.md)
