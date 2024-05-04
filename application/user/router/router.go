package router

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/long250038728/web/application/user/ddd/service"
	user "github.com/long250038728/web/application/user/protoc"
	"github.com/long250038728/web/tool/app"
	"github.com/long250038728/web/tool/auth"
	"github.com/long250038728/web/tool/limiter"
	"github.com/long250038728/web/tool/server/http/tool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
	_ "net/http/pprof"
	"time"
)

//go func() {
//	log.Println(http.ListenAndServe("localhost:6060", nil))  //"net/http/pprof"
//}()

func RegisterHTTPServer(engine *gin.Engine, srv *service.UserService) {
	opts := []tool.MiddlewareOpt{
		tool.Limiter( //设置限流
			limiter.NewCacheLimiter(app.NewUtil().Cache(), time.Second, 10),
		),

		tool.Auth( //设置权限
			auth.NewCacheAuth(app.NewUtil().Cache(), auth.WhiteList([]string{"/", "/hello"})),
		),

		tool.Error( //设置错误
			[]tool.MiddleErr{}, //可以通过数据库处理
		),
	}
	middleware := tool.NewMiddlewarePool(opts...)

	engine.GET("/", func(gin *gin.Context) {
		var request user.RequestHello
		middleware.JSON(gin, &request, func(ctx context.Context) (interface{}, error) {
			return srv.SayHello(ctx, &request)
		})
	})

	engine.POST("/hello", func(gin *gin.Context) {
		var request user.RequestHello
		middleware.JSON(gin, &request, func(ctx context.Context) (interface{}, error) {
			return srv.SayHello(ctx, &request)
		})
	})

	engine.POST("/xls", func(gin *gin.Context) {
		var request user.RequestHello
		middleware.File(gin, &request, func(ctx context.Context) (interface{}, error) {
			return srv.SayHello(ctx, &request)
		})
	})
}

func RegisterGRPCServer(engine *grpc.Server, srv *service.UserService) {
	user.RegisterUserServer(engine, srv)
	grpc_health_v1.RegisterHealthServer(engine, srv)
}
