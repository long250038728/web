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

// UserClaims 外部使用的信息
type UserClaims struct {
	Id   int32  `json:"id" yaml:"id"`
	Name string `json:"name" yaml:"name"`
}

func (c *UserClaims) AuthToken() string {
	hash := sha256.New()
	hash.Write([]byte(fmt.Sprintf("%d", c.Id))) // 向哈希计算对象中写入字符串数据
	return hex.EncodeToString(hash.Sum(nil))
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
