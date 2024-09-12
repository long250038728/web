package session

import (
	"context"
	"github.com/long250038728/web/tool/authorization"
	"github.com/long250038728/web/tool/cache"
	"github.com/long250038728/web/tool/configurator"
	"testing"
)

var c authorization.Store

func init() {
	var redisConfig cache.Config
	configurator.NewYaml().MustLoadConfigPath("redis.yaml", &redisConfig)
	c = cache.NewRedisCache(&redisConfig)
}

func TestSigned(t *testing.T) {
	t.Log(NewAuth(c).Signed(context.Background(),
		&UserInfo{Id: 123456, Name: "john"},
	))
}

func TestParse(t *testing.T) {
	var accessToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjQ0MjU5MDAsImlhdCI6MTcyNDMxNzkwMCwiaWQiOjEyMzQ1NiwibmFtZSI6ImpvaG4ifQ.SlNUIsPcFweo9Abrmr4R_lR7I_GwV1zNDZtfgbKKgwU"
	ctx, err := NewAuth(c).Parse(context.Background(), accessToken)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(GetClaims(ctx))
	t.Log(GetSession(ctx))
}
