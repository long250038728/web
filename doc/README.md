### 登录端
1.获取员工信息
2.根据员工获取角色
3.根据角色获取权限菜单

创建员工表
```sql
CREATE TABLE `project`.`admin`  (
  `id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NULL,
  `password` varchar(255) NULL,
  `telephone` varchar(11) NULL,
  `status` tinyint NULL DEFAULT 0 COMMENT '1:正常  0:删除',
  `role` int NULL,
  `create_time` int NULL,
  `delete_time` int NULL,
  PRIMARY KEY (`id`)
);
```

创建角色表
```sql
CREATE TABLE `project`.`role`  (
  `id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NULL,
  `status` tinyint NULL DEFAULT 0 COMMENT '1:正常  0:删除',
  `rule` int NULL COMMENT '规则列表',
  `create_time` int NULL,
  `delete_time` int NULL,
  PRIMARY KEY (`id`)
);
```

创建权限菜单表
```sql
CREATE TABLE `project`.`rule`  (
  `id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NULL,
  `status` tinyint NULL DEFAULT 0 COMMENT '1:正常  0:删除',
  `path` varchar(255) COMMENT '权限路径',
  `create_time` int NULL,
  `delete_time` int NULL,
  PRIMARY KEY (`id`)
);
```

创建root账号
```sql
INSERT INTO `project`.`admin` (`id`, `name`, `password`, `telephone`, `status`, `role`, `create_time`, `delete_time`) VALUES (1, 'root', '123456', NULL, 1, 1, NULL, NULL);
INSERT INTO `project`.`role` (`id`, `name`, `status`, `rule`, `create_time`, `delete_time`) VALUES (1, '系统管理员', 1, 1, NULL, NULL);
INSERT INTO `project`.`rule` (`id`, `name`, `status`, `path`, `create_time`, `delete_time`) VALUES (1, '客户信息', 1, '/user', NULL, NULL);
```
