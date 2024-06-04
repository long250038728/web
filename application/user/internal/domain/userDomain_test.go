package domain

import (
	"github.com/long250038728/web/application/user/internal/repository"
	"github.com/long250038728/web/protoc"
	"github.com/long250038728/web/protoc/user"
	"github.com/long250038728/web/tool/app"
	"testing"
)

func init() {
	app.InitPathInfo("/Users/linlong/Desktop/web/config", protoc.UserService)
}

func TestUserDomain_Login(t *testing.T) {
	var dm = NewUserDomain(repository.NewUserRepository(app.NewUtil()))
	login, err := dm.Login(authCtx(), &user.LoginRequest{
		Name:     "root",
		Password: "123456",
	})
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(login)
}

func TestUserDomain_Refresh(t *testing.T) {
	var dm = NewUserDomain(repository.NewUserRepository(app.NewUtil()))
	login, err := dm.Refresh(authCtx(), &user.RefreshRequest{
		RefreshToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTc1ODkzNjMsImlhdCI6MTcxNzQ4MTM2MywiaWQiOjEsIm1kNSI6ImlkOjZiODZiMjczZmYzNGZjZTE5ZDZiODA0ZWZmNWEzZjU3NDdhZGE0ZWFhMjJmMWQ0OWMwMWU1MmRkYjc4NzViNGIifQ.XzABUvbuFbt6D5dCfKwAxEhtSCwUUyuvfpIMwH0KSLM",
	})
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(login)
}
