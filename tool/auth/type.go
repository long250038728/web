package auth

import (
	"context"
	"errors"
	"github.com/golang-jwt/jwt"
)

type jwtClaims struct {
	jwt.StandardClaims
	*UserClaims
}

// UserClaims 外部使用的信息
type UserClaims struct {
	Id        int32
	Name      string
	Other     map[string]string
	AuthToken string
}

// UserSession 内部使用的信息
type UserSession struct {
	AuthList []string
}

type Auth interface {
	// Parse 生成accessToken
	Parse(ctx context.Context, accessToken string) (*UserClaims, *UserSession, error)
	// Auth 判断是否有权限
	Auth(ctx context.Context, path string) error
	// Set 生成accessToken  refreshToken
	Set(ctx context.Context, userClaims *UserClaims, userSession *UserSession) (string, error)
}

type claims struct{}
type session struct{}

func SetClaims(ctx context.Context, userClaims *UserClaims) context.Context {
	return context.WithValue(ctx, claims{}, userClaims)
}
func GetClaims(ctx context.Context) (*UserClaims, error) {
	if val, ok := ctx.Value(claims{}).(*UserClaims); ok {
		return val, nil
	}
	return nil, errors.New("claims is null")
}

func SetSession(ctx context.Context, userSession *UserSession) context.Context {
	return context.WithValue(ctx, session{}, userSession)
}
func GetSession(ctx context.Context) (*UserSession, error) {
	if val, ok := ctx.Value(session{}).(*UserSession); ok {
		return val, nil
	}
	return nil, errors.New("session is null")
}
