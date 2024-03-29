package tool

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/long250038728/web/tool/auth"
	"github.com/long250038728/web/tool/limiter"
	"sync"
)

type MiddlewarePool struct {
	error   map[error]MiddleErr
	pool    sync.Pool
	auth    auth.Auth
	limiter limiter.Limiter
}
type HttpFunc func(ctx context.Context) (interface{}, error)

func NewMiddlewarePool(opts ...MiddlewareOpt) *MiddlewarePool {
	middleTool := &MiddlewarePool{}
	for _, opt := range opts {
		opt(middleTool)
	}
	if middleTool.error == nil {
		middleTool.error = map[error]MiddleErr{}
	}
	middleTool.pool = sync.Pool{New: func() any {
		return NewMiddleware(middleTool.error)
	}}
	return middleTool
}

// JSON  json返回
func (m *MiddlewarePool) JSON(gin *gin.Context, request any, function HttpFunc) {
	middleware, _ := m.pool.Get().(*Middleware)
	defer func() {
		middleware.Reset()
		m.pool.Put(middleware)
	}()

	middleware.Set(gin)
	ctx := middleware.Context()

	//限流
	if m.limiter != nil {
		if err := m.limiter.Allow(ctx, "HTTP API"); err != nil {
			middleware.WriteJSON(nil, err)
			return
		}
	}

	//授权
	if m.auth != nil {
		//获取Claims对象
		userClaims, err := auth.Parse(gin.GetHeader("Authorization"))
		if err != nil {
			middleware.WriteJSON(nil, err)
			return
		}

		if err = m.auth.Auth(ctx, userClaims, gin.Request.URL.Path); err != nil {
			middleware.WriteJSON(nil, err)
			return
		}
	}

	//基础处理 （bind绑定  及链路 处理）
	if err := middleware.bind(request); err != nil {
		middleware.WriteJSON(nil, err)
		return
	}

	//处理业务
	res, err := function(ctx)
	middleware.WriteJSON(res, err)
}
