package repository

import (
	"context"
	"errors"
	"github.com/long250038728/web/application/auth/internal/model"
	"github.com/long250038728/web/protoc/auth"
	"github.com/long250038728/web/tool/app"
	"github.com/long250038728/web/tool/authorization"
	"github.com/long250038728/web/tool/sliceconv"
	"time"
)

type Auth struct {
	util *app.Util
}

func NewAuthRepository(util *app.Util) *Auth {
	return &Auth{util: util}
}

func (repository *Auth) Login(ctx context.Context, name, password string) (*auth.UserResponse, error) {
	userInfo, err := repository.GetUser(ctx, 0, name, password)
	if err != nil {
		return nil, err
	}
	return repository.getUserResponse(ctx, userInfo)
}

func (repository *Auth) Refresh(ctx context.Context, refreshToken string) (*auth.UserResponse, error) {
	refreshCla := &authorization.RefreshClaims{}

	auth, err := repository.util.Auth()
	if err != nil {
		return nil, err
	}
	if err = auth.Refresh(ctx, refreshToken, refreshCla); err != nil {
		return nil, err
	}

	refresh := refreshCla.Refresh
	if refresh.Md5 != authorization.GetSessionId(refresh.Id) {
		return nil, errors.New("refresh token error")
	}
	userInfo, err := repository.GetUser(ctx, refresh.Id, "", "")
	if err != nil {
		return nil, err
	}
	resp, err := repository.getUserResponse(ctx, userInfo) //生成新的accessToken及refreshToken

	//如果refreshToken的有效期大于24小时，则返回之前refreshToken，否则返回新的refreshToken
	if refreshCla.ExpiresAt.Time.Unix()-time.Now().Local().Unix() >= 60*60*24 {
		resp.RefreshToken = refreshToken
	}
	return resp, err
}

func (repository *Auth) Logout(ctx context.Context) error {
	claims, err := authorization.GetClaims(ctx)
	if err == nil {
		return err
	}
	auth, err := repository.util.Auth()
	if err != nil {
		return err
	}
	return auth.DeleteSession(ctx, authorization.GetSessionId(claims.Id))
}

//======================================================================================================================

func (repository *Auth) GetUser(ctx context.Context, id int32, name, password string) (*model.User, error) {
	db, err := repository.util.Db(ctx)
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

func (repository *Auth) GetRoles(ctx context.Context, userId int32) ([]*model.Role, error) {
	db, err := repository.util.Db(ctx)
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

func (repository *Auth) GetPermissions(ctx context.Context, roleIds []int32) ([]*model.Permission, error) {
	db, err := repository.util.Db(ctx)
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

// ======================================================================================================================
func (repository *Auth) getUserResponse(ctx context.Context, userInfo *model.User) (*auth.UserResponse, error) {
	//角色
	roles, err := repository.GetRoles(ctx, userInfo.Id)
	if err != nil {
		return nil, err
	}
	roleIds := sliceconv.Extract(roles, func(item *model.Role) int32 { return item.Id })
	roleNames := sliceconv.Extract(roles, func(item *model.Role) string { return item.Name })

	//权限列表
	permissions, err := repository.GetPermissions(ctx, roleIds)
	if err != nil {
		return nil, err
	}
	permissionsPath := sliceconv.Extract(permissions, func(item *model.Permission) string { return item.Path })

	//基本参数
	claims := &authorization.UserInfo{Id: userInfo.Id, Name: userInfo.Name}
	sess := &authorization.UserSession{Id: userInfo.Id, Name: userInfo.Name, AuthList: permissionsPath}

	authorized, err := repository.util.Auth()
	if err != nil {
		return nil, err
	}

	if err = authorized.SetSession(ctx, authorization.GetSessionId(claims.Id), sess); err != nil {
		return nil, err
	}
	accessToken, refreshToken, err := authorized.Signed(ctx, claims)
	if err != nil {
		return nil, err
	}
	return &auth.UserResponse{Id: userInfo.Id, Name: userInfo.Name, Telephone: userInfo.Telephone, Roles: roleNames, Permissions: permissionsPath, AccessToken: accessToken, RefreshToken: refreshToken}, nil
}
