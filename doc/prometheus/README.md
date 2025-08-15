# Prometheus 全面实践教程（含 Docker/Kubernetes 部署与配置）

## 1. 什么是 Prometheus

Prometheus 是一个开源的时序数据库和监控系统，特别适合云原生和微服务环境，核心特性如下：

- **容器与节点监控**：收集 CPU、内存、网络等主机和容器资源指标。
- **业务指标监控**：通过集成客户端 SDK，主动暴露 `/metrics` 接口，Prometheus 拉取业务数据。
- **告警与可视化**：配合 Alertmanager 实现自动告警，结合 Grafana 构建仪表盘。
- **支持自动扩缩容发现**：在 Kubernetes、Docker Swarm 等环境中，实例动态变化自动发现，无需手动配置目标。

---

## 2. 在 Docker 中部署 Prometheus

```bash
docker run -d --name prometheus -p 9090:9090 \
  -v $(pwd)/prometheus.yml:/etc/prometheus/prometheus.yml \
  prom/prometheus:latest
```

> 🔧 将你的 `prometheus.yml` 放在当前目录下，挂载进容器以供加载。

---

## 3. 在 Kubernetes 中部署 Prometheus

### 3.1 创建 ConfigMap 配置

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: prometheus-config
  namespace: monitoring
  labels:
    app: prometheus
data:
  prometheus.yml: |
    global:
      scrape_interval: 15s
      evaluation_interval: 15s

    scrape_configs:
      - job_name: 'prometheus'
        static_configs:
          - targets: ['localhost:9090']

      - job_name: 'kubernetes-pods'
        kubernetes_sd_configs:
          - role: pod
        relabel_configs:
          - source_labels: [__meta_kubernetes_pod_annotation_prometheus_io_scrape] ##从Pod注解中读取 annotations: prometheus.io/scrape: "true"
            action: keep                                                           ## 保持不变
            regex: true
          - source_labels: [__meta_kubernetes_pod_annotation_prometheus_io_path]  ## 从Pod注解中读取 prometheus.io/path 的值  annotations: prometheus.io/path: "/custom-metrics"
            action: replace                                                       ## 替换
            target_label: __metrics_path__
            regex: (.+)
          - source_labels: [__address__, __meta_kubernetes_pod_annotation_prometheus_io_port] ## 提取原始 Pod IP (__address__) 和注解中的端口 (prometheus.io/port)
            action: replace                                                                   ## 替换
            regex: ([^:]+)(?::\d+)?;(\d+)                                                     ## 匹配 IP 地址部分,忽略掉IP后可能带的默认端口（例如 :80）。 从注解中提取端口
            replacement: $1:$2                                                                ## 生成新的地址
            target_label: __address__

      - job_name: 'docker-swarm-nodes'
        dockerswarm_sd_configs:
          - host: unix:///var/run/docker.sock
            role: nodes
        relabel_configs:
          - source_labels: [__meta_dockerswarm_node_address]
            target_label: instance

      - job_name: 'consul-services'
        consul_sd_configs:
          - server: 'consul:8500'
            services: []
        relabel_configs:
          - source_labels: [__meta_consul_tags]   ##tag中包含字符串 prometheus，就匹配成功。
            regex: .*,prometheus,.* 
            action: keep
          - source_labels: [__meta_consul_metadata_metrics_path]    ## 服务注册时加上Meta参数{"Meta": {"metrics_path": "/xxxx/metrics"}} ，抓取是会根据这个metrics_path
            action: replace
            target_label: __metrics_path__
```




### 3.2 创建 Deployment 和 Service

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: prometheus
  namespace: monitoring
spec:
  replicas: 1
  selector:
    matchLabels:
      app: prometheus
  template:
    metadata:
      labels:
        app: prometheus
    spec:
      containers:
        - name: prometheus
          image: prom/prometheus:latest
          args:
            - "--config.file=/etc/prometheus/prometheus.yml"
            - "--storage.tsdb.path=/prometheus"
          ports:
            - containerPort: 9090
          volumeMounts:
            - name: config
              mountPath: /etc/prometheus
            - name: data
              mountPath: /prometheus
      volumes:
        - name: config
          configMap:
            name: prometheus-config
        - name: data
          emptyDir: {}
---
apiVersion: v1
kind: Service
metadata:
  name: prometheus-service
  namespace: monitoring
spec:
  selector:
    app: prometheus
  ports:
    - name: http
      port: 9090
      targetPort: 9090
```

---

## 4. 服务发现与 relabel\_configs

Prometheus 可自动发现动态扩缩容实例，核心在于配置 `source_labels` 来提取元信息并重写为可抓取地址。

### Kubernetes 服务自动发现

```yaml
annotations:
  prometheus.io/scrape: "true"
  prometheus.io/port: "8080"
  prometheus.io/path: "/metrics"
```

