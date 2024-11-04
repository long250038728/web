## Nested Loop Join,Block Nested Loop Join（使用join_buffer版本的Nested Loop Join）
取驱动表的每一行，去驱动表关联的表中查询(如果命中索引O(logN),否则是O(N))，然后和驱动表组成一行放到结果集，以此类推，直到最后一个驱动表，最后返回结果。
* SELECT * FROM t2 join t1 on t1.id=t2.id 驱动表不一定是左边的表（根据效率更高的为驱动表），如果强制则使用`STRAIGHT_JOIN`。
* LEFT JOIN时语义上驱动表是左表。大概率是LEFT但不保证一定(这取决于查询条件、表的大小、索引的使用情况和优化器的决策)

## semi join半连接优化
查询中的in子查询可能会优化改写为INNER JOIN连表，效率更高(5.7版本)
* SELECT * FROM t1 WHERE t1.id in (SELECT id FROM t2)，如不优化每一行都要执行 (SELECT id FROM t2)子表查询判断是否符合
* SELECT t1.* FROM t1 join (SELECT id FROM t2) as tmpTable on t1.id=tmpTable.id 半连接优化后变成把子查询当成一个临时表只查询一次

## join_buffer
如果被驱动表`没有符合索引`，此时需要全表扫描，很大概率会到磁盘里面读取，那可以把`被驱动表`的数据放到join_buffer中减少多次读磁盘的次数

## hash join(8.0版本)
利用join_buffer优化，将数据结构改为`hash`，此时与该buffer的匹配复杂度就变成`O(1)`

## join是否改为应用层实现分析
* 如果被驱动表数据量小时可以用join，如果被驱动表数据量很大，join_buffer放不下，此时会使用临时表，效率会降低
  * 此时把这个join_buffer改为应用层的内存（应用层可以水平扩展无状态）
* 如果使用了分页limit时，如果limit N,M，前面的N条数据会先被join但这数据又不用消耗的时间白白浪费
* 就算被驱动表有索引，由于是树的数据结构，每行的时间复杂度是O(logN)，尽可能改为hash(8.0才支持)，通过应用层实现

## 使用join是性能较差排查思路
* 是否驱动表选择错误: straight_join
* 是否使用了错的索引: force index
* 是否回表获取了不必要的数据且`行数据较大`: 修改索引(最左匹配原则，覆盖索引)