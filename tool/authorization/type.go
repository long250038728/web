package authorization

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/golang-jwt/jwt"
	"time"
)

type TokenType int32

const (
	AccessToken = iota
	RefreshToken
)

//=====================================================================================

// AccessClaims 带有jwt的UserInfo
type AccessClaims struct {
	jwt.StandardClaims
	*UserInfo
}

// RefreshClaims 带有jwt的Refresh
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

func GetSessionId(id int32) string {
	hash := sha256.New()
	hash.Write([]byte(fmt.Sprintf("id:%d", id))) // 向哈希计算对象中写入字符串数据
	return hex.EncodeToString(hash.Sum(nil))
}

// =====================================================================================

type Store interface {
	Get(ctx context.Context, key string) (string, error)
	SetEX(ctx context.Context, key string, value string, expiration time.Duration) (bool, error)
	Del(ctx context.Context, key ...string) (bool, error)
}

type Claims interface {
	Valid() error
}

/**
注意：
	claims Claims 必须传对象指针如: &AccessClaims{} 不能是AccessClaims{}
*/

type Auth interface {
	// Signed 生成accessToken refreshToken
	Signed(ctx context.Context, userClaims *UserInfo) (accessToken string, refreshToken string, err error)

	// Parse 解析accessToken
	// 生成Claims Session存放到ctx中 通过 GetClaims GetSession 获取
	Parse(ctx context.Context, accessToken string) (context.Context, error)

	Refresh(ctx context.Context, refreshToken string, claims Claims) error

	SetSession(ctx context.Context, sessionId string, session *UserSession) (err error)
	DeleteSession(ctx context.Context, sessionId string) error
}
