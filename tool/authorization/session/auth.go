package session

import (
	"context"
	"errors"
	"github.com/golang-jwt/jwt"
	"github.com/long250038728/web/tool/authorization"
	"strings"
	"time"
)

type Opt func(r *CacheAuth)

type CacheAuth struct {
	authorization.Token
	Session
	white authorization.White
}

func SecretKey(secretKey []byte) Opt {
	return func(r *CacheAuth) {
		r.SecretKey = secretKey
	}
}

func WhiteList(white authorization.White) Opt {
	return func(r *CacheAuth) {
		r.white = white
	}
}

func NewAuth(store authorization.Store, opts ...Opt) Auth {
	r := &CacheAuth{}
	r.SecretKey = []byte("secretKey")
	r.Store = store

	for _, opt := range opts {
		opt(r)
	}
	return r
}

// ===============================通过 Claims Session 生成token=============================

// Signed token生成
func (p *CacheAuth) Signed(ctx context.Context, userClaims *UserInfo, session *UserSession) (accessToken string, refreshToken string, err error) {
	if err = p.SetSession(ctx, authorization.GetSessionId(userClaims.Id), session); err != nil {
		return "", "", err
	}
	now := time.Now().Local()
	access := &AccessClaims{StandardClaims: jwt.StandardClaims{ExpiresAt: now.Add(1800 * time.Minute).Unix(), IssuedAt: now.Unix()}, UserInfo: userClaims}
	refresh := &RefreshClaims{StandardClaims: jwt.StandardClaims{ExpiresAt: now.Add(1800 * time.Minute).Unix(), IssuedAt: now.Unix()}, Refresh: &Refresh{Id: userClaims.Id, Md5: authorization.GetSessionId(userClaims.Id)}}

	if accessToken, err = p.SignedToken(access); err != nil {
		return "", "", nil
	}
	if refreshToken, err = p.SignedToken(refresh); err != nil {
		return "", "", nil
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
	if err := p.ParseToken(accessToken, claims, authorization.AccessToken); err != nil {
		return ctx, err
	}
	if err := claims.Valid(); err != nil {
		return ctx, err
	}
	ctx = SetClaims(ctx, claims.UserInfo)

	//获取Session对象
	userSession, err := p.GetSession(ctx, authorization.GetSessionId(claims.UserInfo.Id))
	if err != nil {
		return ctx, err
	}
	ctx = SetSession(ctx, userSession)
	return ctx, nil
}

// ===============================Refresh 生成 Claims Session=============================

func (p *CacheAuth) Refresh(ctx context.Context, refreshToken string, claims authorization.Claims) error {
	if len(refreshToken) == 0 {
		return errors.New("refresh token is null")
	}
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	return p.ParseToken(refreshToken, claims, authorization.RefreshToken)
}

// =================================业务判断===========================

// Auth 判断是否有权限
//
//  1. 判断是否是白名单
//  2. 判断是否是登录
//  3. 判断这个接口有没有权限（从Session中获取）
func (p *CacheAuth) Auth(ctx context.Context, path string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	//转换为小写
	path = strings.ToLower(path)

	//白名单
	if p.whitePath(path) {
		return nil
	}

	//匹配是否登录
	userClaims, err := GetClaims(ctx)
	if err != nil {
		return err
	}
	if userClaims != nil && p.loginPath(path) {
		return nil
	}
	userSession, err := GetSession(ctx)
	if err != nil {
		return err
	}

	//匹配session
	for _, authPath := range userSession.AuthList {
		if authPath == path {
			return nil
		}
	}
	return errors.New("no match path")
}

// whitePath path是否为白名单
func (p *CacheAuth) whitePath(path string) bool {
	if p.white == nil {
		return false
	}
	for _, p := range p.white.WhiteList() {
		if p == path {
			return true
		}
	}
	return false
}

// loginPath path是否为登录
func (p *CacheAuth) loginPath(path string) bool {
	if p.white == nil {
		return false
	}
	for _, p := range p.white.LoginList() {
		if p == path {
			return true
		}
	}
	return false
}
