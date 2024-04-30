```
POST /your_index/_update_by_query
{
  "script": {
    "source": "ctx._source.new_amount = ctx._source.amount",
    "lang": "painless"
  },
  "query": {
    "exists": {
      "field": "amount"
    }
  }
}
```

POST /your_index/_update_by_query: 
    这是一个 HTTP POST 请求，目标是对名为 your_index 的 Elasticsearch 索引执行 Update By Query 操作。

请求的主体包含两个部分：
    "script"：定义了要执行的脚本。在这个例子中，脚本内容是 ctx._source.new_amount = ctx._source.amount，它使用 Painless 脚本语言，将 new_amount 字段的值设置为 ctx._source.amount。
    "query"：指定了要更新的文档范围。在这里，查询条件是 exists，检查文档中是否存在名为 amount 的字段。
综合起来，这个 Update By Query 请求的作用是，对于 your_index 索引中存在 amount 字段的文档，将其每个文档的 new_amount 字段的值设置为 amount。




```
POST _reindex
{
  "source": {
    "index": "sale_index",
    "query": {
      "term": {
        "type": 1
      }
    }
  },
  "dest": {
    "index": "type_1_index"
  }
}
```
使用 Reindex API 迁移数据：
    使用 Reindex API 将符合特定 type 值的文档从原始索引 sale_index 中迁移到新的目标索引中. 
    Reindex API 在迁移数据时默认情况下不会删除原始索引中的文档。它会从源索引中读取数据，并将其复制到目标索引中，保持原始索引的完整性。
    如果您希望在成功迁移数据后删除原始索引中的文档，您可以使用 delete_after_reindex 参数来执行这个操作


setting中重要配置
index.number_of_shards: 索引的主分片数量。
index.number_of_replicas: 每个主分片的副本数量。
index.refresh_interval: 刷新间隔，控制索引的刷新频率。(文档被索引或删除时，这些变更并不会立即对搜索结果产生影响，而是会先被缓存在内存中，然后定期刷新到磁盘上的索引文件)
    较短的刷新间隔会增加索引的写入性能，因为变更会更快地对搜索操作可见，但也会增加系统的负载。
    较长的刷新间隔可以减少系统的负载，但搜索操作可能无法立即看到最新的变更
index.search.idle.after: 设置空闲搜索阶段的时间。(在执行搜索操作后，Elasticsearch 会保持搜索上下文一段时间，以便在未来的搜索请求中重用)
index.max_result_window: 查询结果窗口的最大数量。

index.blocks.read_only: 是否将索引设置为只读模式。
index.blocks.read: 是否允许读取索引。
index.blocks.write: 是否允许写入索引。
index.blocks.metadata: 是否允许修改索引的元数据。
index.mapping.total_fields.limit: 索引中字段的总数限制。
index.auto_expand_replicas: 自动扩展副本。
index.routing_partition_size: 控制路由分片的大小。