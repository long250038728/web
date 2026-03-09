package authorization

//
//import (
//	"context"
//	"errors"
//	"fmt"
//	"github.com/golang-jwt/jwt/v4"
//	"github.com/long250038728/web/tool/store"
//	"time"
//)
//
//type Opt func(r *auth)
//
//type auth struct {
//	Serialization
//	Session
//	accessExpires  time.Duration
//	refreshExpires time.Duration
//}
//
//func SecretKey(secretKey []byte) Opt {
//	return func(r *auth) {
//		r.SecretKey = secretKey
//	}
//}
//
//func AccessExpires(accessExpires time.Duration) Opt {
//	return func(r *auth) {
//		r.accessExpires = accessExpires
//	}
//}
//
//func RefreshExpires(refreshExpires time.Duration) Opt {
//	return func(r *auth) {
//		r.refreshExpires = refreshExpires
//	}
//}
//
//func NewAuth(s store.Store, opts ...Opt) Session {
//	p := &auth{}
//
//	//默认值
//	p.SecretKey = []byte("secretKey")
//	p.accessExpires = 20 * time.Minute
//	p.refreshExpires = 24 * 7 * time.Hour
//	p.store = s
//
//	for _, opt := range opts {
//		opt(p)
//	}
//
//	// 比accessExpires多5s避免获取到accessExpires时未过期，但是获取session已经过期
//	p.Session.accessExpires = p.accessExpires + time.Second*5
//	return p
//}
//
//// ===============================通过 Claims Session 生成token=============================
//
//// Signed token生成
//func (auth *auth) Signed(ctx context.Context, userClaims *UserInfo) (accessToken string, refreshToken string, err error) {
//	now := time.Now().Local()
//	access := &AccessClaims{RegisteredClaims: jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(now.Add(auth.accessExpires)), IssuedAt: jwt.NewNumericDate(now)}, UserInfo: userClaims}
//	refresh := &RefreshClaims{RegisteredClaims: jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(now.Add(auth.refreshExpires)), IssuedAt: jwt.NewNumericDate(now)}, Refresh: &Refresh{Id: userClaims.Id, Md5: GetSessionKey(userClaims.Id)}}
//
//	if accessToken, err = auth.SignedToken(access); err != nil {
//		return "", "", fmt.Errorf("access token signed failed: %w", err)
//	}
//	if refreshToken, err = auth.SignedToken(refresh); err != nil {
//		return "", "", fmt.Errorf("refresh token signed failed: %w", err)
//	}
//	return accessToken, refreshToken, nil
//}
//
//// ===============================解析token 生成 Claims Session=============================
//
//// Parse 通过accessToken转换 userClaims userSession 并存到ctx中
//func (auth *auth) Parse(ctx context.Context, accessToken string) (context.Context, error) {
//	if len(accessToken) == 0 {
//		return ctx, nil
//	}
//	//获取Claims对象
//	claims := &AccessClaims{}
//	if err := auth.ParseToken(accessToken, claims, AccessToken); err != nil {
//		return ctx, err
//	}
//	if err := claims.Valid(); err != nil {
//		return ctx, err
//	}
//	ctx = SetClaims(ctx, claims.UserInfo)
//
//	//获取Session对象
//	userSession, err := auth.GetSession(ctx, GetSessionKey(claims.UserInfo.Id))
//	if err != nil {
//		return ctx, err
//	}
//	ctx = SetSession(ctx, userSession)
//	return ctx, nil
//}
//
//// ===============================Refresh 生成 Claims Session=============================
//
//func (auth *auth) Refresh(ctx context.Context, refreshToken string, claims Claims) error {
//	if len(refreshToken) == 0 {
//		return errors.New("refresh token is null")
//	}
//	select {
//	case <-ctx.Done():
//		return ctx.Err()
//	default:
//	}
//	if err := auth.ParseToken(refreshToken, claims, RefreshToken); err != nil {
//		return err
//	}
//	return claims.Valid()
//}

type TokenType int32

const (
	AccessToken = iota
	RefreshToken
)

// =====================================================================================

//func NewAAA[C Claims, S any](store store.Store, claims Claims, session *S) *AAA[C, S] {
//	return &AAA[C, S]{
//		accessClaims: claims,
//		session:      session,
//		token:        NewToken(),
//		sess:         NewSession(store),
//	}
//}
//
//type AAA[C Claims, S any] struct {
//	sess  Session
//	token Token
//
//	accessClaims Claims
//	session      *S
//}
//
//func (a *AAA[C, S]) AuthorizationFromHeader(ctx context.Context, authorization string) (context.Context, error) {
//	if len(authorization) == 0 {
//		return ctx, errors.New("authorization is null")
//	}
//
//	if err := a.token.Parse(ctx, authorization, a.accessClaims); err != nil {
//		return ctx, err
//	}
//
//	if err := a.sess.GetSession(ctx, GetSessionKey(a.accessClaims.GetSessionId()), a.session); err != nil {
//		return ctx, err
//	}
//
//	return a.setCtx(ctx, a.accessClaims, a.session)
//}
//
//func (a *AAA[C, S]) AuthorizationFromMD(ctx context.Context, md map[string][]string) (context.Context, error) {
//	if claims, ok := md["claims"]; ok && len(claims) == 1 {
//		if err := json.Unmarshal([]byte(claims[0]), &a.accessClaims); err == nil {
//		}
//	}
//	if session, ok := md["session"]; ok && len(session) == 1 {
//		if err := json.Unmarshal([]byte(session[0]), &a.session); err == nil {
//		}
//	}
//	return a.setCtx(ctx, a.accessClaims, a.session)
//}
//
//func (a *AAA[C, S]) setCtx(ctx context.Context, accessClaims any, session *S) (context.Context, error) {
//	sessCtx := NewClaimsSessionContext[C, S]()
//	ctx = sessCtx.SetClaims(ctx, accessClaims.(*C))
//	ctx = sessCtx.SetSession(ctx, session)
//	return ctx, nil
//}
