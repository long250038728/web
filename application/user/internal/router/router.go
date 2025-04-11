package router

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/long250038728/web/application/user/internal/service"
	"github.com/long250038728/web/protoc/user"
	"github.com/long250038728/web/tool/app"
	"github.com/long250038728/web/tool/server/http/gateway"
	"github.com/long250038728/web/tool/server/http/gateway/middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
	"time"
)

func RegisterHTTPServer(engine *gin.Engine, srv *service.UserService) {
	cache, _ := app.NewUtil().Cache()

	// 服务/领域/接口
	userGroup := engine.Group("/user/user/").Use(middleware.BaseHandle(cache))
	{
		userGroup.GET("say", func(c *gin.Context) {
			gateway.Json(c, &user.RequestHello{}).Use(
				middleware.Login(),                    //需要登录
				middleware.Validate([]string{"name"}), //参数必填项
				middleware.Rule(c),                    //权限判断
				middleware.Locker(c, cache, []string{"name"}, middleware.IsClaims(true), middleware.Expiration(time.Second*3)), //分布式锁处理
				middleware.Cache(c, cache, []string{"name"}, middleware.IsClaims(true)),                                        //缓存处理
			).Handle(func(ctx context.Context, req any) (any, error) {
				return srv.SayHello(ctx, req.(*user.RequestHello))
			})
		})

		userGroup.GET("file", func(c *gin.Context) {
			gateway.File(c, &user.RequestHello{}).Handle(func(ctx context.Context, req any) (any, error) {
				return srv.File(ctx, req.(*user.RequestHello))
			})
		})

		userGroup.GET("sse", func(c *gin.Context) {
			gateway.SSE(c, &user.RequestHello{}).Handle(func(ctx context.Context, req any) (any, error) {
				return srv.SendSSE(ctx, req.(*user.RequestHello))
			})
		})
	}
}

func RegisterGRPCServer(engine *grpc.Server, srv *service.UserService) {
	user.RegisterUserServer(engine, srv)
	grpc_health_v1.RegisterHealthServer(engine, srv)
}
