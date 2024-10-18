## 创建账号及权限设置
```
-- 创建用户linl用户，密码是long123456, ip指定%
CREATE user 'linl'@'%' IDENTIFIED by "long123456";

-- 验证用户插件
-- MySQL 5.7 默认mysql_native_password  
-- MySQL 8.0 默认caching_sha2_password  
show variables like 'default_authentication_plugin'; 
CREATE user 'linl'@'%' IDENTIFIED  with 'mysql_native_password' by "long123456";
CREATE user 'linl'@'%' IDENTIFIED  with 'caching_sha2_password' by "long123456";

-- 强制登录使用证书  证书包含xxx字段 issuer
CREATE user 'linl'@'%' identified by 'long123456' require ssl;
CREATE user 'linl'@'%' identified by 'long123456' require subject '/CN=MySQL_Server_8.0.32_Auto_Generated_Client_Certificate';
CREATE user 'linl'@'%' identified by 'long123456' require issuer '/helloworld';

-- 查询用户列表
SELECT * from mysql.`user`;

-- 当前的登录账号
select current_user();

-- 账号的权限
show GRANTS;

-- 指定权限 
-- *.*  表示*库中的*表
GRANT CREATE,INDEX,ALTER,DROP,INSERT,UPDATE,DELETE,SELECT  on *.*  to 'linl'@'%';
GRANT all privileges  on *.*  to 'linl'@'%';
```

## 参数设置
```
-- 查看动态参数(会话生效)
show variables like '%'
-- 查看动态参数(全局生效)
show global variables like '%'
-- 查看其他会话的参数(mysql8.0之后版本)
select * from performance_schema.variables_by_thread


-- 设置动态参数(会话生效)
set xxxx
-- 设置动态参数(全局生效，mysql8.0数据库重启会重置)—— 全局不影响已经创建的会话，新会话才会影响
set global xxxx
-- 设置动态参数(mysql8.0数据库重启依旧生效，保存在mysqld-auto.cnf中，重启先加载这个文件配置)
set persist xxxx


-- 查询行数据最多列表
SELECT * FROM information_schema.tables order by table_rows desc  ;
-- 查询自增id最多列表
SELECT * FROM information_schema.tables order by auto_increment desc ;
-- 查询占用硬盘最多列表
SELECT * FROM information_schema.tables order by data_length + index_length + data_free desc ;
-- 当前的环境变量
show global variables where Variable_name in ('innodb_flush_log_at_trx_commit','sync_binlog','binlog_format','character_set_server','system_time_zone')
```

mysql重要的系统参数
* open_files_limit 限制同时打开的文件数量
* max_connections 最大的连接数
* sort_buffer_size 排序buffer的大小(mysql排序有两种模式，一种是只存id在临时表最后回表，一种是全部字段都在临时表)
* join_buffer_size 连表buffer的大小
* read_rnd_buffer_size MRR大小，由于树是正序，如果数据回表可能先把id排序完回表
* tmp_table_size 临时表大小，太小就会swap
* InnoDB Buffer Pool
  * innodb_buffer_pool_size 池大小
  * innodb_page_size 设置页大小
  * innodb_buffer_pool_chunk_size  每个块大小
  * innodb_read_io_threads innodb_write_io_threads 读写io线程
  * innodb_log_buffer_size redo缓存大小
  * innodb_flush_log_at_trx_commit redo刷盘时机

