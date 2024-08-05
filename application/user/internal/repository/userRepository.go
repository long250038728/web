package repository

import (
	"context"
	"errors"
	"github.com/long250038728/web/application/user/internal/model"
	"github.com/long250038728/web/protoc/user"
	"github.com/long250038728/web/tool/app"
	"github.com/long250038728/web/tool/auth"
	"github.com/long250038728/web/tool/server/http"
	"github.com/olivere/elastic/v7"
)

type UserRepository struct {
	util *app.Util
}

func NewUserRepository(util *app.Util) *UserRepository {
	return &UserRepository{
		util: util,
	}
}

func (r *UserRepository) Login(ctx context.Context, name, password string) (*user.UserResponse, error) {
	userInfo, err := r.GetUser(ctx, 0, name, password)
	if err != nil {
		return nil, err
	}
	return r.getUserResponse(ctx, userInfo)
}

func (r *UserRepository) Refresh(ctx context.Context, refreshToken string) (*user.UserResponse, error) {
	refreshCla, err := auth.NewAuth(r.util.Cache()).Refresh(ctx, refreshToken)
	if err != nil {
		return nil, err
	}
	userInfo, err := r.GetUser(ctx, refreshCla.Id, "", "")
	if err != nil {
		return nil, err
	}
	resp, err := r.getUserResponse(ctx, userInfo) //生成新的accessToken及refreshToken
	return resp, err
}

func (r *UserRepository) getUserResponse(ctx context.Context, userInfo *model.User) (*user.UserResponse, error) {
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
	claims := &auth.UserClaims{Id: userInfo.Id, Name: userInfo.Name}
	session := &auth.UserSession{Id: userInfo.Id, Name: userInfo.Name, AuthList: permissionsPath}
	accessToken, refreshToken, err := auth.NewAuth(r.util.Cache()).Signed(ctx, claims, session)
	if err != nil {
		return nil, err
	}
	return &user.UserResponse{Id: userInfo.Id, Name: userInfo.Name, Telephone: userInfo.Telephone, Roles: roleNames, Permissions: permissionsPath, AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

//======================================================================================================================

func (r *UserRepository) GetUser(ctx context.Context, id int32, name, password string) (*model.User, error) {
	if id == 0 && len(name) == 0 && len(password) == 0 {
		return nil, errors.New("params is empty")
	}

	userInfo := &model.User{}
	db := r.util.Db(ctx).Where("status = 1")
	if id > 0 {
		db = db.Where("id = ?", id)
	}
	if len(name) > 0 && len(password) > 0 {
		db = db.Where("name = ?", name).Where("password = ?", password)
	}
	return userInfo, db.Find(userInfo).Error
}

func (r *UserRepository) GetRoles(ctx context.Context, userId int32) ([]*model.Role, error) {
	var ids []int32
	if err := r.util.Db(ctx).Model(model.UserRole{}).Select("role_id").Where("user_id = ?", userId).Find(&ids).Error; err != nil {
		return nil, err
	}
	var roles []*model.Role
	return roles, r.util.Db(ctx).Where("id in ?", ids).Where("status = 1").Find(&roles).Error
}

func (r *UserRepository) GetPermissions(ctx context.Context, roleIds []int32) ([]*model.Permission, error) {
	var ids []int32
	if err := r.util.Db(ctx).Model(model.RolePermission{}).Select("permission_id").Where("role_id = ?", roleIds).Find(&ids).Error; err != nil {
		return nil, err
	}
	var permissions []*model.Permission
	return permissions, r.util.Db(ctx).Where("id in ?", ids).Where("status = 1").Find(&permissions).Error
}

func (r *UserRepository) GetName(ctx context.Context, request *user.RequestHello) (string, error) {
	c := &model.User{}
	//orm
	r.util.Db(ctx).Select("name").Where("id = ?", 1).Find(c)

	////mq
	//_ = r.util.Mq.Send(ctx, "aaa", "", &mq.Message{Data: []byte("hello")})

	////cache
	//_, _ = r.util.Cache.Set(ctx, "hello", "1")
	//_, _ = r.util.Cache.Get(ctx, "hello")

	////lock
	//lock, err := r.util.Locker("hello", "123", time.Second*5)
	//if err != nil {
	//	return "", err
	//}
	//_ = lock.Lock(ctx)
	//_ = lock.UnLock(ctx)

	//es
	query := elastic.NewBoolQuery().Must(
		elastic.NewTermQuery("merchant_id", 240),
		elastic.NewTermQuery("merchant_shop_id", 867),
		elastic.NewRangeQuery("gold_weight").Gte(0).Lte(10000),
		elastic.NewMatchQuery("admin_user_name", "小刘"),
		elastic.NewMatchPhraseQuery("merchant_shop_name", "大"),
	)
	_, _ = r.util.Es().Search("sale_order_record_report").Query(query).From(0).Size(100).Do(ctx)

	_, _, _ = http.NewClient().Get(ctx, "http://test.zhubaoe.cn:8888/report/sale_report/inventory", map[string]any{
		"merchant_id":      394,
		"merchant_shop_id": 1150,
		"start_date":       "2023-12-01",
		"end_date":         "2023-12-01",
		"field":            "goods_type_id",
		"client_name":      "app",
	})
	_, _, _ = http.NewClient().Post(ctx, "http://test.zhubaoe.cn:9999/", map[string]any{
		"a": "Login",
		"m": "Admin",
		"p": "1",
		"r": "{\"merchant_code\":\"ab190735\",\"user_name\":\"yzt\",\"password\":\"123456\",\"last_admin_id\":\"\",\"last_admin_name\":\"\",\"shift_status\":\"1\"}",
		"t": "00000",
		"v": "2.4.4",
	})

	return "hello:" + request.Name + " " + c.Name, nil
}
