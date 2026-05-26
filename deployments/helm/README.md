# Kubernetes Helm 部署

完整的 Helm Chart 部署配置，支持金丝雀发布、蓝绿部署和滚动更新。

## 目录结构

```
helm/
├── Chart.yaml              # Chart 定义
├── values.yaml            # 默认配置
├── values-prod.yaml      # 生产环境配置
├── values-staging.yaml    # 测试环境配置
├── templates/
│   ├── _helpers.tpl       # 辅助模板
│   ├── deployment.yaml    # Deployment
│   ├── service.yaml       # Service
│   ├── ingress.yaml       # Ingress
│   ├── ingress-canary.yaml  # 金丝雀 Ingress
│   ├── hpa.yaml           # HPA 配置
│   ├── pdb.yaml          # Pod 中断预算
│   ├── configmap.yaml     # 配置 ConfigMap
│   ├── config.yaml        # 应用配置
│   └── tests/
│       └── test-connection.yaml  # 连接测试
└── README.md              # 本文档
```

## 快速开始

### 安装

```bash
# 添加 Helm 仓库
helm repo add lookingcamel https://charts.example.com
helm repo update

# 安装
helm install system-framework ./deployments/helm

# 查看状态
helm status system-framework
```

### 升级

```bash
# 升级到新版本
helm upgrade system-framework ./deployments/helm \
  --set image.tag=v1.1.0

# 使用 values 文件升级
helm upgrade system-framework ./deployments/helm \
  -f values-prod.yaml
```

### 回滚

```bash
# 查看发布历史
helm history system-framework

# 回滚到上一个版本
helm rollback system-framework

# 回滚到指定版本
helm rollback system-framework 3
```

## 配置选项

### 基础配置

```yaml
replicaCount: 3

image:
  repository: lookingcamel/system-framework
  pullPolicy: IfNotPresent
  tag: "latest"

imagePullSecrets: []
nameOverride: ""
fullnameOverride: ""
```

### 服务配置

```yaml
service:
  type: ClusterIP      # ClusterIP / NodePort / LoadBalancer
  port: 8080
  targetPort: 8080
```

### 资源限制

```yaml
resources:
  limits:
    cpu: 500m
    memory: 512Mi
  requests:
    cpu: 100m
    memory: 128Mi
```

### 健康检查

```yaml
livenessProbe:
  httpGet:
    path: /health
    port: 8080
  initialDelaySeconds: 10
  periodSeconds: 30

readinessProbe:
  httpGet:
    path: /ready
    port: 8080
  initialDelaySeconds: 5
  periodSeconds: 10
```

## 部署策略

### 滚动更新

```yaml
updateStrategy:
  type: RollingUpdate
  rollingUpdate:
    maxSurge: 25%        # 最多增加 25% Pod
    maxUnavailable: 0%   # 不允许不可用 Pod
```

### 蓝绿部署

```bash
# 部署绿色版本
helm install system-framework-green ./deployments/helm \
  --set image.tag=v2.0.0

# 切换流量到绿色版本
kubectl patch service system-framework \
  -p '{"spec":{"selector":{"version":"green"}}}'

# 验证后删除蓝色版本
helm uninstall system-framework-blue
```

## 金丝雀发布

### 启用金丝雀

```yaml
traffic:
  canary:
    enabled: true
    weight: 10              # 10% 流量到金丝雀
```

### 按 Header 路由

```yaml
traffic:
  canary:
    enabled: true
    header: "X-Canary"
    headerValue: "always"
```

### 按 Cookie 路由

```yaml
traffic:
  canary:
    enabled: true
    cookie: "canary"
```

### 金丝雀发布流程

```bash
# 1. 部署金丝雀版本
helm upgrade system-framework ./deployments/helm \
  --set traffic.canary.enabled=true \
  --set traffic.canary.weight=10 \
  --set image.tag=canary

# 2. 监控金丝雀指标
kubectl logs -l app.kubernetes.io/canary=true

# 3. 增加流量
helm upgrade system-framework ./deployments/helm \
  --set traffic.canary.weight=50

# 4. 全量发布
helm upgrade system-framework ./deployments/helm \
  --set traffic.canary.enabled=false \
  --set image.tag=canary
```

## 流量策略

### 会话保持

```yaml
traffic:
  sessionAffinity:
    enabled: true
    type: ClientIP
    sessionTimeoutSeconds: 10800
```

### 速率限制

```yaml
traffic:
  rateLimit:
    enabled: true
    requestsPerSecond: 1000
    burst: 2000
```

## 高可用配置

### Pod 中断预算

