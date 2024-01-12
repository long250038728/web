package router

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	user "github.com/long250038728/web/application/user/protoc"
	"github.com/long250038728/web/tool/server/http/tool"
	"sync"
)

func RegisterUserServerServer(engine *gin.Engine, srv user.UserServerServer) {
	//设置错误
	//设置限流
	//设置权限
	var opts = []tool.MiddlewareOpt{
		tool.SetErrData(map[error]tool.Err{errors.New("hello"): tool.Err{Code: "1233333", Message: "错误"}}),
		//http.SetAuth(auth.NewJwtAuth()),
		//http.SetLimiter(redis.NewRedisLimiter()),
	}

	//线程池
	pool := sync.Pool{New: func() any {
		return tool.NewMiddleware(opts...)
	}}

	// ======================================= handle =================================================
	engine.GET("/", func(gin *gin.Context) {
		middle := pool.Get().(*tool.Middleware)
		defer func() {
			middle.Reset()
			pool.Put(middle)
		}()

		req := &user.RequestHello{}
		middle.GinContext(gin).Bind(req).Do(func(ctx context.Context) (interface{}, error) {
			return srv.SayHello(ctx, req)
		})
	})

	engine.GET("/hello", func(gin *gin.Context) {
		middle := pool.Get().(*tool.Middleware)
		defer func() {
			middle.Reset()
			pool.Put(middle)
		}()

		req := &user.RequestHello{}
		middle.GinContext(gin).Bind(req).Do(func(ctx context.Context) (interface{}, error) {
			return srv.SayHello(ctx, req)
		})
	})
}
