package router

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/long250038728/web/application/user/internal/service"
	"github.com/long250038728/web/protoc/user"
	"github.com/long250038728/web/tool/app"
	"github.com/long250038728/web/tool/authorization"
	"github.com/long250038728/web/tool/limiter"
	"github.com/long250038728/web/tool/server/http/tool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
	"time"
)

//go func() {
//	log.Println(http.ListenAndServe("localhost:6060", nil))  //"net/http/pprof"
//}()

func RegisterHTTPServer(engine *gin.Engine, srv *service.UserService) {
	var opts []tool.MiddlewareOpt

	if cache, err := app.NewUtil().Cache(); err == nil {
		opts = append(opts, tool.Limiter( //设置限流
			limiter.NewCacheLimiter(
				cache,
				limiter.SetExpiration(time.Second), limiter.SetTimes(10),
			),
		))
		opts = append(opts, tool.Auth( //设置权限（权限信息可从数据库获取文件获取）
			authorization.NewAuth(
				cache,
				authorization.WhiteList(authorization.NewLocalWhite([]string{"/", "/user/", "/user/hello", "/user/hello2", "/user/hello3"}, []string{})),
			),
		))
	}
	middleware := tool.NewMiddlewarePool(opts...)

	engine.GET("/", func(gin *gin.Context) {
		var request user.RequestHello
		middleware.JSON(gin, &request, func(ctx context.Context) (interface{}, error) {
			return srv.SayHello(ctx, &request)
		})
	})

	engine.GET("/user/", func(gin *gin.Context) {
		var request user.RequestHello
		middleware.JSON(gin, &request, func(ctx context.Context) (interface{}, error) {
			return srv.SayHello(ctx, &request)
		})
	})

	engine.POST("/user/xls", func(gin *gin.Context) {
		var request user.RequestHello
		middleware.File(gin, &request, func(ctx context.Context) (interface{}, error) {
			return srv.SayHello(ctx, &request)
		})
	})

	engine.GET("/user/sse", func(gin *gin.Context) {
		var request user.RequestHello
		middleware.SSE(gin, &request, func(ctx context.Context) (interface{}, error) {
			return srv.SendSSE(ctx, &request)
		})
	})
}

func RegisterGRPCServer(engine *grpc.Server, srv *service.UserService) {
	user.RegisterUserServer(engine, srv)
	grpc_health_v1.RegisterHealthServer(engine, srv)
}
