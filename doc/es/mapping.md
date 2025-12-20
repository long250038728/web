# Elasticsearch Mapping 全参数详解

## 一、通用核心参数（所有字段类型均可配置，含type/index）
| 参数名          | 默认值       | 适用类型       | 含义 & 核心场景                                                                 |
|-----------------|--------------|----------------|--------------------------------------------------------------------------------|
| type            | 无（必配）| 所有类型       | 定义字段的基础数据类型（核心必配参数）：<br>- 数值型：int/long/float/double/short/byte；<br>- 文本型：text（分词）/keyword（精确匹配）；<br>- 日期型：date；<br>- 对象型：object/nested；<br>- 其他：boolean/ip/geo_point等|
| index           | true         | 所有类型       | 控制字段是否生成倒排索引（核心优化参数）：<br>- true：生成倒排索引，支持match/term等查询；<br>- false：不生成倒排索引，仅存储在_source（无法查询，节省内存/磁盘）；<br>- not_analyzed（7.x+废弃）：等价于index:true + type:keyword |
| enabled         | true         | 所有类型       | 控制字段是否「存储 + 索引」：<br>- true：正常存储、解析字段值，参与索引构建；<br>- false：完全忽略该字段（不存储、不索引、不显示在_source），节省磁盘/内存（适合无用冗余字段） |
| null_value      | null         | 非text类型     | 字段值为null/空时的替代值（text不支持）：<br>例：int字段设null_value: 0，可通过term: 0查询所有null值文档。<br>👉 注意：仅影响索引，_source仍显示null。 |
| copy_to         | 无           | 所有类型       | 将当前字段值复制到目标虚拟字段（目标字段无需显式定义），用于多字段聚合/查询。<br>例："copy_to": "all_search"，title/content都复制到all_search，查询该字段可匹配两者。 |
| dynamic         | true         | object/nested/根文档 | 控制动态新增字段规则：<br>- true：自动检测新字段并生成Mapping；<br>- false：忽略新字段；<br>- strict：新增字段报错（生产建议核心索引设strict）。 |
| boost           | 1.0          | 所有类型       | 字段权重系数，影响查询相关性得分：设为2.0则得分是默认值的2倍（适合核心查询字段）。 |

## 二、文本字段（text）专属参数
| 参数名                | 默认值               | 含义 & 核心场景                                                                 |
|-----------------------|----------------------|--------------------------------------------------------------------------------|
| analyzer              | standard             | 索引时分词器（如ik_max_word中文分词、standard英文分词），决定文本拆分为term的规则。 |
| search_analyzer       | 同analyzer           | 查询时分词器（建议与analyzer一致，中文可设ik_smart）。|
| normalizer            | 无                   | 轻量级标准化器（适合keyword）：如小写、去空格（例：{"normalizer": "lowercase"}）。 |
| index_options         | docs                 | 倒排索引存储信息粒度：<br>- docs：仅存文档ID（仅判断匹配）；<br>- freqs：ID+词频（影响评分）；<br>- positions：ID+词频+位置（支持短语查询）；<br>- offsets：全量信息（支持高亮）。<br>👉 无需评分/高亮设docs，节省内存。 |
| norms                 | true                 | 是否存储归一化因子（用于相关性得分）：设false可节省内存（仅过滤的字段）。|
| fielddata             | false                | 是否启用fielddata（文本字段排序/聚合用，默认关闭，启用占大量堆内存，生产慎用）。 |
| position_increment_gap | 100                 | 多值文本字段的位置间隙：避免跨值误匹配短语（如["a","b"]设100，不匹配"a b"）。 |
| similarity            | BM25                 | 相关性评分算法：BM25（默认）、classic（TF-IDF）、boolean（布尔评分）。|
| term_vector           | no                   | 是否存储词向量（分词后的词+位置/偏移）：设with_positions_offsets支持高亮，但增加磁盘占用。 |

