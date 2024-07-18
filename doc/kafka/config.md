### kafka集群配置
#### Broker
* log.dirs: 指定多个文件存储路径（起到提升读写性能，故障转移） 
* log.dir: 补充第一个参数（建议用dirs）

#### Broker连接
* listeners: <协议://主机:端口>
* advertised.listeners:  Broker用于对外发布的 
* host.name/port: 已过期用listeners代替

#### Topic
* auto.create.topics.enable: 不存在的主题是否自动创建（false保证线上不会有垃圾数据） 
* unclean.leader.election.enable:是否允许落后太多的副本成为主副本（false保证数据不会丢失） 
* auto.leader.rebalance.enable:是否定期对主题进行leader重新选举（false这个很消耗性能且无意义）

#### 数据留存
* log.retention.{hour|minutes|ms}: 日志留存多久（log.retention.hour = xx 视日志的重要性） 
* log.retention.bytes:日志能保留多大空间（默认-1） 
* message.max.bytes:每个消息最大的大小

#### Topic级别参数
在创建Topic或修改Topic中进行配置，会覆盖全局配置。
* retention.ms :消息保存时长(默认7天)
* retention.bytes：消息有多少磁盘空间（默认-1，无限制） 
* max.message.bytes：每个消息多大 
* bin/kafka-topics.sh  ##创建指定参数 
* bin/kafka-config.sh  ##修改指定参数

#### JVM参数
* 堆大小建议设置为6G（默认1G） 
* 垃圾回收机制 
  * java7 : cpu资源充裕 -XX:+UseCurrentMarkSweepGC ,否则-XX:+UseParallelGC
  * java8：默认用G1收集器就好-XX:+UseG1GC
* 启动kafka时设置
  * export KAFKA_HEAP_OPTION=--Xms6g  --Xmx6g
  * export KAFKA_JVM_PERFORMANCE_OPTS= -server -XX:useG1GC -xx:MaxGCPauseMillis = 20  -XX:InitiatingHeapOccupancyPercent=45 
  * bin/kafka-server-start.sh  config/server.properties

#### 操作系统参数
* 文件描述符限制
    ulimit -n  100000  //不同担心调大有什么影响，单调小了会经常“Too many open files”
* 文件系统类型
  ext3、ext4、XFS ,官方测试报告提示XFS性能强于ext4
* swappiness
  swap：
  - 0：会禁止kafka使用swap空间，但是物理内存用尽后会触发OOM killer（不推荐）
  - 较小值: 物理内存用尽后使用swap空间时，能观察broker性能急剧下降，能有时间进行处理
  提交时间
* flush时间：默认5s，可以适当的添加时间提高性能
  （但可能有丢失的风险，考虑到kafka本身就有多副本的）

#### 动态配置及静态配置
静态配置
  * kafka路径下有个server.properties，指定给文件启动broker 如果变成broker就需要修改文件后进行重启，对于线上这是个不允许的

