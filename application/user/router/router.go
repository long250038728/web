package router

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/long250038728/web/application/user/ddd/service"
	user "github.com/long250038728/web/application/user/protoc"
	"github.com/long250038728/web/tool/server/http/tool"
	_ "net/http/pprof"
	"sync"
)

func RegisterUserServerServer(engine *gin.Engine, srv *service.UserService) {
	//"net/http/pprof"
	//go func() {
	//	log.Println(http.ListenAndServe("localhost:6060", nil))
	//}()

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