## sql mode
严格模式
* ONLY_FULL_GROUP_BY : SELECT的字段要么也出现GROUP BY中，要么使用聚合函数
* STRICT_TRANS_TABLES:  插入数据时如果类型不一致或超出范围会报错（事务性表）
* STRICT_ALL_TABLES  :  插入数据时如果类型不一致或超出范围会报错 (所有类型的表)
* NO_ZERO_DATE && NO_ZERO_DATE : 对date与datetime类型不允许使用日期0 (需要STRICT_TRANS_TABLES才生效)
* ALLOW_INVALID_DATES : 允许插入不合法的日期会自动变为"00-00-00"(timestamp 无法写入不合法日期)
* ERROR_FOR_DIVISION_BY_ZERO : 除数为0报错
* NO_BACKSLASH_ESCAPES : 反斜杠“\”就变成一个普通的字符
* ANSI_QUOTES : 字符串常量可以使用单引号或双引号来引用 —— 其他数据库``表示相同
* NO_ENGINE_SUBSTITUTION : 建表时如果指定的存储引擎不可用或不存在，SQL 就会报错
* PIPES_AS_CONCAT : 管道符变成连接符 (管道符“||”相当于 OR)
* REAL_AS_FLOAT :  REAL 类型映射为Float类型(不设置Double映射类型)
* IGNORE_SPACE : 多加空格换行符不会报错
* ANSI ： 相当与开启REAL_AS_FLOAT、PIPES_AS_CONCAT、ANSI_QUOTES、IGNORE_SPACE 和 ONLY_FULL_GROUP_BY
* TRADITIONAL ：相当与开启STRICT_TRANS_TABLES、STRICT_ALL_TABLES、NO_ZERO_IN_DATE、NO_ZERO_DATE、ERROR_FOR_DIVISION_BY_ZERO 和 NO_ENGINE_SUBSTITUTION

##  字符集
```
show variables like '%char%';
```
* character_set_system 系统字符集，用于存放元数据
* character_set_server 创建数据库默认字符集
* character_set_database 切换数据库默认字符集
* character_set_client 客户端默认字符集
* character_set_connection 连接字符集（当character_set_connection与character_set_client不一致，会转换为character_set_connection字符集）
* character_set_results 返回客户端结果字符集
* character_set_filesystem 文件名字符集

## 大表DDL
```
alter table salaries modify emp_no int not null comment 'Employee Identity', ALGORITHM=inplace, lock=none;
```
ALGORITHM
* DEFAULT: 默认方式会根据不同的DDL默认选择开销最低的方式

* COPY: 执行期间会进行锁表，Lock为SHARED(共享锁) ——重建表
  * COPY是从server层执行的，所以需要先从存储引擎获取一行数据，再调用存储引擎的接口写入一行数据（性能差，跨层,需要记录 UNDO 日志和 REDO 日志）
  * 整个过程是需要锁表的,因为上面的server层执行产生的 UNDO 日志和 REDO 日志

* INPLACE: 在重建的过程中DML修改数据时记录到在线变更日志中，等到后续这个变更更新到新表时此时会锁表
  * INPLACE DDL 是在 InnoDB 存储引擎的内部复制数据，这个过程中不需要记录 UNDO 日志和 REDO 日志
  * 在内部数据复制是不锁表的，只是等到加载变更日志才锁表(毕竟变更记录的大小比全量的数据要小得多所以锁定时间会小很多)
  * Online DDL 流程
    1. 通过innodb_parallel_read_threads 设置线程数并发读取数据。把数据写入innodb_ddl_buffer_size 缓冲中
    2. 多线程读取的innodb_ddl_buffer_size 合并成一个有序的数据集，并归排序 innodb_ddl_threads 控制线程数
    3. 将有序的数据集写入ibd文件过程中，中间的变更会记录写入在线变更日志buffer中。 innodb_sort_buffer_size 控制缓冲大小 ,innodb_online_alter_log_max_size 变更日志大小 如果变动很大超过这个设置会失败
    4. 有序的数据集写入ibd文件后，在线变更日志加载出来更新（为了避免数据问题此时在这个阶段会进行排他锁处理）
    5. 交换新旧两个的ibd文件名称
  * 会重建表的操作
    1. 主键、修改主键
    2. 修改字段的顺序和 NOT NULL 属性（行的存储格式）
    3. 优化表(optimize table)、修改表的行格式或 key_block_size 属性
  * 不会重建表的操作（元数据修改不影响实际存储）
    1. 添加列到末尾
    2. 添加或删除索引
    3. 重命名表
  * 不支持Online DDL
    1. 删除主键
    2. 修改字符集
    3. 修改字段的类型
    4. varchar长度 如果变动前后都是255字节内或超过255字节没问题，如果变动前后发生255字节的大小变动则无法使用