动态配置
  * 在1.1.0版本引入，2.3版本中broker有200多个参数，可在下面连接进行查看 (https://kafka.apache.org/documentation/#brokerconfigs)
    * read-only ：只有重启才能生效
    * per-broker: 动态配置，只作用于broker
    * cluster-wide: 动态配置，作用于整个集群
    
由于动态配置的特殊性，保存的参数在zookeeper中，znode路径如下:
```
change用来实时监测动态参数的变更,不会保存值
topics用来保存主题的参数
users和clients用于动态调整客户端配额（限制连入集群的客户端吞吐量或cpu等资源）
/config/brokers 才是保存动态broker参数的地方，znode有两大子类节点
第一类，有固定名称<default>保存的事cluster-wide参数
第二类： 已broker.id为名，保存per-broker参数
```

cluster-wide、per-broker和static参数的优先级是这样的
>per-broker参数 > cluster-wide参数 > static参数 > Kafka默认值。

常见设置
* log.retention.ms 日志留存时间（topic级别的）
* num.io.threads和num.network.threads 两组线程组
* num。replica.fetchers 确保有充足线程可以执行follower副本向leader副本的拉取

```
//如果要设置cluster-wide范围的动态字段参数，需要指定entity-default
//如果设置的是per-broker范围的动态字段参数，需要指定 --entity-name  xxx

//broker
kafka-configs.sh --bootstrap-server  kafka-server:9092  --alter --entity-type brokers --entity-name  1   --add-config 'num.io.threads=10'
kafka-configs.sh --bootstrap-server  kafka-server:9092  --alter --entity-type brokers --entity-default   --add-config "num.io.threads=10"
kafka-configs.sh --bootstrap-server  kafka-server:9092  -entity-type brokers --entity-default  --describe

//topic
kafka-configs.sh --bootstrap-server  kafka-server:9092  --alter  --entity-type topics --entity-name  aaa  --add-config 'max.message.bytes=128000'
```
___

### 常用配置

#### Broker配置

|                配置项                 |         说明         |                  默认值                   |
|:----------------------------------:|:------------------:|:--------------------------------------:|
|             `log.dirs`             |      指定文件存储路径      |                                        |
| `offsets.topic.replication.factor` |       分区副本数量       |               broker数量相同               |
|  `unclean.leader.election.enable`  | 落后的分区是否允许unclean选举 |                 false                  |
|       `min.insync.replicas`        |    多少副本成功才能算可见性    | `offsets.topic.replication.factor - 1` |

#### Topic配置

| 配置项                             | 说明                                         | 默认值      |
|:----------------------------------:|:--------------------------------------------:|:-----------:|
| `auto.create.topics.enable`        | 不存在的主题是否自动创建                    |             |
| `unclean.leader.election.enable`   | 落后的分区是否允许unclean选举                | false       |
| `auto.leader.rebalance.enable`     | 是否定期重新leader选举                      | false       |
| `message.max.bytes`                | 每个消息最大的大小                          |             |
| `retention.bytes`                  | 日志保留最大空间                             | 默认-1不限制 |
| `retention.{hour|minutes|ms}`      | 日志保留多长时间                             |             |

#### Producer配置

| 配置项                            | 说明                                         | 默认值       |
|:---------------------------------:|:--------------------------------------------:|:------------:|
| `request.required.acks`           | 生产者消息确认方式                           | 0:异步 1:leader确认 -1:全部确认 |
| `message.send.max.retries`        | 消息发送重试次数                             |              |
| `buffer.memory`                   | 缓冲消息的缓冲区大小                         |              |
| `compression.type`                | 压缩方式                                     | none、gzip、snappy 和 lz4 |
| `batch.size`                      | 批大小                                       |              |
| `linger.ms`                       | 批时间                                       |              |
| `request.timeout.ms`              | 等待请求响应超时时间                         |              |

#### Consumer配置

|               配置项               | 说明                                         | 默认值      |
|:----------------------------------:|:--------------------------------------------:|:-----------:|
| `enable.auto.commit`               | 消息是否自动提交                            | false       |
| `auto.commit.interval.ms`          | 自动提交时的提交时间                        |             |
| `request.timeout.ms`               | 等待请求响应超时时间                        |             |


___

### 参数优化
优化漏斗，我们可以在每一层执行对应的调整，层级越上优化的效果越明显（越下层越没有优化的可能）
![config1.png](pic%2Fconfig1.png)

* 应用层：使用合理的算法，数据结构，缓存等开销计算 
* 框架层：合理设置kafka的各种参数进行调优 *
* JVM层：broker是java进程，所有jvm优化的效果虽然比不上上两层，但是有时有巨大的改善 *
* 系统层：由于系统已经优化了很多，可能在一些配置设置有误

#### 操作系统层优化
1. 挂载文件系统禁用atime更新，记录文件最后访问时间，由于会访问inode资源，禁用会减少系统写操作
mount -o noatime
2. 文件系统建议选择ext4或XFS,能帮助kafka改善I/O性能。
3. swap空间设置为一个较小值，防止linux oom killer进程
sudo sysctl.vm.swappiness=N   或/etc/sysctl.confi增加vm.swappiness=N
4. ulimit -n文件打开数量    vm.max_map_count最大内存映射数
/etc/sysctl.conf增加vm.max_map_count=655360保存后执行  sysctl -p

#### JVM优化
1. JVM堆大小设置为6-8G（如果精确的话可以查看GC log，关注Full GC之后堆存活对象的大小，设置2倍）
2. GC收集器的选择（建议使用G1,方便省事），但要竭力的避免Full GC，由于是单线程运行，非常慢
3. 大对象，指占用至少半个区域大小的对象，可以调大区域的大小 
   * JVM启动参数 -XX:+G1HeapRegionSize=N

#### Broker优化
1. 保持客户端版本和broker版本一致（减少压缩/解压的开销，zero copy等）
2. 参数的调整

#### 应用层优化
1. 不要频繁创建Producer和Consumer对象，尽可能复用（metadata及tcp连接）
2. 用完及时关闭，如不及时关闭必定为造成资源的浪费
3. 合理使用多线程，kafka的java producer是线程安全，可以在多个线程中共享同一个实例，



吞吐量优化
```
假如一条消息花费2ms，
延迟：2ms 吞吐量：500/s  	   （1000ms / 2ms = 500/s）
如果producer每次不是发一条信息，发送前等待一段时间，统一一批发送，比如等待8ms，共缓存了1000条消息
延迟： 2ms + 8ms = 10ms
吞吐量：100000/s       (1000ms / (10ms / 1000)  = 100000/s)
实际上愿意用较少的延迟代价换取TPS的提升
```
#### broker ：
1. 增加num.replica.fetchers参数，但不超过cpu核数（副本用多少线程来拉取消息）
2. GC避免经常性的Full GC

#### Producer：
1. 增加batch.size大小，默认16kb，可加到512kb或1MB（批次大小）
2. 加大linger.ms的值 （批次时间）
3. 设置compression.type=lz4或zstd （压缩方式）
4. acks=0或1 （提交方式）
5. retries = 0  （重试次数）
6. 如果多线程共享一个producer就增加buffer.memory （消息缓冲区）

#### Consumer
1. 采用多进程/线程同时消费
2. 增加fetch.min.bytes参数，如1k或更大 (broker积累多少字节再发给Consumer)

#### 压缩
尽可能做到producer端压缩，broker端保持，consumer端解压

消息层次（一个消息集合中包含若干条日志项，日志项才是封装消息的地方）
1. 消息集合		
2. 消息

v1版本对每条消息进行CRC校验，某些情况CRC值是会发生改变，鉴于这些情况对每条消息CRC校验没必要
v2版本针对v1，把消息的公共部分抽取出来放到外层的消息集合中，就不用每条消息都保存这个信息 
  * broker对消息时间戳更新，计算后CRC会改变 
  * broker在执行消息转换时（为了兼容老版本客户端），CRC也会改变
v2压缩能力更强，把多条信息进行压缩放到外层消息字段中，是对整个消息进行压缩


何时压缩： 
* 生产者
  配置了compersion.type参数指定压缩算法 
* broker（大部分原封不动的保存）
  1. broker端与producer端采用不同的压缩算法（broker默认为producer即跟随生产者统一算法）
  2. 为了兼顾老版本的消费者程序（kafka集群同时保存多个版本非常常见）
  
何时解压
  1. 消费者接收到消息时，会根据消息中启动了什么算法进行解压


压缩对比： GZIP、Snapppy、LZ4、 zstd(2.1.1之后) 压缩比越好，当然吞吐量是越差（压缩越小越消耗性能）
压缩比（压缩后能缩小多少）： zstd> LZ4 > GZIP >Snappy
压缩吞吐量（每秒能压缩多少）：LZ4> Snappy >zstd和GZIP

实践
1. Producer端CPU资源好开启，CPU资源差关闭 
2. 带宽即便足够也建议开启压缩 
3. 尽可能避免意料之外的解压缩（兼容老版本引入的解压缩，尽可能不要出现消息格式转换）