### Docker Swarm 自动发现

```bash
--label prometheus-job=web
--label prometheus-port=8080
```

### Consul 注册服务方式

```json
{
  "name": "web-service",
  "tags": ["prometheus", "v1"],
  "address": "10.0.0.1",
  "port": 8080
}
```

---

## 5. Prometheus 指标类型

### Counter（只增不减）

- 示例：请求总数、错误数
- 特点：只能递增，适用于速率、累计分析
- PromQL:

```promql
http_requests_total{method="POST", handler="/login"} 
# 查看 POST /login 请求总数

rate(http_requests_total[5m])                         
# 每秒请求速率（增长速率）

rate(http_requests_total{method="POST"}[5m])          
# 每秒 POST 请求速率（带条件）

sum by (handler) (rate(http_requests_total[5m]))      
# 请求速率按 handler 分组统计
```

---

### Gauge（可增可减）

- 示例：当前连接数、温度、内存使用率
- 特点：表示当前状态，可增可减
- PromQL:

```promql
avg_over_time(cpu_usage[5m])                
# CPU 5 分钟平均使用率

max_over_time(go_memstats_alloc_bytes[5m])  
# 5 分钟内内存最大使用量
```

---

### Histogram（桶统计，适合延迟分析）

- 示例：请求延迟、响应时间分布
- 特点：
    - 写入时通过 `Observe(value)` 将值分类到桶中
    - 自动生成以下三个指标：
        - `your_histogram_bucket{le="0.5"}`：小于等于该值的请求数
        - `your_histogram_sum`：总耗时
        - `your_histogram_count`：总请求数
- PromQL:

```promql
histogram_quantile(0.9, rate(http_request_duration_seconds_bucket[5m]))  
# 过去 5 分钟估算 P90 响应时间

histogram_quantile(0.99, rate(http_request_duration_seconds_bucket[5m])) 
# 过去 5 分钟估算 P99 响应时间
```

---

### Summary（百分位统计）

- 示例：客户端估算的请求延迟（不适用于跨实例聚合）
- 特点：
    - 使用 `Observe(value)` 上报
    - 自动生成：
        - `your_summary{quantile="0.99"}`：客户端 P99
        - `your_summary_sum`：总值
        - `your_summary_count`：次数
    - 示例定义方式：

```
Objectives: map[float64]float64{
  0.5:  0.05,  # P50，误差容忍 ±5%
  0.9:  0.01,  # P90，误差容忍 ±1%
}
```

- PromQL:

```promql
http_request_duration_seconds{quantile="0.99"}  
# 客户端上报的 P99

http_request_duration_seconds_sum / http_request_duration_seconds_count  
# 平均响应时间
```

---

### Alertmanager 报警规则示例

将 PromQL 表达式封装为报警规则：

```yaml
groups:
- name: example.rules
  rules:
  - alert: HighCpuUsage
    expr: 100 - (avg by(instance)(rate(node_cpu_seconds_total{mode="idle"}[5m])) * 100) > 90
    for: 5m
    labels:
      severity: warning
    annotations:
      summary: "CPU usage high on {{ $labels.instance }}"
      description: "CPU 使用率高于 90%，已持续 5 分钟。"
```

- `expr`：报警表达式
- `for`：持续触发多久才报警
- `labels`：报警等级标签
- `annotations`：报警内容，支持模板变量

---
---

## 6. 集群与微服务监控实践

- 微服务应引入 Prometheus SDK 暴露 `/metrics` 接口。
- 每台实例会通过 Kubernetes 注解或 Docker Swarm label 自动被 Prometheus 发现。
- 实例扩容时，无需修改 Prometheus 配置。

> ✅ 你可以通过 `kubectl get endpoints` 验证哪些服务实例已经被 Prometheus 抓取。

---

## 7. 常用 Exporter 示例（Docker 环境）

```bash
# 宿主机资源监控
$ docker run -d --name node-exporter -p 9100:9100 \
  -v /proc:/host/proc:ro -v /sys:/host/sys:ro \
  prom/node-exporter --path.procfs /host/proc --path.sysfs /host/sys

# 容器级别监控
$ docker run -d --name cadvisor -p 8080:8080 \
  -v /:/rootfs:ro -v /var/run:/var/run:ro \
  -v /sys:/sys:ro -v /var/lib/docker:/var/lib/docker:ro \
  google/cadvisor
```

---

## 8. Prometheus Go 客户端集成

```bash
go get github.com/prometheus/client_golang
```

在业务中暴露 `/metrics`：

```
http.Handle("/metrics", promhttp.Handler())
http.ListenAndServe(":8080", nil)
```

---

> 📘 本文详细讲解了 Prometheus 的核心概念、部署方式、服务发现机制与业务接入实践，适用于 Kubernetes 与 Docker 混合环境下的微服务监控体系。

