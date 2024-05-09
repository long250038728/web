package auth

import (
	"context"
	"github.com/long250038728/web/tool/cache"
	"github.com/long250038728/web/tool/configurator"
	"testing"
)

var c cache.Cache

func init() {
	var conf cache.Config
	configurator.NewYaml().MustLoad("/Users/linlong/Desktop/web/config/redis.yaml", &conf)
	c = cache.NewRedisCache(&conf)
}

func TestSet(t *testing.T) {
	auth := NewCacheAuth(c)
	t.Log(auth.Set(context.Background(),
		&UserClaims{Id: 123456, Name: "john", Other: map[string]string{"size": "11111"}},
		&UserSession{AuthList: []string{"123", "456", "789"}},
	))
}

func TestParse(t *testing.T) {
	var accessToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTUyNjc0MTksImlhdCI6MTcxNTE1OTQxOSwiSWQiOjEyMzQ1NiwiTmFtZSI6ImxpbmwiLCJPdGhlciI6eyJzaXplIjoiMTExMTEifSwiQXV0aFRva2VuIjoiMTExNzgzN2IxYjFmZjExZmExN2YwZDhhNjIwOWI1M2ZlNTk2Mjk5MWZhZmQxMDczM2MxY2NkYTY0ZTg2ZTEwZSJ9.R4LqDkWvMcHHJFqg8FkiK2_8Ye0Lk01behcxFYrgCZs"
	auth := NewCacheAuth(c)
	t.Log(auth.Parse(context.Background(), accessToken))
}