* INSTANT: 8.0新增的且默认，在添加/删除字段不需要重建表（不能再指定Lock关键字），通过内部version版本号及元数据的方式
  * INSTANT DDL 是在 InnoDB 存储引擎的内部复制数据，这个过程中不需要记录 UNDO 日志和 REDO 日志
  * 不支持 INSTANT DDL
    1. 修改类型
    2. 是否为null
    3. 删除枚举值
    4. varchar长度 如果变动前后都是255字节内或超过255字节没问题，如果变动前后发生255字节的大小变动则无法使用
  * 注意: 同一个表指定instant操作次数(64)是有限制的，超过会报错，需要使用COPY或INPLACE重建表(INSTANT是通过版本号控制)

Lock
* NONE: 无锁
* SHARED: 共享锁
* EXCLUSIVE: 排他锁


Metadata(元数据修改不影响实际存储不会导致大DDL)
1. 修改表名、字段名或索引名
2. 删除表/索引 (DROP TABLE、DROP INDEX 但是如果开启了innodb_file_per_table，DROP 表时需要删除对应的 ibd 文件)
3. 修改表、字段的注释

注意:
* 删除索引时需要有可能会导致有些SQL原先命中的现在无法命中导致慢SQL.此时恢复之前索引就会长时间的锁。在mysql 8.0之后建议设置索引不可见后再删除
* 删除表时需要确定确实没有任何业务会访问这个表了再删除，优先建议先改表名后再删除
* 即使你执行的 DDL 只需要修改元数据，在 DDL 执行开始和执行结束的时候，也是需要短暂地获取元数据锁的，如果数据库中有别的长事务提前获取了元数据锁，那么 DDL 就会被阻塞，而 DDL 被阻塞后，后续其他会话访问同一个表时，也会被阻塞。因此在 DDL 执行的过程中，需要注意观察数据库的整体状况，特别是要注意有没有会话在等待元数据锁。

## 执行计划
执行流程
  1. 根据sql文件进行解析，生成sql语法树
  2. 优化器根据sql语法树，表和索引的结构跟统计信息，生成执行计划
  3. sql执行引擎根据执行计划。按一定的步骤，调用存储引擎接口获取数据，执行表连接、排序等操作，
  4. 生成结果集返回

b+树
  * 数据是由page组成 page大小为16k(由参数innodb_page_size设置)，同时innodb对每行长度有一定限制不可超过page的一半。
  * b+树中root节点及中间节点只存放key值(主键id或索引key)，key值是有序page树排列的。
  * 两个相邻叶子节点是通过链式连接（为了解决区间读取的数据）
  * 主键索引叶子节点存放整行数据，二级索引叶子节点存放主键数据（目的是为了索引小可快速找到行数据，可通过二级索引快速找到主键，通过一级索引可以快速找到行数据）
  * 索引指导
    * 通过索引的有序性避免需要重新排序
    * 使用覆盖索引避免回表，通过索引下推的提高server层条件判断
    * 表连接时被驱动表最好是主键索引，其次是二级索引，避免全表匹配
    * 当explain时出现merge_index时代表使用了多个索引查询数据后进行合并，应该优化索引
  * 索引失效/性能差的场景: 
    1. 不满足最左匹配原则
    2. 索引key值进行了函数运算
    3. 索引key值被隐射转换(定义为string类型，条件值为int，或连表时表字符集类型不一致(建议都用utf8mb4))
       * 
       * 索引值由于sting转int后，索引排序是不一致可能会破坏索引的有序性。（"109"跟"11"字符串比较会觉得"11"较大）
    4. 普通索引(唯一索引除外)使用了不支持的运算符(!= , or 等)
    5. mysql8.0之前不支持exists子查询无法进行半连接优化(性能差)
    6. in条件值过多时会导致全表扫描（消耗的内存太多超出了限制，因此放弃使用range执行计划最终使用了全表扫描）
       * 当使用多个in时如a in (1,2) ,b in (4,5) 优化器会改为 (a=1 and b=4) or (a=1 and b=5) or (a=2 and b=4) or (a=2 and b=5)

