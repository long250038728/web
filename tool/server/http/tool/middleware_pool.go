package tool

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"sync"
)

type HttpFunc func(ctx context.Context) (interface{}, error)

type FileInterface interface {
	FileName() string
	FileData() []byte
}

type MiddlewarePool struct {
	pool sync.Pool
}

func NewMiddlewarePool(opts ...MiddlewareOpt) *MiddlewarePool {
	middleTool := &MiddlewarePool{
		pool: sync.Pool{New: func() any {
			return NewMiddleware(opts...)
		}},
	}
	return middleTool
}

// JSON  json返回
func (m *MiddlewarePool) JSON(gin *gin.Context, request any, function HttpFunc) {
	middleware, _ := m.pool.Get().(*Middleware)
	defer func() {
		middleware.Reset()
		m.pool.Put(middleware)
	}()

	ctx, err := middleware.Context(gin)

	if err != nil {
		middleware.WriteJSON(nil, err)
		return
	}

	//基础处理 （bind绑定  及链路 处理）
	if err := middleware.Bind(request); err != nil {
		middleware.WriteJSON(nil, err)
		return
	}

	//处理业务
	res, err := function(ctx)
	middleware.WriteJSON(res, err)
}

// File  File返回
//
//	response 必须实现 FileInterface 接口
func (m *MiddlewarePool) File(gin *gin.Context, request any, function HttpFunc) {
	middleware, _ := m.pool.Get().(*Middleware)
	defer func() {
		middleware.Reset()
		m.pool.Put(middleware)
	}()

	ctx, err := middleware.Context(gin)
	if err != nil {
		middleware.WriteJSON(nil, err)
		return
	}

	//基础处理 （bind绑定  及链路 处理）
	if err := middleware.Bind(request); err != nil {
		middleware.WriteJSON(nil, err)
		return
	}

	//处理业务
	res, err := function(ctx)
	if err != nil {
		middleware.WriteJSON(nil, err)
		return
	}

	if file, ok := res.(FileInterface); !ok {
		middleware.WriteJSON(nil, errors.New("the file is not interface to FileInterface"))
	} else {
		middleware.WriteFile(file.FileName(), file.FileData())
	}
}

// SSE  SSE返回
//
// response 必须是<-chan string
func (m *MiddlewarePool) SSE(gin *gin.Context, request any, function HttpFunc) {
	middleware, _ := m.pool.Get().(*Middleware)
	defer func() {
		middleware.Reset()
		m.pool.Put(middleware)
	}()

	ctx, err := middleware.Context(gin)
	if err != nil {
		middleware.WriteJSON(nil, err)
		return
	}

	//基础处理 （bind绑定  及链路 处理）
	if err := middleware.Bind(request); err != nil {
		middleware.WriteJSON(nil, err)
		return
	}

	// 检查是否支持SSE
	if err := middleware.CheckSupportSEE(); err != nil {
		middleware.WriteJSON(nil, err)
		return
	}

	//处理业务
	res, err := function(ctx)
	if err != nil {
		middleware.WriteJSON(nil, err)
		return
	}

	if ch, ok := res.(<-chan string); !ok {
		middleware.WriteJSON(nil, errors.New("the response is not chan string"))
	} else {
		middleware.WriteSSE(ch)
	}
}
