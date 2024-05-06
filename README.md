# web


// GMP 还有 饥饿模式，正常模式 （第九周最后一节课加餐没听，第十周第一节）
// 服务治理总结没听  ， 可观测性可以再粗略看一下

// GC学习一下




所有的服务/中间件都应该在服务器集群中且不可暴露，仅提供少量的端口对外暴露(如网关入口)。保证了服务/中间件的安全不被恶意攻击且合理有效的控制内部人员的使用权限
    1.公/私有云集群
    2.docker network

暴露端口
    1.consul 8500 为了可以观察服务注册发现相关的信息
    2.kong 8000 通过网关入口可以访问到后端web服务
    3.konga 1337 通过konga配置服务（内部调用kong admin 端口8001）

固定ip
    172.22.0.2 consul
    172.22.0.3 kong-database
    172.22.0.4 kong
    172.22.0.4 konga


### docker运行
1.docker network 创建
```
docker network create my-network
docker network create --driver bridge --subnet 172.22.0.0/24 my-service-network
```

2.consul 创建
```
docker pull consul:1.15
docker run --name=consul \
--ip=172.22.0.2 \
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
--ip=172.22.0.3 \
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
--ip=172.22.0.4 \
--network=my-service-network \
-e "KONG_DATABASE=postgres" \
-e "KONG_PG_HOST=172.22.0.3" \
-e "KONG_PG_USER=kong" \
-e "KONG_PG_PASSWORD=kong" \
-e "KONG_PROXY_ACCESS_LOG=/dev/stdout" \
-e "KONG_ADMIN_ACCESS_LOG=/dev/stdout" \
-e "KONG_PROXY_ERROR_LOG=/dev/stderr" \
-e "KONG_ADMIN_ERROR_LOG=/dev/stderr" \
-e "KONG_DNS_RESOLVER=172.22.0.2:8600" \
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
--ip=172.22.0.5 \
--network=my-service-network \
-p 1337:1337 \
pantsel/konga
```

4.web服务应用
```
docker pull golang:1.20 
docker run --network=my-service-network --name=app -itd -v /Users/linlong/Desktop/web:/app golang:1.20 
export GOPROXY=https://goproxy.cn,direct
cd /app
go run application/user/cmd/main.go
```

### konga配置 127.0.0.1:1337
创建
    admin api: 172.22.0.4:8001

创建service配置信息
    Protocol: http                          //指定发送http请求
    Host: xxxx.service.consul               //注册到consul的服务名.service.consul
    Port: 8002                              //服务暴露的端口号

创建router配置信息
    Paths: /user                            //指定某个路径调用那个服务(记得回车确认)
    Strip Path: false                       //从上游请求URL中删除匹配的前缀。（是否有删除Paths前缀）
        true:   http://127.0.0.1:8000/user/hello  =>  后端path "/hello"
        false:  http://127.0.0.1:8000/user/hello  =>  后端path "/user/hello"