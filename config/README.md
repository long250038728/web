获取配置信息加载中间件
根配置文件路径下必须的yaml文件
1. config.yaml 获取各个服务的配置信息（http端口，grpc端口等）


配置分为两种
1. 根据文件获取配置信息
    > app.InitPathInfo(根配置文件路径, 服务名)
2. 根据配置中心获取配置中心
    > app.InitCenterInfo(根配置文件路径, 服务名)


文件
```
db.yaml           // 数据库       
db_read.yaml      // 数据库(只读)                    
es.yaml           // ES     
mq.yaml           // 消息队列   
redis.yaml        // 缓存     
register.yaml     // 服务注册 
tracing.yaml      // 链路    
```   

配置中心
```
center.yaml        //配置中心的地址
```
