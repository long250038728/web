package router

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/long250038728/web/application/auth/internal/service"
	"github.com/long250038728/web/protoc/auth_rpc"
	"github.com/long250038728/web/tool/app"
	"github.com/long250038728/web/tool/server/http"
	"github.com/long250038728/web/tool/server/http/gateway"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func RegisterHTTPServer(engine *gin.Engine, srv *service.Service) {
	cache, _ := app.NewUtil().Cache()
	userGroup := engine.Group("/authorization/user/").Use(http.BaseHandle(cache), http.LimitHandle(cache))
	{
		userGroup.POST("login", func(c *gin.Context) {
			gateway.Json(c, &auth_rpc.LoginRequest{}).Use(gateway.Limit(), gateway.Cache()).Handle(func(ctx context.Context, req any) (any, error) {
				return srv.Login(ctx, req.(*auth_rpc.LoginRequest))
			})
		})

		userGroup.POST("refresh", func(c *gin.Context) {
			gateway.Json(c, &auth_rpc.RefreshRequest{}).Use(gateway.Limit(), gateway.Cache()).Handle(func(ctx context.Context, req any) (any, error) {
				return srv.Refresh(ctx, req.(*auth_rpc.RefreshRequest))
			})
		})
	}
}

func RegisterGRPCServer(engine *grpc.Server, srv *service.Service) {
	auth_rpc.RegisterAuthServer(engine, srv)
	grpc_health_v1.RegisterHealthServer(engine, srv)
}
