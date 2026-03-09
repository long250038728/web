package authorization

import (
	"context"
	"github.com/long250038728/web/tool/app_error"
)

type claimsKey struct{}
type sessionKey struct{}

func SetClaims(ctx context.Context, userClaims *UserInfo) context.Context {
	return context.WithValue(ctx, claimsKey{}, userClaims)
}
func GetClaims(ctx context.Context) (*UserInfo, error) {
	if val, ok := ctx.Value(claimsKey{}).(*UserInfo); ok {
		return val, nil
	}
	return nil, app_error.ClaimsNull
}

func SetSession(ctx context.Context, userSession *UserSession) context.Context {
	return context.WithValue(ctx, sessionKey{}, userSession)
}
func GetSession(ctx context.Context) (*UserSession, error) {
	if val, ok := ctx.Value(sessionKey{}).(*UserSession); ok {
		return val, nil
	}
	return nil, app_error.SessionNull
}
