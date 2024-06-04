### 登录端
1.获取员工信息
2.根据员工获取角色
3.根据角色获取权限菜单


```sql
-- 创建员工表
CREATE TABLE `project`.`admin_user`  (
`id` int NOT NULL AUTO_INCREMENT,
`name` varchar(255) NULL,
`password` varchar(255) NULL,
`telephone` varchar(11) NULL,
`status` tinyint NULL DEFAULT 0 COMMENT '1:正常  0:删除',
`create_time` int NULL,
`delete_time` int NULL,
PRIMARY KEY (`id`)
);

-- 创建角色表
CREATE TABLE `project`.`admin_role`  (
`id` int NOT NULL AUTO_INCREMENT,
`name` varchar(255) NULL,
`status` tinyint NULL DEFAULT 0 COMMENT '1:正常  0:删除',
`create_time` int NULL,
`delete_time` int NULL,
PRIMARY KEY (`id`)
);

-- 创建权限菜单表
CREATE TABLE `project`.`admin_permission`  (
`id` int NOT NULL AUTO_INCREMENT,
`name` varchar(255) NULL,
`status` tinyint NULL DEFAULT 0 COMMENT '1:正常  0:删除',
`path` varchar(255) COMMENT '权限路径',
`create_time` int NULL,
`delete_time` int NULL,
PRIMARY KEY (`id`)                                        
);

-- 创建员工角色关联表
CREATE TABLE `project`.`admin_user_role`  (
`id` int NOT NULL AUTO_INCREMENT,
`user_id` int NULL,
`role_id` int NULL,
PRIMARY KEY (`id`)
);

-- 创建角色权限关联表
CREATE TABLE `project`.`admin_role_permission`  (
`id` int NOT NULL AUTO_INCREMENT,
`role_id` int NULL,
`permission_id` int NULL,
PRIMARY KEY (`id`)
);
```

创建root账号
```sql
INSERT INTO `project`.`admin_user` (`id`, `name`, `password`, `telephone`, `status`, `create_time`, `delete_time`) VALUES (1, 'root', '123456', NULL, 1, UNIX_TIMESTAMP(), NULL);
INSERT INTO `project`.`admin_role` (`id`, `name`, `status`, `create_time`, `delete_time`) VALUES (1, '系统管理员', 1, UNIX_TIMESTAMP(), NULL);
INSERT INTO `project`.`admin_permission` (`id`, `name`, `status`, `path`, `create_time`, `delete_time`) VALUES (1, '客户信息', 1, '/user', UNIX_TIMESTAMP(), NULL);

INSERT INTO `project`.`admin_user_role` (`user_id`, `role_id`) VALUES ( 1, 1);
INSERT INTO `project`.`admin_role_permission` (`role_id`, `permission_id`) VALUES ( 1, 1);
```
