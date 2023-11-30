package repository

import (
	"context"
	user "github.com/long250038728/web/application/user/protoc"
	"github.com/long250038728/web/tool/app"
)

type UserRepository struct {
	util *app.Util
}

func NewUserRepository(util *app.Util) *UserRepository {
	return &UserRepository{
		util: util,
	}
}

func (r *UserRepository) GetName(ctx context.Context, request *user.RequestHello) string {
	type customer struct {
		Name string `json:"name"`
	}
	c := &customer{}
	r.util.Db.Table("zby_customer").Select("name").Where("id = ?", 1).Find(c)

	return "hello:" + c.Name
}
