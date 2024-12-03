package authorization

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

type Opt func(r *CacheAuth)

type CacheAuth struct {
	Serialization
	Session
	accessExpires  time.Duration
	refreshExpires time.Duration
}

func SecretKey(secretKey []byte) Opt {
	return func(r *CacheAuth) {
		r.SecretKey = secretKey
	}
}

func AccessExpires(accessExpires time.Duration) Opt {
	return func(r *CacheAuth) {
		r.accessExpires = accessExpires
	}
}

func RefreshExpires(refreshExpires time.Duration) Opt {
	return func(r *CacheAuth) {
		r.refreshExpires = refreshExpires
	}
}

func LocalStore(localStore Store) Opt {
	return func(r *CacheAuth) {
		r.LocalStore = localStore
	}
}

func NewAuth(store Store, opts ...Opt) Auth {
	p := &CacheAuth{}

	//默认值
	p.SecretKey = []byte("secretKey")
	p.accessExpires = 5 * time.Minute
	p.refreshExpires = 60 * 24 * 7 * time.Minute

	p.Store = store
	if localStore, err := NewLocalStore(10000); err == nil {
		p.LocalStore = localStore
	}

	for _, opt := range opts {
		opt(p)
	}

	// 比accessExpires多5s避免获取到accessExpires时未过期，但是获取session已经过期
	p.Session.accessExpires = p.accessExpires + time.Second*5
	return p
}

// ===============================通过 Claims Session 生成token=============================

// Signed token生成
func (p *CacheAuth) Signed(ctx context.Context, userClaims *UserInfo) (accessToken string, refreshToken string, err error) {
	now := time.Now().Local()
	access := &AccessClaims{StandardClaims: jwt.StandardClaims{ExpiresAt: now.Add(p.accessExpires).Unix(), IssuedAt: now.Unix()}, UserInfo: userClaims}
	refresh := &RefreshClaims{StandardClaims: jwt.StandardClaims{ExpiresAt: now.Add(p.refreshExpires).Unix(), IssuedAt: now.Unix()}, Refresh: &Refresh{Id: userClaims.Id, Md5: GetSessionId(userClaims.Id)}}

	if accessToken, err = p.SignedToken(access); err != nil {
		return "", "", fmt.Errorf("access token signed failed: %w", err)
	}
	if refreshToken, err = p.SignedToken(refresh); err != nil {
		return "", "", fmt.Errorf("refresh token signed failed: %w", err)
	}
	return accessToken, refreshToken, nil
}

// ===============================解析token 生成 Claims Session=============================

// Parse 通过accessToken转换 userClaims userSession 并存到ctx中
func (p *CacheAuth) Parse(ctx context.Context, accessToken string) (context.Context, error) {
	if len(accessToken) == 0 {
		return ctx, nil
	}
	//获取Claims对象
	claims := &AccessClaims{}
	if err := p.ParseToken(accessToken, claims, AccessToken); err != nil {
		return ctx, err
	}
	if err := claims.Valid(); err != nil {
		return ctx, err
	}
	ctx = SetClaims(ctx, claims.UserInfo)

	//获取Session对象
	userSession, err := p.GetSession(ctx, GetSessionId(claims.UserInfo.Id))
	if err != nil {
		return ctx, err
	}
	ctx = SetSession(ctx, userSession)
	return ctx, nil
}

// ===============================Refresh 生成 Claims Session=============================

func (p *CacheAuth) Refresh(ctx context.Context, refreshToken string, claims Claims) error {
	if len(refreshToken) == 0 {
		return errors.New("refresh token is null")
	}
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	if err := p.ParseToken(refreshToken, claims, RefreshToken); err != nil {
		return err
	}
	return claims.Valid()
}
