package auth

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/long250038728/web/tool/struct_map"
	"strings"
	"time"
)

type Opt func(r *cacheAuth)

type cacheAuth struct {
	store Store
	white White

	secretKey   []byte
	userClaims  *UserClaims
	userSession *UserSession
}

func SecretKey(secretKey []byte) Opt {
	return func(r *cacheAuth) {
		r.secretKey = secretKey
	}
}

func WhiteList(white White) Opt {
	return func(r *cacheAuth) {
		r.white = white
	}
}

func NewAuth(cache Store, opts ...Opt) Auth {
	r := &cacheAuth{
		store:       cache,
		userClaims:  &UserClaims{},
		userSession: &UserSession{AuthList: []string{}},
		secretKey:   []byte("secret_key"),
	}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

// ===============================通过 Claims Session 生成token=============================

// Signed token生成
func (p *cacheAuth) Signed(ctx context.Context, userClaims *UserClaims, session *UserSession) (accessToken string, refreshToken string, err error) {
	b, err := json.Marshal(session)
	if err != nil {
		return "", "", err
	}

	// session信息存放到中间件中
	ok, err := p.store.Set(ctx, authToken(userClaims.Id), string(b))
	if err != nil {
		return "", "", err
	}
	if !ok {
		return "", "", errors.New("session setting is err")
	}

	claims := &jwtClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(1800 * time.Minute).Unix(),
			IssuedAt:  time.Now().Local().Unix(),
		},
		UserClaims: userClaims,
	}

	refreshClaims := &RefreshClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(1800 * time.Minute).Unix(),
			IssuedAt:  time.Now().Local().Unix(),
		},
		Id:  userClaims.Id,
		Md5: authToken(userClaims.Id),
	}

	accessToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(p.secretKey)
	if err != nil {
		return "", "", nil
	}

	refreshToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString(p.secretKey)
	if err != nil {
		return "", "", nil
	}

	return accessToken, refreshToken, nil
}

// ===============================Refresh 生成 Claims Session=============================

func (p *cacheAuth) Refresh(ctx context.Context, refreshToken string) (*RefreshClaims, error) {
	if len(refreshToken) == 0 {
		return nil, errors.New("refresh token is null")
	}

	// 解析JWT字符串
	token, err := jwt.ParseWithClaims(refreshToken, &RefreshClaims{}, func(token *jwt.Token) (interface{}, error) {
		return p.secretKey, nil // 这里你需要提供用于签名的密钥
	})
	if err != nil {
		if validationErr, ok := err.(*jwt.ValidationError); ok && validationErr.Errors == jwt.ValidationErrorExpired {
			err = errors.New("token is Disabled")
		}
		return nil, err
	}

	//获取Claims对象
	claims := token.Claims.(*RefreshClaims)
	if claims.Md5 != authToken(claims.Id) {
		return nil, errors.New("refresh token invalid")
	}

	return claims, nil
}

// ===============================解析token 生成 Claims Session=============================

// Parse 通过accessToken转换 userClaims userSession 并存到ctx中
func (p *cacheAuth) Parse(ctx context.Context, accessToken string) (context.Context, error) {
	userClaims, userSession, err := p.parse(ctx, accessToken)
	if err != nil {
		return ctx, err
	}
	ctx = SetClaims(ctx, userClaims)
	ctx = SetSession(ctx, userSession)
	return ctx, nil
}

func (p *cacheAuth) parse(ctx context.Context, accessToken string) (*UserClaims, *UserSession, error) {
	if len(accessToken) == 0 {
		return nil, nil, nil
	}

	//获取Claims对象
	userClaims, err := p.getClaims(accessToken)
	if err != nil {
		return nil, nil, err
	}

	//获取Session对象
	userSession, err := p.getSession(ctx, authToken(userClaims.Id))
	if err != nil {
		return nil, nil, err
	}
	p.userClaims = userClaims
	p.userSession = userSession
	return p.userClaims, p.userSession, nil
}

func (p *cacheAuth) getClaims(signedString string) (*UserClaims, error) {
	userClaims := &UserClaims{}

	// 解析JWT字符串
	token, err := jwt.ParseWithClaims(signedString, &jwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		return p.secretKey, nil // 这里你需要提供用于签名的密钥
	})
	if err != nil {
		if validationErr, ok := err.(*jwt.ValidationError); ok && validationErr.Errors == jwt.ValidationErrorExpired {
			err = errors.New("token is Disabled")
		}
		return nil, err
	}

	//获取Claims对象
	claims := token.Claims.(*jwtClaims)
	err = struct_map.Map(claims.UserClaims, userClaims) //带有jwt.StandardClaims 的对象 转换为 外部不带有 jwt.StandardClaims 的对象
	if err != nil {
		return nil, err
	}
	return userClaims, nil
}

func (p *cacheAuth) getSession(ctx context.Context, token string) (session *UserSession, err error) {
	//检查authToken
	if token == "" {
		return nil, errors.New("token is empty")
	}
	sessionStr, err := p.store.Get(ctx, token)
	if err != nil {
		return nil, err
	}

	if len(sessionStr) == 0 {
		return nil, errors.New("token is empty")
	}
	return session, json.Unmarshal([]byte(sessionStr), &session)
}

func authToken(id int32) string {
	hash := sha256.New()
	hash.Write([]byte(fmt.Sprintf("id:%d", id))) // 向哈希计算对象中写入字符串数据
	return hex.EncodeToString(hash.Sum(nil))
}

// =================================业务判断===========================

// Auth 判断是否有权限
//
//  1. 判断是否是白名单
//  2. 判断是否是登录
//  3. 判断这个接口有没有权限（从Session中获取）
func (p *cacheAuth) Auth(ctx context.Context, path string) error {
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
	if p.userClaims != nil && p.loginPath(path) {
		return nil
	}

	//匹配session
	for _, authPath := range p.userSession.AuthList {
		if authPath == path {
			return nil
		}
	}
	return errors.New("no match path")
}

// whitePath path是否为白名单
func (p *cacheAuth) whitePath(path string) bool {
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
func (p *cacheAuth) loginPath(path string) bool {
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
