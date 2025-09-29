
```
// 查看cat命令下有什么字命令
GET /_cat
```


集群
```
// 查看集群健康状态
GET /_cat/health?v

// 查看集群状态
GET /_cluster/state

// 查看集群统计信息
GET /_cluster/stats

// 查看统计信息
GET _stats/search?pretty
GET _stats/indexing
```

节点
```
// 节点基本信息
GET /_cat/nodes?v
GET /_nodes/1714387876000164332
GET /_nodes/stats/jvm?pretty
```

索引&分片
```
// 查看索引
GET /_cat/indices?v
GET /_cat/indices/goods_stock_history_2024_10?v

// 查看索引的设置及字段等详情
GET /goods_stock_history_2024_10

// 查看分片
GET /_cat/shards?v

// 查看索引分片信息
GET /_cat/shards/zby_stock_allocation_record_2?v

```