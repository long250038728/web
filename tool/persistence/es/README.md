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