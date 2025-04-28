package authorization

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/long250038728/web/tool/store"
	"time"
)

type Opt func(r *auth)

type auth struct {
	Serialization
	Session
	accessExpires  time.Duration
	refreshExpires time.Duration
}

func SecretKey(secretKey []byte) Opt {
	return func(r *auth) {
		r.SecretKey = secretKey
	}
}

func AccessExpires(accessExpires time.Duration) Opt {
	return func(r *auth) {
		r.accessExpires = accessExpires
	}
}

func RefreshExpires(refreshExpires time.Duration) Opt {
	return func(r *auth) {
		r.refreshExpires = refreshExpires
	}
}

func AddStore(s store.Store) Opt {
	return func(r *auth) {
		r.Stores = append(r.Stores, s)
	}
}

func NewAuth(s store.Store, opts ...Opt) Auth {
	p := &auth{}

	//默认值
	p.SecretKey = []byte("secretKey")
	p.accessExpires = 20 * time.Minute
	p.refreshExpires = 24 * 7 * time.Hour
	p.Stores = []store.Store{s}

	for _, opt := range opts {
		opt(p)
	}

	// 比accessExpires多5s避免获取到accessExpires时未过期，但是获取session已经过期
	p.Session.accessExpires = p.accessExpires + time.Second*5
	return p
}

// ===============================通过 Claims Session 生成token=============================

// Signed token生成
func (auth *auth) Signed(ctx context.Context, userClaims *UserInfo) (accessToken string, refreshToken string, err error) {
	now := time.Now().Local()
	access := &AccessClaims{RegisteredClaims: jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(now.Add(auth.accessExpires)), IssuedAt: jwt.NewNumericDate(now)}, UserInfo: userClaims}
	refresh := &RefreshClaims{RegisteredClaims: jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(now.Add(auth.refreshExpires)), IssuedAt: jwt.NewNumericDate(now)}, Refresh: &Refresh{Id: userClaims.Id, Md5: GetSessionId(userClaims.Id)}}

	if accessToken, err = auth.SignedToken(access); err != nil {
		return "", "", fmt.Errorf("access token signed failed: %w", err)
	}
	if refreshToken, err = auth.SignedToken(refresh); err != nil {
		return "", "", fmt.Errorf("refresh token signed failed: %w", err)
	}
	return accessToken, refreshToken, nil
}

// ===============================解析token 生成 Claims Session=============================

// Parse 通过accessToken转换 userClaims userSession 并存到ctx中
func (auth *auth) Parse(ctx context.Context, accessToken string) (context.Context, error) {
	if len(accessToken) == 0 {
		return ctx, nil
	}
	//获取Claims对象
	claims := &AccessClaims{}
	if err := auth.ParseToken(accessToken, claims, AccessToken); err != nil {
		return ctx, err
	}
	if err := claims.Valid(); err != nil {
		return ctx, err
	}
	ctx = SetClaims(ctx, claims.UserInfo)

	//获取Session对象
	userSession, err := auth.GetSession(ctx, GetSessionId(claims.UserInfo.Id))
	if err != nil {
		return ctx, err
	}
	ctx = SetSession(ctx, userSession)
	return ctx, nil
}

// ===============================Refresh 生成 Claims Session=============================

func (auth *auth) Refresh(ctx context.Context, refreshToken string, claims Claims) error {
	if len(refreshToken) == 0 {
		return errors.New("refresh token is null")
	}
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	if err := auth.ParseToken(refreshToken, claims, RefreshToken); err != nil {
		return err
	}
	return claims.Valid()
}
