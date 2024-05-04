### ElasticSearch
Elastic Stack （ELK）
ElasticSearch  存储（倒排索引）
Logstash       日志收集（强大数据处理,传输）
Kibana         视图
Beats          日志收集（轻量数据收集,有日志采集器Filebeat,指标采集器Metricbeat,网络数据采集器Packetbeat）

注：
数据需要转换用Logstash
只需收集Beats

使用用途
1.全文检索    2.日志分析   3.商业智能


特点：
1.采用了倒排索引可以快速找到对应的数据（缺点：资源占用大）,可以对数据快速检索、聚合、统计、分析
2.采用API的方式,可以适配各种编程语言及客户端
3.采用集群副本等分布式理念。可以快速写入/查询。确保数据的可用性  
（分片：可以把数据写到不同的Node提高io写   副本：读取可以通过不同的Node提高io读,同时可以解决单一Node挂后数据无法找回问题）
```
//倒排索引  
{id：1 ,doc： 今天天气晴朗,这天气真舒服}
{id：2 ,doc： 今天是星期二}

转换为
    今天      (1:1) (2:1)    ===>   (文档id:当前文档出现次数)
    天气      (1:2)
    晴朗      (1:1)
    这        (1:1)
    真        (1:1)
    舒服      (1:1)
    星期      (2:1)
    二        (2:1)
```

### 索引
```
//索引创建
PUT 索引名　　　
{
   "aliases": {},             //设置别名
   "setting":{
     "number_of_shares": 2,   //*两个分片(静态设置——只有在创建的时候指定,后续修改无效。只能通过reindex迁移实现)
     "number_of_relicas": 1,  //一个副本 (动态设置),
     "analysis": {
        "char_filter": {                     //对应的字符过滤部分
            "my_char_filter":{"type":"mapping","mappings":[",=> "]}                                                          //把包含,的过滤掉
        },                 
        "tokenizer": {                       //对应文本切分为分词部分
            "my_tokenizer":{"my_tokenier":{"type":"pattern","pattern":"""\;"""}}                                             //将";"作为自定义分词
        }, 
        "filter": {                         //对应分词后再过滤部分
            "my_filter":{"type":"synonym","expand":true,"synonyms":["leileili => lileilei", "meimeihan => hanmeimei"]}      //添加同义词"leileili => lileilei", "meimeihan => hanmeimei"
        },    
        "analyzer": {                       //对应分词器,包含上述三种
            "my_analyzer":{"tokenizer":"my_tokenizer","char_filter":["my_char_filter"],"filter":["filter","lowercase"]}
        },
     },
     "index.default_pipeline" : "预处理名"    //设置预处理
   },
   "mappings":{
     "_source": {"enabled": false},     //默认值true,如果false时执行update会报错
     "dynamic": false,                  //*true:会采用动态映射进行生成     false:忽略新字段类型(插入的key不在mappings中则忽略不插入)     strict:会报错
     "propertes": {
        "name":{
           "type":"text",               //定义为text文本类型（会分词）
           "analyzer": "my_analyzer",   //可使用ik_max_word分词器或自定义,分词是在数据写入阶段
           "fields": {
              "keyword": {
                 "type" : "keyword"     //name.keyword 则是一个不分词的字符串 （多字段类型）
              }
           }
        } 
     }
   }
}


//索引参数动态修改
PUT 索引名/_settings
{
    "number_of_relicas":3,      //修改为3副本 (动态设置)
    "refresh_interval": "1s",   //刷新频率
    "max_result_window": 50000, //最大窗口大小(搜索时最大返回的数据条数)
}

//删除索引(物理删除,速度快,但是无法恢复)
DELETE 索引名
//删除索引数据(逻辑删除,不会马上释放空间)
DELETE 索引名/_delete_by_query
{
    "query":{
       "match_all":{}
    }
}
```

### 映射
基本数据类型
binary
boolean
keyword
number
date
alisa   
text    
复杂数据类型
array 数组
object json对象
Nested  嵌套类型
Join    父子关联
Flattened 将复杂的object或Nested统一映射为扁平字段

### 文档
索引是一组相关的文档的集合体,文档存储在索引中。每个文档在索引中都有一个唯一的id,每个文档都是一组字段组合。字段是可以任意类型的。
结构体有几个字段
_index: 文档索引
_id: 文档id
_score: 得分
_source: 字段

```
//指定文档id插入数据或修改
PUT 索引/_doc/1
{
    "name":"hanmeimei"
}
//不指定文档id插入数据
PUT 索引/_doc/
{
    "name":"hanmeimei"
}

//批量插入数据
POST 索引/_bulk
{"index":{"_index":"索引名","_id":1}}  //index单不存在时创建,存在时更新
{"field":"value"}
{"delete":{"_index":"索引名","_id":1}}

//删除文档
DELETE 索引名/_doc/1

//批量删除文档
POST 索引/_delete_by_query
{
    "query":{"match":{ "message":"hello" }}
}

//修改文档
POST 索引/_update/1
{
   "doc":{
      "name":"limeimei"
   }
}   

//修改文档通过script脚本
POST 索引/_update/1
{   
   "upsert":{
       "counter": 1       //当文档不存在是counter默认值为1
   },
   "script":{
      "lang": "painless",
      "source":"ctx._source.counter += params.count; ctx._source.tags.add(params.tag)",
      "params": {
        "count":4,
        "tag": "blue"
      }
   }
}   


//批量修改文档通过script脚本
POST 索引/_update_by_query
{   
   "upsert":{
       "counter": 1       //当文档不存在是counter默认值为1
   },
   "script":{
      "lang": "painless",
      "source":"ctx._source.counter += params.count; ctx._source.tags.add(params.tag)",
      "params": {
        "count":4,
        "tag": "blue"
      }
   },
   "query": {
      "match":{ "message":"hello" }
   }
}   

//批量修改文档通过预处理
POST 索引/_update_by_query?pipeline=预处理名



PUT _cluster/settings
{
    "persistent": { "action.auto_create_index" : false}    //插入文档时如果索引不存在是否自动创建
}


```

