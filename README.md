# 透明多级分流

### 分流组件
1. **本地缓存、内容分发网络、反向代理**：利用本地缓存、CDN和反向代理，实现请求的初步分流，减轻后端服务压力。
2. **可自动扩缩的服务节点**：使用较小的代价堆叠机器，构建能够自动扩缩的服务节点，以应对不同的流量需求。
3. **传统关系数据库**：作为调用链末端的数据存储，处理核心业务数据。

### 高可用设计目标
尽可能减少单点部件，如果某些单点是无可避免的，则应尽最大限度减少到达单点部件的流量。引导请求分流至最合适的组件，降低单点的压力。具体措施如下：
1. **客户端缓存**：采用空间换时间的概念，把部分数据缓存在客户端本地，无需每次都调用获取，加快速度并减少服务器压力。
2. **DNS**：可通过不同运营商、不同区域获取最近的节点，提高响应速度，避免单点服务和负载均衡的压力。
3. **CDN**：对于静态资源，如图片、CSS、JS等，可减少对服务的获取，利用缓存机制提高响应速度。
4. **负载均衡**：将请求流量指向不同的服务器集群，实现负载均衡。可分为四层（根据IP等报文即可）及七层（读取所有的数据后再转发）。
    - **避免单点故障**：引入`KeepAlived`，防止负载均衡成为单点故障。
    - **调度算法**：包括轮训、加权轮训、最少连接、加权最少连接、IP hash等。
    - **健康检测**：为避免把请求流量指向已经不可用的服务器，需要进行`健康检测`。
5. **服务器集群**：无状态服务可快速扩展，避免单台服务的计算压力。
6. **参数校验**：检验参数的合理性，减少多余的查询和处理。
7. **缓存**：由于数据库QPS较低，将频繁请求的数据缓存到内存中。如有全局概念，可使用分布式缓存。
8. **NoSQL非事务型数据库**：对非事务要求的操作，存储和读取到非事务型数据库。
9. **异步队列/事务数据库**：让一些实时性不高的操作进行异步处理。
10. **同步事务数据库**：最终将数据落到数据库。

---

# Web服务安全与部署

### 安全策略
所有的服务/中间件都应该在服务器集群中且不可暴露，仅提供少量的端口对外暴露（如网关入口）。这样可以保证服务/中间件的安全，不被恶意攻击，同时合理有效地控制内部人员的使用权限。

### 部署方案
1. **公/私有云集群**：利用公/私有云的资源，构建高可用的服务集群。
2. **docker network**：使用Docker网络，实现容器之间的通信和隔离。
3. **k8s**：借助Kubernetes进行容器编排和管理，提高服务的可扩展性和可靠性。

---

# 运行场景

### 本地调试
在本地运行时，一般不会把所有的服务启动运行，通常是针对少数的服务进行跨服务的调用。使用gRPC通过本地IP:port获取信息进行调用，无需注册中心。

### 对外服务（k8s）
在Kubernetes环境中，每个服务都应该创建对应的server进行对外访问。gRPC可使用k8s DNS的机制（server - name:port）进行调用，无需注册中心。

### 对外服务（注册中心）
程序启动时，需要把服务注册到注册中心，提供访问。可采用以下两种方式：
1. **DNS形式**：通过注册中心DNS进行地址解析获取服务地址。
2. **HTTP形式**：通过注册中心获取服务的对应列表，进行细粒度的地址获取。

### 常用暴露端口
| 服务类型 | 服务名称 | 端口号 | 用途 |
| ---- | ---- | ---- | ---- |
| 服务注册与发现 | consul | 8500 | 观察服务注册发现相关的信息 |
| 服务网关 | kong | 8000 | 通过网关入口访问后端web服务 |
| 服务网关 | konga | 1337 | 通过konga配置服务（内部调用kong admin端口8001） |
| 配置中心 | etcd | 2379 | 存储和管理配置信息 |
| SQL/NoSQL | mysql | 3306 | 关系型数据库 |
| SQL/NoSQL | redis | 6379 | 非关系型数据库 |
| 消息队列 | kafka | 9092 | 消息队列服务 |
| 消息队列 | rocketmq | 9876 | 消息队列服务 |
| 服务检测 | openTelemetry | 4317 | 服务检测和追踪 |
| 服务检测 | jaeger webUI | 16686 | 可视化服务调用链 |

### 固定IP
| 服务名称 | IP地址        |
| ---- |-------------|
| consul | 172.40.0.11 |
| etcd | 172.40.0.12 |
| kong - database | 172.40.0.21 |
| kong | 172.40.0.22 |
| konga | 172.40.0.23 |
| mysql | 172.40.0.31 |
| redis | 172.40.0.32 |
| zookeeper | 172.40.0.33 |
| kafka | 172.40.0.34 |
| canal | 172.40.0.41 |
| kafka - ui | 172.40.0.42 |

