package domain

import (
	"context"
	"github.com/long250038728/web/tool/app"
	"github.com/long250038728/web/tool/auth/auth"
)

func authCtx() context.Context {
	cache, err := app.NewUtil().Cache()
	if err != nil {
		panic(err)
		return nil
	}
	var accessToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTU5NDg4NzAsImlhdCI6MTcxNTg0MDg3MCwiaWQiOjEyMzQ1NiwibmFtZSI6ImpvaG4ifQ.vk7CR288G1s5a8ky5gV2iUtmbzxyz1LYRT5eJSIpnqE"
	ctx, _ := auth.NewAuth(cache).Parse(context.Background(), accessToken)
	return ctx
}
