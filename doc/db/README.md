创建账号及权限设置
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