### docker - composer
[基本环境搭建(consul + kong + mysql + redis + kafka + other)](docker - compose.yaml)

### canal
- [Docker快速启动](https://github.com/alibaba/canal/wiki/Docker-QuickStart)
- [Canal Kafka/RocketMQ快速启动](https://github.com/alibaba/canal/wiki/Canal-Kafka-RocketMQ-QuickStart)

配置示例：
```properties
# mysql设置
canal.instance.master.address=172.40.0.7:3306
canal.instance.dbUsername = root
canal.instance.dbPassword = root123456
# mq主题名称
canal.mq.topic=canal
```
```properties
# 设置mq为kafka，格式为json格式，kafka地址
canal.serverMode = kafka
canal.mq.flatMessage = true
kafka.bootstrap.servers = 159.75.1.200:9093
```


### 其他注意事项
1. 在配置kong时指定`KONG_DNS_RESOLVER`，可以通过consul的DNS使用到服务注册与发现，无需手动维护服务列表。注册到consul的服务名格式为`服务名.service.consul`。
2. `consul + kong + 微服务`运行在同一个docker network中，以实现互相访问。
3. 在docker中暴露端口是为了给宿主机使用，docker内部可以互相访问，无需暴露端口，微服务应用无需暴露docker端口到宿主机中。
4. 可以通过 dig 172.40.0.2 -p 8600 user - HTTP.service.consul SRV 验证consul的SRV

## web服务应用启动实例
```bash
docker run --network=web_my - service - network --name=user -e WEB="/app" -itd -v /Users/linlong/Desktop/web:/app golang:1.23.8 
export GOPROXY=https://goproxy.cn,direct
cd /app
go run application/user/cmd/main.go -path /app
```
```bash
docker run --network=web_my - service - network --name=order -e WEB="/app" -itd -v /Users/linlong/Desktop/web:/app golang:1.23.8 
export GOPROXY=https://goproxy.cn,direct
cd /app
go run application/order/cmd/main.go -path /app
```

## konga配置（127.0.0.1:1337）
### 微服务配置信息
```bash
admin api: http://172.40.0.4:8001

创建service配置信息
    Protocol: http                          # 指定发送http请求
    Host: xxxx.service.consul               # 注册到consul的服务名.service.consul
    Port: 8002                              # 服务暴露的端口号

创建router配置信息
    Paths: /user                            # 指定某个路径调用那个服务(记得回车确认)
    Strip Path: false                       # 从上游请求URL中删除匹配的前缀。（是否有删除Paths前缀）
        true:   http://127.0.0.1:8000/user/hello  =>  后端path "/hello"
        false:  http://127.0.0.1:8000/user/hello  =>  后端path "/user/hello"
```
### 负载均衡配置
```bash
创建unstreams配置信息
    Name: 负载均衡的名称
    Target: 具体的服务地址:端口号

创建service配置信息
    Protocol: http                          # 指定发送http请求
    Host: unstreams                         # unstreams配置的名称 （负载均衡与正常服务的区别）

创建router配置信息
    Host: 域名地址                           # 相当与nginx中的server，根据传入的header中的Host判断
    Paths: /                                # 指定某个路径调用那个服务(记得回车确认)
    Strip Path: false                       # 从上游请求URL中删除匹配的前缀。（是否有删除Paths前缀）
```
## admin api配置
```bash
创建service配置信息
curl -X POST http://192.168.0.74:8001/services \
  --data "name=order-http" \
  --data "host=order-HTTP.service.consul" \
  --data "protocol=http" 

# 返回的json字段，设置也使用相同的字段 
{"created_at":1733465542,"updated_at":1733465542,"path":null,"host":"test.zhubaoe.cn","retries":5,"write_timeout":60000,"enabled":true,"port":80,"tags":null,"ca_certificates":null,"client_certificate":null,"read_timeout":60000,"connect_timeout":60000,"name":"test_srv","protocol":"http","tls_verify":null,"id":"9d3f5b04 - d0de - 46f4 - 8c5d - 5e7238220658","tls_verify_depth":null}

创建router配置信息
curl -X POST http://192.168.0.74:8001/routes \
  --data "name=order_routes" \
  --data "paths[]=/order" \
  --data "service.name=order-http" \
  --data "strip_path=false"

# 返回的json字段，设置也使用相同的字段  
{"created_at":1733465781,"updated_at":1733465781,"service":{"id":"9d3f5b04 - d0de - 46f4 - 8c5d - 5e7238220658"},"path_handling":"v0","methods":null,"hosts":["www.xx1.com"],"request_buffering":true,"response_buffering":true,"strip_path":true,"snis":null,"regex_priority":0,"tags":null,"paths":null,"protocols":["http","https"],"name":"test_routes","headers":null,"https_redirect_status_code":426,"id":"6af6e619 - fac8 - 4ec0 - 8908 - 751f64a75773","preserve_host":false,"sources":null,"destinations":null}
```