package domain

import (
	"context"
	"github.com/long250038728/web/application/auth/internal/repository"
	"github.com/long250038728/web/protoc/auth_rpc"
	"github.com/long250038728/web/tool/app"
	"github.com/long250038728/web/tool/authorization/session"
	"testing"
)

func authCtx() context.Context {
	cache, err := app.NewUtil().Cache()
	if err != nil {
		panic(err)
		return nil
	}
	var accessToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTU5NDg4NzAsImlhdCI6MTcxNTg0MDg3MCwiaWQiOjEyMzQ1NiwibmFtZSI6ImpvaG4ifQ.vk7CR288G1s5a8ky5gV2iUtmbzxyz1LYRT5eJSIpnqE"
	ctx, _ := session.NewAuth(cache).Parse(context.Background(), accessToken)
	return ctx
}

func TestUserDomain_Login(t *testing.T) {
	var dm = NewDomain(repository.NewRepository(app.NewUtil()))
	login, err := dm.Login(authCtx(), &auth_rpc.LoginRequest{
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
	var dm = NewDomain(repository.NewRepository(app.NewUtil()))
	login, err := dm.Refresh(authCtx(), &auth_rpc.RefreshRequest{
		RefreshToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTc1ODkzNjMsImlhdCI6MTcxNzQ4MTM2MywiaWQiOjEsIm1kNSI6ImlkOjZiODZiMjczZmYzNGZjZTE5ZDZiODA0ZWZmNWEzZjU3NDdhZGE0ZWFhMjJmMWQ0OWMwMWU1MmRkYjc4NzViNGIifQ.XzABUvbuFbt6D5dCfKwAxEhtSCwUUyuvfpIMwH0KSLM",
	})
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(login)
}
