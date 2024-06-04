package model

type User struct {
	CreateTime int32  `json:"create_time" yaml:"create_time" form:"create_time"`
	DeleteTime int32  `json:"delete_time" yaml:"delete_time" form:"delete_time"`
	Id         int32  `json:"id" yaml:"id" form:"id"`
	Name       string `json:"name" yaml:"name" form:"name"`
	Password   string `json:"password" yaml:"password" form:"password"`
	Status     int32  `json:"status" yaml:"status" form:"status"` // 1:正常  0:删除
	Telephone  string `json:"telephone" yaml:"telephone" form:"telephone"`
}

type Permission struct {
	CreateTime int32  `json:"create_time" yaml:"create_time" form:"create_time"`
	DeleteTime int32  `json:"delete_time" yaml:"delete_time" form:"delete_time"`
	Id         int32  `json:"id" yaml:"id" form:"id"`
	Name       string `json:"name" yaml:"name" form:"name"`
	Path       string `json:"path" yaml:"path" form:"path"`       // 权限路径
	Status     int32  `json:"status" yaml:"status" form:"status"` // 1:正常  0:删除

}

type Role struct {
	CreateTime int32  `json:"create_time" yaml:"create_time" form:"create_time"`
	DeleteTime int32  `json:"delete_time" yaml:"delete_time" form:"delete_time"`
	Id         int32  `json:"id" yaml:"id" form:"id"`
	Name       string `json:"name" yaml:"name" form:"name"`
	Status     int32  `json:"status" yaml:"status" form:"status"` // 1:正常  0:删除
	
}

type RolePermission struct {
	Id           int32 `json:"id" yaml:"id" form:"id"`
	PermissionId int32 `json:"permission_id" yaml:"permission_id" form:"permission_id"`
	RoleId       int32 `json:"role_id" yaml:"role_id" form:"role_id"`
}

type UserRole struct {
	Id     int32 `json:"id" yaml:"id" form:"id"`
	RoleId int32 `json:"role_id" yaml:"role_id" form:"role_id"`
	UserId int32 `json:"user_id" yaml:"user_id" form:"user_id"`
}