## 三、数值/日期/布尔等基础类型专属参数
| 参数名            | 默认值                | 适用类型               | 含义 & 核心场景                                                                                                                      |
|-------------------|-----------------------|------------------------|--------------------------------------------------------------------------------------------------------------------------------|
| doc_values        | true                  | 数值/date/keyword/boolean等（text除外） | 控制是否启用列式存储（核心性能参数）：<br>- true：启用列式存储，支持排序/聚合/脚本操作（默认开启）；<br>- false：禁用列式存储，节省磁盘/内存，但无法排序/聚合                                    |
| format            | 类型专属默认值        | date/数值              | 字段值格式规则：<br>- 日期：默认strict_date_optional_time or epoch_millis，可自定义yyyy-MM-dd HH:mm:ss；<br>- 数值：仅控制展示，不影响存储/索引。<br>👉 int字段无需配置。 |
| ignore_malformed  | false                 | 数值/date/ip/keyword   | 是否忽略格式错误值：<br>- false：值错误（int传字符串）报错；<br>- true：忽略错误值，文档正常写入。<br>👉 int字段建议设true，防脏数据导致写入失败。                                   |
| coerce            | true                  | 数值/date              | 是否自动类型转换：<br>- true：字符串"123"转int、123.9转int（取整）；<br>- false：严格匹配类型，转换失败报错。<br>👉 核心int字段设false，避免隐式转换。                          |
| precision_step    | int/long=16；float/double=8 | 数值类型       | 数值倒排索引精度步长：步长越小范围查询越快，但索引体积越大（int默认16足够）。                                                                                      |
| ignore_above      | null（不限制）| keyword/text           | 忽略超过指定长度的值（仅索引≤长度的值）：例ignore_above:256，超256字符的keyword不索引。<br>👉 int字段无需配置。                                                     |

## 四、对象/嵌套/特殊类型专属参数
| 参数名          | 默认值       | 适用类型       | 含义 & 核心场景                                                                 |
|-----------------|--------------|----------------|--------------------------------------------------------------------------------|
| properties      | 无           | object/nested  | 定义对象/嵌套字段的子字段（必选）：<br>例："user": {"type":"object","properties":{"name":{"type":"keyword"}}}。 |
| include_in_all  | true         | 7.x+废弃       | 旧版参数（被copy_to替代），控制字段是否包含在_all虚拟字段中。|
| ignore_z_value  | true         | geo_point      | 地理坐标字段是否忽略Z轴/高度（仅保留经纬度）。|
| max_shingle_size | 无          | text/keyword   | 分词器shingle大小：设2则索引相邻两个词，提升短语查询速度。|

## 五、存储与返回控制参数
| 参数名          | 默认值       | 适用类型       | 含义 & 核心场景                                                                 |
|-----------------|--------------|----------------|--------------------------------------------------------------------------------|
| store           | false        | 所有类型       | 是否单独存储字段值（脱离_source）：<br>- false：仅存在于_source；<br>- true：单独存储，可通过fields参数返回（大字段适用）。 |
| _source         | true         | 文档级（非字段） | 控制是否存储原始文档：<br>- false：节省磁盘，但无法reindex/更新字段；<br>- 可通过includes/excludes指定存储/排除字段：<br>"_source": {"excludes":["useless_field1"]}。 |
| fields          | 无           | 所有类型       | 多字段定义：一个字段支持多种索引方式（如text分词 + keyword精确匹配）：<br>例："title":{"type":"text","fields":{"keyword":{"type":"keyword"}}}。 |
| alias           | 无           | 7.x+支持       | 字段别名：映射到真实字段（例：{"alias":{"type":"alias","path":"real_int_field"}}）。 |

## 六、性能优化参数（腾讯云ES集群重点关注）
| 参数名                  | 默认值       | 适用类型       | 含义 & 核心场景                                                                 |
|-------------------------|--------------|----------------|--------------------------------------------------------------------------------|
| eager_global_ordinals   | false        | keyword/text   | 是否提前加载全局序号（聚合用）：<br>- true：索引刷新时加载，提升聚合速度，消耗更多内存；<br>👉 高频聚合的int/keyword字段可设true。 |
| index_phrases           | false        | text           | 是否索引常用双词短语：提升短语查询速度，增加索引体积。|
| index_prefixes          | false        | keyword/text   | 是否索引前缀：设{"min_chars":1,"max_chars":5}，提升前缀查询（如term: abc*）速度。 |
| fielddata_frequency_filter | 无      | text           | fielddata频率过滤：仅加载高频term（如{"min":0.01,"max":1.0}），减少内存占用。 |

# 核心优化建议（适配你的100个int字段场景）
1. type: integer：所有int字段必配基础类型；
2. index: false：无需查询的字段必配，减少倒排索引（核心降内存手段）；
3. doc_values: true/false：需排序/聚合则保留true，无需则设false；
4. ignore_malformed: true：int字段防脏数据，避免写入失败；
5. coerce: false：核心int字段关闭隐式转换，保证数据准确性；
6. enabled: false：完全无用的int字段直接禁用，比index: false更彻底；
7. _source.excludes：剔除无需返回的int字段，节省磁盘和网络传输。