### 别名
别名的意义/作用
1.可以指向多个索引（如log2023,log2024。无需切换即可快速查找）
POST log2023/_search POST log2024/_search  === 改为 ===> POST log/_search
2.索引与另外一个索引的切换
POST 别名/search    别名指向索引1,此时需要指向索引2,可通过修改别名指向的索引无需修改所有业务逻辑代码

注意事项
1.写入数据是只能指定具体的索引,如果指定了别名会报错,除非在别名中设置is_write_index(但不推荐)
```
//别名
POST /_aliases
{
    "actions":{
        "add": {
            "index":"索引名",
            "alias":"别名名",
            "is_write_index": true, //允许往别名中写入数据
        }
    }
}

//查询数据
GET 索引名/_search 
{
    highight: {                         //字段高亮,在相应中有highight对象
        "fields": {"字段名":{}}
    },
    "query":{
       "match_all":{}
    }
}

//获取别名列表
GET _cat/aliases?v
```

### 模版
模版的意义/作用
需要根据时间创建不同的索引,每次手动创建很麻烦,先创建一个模版,后续使用该模版可以快速创建索引
多个索引中用到相同的字段或配置
```
//模版创建(普通模版定义)
PUT _index_template/模版名
{
    "index_patterns":{
        "logs*",                //需要匹配到什么索引
    },
    "template":{                //模版数据
        "aliases": {
            "alias1": {},                            //设置别名
        },
        "settings":{
            "number_of_shareds": 2,                  //分区数量2
            "index.default_pipeline" : "预处理名",    //设置预处理
        },
        "mappings":{
            "properties": {
                "name": {
                    "type":"keyword",               //name字段为keyword
                }   
            }
        }
    }
}

//组合模版创建(将mappings,settings等以组件方式分隔)
PUT _component_template/模版名1
{
    "template":{
        "mappings": {
            "properties": {
                "name": {
                    "type":"keyword",    //name字段为keyword
                }   
            }
        }
    }
}
PUT _component_template/模版名2
{
    "template":{
        "settings":{
            "number_of_shareds": 2,     //分区数量2
        }
    }
}
PUT _index_template/模版名3
{
    "priority": 5000,
    "composed_of": [
        "模版1","模版2"
    ],
    "version": 1,
    "_meta": {
        "description": "my custom template"
    }
}

//删除模版
DELETE _index_template/模版名
//查询模版
GET _index_template/模版名
```

### 预处理
在写入数据之前会对数据预处理操作
1.数据清洗:去除重复数据,默认值等
2.数据集成:将多个数据源放到统一的存储平台
3.数据转换:将数据转换为可以挖掘或分析的形式中

ingest节点本质是在实际文档建立索引之前使用ingest节点进行预处理。然后再通过文档的api写入。预处理的步骤有
1.定义预处理管道
2.将预处理管道与索引进行关联
3.写入数据

```
//创建预处理
PUT _ingest/pipeline/预处理管道名称
{
    "version": 1,
    "processors": [
    {
        "script": {
            "lang": "painless",
            "sourece": "ctx.mid_prefix = ctx.mid.substring(0,2)"   //把mid提取前两个字符到mid_prefix
        }
    },
    {
        "json": {
            "field" : "header.userInfo"
            "target_field" : "header.userInfoJson"                //将json字符串转换为json（object对象）
        }
    },
    {
        "script": {
            "lang" : "painless"                                    //将数组tag中的值+2
            "sourece" : """
                for (int i=0 ; i < ctx.tag.length;i++){
                    ctx.tag[i]=ctx.tag[i] + "2";
                }
            """   
        }
    },
    
    ]
}

//使用预处理
POST 索引名/_update_by_query?pipline=预处理管道名称
{
    "query": {
        "match_all" : {}
    }
}
```

### 文档迁移reindex
把一个索引迁移到另外一个索引，Reindex API 在迁移数据时默认情况下不会删除原始索引中的文档。它会从源索引中读取数据，并将其复制到目标索引中，保持原始索引的完整性。
如果您希望在成功迁移数据后删除原始索引中的文档，您可以使用 delete_after_reindex 参数来执行这个操作
```
POST _reindex
{   
    "confilicts": "proceed",                    //如果版本冲突时继续执行
    "source" : {
        "index":"source索引名",
        "query": {
            "term": {"user":"meimeihan"},       //指定条件迁移
        }
    },
    "dest": {
        "index": "dest索引名",
        "pipeline": "预处理名"                   //迁移后的索引指定使用预处理
    },
    "script": {
        "lang": "painless",
        "source": "if (ctx._source.user = "meimeihan") {ctx._source.remove('user')}"
    }
}
```

reindex一般以任务的形式
```
//获取任务列表
GET _task?detailed=true&action=*reindex

//获取任务详情
GET _task/任务id

//取消任务
POST _tasks/任务id/_cancel
```

### 脚本
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

### 重要参数
setting中重要配置
```
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
```