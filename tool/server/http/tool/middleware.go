package tool

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/long250038728/web/tool/auth"
	"github.com/long250038728/web/tool/limiter"
	"github.com/long250038728/web/tool/tracing/opentelemetry"
	"sync"
)

type Middleware struct {
	error   map[error]MiddleErr
	pool    sync.Pool
	auth    auth.Auth
	limiter limiter.Limiter
}

type HttpFunc func(ctx context.Context) (interface{}, error)

func NewMiddle(opts ...MiddlewareOpt) *Middleware {
	middleTool := &Middleware{
		pool: sync.Pool{New: func() any {
			return &MiddleItem{}
		}},
	}
	for _, opt := range opts {
		opt(middleTool)
	}
	return middleTool
}

func (m *Middleware) JSON(gin *gin.Context, function HttpFunc) {
	//基本处理
	item, _ := m.pool.Get().(*MiddleItem)
	item.Set(gin)
	ctx := item.Ctx()

	span := opentelemetry.NewSpan(ctx, "HTTP Request")
	defer func() {
		//返回线程池
		span.Close()
		item.Reset()
		m.pool.Put(item)
	}()

	//限流
	if m.limiter != nil {
		if err := m.limiter.Allow(ctx, "hello"); err != nil {
			item.WriteJSON(nil, err)
			return
		}
	}

	//授权
	if m.auth != nil {
		//获取Claims对象
		userClaims, err := auth.Parse(gin.GetHeader("Authorization"))
		if err != nil {
			item.WriteJSON(nil, err)
			return
		}

		if err = m.auth.Auth(ctx, userClaims, gin.Request.RequestURI); err != nil {
			item.WriteJSON(nil, err)
			return
		}
	}

	//处理业务
	res, err := function(ctx)
	item.WriteJSON(res, err)
}
