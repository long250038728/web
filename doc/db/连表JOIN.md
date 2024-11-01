## Nested Loop Join
取驱动表的每一行，去驱动表关联的表中查询(如果命中索引O(logN),否则是O(N))，然后和驱动表组成一行放到结果集，以此类推，直到最后一个驱动表，最后返回结果。
* SELECT * FROM t2 join t1 on t1.id=t2.id 驱动表不一定是左边的表（根据效率更高的为驱动表），如果强制则使用`STRAIGHT_JOIN`。
* LEFT JOIN时语义上驱动表是左表。大概率是LEFT但不保证一定

## semijoin半连接优化
查询中的in子查询可能会优化为INNER JOIN连表，效率更高(5.7版本)
* SELECT * FROM t1 WHERE t1.id in (SELECT id FROM t2)，如不优化每一行取子表查询判断是否符合

## join_buffer
如果被驱动表没有符合索引，此时需要全表扫描，很大概率会到磁盘里面读取，那可以把一批数据放到join_buffer中减少多次读磁盘的次数

## hash join(8.0版本)
利用join_buffer优化，将数据结构改为hash，此时与该buffer的匹配复杂度就变成O(1)

## 是否使用join或代码拆成几个处理
* 如果被驱动表数据量小时可以用join，如果被驱动表数据量很大，join_buffer放不下，此时会使用临时表，效率会降低
  * 此时把这个join_buffer改为应用层的内存（应用层可以水平扩展无状态）
* 如果使用了分页limit时，如果limit 100,10，前面的100条数据会先被join但这数据又不用消耗的时间白白浪费