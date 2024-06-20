![图片](doc/microservices/microservices.png)

# web
所有的服务/中间件都应该在服务器集群中且不可暴露，仅提供少量的端口对外暴露(如网关入口)。保证了服务/中间件的安全不被恶意攻击且合理有效的控制内部人员的使用权限
1. 公/私有云集群
2. docker network

暴露端口
1. consul 8500 为了可以观察服务注册发现相关的信息
2. kong 8000 通过网关入口可以访问到后端web服务
3. konga 1337 通过konga配置服务（内部调用kong admin 端口8001）
4. mysql 3306 外部mysql数据访问


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

web服务应用
```
docker pull golang:1.20 

docker run --network=my-service-network --name=user -p 8001:8001  -p 6060:6060 -itd -v /Users/linlong/Desktop/web:/app golang:1.20 
export GOPROXY=https://goproxy.cn,direct
cd /app
go run application/user/cmd/main.go -config /app/config


docker run --network=my-service-network --name=order -p 6060:6060  -itd -v /Users/linlong/Desktop/web:/app golang:1.20 
export GOPROXY=https://goproxy.cn,direct
cd /app
go run application/order/cmd/main.go -config /app/config

curl -H "Authorization:eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTYwMzQxMTksImlhdCI6MTcxNTkyNjExOSwiaWQiOjEyMzQ1NiwibmFtZSI6ImpvaG4ifQ.FzmTyzp3TK1cLiZnuv0xMQeXK01e-IlMAdOJgW3uKNU" http://172.40.0.4:8002/
```

## konga配置 127.0.0.1:1337
```
创建
    admin api: 172.40.0.4:8001

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

---

## 扩展
### 创建一个DNS解析器
```
import (
	"fmt"
	"net"
	"os"
	"github.com/miekg/dns"
)

func handleRequest(w dns.ResponseWriter, r *dns.Msg) {
	m := new(dns.Msg)
	m.SetReply(r)

	// 获取查询的域名
	name := r.Question[0].Name
	if name == "user.service.consul." {  // 对于 user.service.consul 的查询，返回一个虚构的 IP 地址
		a := &dns.A{
			Hdr: dns.RR_Header{Name: name, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 60},
			A:   net.ParseIP("192.168.1.100"),  //这里可以获取数据库或其他得到具体ip，这里是伪代码
		}
		m.Answer = append(m.Answer, a)
	}
	w.WriteMsg(m)
}

func main() {
	addr := ":53" // 定义 DNS 服务器监听的地址和端口
	server := dns.Server{Addr: addr, Net: "udp"}  // 创建一个 DNS 服务器
	dns.HandleFunc(".", handleRequest) // 定义 DNS 请求处理函数
	if err := server.ListenAndServe(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to start DNS server: %s\n", err.Error())
		os.Exit(1)
	}
}
```
使用 dig @127.0.0.1 -p 53 user.service.consul SRV 命令进行输出
```
linlong@linlongdeMacBook-Pro-2 ~ % dig @127.0.0.1 -p 53 user.service.consul
; (1 server found)
;; global options: +cmd
;; Got answer:
;; ->>HEADER<<- opcode: QUERY, status: NOERROR, id: 37398
;; flags: qr rd; QUERY: 1, ANSWER: 1, AUTHORITY: 0, ADDITIONAL: 0
;; WARNING: recursion requested but not available

;; QUESTION SECTION:
;user.service.consul.		IN	SRV

;; ANSWER SECTION:
user.service.consul.	60	IN	A	192.168.1.100

;; Query time: 0 msec
;; SERVER: 127.0.0.1#53(127.0.0.1)
;; WHEN: Tue May 07 09:45:40 CST 2024
;; MSG SIZE  rcvd: 72

linlong@linlongdeMacBook-Pro-2 ~ % dig @127.0.0.1 -p 53 user.service.consul SRV

; <<>> DiG 9.10.6 <<>> @127.0.0.1 -p 53 user.service.consul SRV
; (1 server found)
;; global options: +cmd
;; connection timed out; no servers could be reached
```

---

### go项目依赖生成图
安装dot 工具 && 安装godepgraph 工具  && 图生成
```
brew install graphviz
go install github.com/kisielk/godepgraph@latest
godepgraph -s  ./application/user/cmd/ | dot -Tpng -o godepgraph.png
```

---

###  是否数据竞争
go run -race main.go


---

### 分析
代码中加入 pprof包，启动6060端口的访问，然后执行 go tool pprof http://localhost:6060/debug/pprof/profile
```
import (
    _ "net/http/pprof"
    "net/http"
    "log"
)

// 启动 HTTP 服务器以便访问 pprof 数据
go func() {
    log.Println(http.ListenAndServe("localhost:6060", nil))
}()
```
1. 执行 go tool pprof http://localhost:6060/debug/pprof/profile  //获取 CPU 性能分析数据。具体来说，它会收集程序在一段时间内的 CPU 使用情况，包括函数调用的频率和耗时
2. 执行 go tool pprof http://localhost:6060/debug/pprof/heap     //获取堆内存的使用情况。它会收集当前程序的堆内存分配数据，包括哪些对象占用了多少内存
3. 访问 web  http://localhost:6060/debug/pprof/
   * http.HandleFunc("/debug/pprof/", Index)
   * http.HandleFunc("/debug/pprof/cmdline", Cmdline)
   * http.HandleFunc("/debug/pprof/profile", Profile)
   * http.HandleFunc("/debug/pprof/symbol", Symbol)
   * http.HandleFunc("/debug/pprof/trace", Trace)

---

### 查看内存逃逸
go build -gcflags "-m" main.go 返回结果如下
1. inlining call    内联优化（将函数调用替换为函数实际代码，减少调用函数的开销，减少栈空间的使用）
2. does not escape  没有逃逸到堆，意味着它被分配在栈上
3. escapes to heap  逃逸到堆
```
这三个并不矛盾分别描述了内联优化和逃逸分析的结果

./main.go:35:14: inlining call to app.Servers
./main.go:35:14: func literal does not escape
./main.go:35:14: ... argument escapes to heap

1. 对app.servers进行内联优化,
2. 同时函数内的变量在栈上分配
3. 然后在函数外使用了...arg此时arg变量才逃逸到堆上
```








---



// GMP 还有 饥饿模式，正常模式 （第九周最后一节课加餐没听，第十周第一节）
// 服务治理总结没听  ， 可观测性可以再粗略看一下

// GC学习一下
