package domain

import (
	"context"
	"github.com/long250038728/web/tool/app"
	"github.com/long250038728/web/tool/auth"
)

func authCtx() context.Context {
	var accessToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTU5NDg4NzAsImlhdCI6MTcxNTg0MDg3MCwiaWQiOjEyMzQ1NiwibmFtZSI6ImpvaG4ifQ.vk7CR288G1s5a8ky5gV2iUtmbzxyz1LYRT5eJSIpnqE"
	ctx, _ := (auth.NewCacheAuth(app.NewUtil().Cache())).Parse(context.Background(), accessToken)
	return ctx
}
