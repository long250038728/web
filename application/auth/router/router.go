package router

import (
	"github.com/gin-gonic/gin"
	"github.com/long250038728/web/application/auth/internal/service"
	"github.com/long250038728/web/protoc/auth_rpc"
	"github.com/long250038728/web/tool/app"
	"github.com/long250038728/web/tool/server/http/tool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func RegisterHTTPServer(engine *gin.Engine, srv *service.Service) {
	cache, _ := app.NewUtil().Cache()
	user := engine.Group("/auth/user/").Use(tool.BaseHandle(cache), tool.LimitHandle(cache))
	{
		user.POST("login", func(c *gin.Context) {
			var request auth_rpc.LoginRequest
			tool.NewHttpTools().JSON(c, &request, func() (interface{}, error) {
				return srv.Login(c.Request.Context(), &request)
			})
		})
	}

	//var opts []tool.MiddlewareOpt
	//if cache, err := app.NewUtil().Cache(); err == nil {
	//	opts = append(opts, tool.Limiter( //设置限流
	//		limiter.NewCacheLimiter(
	//			cache,
	//			limiter.SetExpiration(time.Second), limiter.SetTimes(10),
	//		),
	//	))
	//	opts = append(opts, tool.Auth( //设置权限（权限信息可从数据库获取文件获取）
	//		authorization.NewAuth(
	//			cache,
	//			authorization.WhiteList(authorization.NewLocalWhite([]string{"/", "/user/", "/user/hello", "/user/hello2", "/user/hello3"}, []string{})),
	//		),
	//	))
	//}
	//middleware := tool.NewHttpTools(opts...)
	//
	//engine.GET("/authorization/login", func(gin *gin.Context) {
	//	var request auth_rpc.LoginRequest
	//	middleware.JSON(gin, &request, func(ctx context.Context) (interface{}, error) {
	//		return srv.Login(ctx, &request)
	//	})
	//})
	//
	//engine.GET("/authorization/refresh", func(gin *gin.Context) {
	//	var request auth_rpc.RefreshRequest
	//	middleware.JSON(gin, &request, func(ctx context.Context) (interface{}, error) {
	//		return srv.Refresh(ctx, &request)
	//	})
	//})
}

//func RegisterHTTPServer(engine *gin.Engine, srv *service.Service) {
//
//	var opts []tool.MiddlewareOpt
//	if cache, err := app.NewUtil().Cache(); err == nil {
//		opts = append(opts, tool.Limiter( //设置限流
//			limiter.NewCacheLimiter(
//				cache,
//				limiter.SetExpiration(time.Second), limiter.SetTimes(10),
//			),
//		))
//		opts = append(opts, tool.Auth( //设置权限（权限信息可从数据库获取文件获取）
//			authorization.NewAuth(
//				cache,
//				authorization.WhiteList(authorization.NewLocalWhite([]string{"/", "/user/", "/user/hello", "/user/hello2", "/user/hello3"}, []string{})),
//			),
//		))
//	}
//	middleware := tool.NewHttpTools(opts...)
//
//	engine.GET("/authorization/login", func(gin *gin.Context) {
//		var request auth_rpc.LoginRequest
//		middleware.JSON(gin, &request, func(ctx context.Context) (interface{}, error) {
//			return srv.Login(ctx, &request)
//		})
//	})
//
//	engine.GET("/authorization/refresh", func(gin *gin.Context) {
//		var request auth_rpc.RefreshRequest
//		middleware.JSON(gin, &request, func(ctx context.Context) (interface{}, error) {
//			return srv.Refresh(ctx, &request)
//		})
//	})
//}

func RegisterGRPCServer(engine *grpc.Server, srv *service.Service) {
	auth_rpc.RegisterAuthServer(engine, srv)
	grpc_health_v1.RegisterHealthServer(engine, srv)
}
