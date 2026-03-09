package authorization

import (
	"context"
)

type Claims interface {
	Valid() error
}

//================================ token ================================

type Token interface {
	// Signed 生成token
	Signed(ctx context.Context, claims Claims) (string, error)
	// Parse 解析token转换为Claims
	Parse(ctx context.Context, token string, claims Claims) error
}

//================================ session ================================

type Session interface {
	// GetSession 获取session
	GetSession(ctx context.Context, sessionId string, session any) error
	// SetSession 设置session
	SetSession(ctx context.Context, sessionId string, session any) (err error)
	// DeleteSession 删除session
	DeleteSession(ctx context.Context, sessionId string) error
}
