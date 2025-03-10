![图片](doc/microservices/microservices.png)

# 透明多级分流
1. 本地缓存、内容分发网络、反向代理
2. 使用较小的代价堆叠机器能够自动扩缩的服务节点
3. 调用链末端的传统关系数据库

### 高可用设计目标
尽可能减少单点部件，如果某些单点是无可避免的，则应尽最大限度减少到达单点部件的流量。此时引导请求分流至最合适的组件降低单点的压力。
1. 客户端缓存 (空间换时间的概念，把部分数据缓存在客户端本地无需每次都调用获取加快速度&减少服务器压力)
2. DNS (可通过不同运营商，不同区域获取最近的节点。提高响应速度及避免单点服务/负载均衡的压力)
3. CDN (对于静态资源，如图片、css、js等，可减少对服务的获取，缓存的概念)
4. 负载均衡 (可以把请求流量指向不同的服务器集群) ————四层(根据IP等报文即可)及七层(读取所有的数据后再转发)
   * 为了避免负载均衡成了单点故障需要引入`KeepAlived`
   * 调度算法（轮训，加权轮训，最少连接，加权最少连接，ip hash）
   * 为了避免把请求流量指向已经不可以的服务器上，所以需要有`健康检测`
5. 服务器集群 (无状态服务客户快速扩展，避免单台服务的计算压力)
6. 参数校验 (检验参数的合理性可以减少多余的查询/处理)
7. 缓存 (数据库QPS较低，把频繁请求的数据缓存到内存中。如果有全局的概念可用分布式缓存)
8. NoSQL非事务型数据库 (对非事务要求的操作存储/读取到非事务型数据库)
9. 异步队列/事务数据库 (让一些实时性不高的异步处理)
10. 同步事务数据库 (最终落到数据库)


---

# web
所有的服务/中间件都应该在服务器集群中且不可暴露，仅提供少量的端口对外暴露(如网关入口)。保证了服务/中间件的安全不被恶意攻击且合理有效的控制内部人员的使用权限
1. 公/私有云集群
2. docker network
3. k8s

# 运行场景
1. 本地调试:  在本地运行时一般不会把所有的服务启动运行一般是针对少数的服务跨服务的调用， grpc使用本地IP:port获取信息进行调用 (无需注册中心)
2. 对外服务(k8s): 在k8s每个服务都应该创建对应的server进行对外访问，grpc可使用k8sDNS的机制(server-name:port)进行调用  (无需注册中心)、
3. 对外服务(注册中心): 程序启动时需要把服务注册到注册中心提供访问，
   * 使用DNS形式(通过注册中心DNS进行地址解析获取)
   * 使用HTTP形式的(通过注册中心获取服务的对应列表进行细粒度的地址获取)

  
常用暴露端口
* 服务注册与发现
  * consul 8500 为了可以观察服务注册发现相关的信息 
* 服务网关
  * kong 8000 通过网关入口可以访问到后端web服务 
  * konga 1337 通过konga配置服务（内部调用kong admin 端口8001）
* 配置中心
  * etcd 2379
* SQL\NoSQL
  * mysql 3306 
  * redis 6379 
* 消息队列
  * kafka 9092
  * rocketmq 9876
* 服务检测
  * openTelemetry 4317
  * jaeger webUI 16686

固定ip
1. 172.40.0.2 consul
2. 172.40.0.3 kong-database
3. 172.40.0.4 kong
4. 172.40.0.5 konga
5. 172.40.0.6 etcd
6. 172.40.0.7 mysql
7. 172.40.0.8 canal

## docker运行
1.docker network 创建
```
docker network create my-network
docker network create --driver bridge --subnet 172.40.0.0/24 my-service-network
```

2.consul 创建
```
docker pull consul:1.15

docker run --name=consul \
--ip=172.40.0.2 \
--network=my-service-network \
-d -p 8500:8500  \
consul:1.15 agent -dev -ui -client='0.0.0.0'
```

