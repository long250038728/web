package auth

import (
	"context"
	"github.com/long250038728/web/tool/cache"
	"github.com/long250038728/web/tool/configurator"
	"testing"
)

var c Store

func init() {
	var conf cache.Config
	configurator.NewYaml().MustLoad("/Users/linlong/Desktop/web/config/redis.yaml", &conf)
	c = cache.NewRedisCache(&conf)
}

func TestSet(t *testing.T) {
	auth := NewAuth(c)
	t.Log(auth.Signed(context.Background(),
		&UserClaims{Id: 123456, Name: "john"},
		&UserSession{AuthList: []string{"123", "456", "789"}},
	))
}

func TestParse(t *testing.T) {
	var accessToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTU5NDg4NzAsImlhdCI6MTcxNTg0MDg3MCwiaWQiOjEyMzQ1NiwibmFtZSI6ImpvaG4ifQ.vk7CR288G1s5a8ky5gV2iUtmbzxyz1LYRT5eJSIpnqE"
	auth := NewAuth(c)
	ctx, err := auth.Parse(context.Background(), accessToken)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(GetClaims(ctx))
	t.Log(GetSession(ctx))
}
