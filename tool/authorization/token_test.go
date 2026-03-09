package authorization

import (
	"context"
	"github.com/golang-jwt/jwt/v4"
	"testing"
	"time"
)

//========================================== claims ===========================================

// AccessClaims 带有jwt的UserInfo
type AccessClaims struct {
	jwt.RegisteredClaims
	*UserInfo
}

func (a AccessClaims) GetSessionId() int32 {
	return a.UserInfo.Id
}

// RefreshClaims 带有jwt的Refresh
type RefreshClaims struct {
	jwt.RegisteredClaims
	*Refresh
}

func (a RefreshClaims) GetSessionId() int32 {
	return a.Refresh.Id
}

// UserInfo 外部使用的信息
type UserInfo struct {
	Id   int32  `json:"id" yaml:"id"`
	Name string `json:"name" yaml:"name"`
}

type Refresh struct {
	Id  int32  `json:"id" yaml:"id"`
	Md5 string `json:"md5" yaml:"md5"`
}

func Test_token_Signed(t *testing.T) {
	ctx := context.Background()
	token := NewToken()

	now := time.Now().Local()
	userClaims := &UserInfo{Id: 1, Name: "test"}

	// 生成token && token转换为claims
	accessClaims := &AccessClaims{RegisteredClaims: jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour)), IssuedAt: jwt.NewNumericDate(now)}, UserInfo: userClaims}
	accessToken, _ := token.Signed(ctx, accessClaims)

	newAccessClaims := &AccessClaims{}
	_ = token.Parse(ctx, accessToken, newAccessClaims)
	t.Log(newAccessClaims.UserInfo)

	// 生成token && token转换为claims
	refreshClaims := &RefreshClaims{RegisteredClaims: jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour)), IssuedAt: jwt.NewNumericDate(now)}, Refresh: &Refresh{Id: userClaims.Id, Md5: GetSessionKey(userClaims.Id)}}
	refreshToken, _ := token.Signed(ctx, refreshClaims)

	newRefreshClaims := &RefreshClaims{}
	_ = token.Parse(ctx, refreshToken, newRefreshClaims)
	t.Log(newRefreshClaims.Refresh)
}
