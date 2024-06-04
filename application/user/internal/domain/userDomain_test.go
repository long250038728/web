package domain

import (
	"context"
	"github.com/long250038728/web/application/user/internal/repository"
	"github.com/long250038728/web/protoc"
	"github.com/long250038728/web/protoc/user"
	"github.com/long250038728/web/tool/app"
	"github.com/long250038728/web/tool/auth"
	"testing"
)

func TestUserDomain_Login(t *testing.T) {
	app.InitPathInfo("/Users/linlong/Desktop/web/config", protoc.UserService)
	dm := NewUserDomain(repository.NewUserRepository(app.NewUtil()))

	login, err := dm.Login(context.Background(), &user.LoginRequest{
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
	app.InitPathInfo("/Users/linlong/Desktop/web/config", protoc.UserService)
	dm := NewUserDomain(repository.NewUserRepository(app.NewUtil()))

	login, err := dm.Refresh(context.Background(), &user.RefreshRequest{
		RefreshToken: "",
	})
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(login)
}

func TestParse(t *testing.T) {
	app.InitPathInfo("/Users/linlong/Desktop/web/config", protoc.UserService)
	var accessToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTU5NDg4NzAsImlhdCI6MTcxNTg0MDg3MCwiaWQiOjEyMzQ1NiwibmFtZSI6ImpvaG4ifQ.vk7CR288G1s5a8ky5gV2iUtmbzxyz1LYRT5eJSIpnqE"
	ctx, err := (auth.NewCacheAuth(app.NewUtil().Cache())).Parse(context.Background(), accessToken)
	if err != nil {
		t.Error(err)
		return
	}

	dm := NewUserDomain(repository.NewUserRepository(app.NewUtil()))
	login, err := dm.Refresh(ctx, &user.RefreshRequest{
		RefreshToken: "",
	})
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(login)
}