3.kong 创建
```
docker pull postgres
docker pull kong
docker pull pantsel/konga

#这里指定ip是因为kong需要用到，同时还需要暴露给consul的dns使用
docker run -d --name kong-database \
--ip=172.40.0.3 \
--network=my-service-network \
-p 5432:5432 \
-e "POSTGRES_USER=kong" \
-e "POSTGRES_DB=kong" \
-e "POSTGRES_PASSWORD=kong" \
postgres

docker run --rm \
--network=my-service-network \
-e "KONG_DATABASE=postgres" \
-e "KONG_PG_HOST=kong-database" \
-e "KONG_PG_USER=kong" \
-e "KONG_PG_PASSWORD=kong" \
-e "KONG_CASSANDRA_CONTACT_POINTS=kong-database" \
kong kong migrations bootstrap

#这里KONG_DNS_RESOLVER是为了可以通过consul的dns使用到服务注册与发现（不用手动维护服务列表）
docker run -d --name kong \
--ip=172.40.0.4 \
--network=my-service-network \
-e "KONG_DATABASE=postgres" \
-e "KONG_PG_HOST=172.40.0.3" \
-e "KONG_PG_USER=kong" \
-e "KONG_PG_PASSWORD=kong" \
-e "KONG_PROXY_ACCESS_LOG=/dev/stdout" \
-e "KONG_ADMIN_ACCESS_LOG=/dev/stdout" \
-e "KONG_PROXY_ERROR_LOG=/dev/stderr" \
-e "KONG_ADMIN_ERROR_LOG=/dev/stderr" \
-e "KONG_DNS_RESOLVER=172.40.0.2:8600" \
-e "KONG_DNS_ORDER=SRV,LAST,A,CNAME" \
-e "KONG_ADMIN_LISTEN=0.0.0.0:8001, 0.0.0.0:8444 ssl" \
-e "KONG_PROXY_LISTEN=0.0.0.0:8000, 0.0.0.0:9080 http2, 0.0.0.0:9081 http2 ssl" \
-p 8000:8000 \
-p 9080:9080 \
-p 8443:8443 \
-p 8001:8001 \
-p 8444:8444 \
kong

docker run -d --name konga \
--ip=172.40.0.5 \
--network=my-service-network \
-p 1337:1337 \
pantsel/konga


DNS验证
1. dig @127.0.0.1 -p 8600 user.service.consul  SRV          //在consul服务器上
2. dig $KONG_DNS_RESOLVER -p 8600 user.service.consul  SRV  //在kong服务器上
```


4.etcd 创建
```
docker pull bitnami/etcd:latest

docker run -d \
  --ip=172.40.0.6 \
  --network=my-service-network \
  --name etcd \
  --restart always \
  -p 2379:2379 \
  -p 2380:2380 \
  -e ALLOW_NONE_AUTHENTICATION=yes \
  bitnami/etcd:latest
```

5.mysql 创建
```
docker run --name mysql \
 --ip=172.40.0.7 \
 --network=my-service-network \
 -e MYSQL_ROOT_PASSWORD=root123456 \
 -p 3306:3306 -itd \
 mysql:8.0
```

6.canal 创建
>https://github.com/alibaba/canal/wiki/Docker-QuickStart
https://github.com/alibaba/canal/wiki/Canal-Kafka-RocketMQ-QuickStart
```
docker pull canal/canal-server:latest

docker run -itd \
  --ip=172.40.0.8 \
  --name canal-server \
  --network=my-service-network \
  -p 11111:11111 \
  -e canal.instance.master.address=172.40.0.7:3306 \
  -e canal.instance.dbUsername=root \
  -e canal.instance.dbPassword=root123456 \
  -e canal.mq.topic=canal \
  -e canal.serverMode=kafka \
  -e canal.mq.flatMessage=true \
  -e kafka.bootstrap.servers=159.75.1.200:9093 \
  canal/canal-server:latest
  
或是通过配置文件指定参数
vi conf/example/instance.properties
//mysql 设置
canal.instance.master.address=172.40.0.7:3306
canal.instance.dbUsername = root
canal.instance.dbPassword = root123456
//mq 主题名称
canal.mq.topic=canal


vi /usr/local/canal/conf/canal.properties 
//设置mq为kafka  格式为json格式  kafka地址
canal.serverMode = kafka
canal.mq.flatMessage = true
kafka.bootstrap.servers = 159.75.1.200:9093
```

