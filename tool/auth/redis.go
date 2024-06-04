package auth

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/golang-jwt/jwt"
	"github.com/long250038728/web/tool/cache"
	"github.com/long250038728/web/tool/struct_map"
	"strings"
	"time"
)

type Opt func(r *cacheAuth)

type cacheAuth struct {
	cache       cache.Cache
	secretKey   []byte
	whiteList   []string
	userClaims  *UserClaims
	userSession *UserSession
}

func SecretKey(secretKey []byte) Opt {
	return func(r *cacheAuth) {
		r.secretKey = secretKey
	}
}

func WhiteList(list []string) Opt {
	return func(r *cacheAuth) {
		r.whiteList = list
	}
}

func NewCacheAuth(cache cache.Cache, opts ...Opt) Auth {
	r := &cacheAuth{
		cache:       cache,
		userClaims:  &UserClaims{},
		userSession: &UserSession{},
		secretKey:   []byte("secret_key"),
	}
	for _, opt := range opts {
		opt(r)
	}

	if r.whiteList == nil {
		r.whiteList = make([]string, 0, 0)
	}

	return r
}

func (p *cacheAuth) Parse(ctx context.Context, accessToken string) (context.Context, error) {
	userClaims, userSession, err := p.parse(ctx, accessToken)
	if err != nil {
		return nil, err
	}
	ctx = SetClaims(ctx, userClaims)
	ctx = SetSession(ctx, userSession)
	return ctx, nil
}

// Parse 解析 signed
func (p *cacheAuth) parse(ctx context.Context, accessToken string) (*UserClaims, *UserSession, error) {
	if len(accessToken) == 0 {
		return p.userClaims, p.userSession, nil
	}

	//获取Claims对象
	userClaims, err := p.Claims(accessToken)
	if err != nil {
		return nil, nil, err
	}

	//获取Session对象
	userSession, err := p.Session(ctx, userClaims.AuthToken())
	if err != nil {
		return nil, nil, err
	}
	p.userClaims = userClaims
	p.userSession = userSession
	return p.userClaims, p.userSession, nil
}

// Refresh 续
func (p *cacheAuth) Refresh(ctx context.Context, refreshToken string) (*RefreshClaims, error) {
	if len(refreshToken) == 0 {
		return nil, errors.New("refresh token is null")
	}

	cla := &RefreshClaims{}

	// 解析JWT字符串
	token, err := jwt.ParseWithClaims(refreshToken, &jwtClaims{}, func(token *jwt.Token) (interface{}, error) {
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
	err = struct_map.Map(claims.UserClaims, cla) //带有jwt.StandardClaims 的对象 转换为 外部不带有 jwt.StandardClaims 的对象
	if err != nil {
		return nil, err
	}

	if cla.Md5 != "1234567890" {
		return nil, errors.New("refresh token invalid")
	}

	return cla, nil
}

// Auth 判断是否有权限
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

	//匹配
	for _, authPath := range p.userSession.AuthList {
		if authPath == path {
			return nil
		}
	}
	return errors.New("no match path")
}

func (p *cacheAuth) Claims(signedString string) (*UserClaims, error) {
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

func (p *cacheAuth) Session(ctx context.Context, token string) (session *UserSession, err error) {
	//检查authToken
	if token == "" {
		return nil, errors.New("token is empty")
	}
	sessionStr, err := p.cache.Get(ctx, token)
	if err != nil {
		return nil, err
	}

	if len(sessionStr) == 0 {
		return nil, errors.New("token is empty")
	}
	return session, json.Unmarshal([]byte(sessionStr), &session)
}

// whitePath path是否为白名单
func (p *cacheAuth) whitePath(path string) bool {
	for _, whitePath := range p.whiteList {
		if whitePath == path {
			return true
		}
	}
	return false
}

// Signed 用户内部信息生产token
func (p *cacheAuth) Signed(ctx context.Context, userClaims *UserClaims, userSession *UserSession) (string, string, error) {
	b, err := json.Marshal(userSession)
	if err != nil {
		return "", "", err
	}

	ok, err := p.cache.Set(ctx, userClaims.AuthToken(), string(b))
	if err != nil {
		return "", "", err
	}
	if !ok {
		return "", "", errors.New("session setting is err")
	}

	claims := &jwtClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(1800 * time.Minute).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		UserClaims: userClaims,
	}

	refreshClaims := &RefreshClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(1800 * time.Minute).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		Id:  userClaims.Id,
		Md5: "1234567890",
	}

	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(p.secretKey)
	if err != nil {
		return "", "", nil
	}

	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString(p.secretKey)
	if err != nil {
		return "", "", nil
	}

	return accessToken, refreshToken, nil
}
