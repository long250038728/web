# web


// GMP 还有 饥饿模式，正常模式 （第九周最后一节课加餐没听，第十周第一节）
// 服务治理总结没听  ， 可观测性可以再粗略看一下

// GC学习一下




所有的服务/中间件都应该在服务器集群中且不可暴露，仅提供少量的端口对外暴露(如网关入口)。保证了服务/中间件的安全不被恶意攻击且合理有效的控制内部人员的使用权限
    1.公/私有云集群
    2.docker network

暴露端口
    1.consul 8500 为了可以观察服务注册发现相关的信息
    2.kong 通过网关入口可以访问到后端web服务
    



通过docker network 内部集群演示（由于演示内存会设置较小。`--memory`或`-m`参数内存限制的设置）：
    1.docker network 创建
        docker network create my-network
    2.consul 创建
        docker pull consul:1.15
        docker run --name=consul -d -p 8500:8500 -m 128m consul:1.15 agent -dev -ui -client='0.0.0.0'
    3.kong 创建
            
    4.web服务应用
        docker pull golang:1.20 
        docker run --name=app -d-v /Users/linlong/Desktop/web:/app  -m 128m golang:1.20 



docker run -d --name kong-database \
--network=my-network \
-p 5432:5432 \
-e "POSTGRES_USER=kong" \
-e "POSTGRES_DB=kong" \
-e "POSTGRES_PASSWORD=kong" \
postgres

docker run --rm \
--network=my-network \
-e "KONG_DATABASE=postgres" \
-e "KONG_PG_HOST=kong-database" \
-e "KONG_PG_USER=kong" \
-e "KONG_PG_PASSWORD=kong" \
-e "KONG_CASSANDRA_CONTACT_POINTS=kong-database" \
kong kong migrations bootstrap


docker run -d --name kong \
--network=my-network \
-e "KONG_DATABASE=postgres" \
-e "KONG_PG_HOST=kong-database" \
-e "KONG_PG_USER=kong" \
-e "KONG_PG_PASSWORD=kong" \
-e "KONG_PROXY_ACCESS_LOG=/dev/stdout" \
-e "KONG_ADMIN_ACCESS_LOG=/dev/stdout" \
-e "KONG_PROXY_ERROR_LOG=/dev/stderr" \
-e "KONG_ADMIN_ERROR_LOG=/dev/stderr" \
-e "KONG_ADMIN_LISTEN=0.0.0.0:8001, 0.0.0.0:8444 ssl" \
-e "KONG_PROXY_LISTEN=0.0.0.0:8000, 0.0.0.0:9080 http2, 0.0.0.0:9081 http2 ssl" \
-p 8000:8000 \
-p 9080:9080 \
-p 8443:8443 \
-p 8001:8001 \
-p 8444:8444 \
kong


docker run -d --name konga \
-p 1337:1337 \
--network my-network \
pantsel/konga