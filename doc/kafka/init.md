### 创建网络
```
docker network create app-tier --driver bridge  
```

### 创建zookeeper
```
docker run -d --name zookeeper-server \
--network app-tier \    
-e ALLOW_ANONYMOUS_LOGIN=yes \     
bitnami/zookeeper:latest
```

### /创建kafka
```
docker run -d --name kafka-server \
--network app-tier \
-e ALLOW_PLAINTEXT_LISTENER=yes \
-e KAFKA_CFG_ZOOKEEPER_CONNECT=zookeeper-server:2181 \
bitnami/kafka:latest
```