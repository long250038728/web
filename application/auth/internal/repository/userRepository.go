package repository

import (
	"context"
	"errors"
	"github.com/long250038728/web/application/auth/internal/model"
	"github.com/long250038728/web/protoc/auth_rpc"
	"github.com/long250038728/web/tool/app"
	"github.com/long250038728/web/tool/authorization"
	"github.com/long250038728/web/tool/authorization/session"
)

type Repository struct {
	util *app.Util
}

func NewRepository(util *app.Util) *Repository {
	return &Repository{
		util: util,
	}
}

func (r *Repository) Login(ctx context.Context, name, password string) (*auth_rpc.UserResponse, error) {
	userInfo, err := r.GetUser(ctx, 0, name, password)
	if err != nil {
		return nil, err
	}
	return r.getUserResponse(ctx, userInfo)
}

func (r *Repository) Refresh(ctx context.Context, refreshToken string) (*auth_rpc.UserResponse, error) {
	cache, err := r.util.Cache()
	if err != nil {
		return nil, err
	}

	refreshCla := session.RefreshClaims{}
	if err = session.NewAuth(cache).Refresh(ctx, refreshToken, refreshCla); err != nil {
		return nil, err
	}
	if err := refreshCla.Valid(); err != nil {
		return nil, err
	}

	refresh := refreshCla.Refresh

	if refresh.Md5 != authorization.GetSessionId(refresh.Id) {
		return nil, errors.New("refresh token error")
	}
	userInfo, err := r.GetUser(ctx, refresh.Id, "", "")
	if err != nil {
		return nil, err
	}
	resp, err := r.getUserResponse(ctx, userInfo) //生成新的accessToken及refreshToken
	return resp, err
}

func (r *Repository) getUserResponse(ctx context.Context, userInfo *model.User) (*auth_rpc.UserResponse, error) {
	cache, err := r.util.Cache()
	if err != nil {
		return nil, err
	}

	//角色
	roles, err := r.GetRoles(ctx, userInfo.Id)
	if err != nil {
		return nil, err
	}
	roleIds := make([]int32, 0, len(roles))
	roleNames := make([]string, 0, len(roles))
	for _, role := range roles {
		roleIds = append(roleIds, role.Id)
		roleNames = append(roleNames, role.Name)
	}

	//权限列表
	permissions, err := r.GetPermissions(ctx, roleIds)
	if err != nil {
		return nil, err
	}
	permissionsPath := make([]string, 0, len(roles))
	for _, permission := range permissions {
		permissionsPath = append(permissionsPath, permission.Path)
	}

	//基本参数
	claims := &session.UserInfo{Id: userInfo.Id, Name: userInfo.Name}
	sess := &session.UserSession{Id: userInfo.Id, Name: userInfo.Name, AuthList: permissionsPath}
	accessToken, refreshToken, err := session.NewAuth(cache).Signed(ctx, claims, sess)
	if err != nil {
		return nil, err
	}
	return &auth_rpc.UserResponse{Id: userInfo.Id, Name: userInfo.Name, Telephone: userInfo.Telephone, Roles: roleNames, Permissions: permissionsPath, AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

//======================================================================================================================

func (r *Repository) GetUser(ctx context.Context, id int32, name, password string) (*model.User, error) {
	if id == 0 && len(name) == 0 && len(password) == 0 {
		return nil, errors.New("params is empty")
	}

	db, err := r.util.Db(ctx)
	if err != nil {
		return nil, err
	}

	userInfo := &model.User{}
	dao := db.Where("status = 1")
	if id > 0 {
		dao = dao.Where("id = ?", id)
	}
	if len(name) > 0 && len(password) > 0 {
		dao = dao.Where("name = ?", name).Where("password = ?", password)
	}
	return userInfo, dao.Find(userInfo).Error
}

func (r *Repository) GetRoles(ctx context.Context, userId int32) ([]*model.Role, error) {
	db, err := r.util.Db(ctx)
	if err != nil {
		return nil, err
	}

	var ids []int32
	if err := db.Model(model.UserRole{}).Select("role_id").Where("user_id = ?", userId).Find(&ids).Error; err != nil {
		return nil, err
	}
	var roles []*model.Role
	return roles, db.Where("id in ?", ids).Where("status = 1").Find(&roles).Error
}

func (r *Repository) GetPermissions(ctx context.Context, roleIds []int32) ([]*model.Permission, error) {
	db, err := r.util.Db(ctx)
	if err != nil {
		return nil, err
	}

	var ids []int32
	if err := db.Model(model.RolePermission{}).Select("permission_id").Where("role_id = ?", roleIds).Find(&ids).Error; err != nil {
		return nil, err
	}
	var permissions []*model.Permission
	return permissions, db.Where("id in ?", ids).Where("status = 1").Find(&permissions).Error
}
