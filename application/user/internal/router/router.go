package router

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/long250038728/web/application/user/internal/service"
	"github.com/long250038728/web/protoc/user"
	"github.com/long250038728/web/tool/app"
	"github.com/long250038728/web/tool/server/http"
	gateway2 "github.com/long250038728/web/tool/server/http/gateway"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
)

//go func() {
//	log.Println(http.ListenAndServe("localhost:6060", nil))  //"net/http/pprof"
//}()

func RegisterHTTPServer(engine *gin.Engine, srv *service.UserService) {
	cache, _ := app.NewUtil().Cache()

	userGroup := engine.Group("/user/user/").Use(http.BaseHandle(cache), http.LimitHandle(cache))
	{
		userGroup.GET("say", func(c *gin.Context) {
			gateway2.Json(c, &user.RequestHello{}).Use(gateway2.Limit(), gateway2.Cache()).Handle(func(ctx context.Context, req any) (any, error) {
				return srv.SayHello(ctx, req.(*user.RequestHello))
			})
		})

		userGroup.GET("xls", func(c *gin.Context) {
			gateway2.File(c, &user.RequestHello{}).Handle(func(ctx context.Context, req any) (any, error) {
				return srv.SayHello(ctx, req.(*user.RequestHello))
			})
		})

		userGroup.GET("sse", func(c *gin.Context) {
			gateway2.SSE(c, &user.RequestHello{}).Handle(func(ctx context.Context, req any) (any, error) {
				return srv.SendSSE(ctx, req.(*user.RequestHello))
			})
		})
	}
}

func RegisterGRPCServer(engine *grpc.Server, srv *service.UserService) {
	user.RegisterUserServer(engine, srv)
	grpc_health_v1.RegisterHealthServer(engine, srv)
}
