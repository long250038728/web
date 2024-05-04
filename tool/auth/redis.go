package auth

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/long250038728/web/tool/cache"
	"strings"
)

//通过token字符串获取授权信息

//根据传入的地址及参数与授权信息继续匹配

//匹配到则代表有权限

type cacheAuth struct {
	cache     cache.Cache
	whiteList []string
}

type Opt func(r *cacheAuth)

func WhiteList(list []string) Opt {
	return func(r *cacheAuth) {
		r.whiteList = list
	}
}

func NewCacheAuth(cache cache.Cache, opts ...Opt) Auth {
	r := &cacheAuth{
		cache: cache,
	}
	for _, opt := range opts {
		opt(r)
	}

	if r.whiteList == nil {
		r.whiteList = make([]string, 0, 0)
	}

	return r
}

// Set 用户内部信息生产token
func (p *cacheAuth) Set(ctx context.Context, userToken *TokenInfo, token string) error {
	if len(token) == 0 {
		return errors.New("token str is err")
	}

	b, err := json.Marshal(userToken)
	if err != nil {
		return err
	}

	ok, err := p.cache.Set(ctx, token, string(b))
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("auth to token is err")
	}
	return nil
}

// Auth 判断是否有权限
func (p *cacheAuth) Auth(ctx context.Context, userClaims *UserClaims, path string) error {
	//转换为小写
	path = strings.ToLower(path)

	//白名单
	if p.whitePath(path) {
		return nil
	}

	//检查authToken
	if userClaims.AuthToken == "" {
		return errors.New("token is empty")
	}

	token, err := p.cache.Get(ctx, userClaims.AuthToken)
	if err != nil {
		return err
	}

	var UserToken TokenInfo
	err = json.Unmarshal([]byte(token), &UserToken)
	if err != nil {
		return err
	}

	//匹配
	for _, authPath := range UserToken.AuthList {
		if authPath == path {
			return nil
		}
	}

	return errors.New("no match path")
}

// whitePath path是否为白名单
func (p *cacheAuth) whitePath(path string) bool {
	for _, whitePath := range p.whiteList {
		if whitePath == path {
			return true
		}
	}
	return false
}
