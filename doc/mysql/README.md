创建账号及权限设置
```
-- 创建用户linl用户，密码是long123456, ip指定%
CREATE user 'linl'@'%' IDENTIFIED by "long123456";

-- 查询用户列表
SELECT * from mysql.`user`;

-- 当前的登录账号
select current_user();

-- 账号的权限
show GRANTS;

-- 指定权限 *.*  表示*库中的*表
GRANT CREATE,INDEX,ALTER,DROP,INSERT,UPDATE,DELETE,SELECT  on *.*  to 'linl'@'%';
GRANT all privileges  on *.*  to 'linl'@'%';
```