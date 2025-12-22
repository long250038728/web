## fink && doris
通过fink进行ETL数据转换写入到doris数据中(OLAP)


## docker 环境搭建检查
1. http://localhost:8081 ui可以访问
2. http://localhost:8030 Doris FE Web UI
3.  mysql -h 127.0.0.1 -P 9030 -u root  登录数据库，状态
    * SHOW FRONTENDS;
    * SHOW BACKENDS;


## fink && doris
docker exec -it flink-jobmanager ./bin/sql-client.sh
1. 创建 Source 表，类似sql语句，但是有区别（这个表用于接收kafka的主题）
   event_time 为fink理解的时间    WATERMARK最多延迟5秒，超过这个时间我就当做丢了（并不是“字段”，而是Flink Table/Stream的事件时间机制配置）
2. 创建 Sink 表，用于把数据同步到doris的表格中
3. 创建 doris真实表 （根据AGGREGATE KEY = 防重复写 计算）
4. insert into select（TUMBLE(...)是一个动态的表），插入到Sink 表，映射到doris的真实表（往doris插入数据）

---

## 流程图
```
 Kafka user_events
       │
       ▼
Flink Source 表 user_events
       │
       ▼
Flink 窗口聚合计算 (TUMBLE + COUNT/SUM)
       │
       ▼
Flink Sink 表 doris_buy_stats_1m
       │
       ▼
Doris 物理表 buy_stats_1m (AGGREGATE KEY)

```


---

## sql示例
读取kafka中user_events主题的数据，形成fink的source表
```fink表 
CREATE TABLE user_events (
    user_id INT,
    event STRING,
    amount DOUBLE,
    ts BIGINT,
    event_time AS TO_TIMESTAMP_LTZ(ts * 1000, 3),
    WATERMARK FOR event_time AS event_time - INTERVAL '5' SECOND
) WITH (
    'connector' = 'kafka',
    'topic' = 'user_events',
    'properties.bootstrap.servers' = 'kafka:9092',
    'properties.group.id' = 'flink-etl',
    'scan.startup.mode' = 'earliest-offset',
    'format' = 'json'
);
```

把fink中的数据，同步到doris实际存储表
```fink表
CREATE TABLE doris_buy_stats_1m (
    window_start TIMESTAMP(3),
    cnt BIGINT,
    total_amount DOUBLE
)
WITH (
    'connector' = 'doris',
    'fenodes' = 'doris-fe:8030',
    'table.identifier' = 'realtime.buy_stats_1m',
    'username' = 'root',
    'password' = '',
    'sink.label-prefix' = 'flink_buy_1m'
);
```

doris实际存储表
```doris表 
CREATE TABLE buy_stats_1m (
    window_start DATETIME,
    cnt BIGINT,
    total_amount DOUBLE
)
AGGREGATE KEY(window_start)
DISTRIBUTED BY HASH(window_start) BUCKETS 1
PROPERTIES (
    "replication_num" = "1"
);
```

根据fink的source数据，插入到fink的sink表，sink表同步到doris的实际存储表
```
INSERT INTO doris_buy_stats_1m
    SELECT
        window_start,
        COUNT(*) AS cnt,
        SUM(amount) AS total_amount 
    FROM TABLE( 
        TUMBLE(
            TABLE user_events,
            DESCRIPTOR(event_time),
            INTERVAL '1' MINUTE
        )
    )
    WHERE event = 'buy'
    GROUP BY window_start;
```