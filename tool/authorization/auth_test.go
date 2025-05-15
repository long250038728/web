package authorization

import (
	"context"
	"github.com/long250038728/web/tool/configurator"
	"github.com/long250038728/web/tool/persistence/cache"
	"github.com/long250038728/web/tool/store"
	"testing"
	"time"
)

var c store.Store

func init() {
	var redisConfig cache.Config
	configurator.NewYaml().MustLoadConfigPath("redis.yaml", &redisConfig)
	c = store.NewRedisStore(cache.NewRedis(&redisConfig))
}

func TestSigned(t *testing.T) {
	access, refresh, _ := NewAuth(c).Signed(context.Background(),
		&UserInfo{Id: 123456, Name: "john"},
	)
	t.Log(access)
	t.Log(refresh)
}

func TestParse(t *testing.T) {
	var accessToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MzI3OTM4NTEsImlhdCI6MTczMjY4NTg1MSwiaWQiOjEyMzQ1NiwibmFtZSI6ImpvaG4ifQ.qOVreMOtfxARUGyrlOJTBI47i1YLx09kWQKL6dDZXfQ"
	ctx, err := NewAuth(c).Parse(context.Background(), accessToken)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(GetClaims(ctx))
	t.Log(GetSession(ctx))
}

func TestRefresh(t *testing.T) {
	refreshToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MzMyOTI1MDcsImlhdCI6MTczMjY4NzcwNywiaWQiOjEyMzQ1NiwibWQ1IjoiMmMxYWNiMWRmZjRhNTk2MmFhYWQzOWE1ZWNjMWJlNzQ3MjQxYjY2NTlhMTIwYzBmZTQzMDRmNDQ3ODQyMTU2ZCJ9.4-kwPx7ASS7V7XWEC6XQsgUyA9_i9WA_AWTdpJMpqEo"
	refreshCla := &RefreshClaims{}
	if err := NewAuth(c).Refresh(context.Background(), refreshToken, refreshCla); err != nil {
		t.Error(err)
		return
	}

	t.Log((refreshCla.ExpiresAt.Time.Unix() - time.Now().Local().Unix()) / 60 / 60 / 24)
	t.Log(refreshCla.ExpiresAt.Time.Unix()-time.Now().Local().Unix() <= 60*60*24)
}
