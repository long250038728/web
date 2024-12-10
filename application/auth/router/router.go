package router

import (
	"github.com/gin-gonic/gin"
	"github.com/long250038728/web/application/auth/internal/service"
	"github.com/long250038728/web/protoc/auth_rpc"
	"github.com/long250038728/web/tool/app"
	"github.com/long250038728/web/tool/server/http"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func RegisterHTTPServer(engine *gin.Engine, srv *service.Service) {
	cache, _ := app.NewUtil().Cache()
	userGroup := engine.Group("/authorization/user/").Use(http.BaseHandle(cache), http.LimitHandle(cache))
	{
		userGroup.POST("login", func(c *gin.Context) {
			var request auth_rpc.LoginRequest
			http.NewGateway().JSON(c, &request, func() (interface{}, error) {
				return srv.Login(c.Request.Context(), &request)
			})
		})

		userGroup.POST("refresh", func(c *gin.Context) {
			var request auth_rpc.RefreshRequest
			http.NewGateway().JSON(c, &request, func() (interface{}, error) {
				return srv.Refresh(c.Request.Context(), &request)
			})
		})
	}
}

func RegisterGRPCServer(engine *grpc.Server, srv *service.Service) {
	auth_rpc.RegisterAuthServer(engine, srv)
	grpc_health_v1.RegisterHealthServer(engine, srv)
}
