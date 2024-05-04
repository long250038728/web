package auth

import (
	"context"
	"github.com/long250038728/web/tool/cache"
	"github.com/long250038728/web/tool/configurator"
	"testing"
)

var conf cache.Config
var c cache.Cache

func init() {
	configurator.NewYaml().MustLoad("/Users/linlong/Desktop/web/application/user/config/redis.yaml", &conf)
	c = cache.NewRedisCache(&conf)
}

var authToken = "12345678910"
var at = &TokenInfo{AuthList: []string{"/ok", "/true", "/1"}}
var whiteList = []string{"/"}

func TestRedis_Set(t *testing.T) {
	type fields struct {
		cache     cache.Cache
		whiteList []string
	}
	type args struct {
		ctx       context.Context
		userToken *TokenInfo
		token     string
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "auth is ok",
			fields:  fields{cache: c},
			args:    args{ctx: context.Background(), userToken: at, token: authToken},
			wantErr: false,
		},
		{
			name:    "token is empty",
			fields:  fields{cache: c},
			args:    args{ctx: context.Background(), userToken: at, token: ""},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewCacheAuth(tt.fields.cache)
			if err := p.Set(tt.args.ctx, tt.args.userToken, tt.args.token); (err != nil) != tt.wantErr {
				t.Errorf("Set() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRedis_Auth(t *testing.T) {
	type fields struct {
		cache     cache.Cache
		whiteList []string
	}
	type args struct {
		ctx        context.Context
		userClaims *UserClaims
		path       string
	}

	u := &UserClaims{AuthToken: authToken}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "white list",
			fields:  fields{cache: c, whiteList: whiteList},
			args:    args{ctx: context.Background(), userClaims: u, path: "/"},
			wantErr: false,
		},
		{
			name:    "auth is ok",
			fields:  fields{cache: c, whiteList: whiteList},
			args:    args{ctx: context.Background(), userClaims: u, path: "/ok"},
			wantErr: false,
		},
		{
			name:    "auth is ok path to Lower",
			fields:  fields{cache: c, whiteList: whiteList},
			args:    args{ctx: context.Background(), userClaims: u, path: "/TRUE"},
			wantErr: false,
		},
		{
			name:    "auth is not ok",
			fields:  fields{cache: c, whiteList: whiteList},
			args:    args{ctx: context.Background(), userClaims: u, path: "/not ok"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &cacheAuth{
				cache:     tt.fields.cache,
				whiteList: tt.fields.whiteList,
			}
			if err := p.Auth(tt.args.ctx, tt.args.userClaims, tt.args.path); (err != nil) != tt.wantErr {
				t.Errorf("Auth() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
