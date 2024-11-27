package router

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/long250038728/web/application/auth/internal/service"
	"github.com/long250038728/web/protoc/auth_rpc"
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

func RegisterHTTPServer(engine *gin.Engine, srv *service.Service) {
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

	engine.GET("/authorization/login", func(gin *gin.Context) {
		var request auth_rpc.LoginRequest
		middleware.JSON(gin, &request, func(ctx context.Context) (interface{}, error) {
			return srv.Login(ctx, &request)
		})
	})

	engine.GET("/authorization/refresh", func(gin *gin.Context) {
		var request auth_rpc.RefreshRequest
		middleware.JSON(gin, &request, func(ctx context.Context) (interface{}, error) {
			return srv.Refresh(ctx, &request)
		})
	})
}

func RegisterGRPCServer(engine *grpc.Server, srv *service.Service) {
	auth_rpc.RegisterAuthServer(engine, srv)
	grpc_health_v1.RegisterHealthServer(engine, srv)
}
