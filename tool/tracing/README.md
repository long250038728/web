## tracing
日志：来记录发生的请求及整个过程。
链路：了解服务之间的调用关系

opentracing通过日志，链路，监控三个整合在一起，通过一个链路就能排查服务之间联系，中间处理的事项信息


---
#### 加载opentracing jaeger库
```
go get github.com/opentracing/opentracing-go
go get github.com/uber/jaeger-client-go
go get github.com/uber/jaeger-lib
```
