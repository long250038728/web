package router

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/long250038728/web/application/user/internal/service"
	"github.com/long250038728/web/protoc/user"
	"github.com/long250038728/web/tool/app"
	"github.com/long250038728/web/tool/server/http/gateway"
	"github.com/long250038728/web/tool/server/http/gateway/interceptor"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
	"time"
)

//go func() {
//	log.Println(http.ListenAndServe("localhost:6060", nil))  //"net/http/pprof"
//}()

func RegisterHTTPServer(engine *gin.Engine, srv *service.UserService) {
	cache, _ := app.NewUtil().Cache()

	userGroup := engine.Group("/user/user/").Use(interceptor.BaseHandle(cache))
	{
		userGroup.GET("say", func(c *gin.Context) {
			gateway.Json(c, &user.RequestHello{}).Use(
				interceptor.Login(),
				interceptor.Validate([]string{"name"}),
				interceptor.Cache(c, cache, []string{}, interceptor.SetIsClaims(true), interceptor.SetExpiration(time.Second*10)),
			).Handle(func(ctx context.Context, req any) (any, error) {
				return srv.SayHello(ctx, req.(*user.RequestHello))
			})
		})

		userGroup.GET("xls", func(c *gin.Context) {
			gateway.File(c, &user.RequestHello{}).Handle(func(ctx context.Context, req any) (any, error) {
				return srv.SayHello(ctx, req.(*user.RequestHello))
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
