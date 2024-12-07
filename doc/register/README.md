
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