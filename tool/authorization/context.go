package authorization

import (
	"context"
	"errors"
)

type claims struct{}
type session struct{}

func SetClaims(ctx context.Context, userClaims *UserInfo) context.Context {
	return context.WithValue(ctx, claims{}, userClaims)
}
func GetClaims(ctx context.Context) (*UserInfo, error) {
	if val, ok := ctx.Value(claims{}).(*UserInfo); ok {
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
