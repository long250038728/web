package auth

import (
	"context"
	"errors"
)

// UserClaims 返回给客户的token信息，内部使用的放在other里面的auth_token
type UserClaims struct {
	ClaimsId  int32
	Name      string
	Other     map[string]string
	AuthToken string
}

// TokenInfo 内部使用的信息
type TokenInfo struct {
	AuthList []string
}

type Auth interface {
	Set(ctx context.Context, userToken *TokenInfo, token string) error
	Auth(ctx context.Context, userClaims *UserClaims, path string) error
}

type claims struct{}

func SetClaims(ctx context.Context, userClaims *UserClaims) context.Context {
	return context.WithValue(ctx, claims{}, userClaims)
}
func GetClaims(ctx context.Context) (*UserClaims, error) {
	if val, ok := ctx.Value(claims{}).(*UserClaims); ok {
		return val, nil
	}
	return nil, errors.New("claims is null")
}