## web服务应用
```
docker pull golang:1.20 

docker run --network=my-service-network --name=user -e WEB="/app" -p #http_port:#http_port #grpc_port:#grpc_port -itd -v /Users/linlong/Desktop/web:/app golang:1.20 
export GOPROXY=https://goproxy.cn,direct
cd /app
go run application/user/cmd/main.go -path /app


docker run --network=my-service-network --name=order -e WEB="/app" -p #http_port:#http_port #grpc_port:#grpc_port -itd -v /Users/linlong/Desktop/web:/app golang:1.20 
export GOPROXY=https://goproxy.cn,direct
cd /app
go run application/order/cmd/main.go -path /app
```

## konga配置 127.0.0.1:1337
指定kong地址
```
创建
    admin api: 172.40.0.4:8001
```
微服务网关配置
```    
创建service配置信息
    Protocol: http                          //指定发送http请求
    Host: xxxx.service.consul               //注册到consul的服务名.service.consul
    Port: 8002                              //服务暴露的端口号

创建router配置信息
    Paths: /user                            //指定某个路径调用那个服务(记得回车确认)
    Strip Path: false                       //从上游请求URL中删除匹配的前缀。（是否有删除Paths前缀）
        true:   http://127.0.0.1:8000/user/hello  =>  后端path "/hello"
        false:  http://127.0.0.1:8000/user/hello  =>  后端path "/user/hello"
```

负载均衡配置
```
// 模拟请求域名为"xxxx"
// curl -i http://0.0.0.0 --header "Host: 域名地址"
创建unstreams配置信息
    Name: 负载均衡的名称
    Target: 具体的服务地址:端口号

创建service配置信息
    Protocol: http                          //指定发送http请求
    Host: unstreams                         //unstreams配置的名称

创建router配置信息
    Host: 域名地址                           //相当与nginx中的server，根据传入的header中的Host判断
    Paths: /                                //指定某个路径调用那个服务(记得回车确认)
    Strip Path: false                       //从上游请求URL中删除匹配的前缀。（是否有删除Paths前缀）
```

admin api配置
```
创建service配置信息
curl -X POST http://0.0.0.0:8001/services \
  --data "name=test_srv" \
  --data "host=test.zhubaoe.cn" \
  --data "protocol=http" 

# 返回的json字段，设置也使用相同的字段 
{"created_at":1733465542,"updated_at":1733465542,"path":null,"host":"test.zhubaoe.cn","retries":5,"write_timeout":60000,"enabled":true,"port":80,"tags":null,"ca_certificates":null,"client_certificate":null,"read_timeout":60000,"connect_timeout":60000,"name":"test_srv","protocol":"http","tls_verify":null,"id":"9d3f5b04-d0de-46f4-8c5d-5e7238220658","tls_verify_depth":null} 

创建router配置信息
curl -X POST http://0.0.0.0:8001/routes \
  --data "name=test_routes" \
  --data "hosts[]=www.xx1.com" \
  --data "paths[]=/hello" \
  --data "service.name=test_srv" \
  --data "strip_path=false"

# 返回的json字段，设置也使用相同的字段  
{"created_at":1733465781,"updated_at":1733465781,"service":{"id":"9d3f5b04-d0de-46f4-8c5d-5e7238220658"},"path_handling":"v0","methods":null,"hosts":["www.xx1.com"],"request_buffering":true,"response_buffering":true,"strip_path":true,"snis":null,"regex_priority":0,"tags":null,"paths":null,"protocols":["http","https"],"name":"test_routes","headers":null,"https_redirect_status_code":426,"id":"6af6e619-fac8-4ec0-8908-751f64a75773","preserve_host":false,"sources":null,"destinations":null}  
```