```yaml
podDisruptionBudget:
  enabled: true
  minAvailable: 1
  # maxUnavailable: 1
```

### HPA 自动扩缩容

```yaml
autoscaling:
  enabled: true
  minReplicas: 2
  maxReplicas: 10
  targetCPUUtilizationPercentage: 80
  # targetMemoryUtilizationPercentage: 80
```

## 配置管理

### ConfigMap 挂载

```yaml
configmap:
  enabled: true

# 挂载到容器
volumeMounts:
  - name: config
    mountPath: /config

# 添加环境变量
env:
  - name: CONFIG_PATH
    value: /config/config.yaml
```

### 敏感信息

```yaml
# 使用 Kubernetes Secret
secretName: system-framework-secrets

env:
  - name: DB_PASSWORD
    valueFrom:
      secretKeyRef:
        name: system-framework-secrets
        key: db-password
```

## 多环境配置

### 开发环境

```yaml
# values-dev.yaml
replicaCount: 1
image:
  tag: "dev"

resources:
  limits:
    cpu: 200m
    memory: 256Mi

autoscaling:
  enabled: false
```

### 生产环境

```yaml
# values-prod.yaml
replicaCount: 3
image:
  tag: "latest"

resources:
  limits:
    cpu: 1000m
    memory: 1Gi

autoscaling:
  enabled: true
  minReplicas: 3
  maxReplicas: 20

podDisruptionBudget:
  enabled: true
  minAvailable: 2
```

### 部署命令

```bash
# 开发环境
helm install system-framework ./deployments/helm -f values-dev.yaml

# 测试环境
helm install system-framework ./deployments/helm -f values-staging.yaml

# 生产环境
helm install system-framework ./deployments/helm -f values-prod.yaml
```

## Ingress 配置

```yaml
ingress:
  enabled: true
  className: nginx
  annotations:
    nginx.ingress.kubernetes.io/rewrite-target: /
  hosts:
    - host: system.example.com
      paths:
        - path: /
          pathType: Prefix
  tls:
    - secretName: system-framework-tls
      hosts:
        - system.example.com
```

## 服务网格集成

### Istio 虚拟服务

```yaml
apiVersion: networking.istio.io/v1beta1
kind: VirtualService
metadata:
  name: system-framework
spec:
  hosts:
    - system.example.com
  http:
    - route:
        - destination:
            host: system-framework
            subset: stable
          weight: 90
        - destination:
            host: system-framework
            subset: canary
          weight: 10
```

### Linkerd 加权路由

```yaml
apiVersion: split.smi-spec.io/v1alpha1
kind: TrafficSplit
metadata:
  name: system-framework
spec:
  service: system-framework
  backends:
    - service: system-framework-stable
      weight: 90
    - service: system-framework-canary
      weight: 10
```

## 监控集成

### Pod 注解

```yaml
podAnnotations:
  prometheus.io/scrape: "true"
  prometheus.io/port: "8080"
  prometheus.io/path: "/metrics"
```

### ServiceMonitor

```yaml
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  name: system-framework
spec:
  selector:
    matchLabels:
      app: system-framework
  endpoints:
    - port: http
      path: /metrics
```

## 最佳实践

1. **资源限制**
   - 设置合理的 CPU 和内存限制
   - 使用 requests 确保最小资源
   - 根据实际负载调整

2. **健康检查**
   - 配置存活和就绪探针
   - 设置合理的超时时间
   - 监控探针失败

3. **滚动更新**
   - 设置 maxUnavailable 为 0
   - 使用 maxSurge 控制更新速度
   - 监控滚动更新过程

4. **安全配置**
   - 使用 ServiceAccount
   - 配置 SecurityContext
   - 限制 Pod 特权

## 故障排查

### 查看 Pod 状态

```bash
kubectl get pods -l app=system-framework
kubectl describe pod system-framework-xxx
kubectl logs system-framework-xxx
```

### 查看 Helm 状态

```bash
helm status system-framework
helm history system-framework
helm get values system-framework
```

### 调试技巧

```bash
# 进入容器调试
kubectl exec -it system-framework-xxx -- /bin/sh

# 转发端口到本地
kubectl port-forward system-framework-xxx 8080:8080

# 查看资源使用
kubectl top pods
kubectl top nodes
```

## 相关文档

- [Helm 官方文档](https://helm.sh/docs/)
- [Kubernetes 部署](https://kubernetes.io/docs/concepts/workloads/controllers/deployment/)
- [Ingress 配置](https://kubernetes.io/docs/concepts/services-networking/ingress/)
- [项目总览](../../README.md)
