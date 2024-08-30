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