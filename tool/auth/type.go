package auth

import (
	"context"
)

type Auth interface {
	Set(ctx context.Context, userToken *UserToken, token string) error
	Auth(ctx context.Context, userClaims *UserClaims, path string) error
}

// UserClaims 返回给客户的token信息，内部使用的放在other里面的auth_token
type UserClaims struct {
	ClaimsId int32
	Name     string
	Other    map[string]string
}

func (u *UserClaims) SetAuthToken(token string) {
	if u.Other == nil {
		u.Other = make(map[string]string, 1024)
	}
	u.Other["auth_token"] = token
}

func (u *UserClaims) AuthToken() string {
	if u.Other == nil {
		u.Other = make(map[string]string, 1024)
	}
	if token, ok := u.Other["auth_token"]; ok {
		return token
	}
	return ""
}

// UserToken 内部使用的信息
type UserToken struct {
	AuthList []string
}
