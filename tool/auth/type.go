package auth

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
)

type jwtClaims struct {
	jwt.StandardClaims
	*UserClaims
}

type RefreshClaims struct {
	jwt.StandardClaims
	Id  int32  `json:"id" yaml:"id"`
	Md5 string `json:"md5" yaml:"md5"`
}

// UserClaims 外部使用的信息
type UserClaims struct {
	Id   int32  `json:"id" yaml:"id"`
	Name string `json:"name" yaml:"name"`
}

// UserSession 内部使用的信息
type UserSession struct {
	Id       int32  `json:"id" yaml:"id"`
	Name     string `json:"name" yaml:"name"`
	AuthList []string
}

type Auth interface {
	// Parse 解析accessToken
	// 生成Claims Session存放到ctx中 通过 GetClaims GetSession 获取
	Parse(ctx context.Context, accessToken string) (context.Context, error)

	// Auth 判断是否有权限 判断path是否在GetSession中
	Auth(ctx context.Context, path string) error

	// Signed 生成accessToken refreshToken
	Signed(ctx context.Context, userClaims *UserClaims, userSession *UserSession) (string, string, error)

	Refresh(ctx context.Context, refreshToken string) (*RefreshClaims, error)
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

func authToken(id int32) string {
	hash := sha256.New()
	hash.Write([]byte(fmt.Sprintf("id:%d", id))) // 向哈希计算对象中写入字符串数据
	return hex.EncodeToString(hash.Sum(nil))
}
