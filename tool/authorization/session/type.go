package session

import (
	"context"
	"errors"
	"github.com/golang-jwt/jwt"
	"github.com/long250038728/web/tool/authorization"
)

//=====================================================================================

type AccessClaims struct {
	jwt.StandardClaims
	*UserInfo
}

type RefreshClaims struct {
	jwt.StandardClaims
	*Refresh
}

//=====================================================================================

// UserInfo 外部使用的信息
type UserInfo struct {
	Id   int32  `json:"id" yaml:"id"`
	Name string `json:"name" yaml:"name"`
}

// UserSession 内部使用的信息
type UserSession struct {
	Id       int32  `json:"id" yaml:"id"`
	Name     string `json:"name" yaml:"name"`
	AuthList []string
}

type Refresh struct {
	Id  int32  `json:"id" yaml:"id"`
	Md5 string `json:"md5" yaml:"md5"`
}

// =====================================================================================

type Auth interface {
	// Signed 生成accessToken refreshToken
	Signed(ctx context.Context, userClaims *UserInfo) (accessToken string, refreshToken string, err error)

	// Parse 解析accessToken
	// 生成Claims Session存放到ctx中 通过 GetClaims GetSession 获取
	Parse(ctx context.Context, accessToken string) (context.Context, error)

	Refresh(ctx context.Context, refreshToken string, claims authorization.Claims) error

	SetSession(ctx context.Context, sessionId string, session *UserSession) (err error)
	DeleteSession(ctx context.Context, sessionId string) error

	// Auth 判断是否有权限 判断path是否在GetSession中
	Auth(ctx context.Context, path string) error
}

type claims struct{}
type session struct{}

func SetClaims(ctx context.Context, userClaims *UserInfo) context.Context {
	return context.WithValue(ctx, claims{}, userClaims)
}
func GetClaims(ctx context.Context) (*UserInfo, error) {
	if val, ok := ctx.Value(claims{}).(*UserInfo); ok {
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