```
-- 需要特别留意type，key，rows，filtered，extra 这几个字段
explain [format=json] sql
-- 在MySQL 8.0 及更高版本提供的功能。它不仅显示查询的执行计划，还实际执行查询并给出运行时的分析结果
explian ANALYZE sql  

-- 打开优化追踪
SET optimizer_trace="enabled=on";
-- 执行sql
SELECT * FROM sql;
-- 查看SQL执行优化
SELECT * FROM information_schema.OPTIMIZER_TRACE;
```
### 返回字段及含义
* select_type :
  * SIMPLE 表示查询中没有使用任何复杂的子查询或联合(一般是这个)
     * 优化器对子连接可能会进行半连接优化sql,需要满足一下条件
       1. 子查询没有union
       2. 子查询没有having
       3. 子查询没有使用函数（avg，sum等）
       4. 子查询不允许使用limit
       5. 主查询与子查询没有使用STRAIGHT_JOIN（强制指定左表连接右表）
       6. 主查询及子查询不超过表最大连接数量（61个）
  * PRIMARY 表示主查询，即包含其他子查询的查询。(可能被优化器优化为SIMPLE)
  * SUBQUERY 表示一个子查询。(可能被优化器优化为SIMPLE)
  * DERIVED 表示派生表，通常是从子查询中生成的临时表。(可能被优化器优化为SIMPLE)
  * UNION 表示联合查询的结果
  * UNION RESULT 联合查询的最终结果集
* table : 指定的表名或别名
* partitions : 表分区
* type :
  * const: 唯一索引 等值匹配
  * eq_ref: 唯一索引 (匹配索引字段的值来自驱动表，不是固定的常量) 等值匹配
  * ref: 普通索引字段
  * ref_or_null: 索引字段的条件使用了 or 或 in  空值null
  * range: 索引字段上的范围条件查询数据
  * index: 查找索引中的每一行数据
  * All: 查找表中的每一行数据
  * index_merge: 使用多个索引来查询数据
* possible_keys : 可能会使用到的索引
* key : 最终使用的索引
* key_len : 索引的长度信息 （varchar、char = 长度 * 字符集长度 ， 如果是varchar会额外加2 ，如果字段可以为null则再加1）
* ref: 显示所以查找的值
  * const 常量匹配
  * 表.xx 使用驱动表的某个字段匹配
  * func 使用某个函数计算结果匹配
* rows : 预估扫描的行数
* filtered : 结果行占扫描行的比例
* extra : 其他额外信息
  * Using index 使用了覆盖索引 (无需回表)
  * Using temporary 使用临时表
  * Using index for skip scan 查询条件没有传入索引的前缀字段，又用到了覆盖索引时
  * Using index condition 使用索引下推
  * Using filesort 使用文件排序
  * Using join buffer (xxxx) ———— 优化器一般会使用BNL/hash，BKA
    * Nested-Loop Join: (NLJ)性能最差，对每一行外表的记录都去表匹配查找  (低版本就有)
    * Block Nested-Loop Join (BNL): Nested-Loop Join 的优化版本，它会把外表的数据块存储在内存中，并在内存中逐块处理内表的查找  (低版本就有)
    * Hash Join：以哈希表的形式存储，并根据哈希表快速匹配另一张表的数据   (8.0版本之后，用于替换BNL算法，hash从O(N)改为O(1))
    * Index Nested-Loop Join (INLJ): 每当从外表中取得一行记录时，直接利用索引在内表中进行查找  (低版本就有)
    * Batched Key Access (BKA): 对 INLJ 的进一步优化.收集一批键（即多条记录），内表中进行查找 (低版本就有)
    * Sort-Merge Join (SMJ) 通过排序两个结果集有序后合并  (8.0版本之后)
  * Using MRR 减少回表查询数据时随机 IO，对主键id进行排序回表

## binlog
```
-- 查看binlog 信息
show master logs;
show master status;
show binlog events;

-- 查看binlog_format的格式
SHOW VARIABLES LIKE 'binlog_format';

-- --verbose 查询binlog中执行内容
mysqlbinlog --verbose  ./binlog.000011
```

## 常见命令
```
-- 整理表的碎片
OPTIMIZE TABLE table_name;
-- 采集数据重新刷新
ANALYZE TABLE;
```

top 
mpstat  -P ALL 3

free -m
cat /proc/meminfo
ps aux | head -1; ps aux | sort -nr -k +6 | head
vmstat 3

iostat
iotop

