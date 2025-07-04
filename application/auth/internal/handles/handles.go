package handles

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/long250038728/web/application/auth/internal/service"
	"github.com/long250038728/web/protoc/auth"
	"github.com/long250038728/web/tool/app"
	"github.com/long250038728/web/tool/server/http/gateway"
	"github.com/long250038728/web/tool/server/http/gateway/middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
)

type Handles struct {
	util *app.Util
}

func NewHandles(util *app.Util) *Handles {
	return &Handles{util: util}
}

func (r *Handles) RegisterHTTPServer(engine *gin.Engine, srv *service.Auth) {
	authorized, err := r.util.Auth()
	if err != nil {
		panic(err)
	}

	// 服务/领域/接口
	userGroup := engine.Group("/auth/user/").Use(middleware.BaseHandle(authorized))
	{
		userGroup.POST("login", func(c *gin.Context) {
			gateway.Json(c, &auth.LoginRequest{}).Use(middleware.Limit()).Handle(func(ctx context.Context, req any) (any, error) {
				return srv.Login(ctx, req.(*auth.LoginRequest))
			})
		})

		userGroup.POST("refresh", func(c *gin.Context) {
			gateway.Json(c, &auth.RefreshRequest{}).Use(middleware.Limit()).Handle(func(ctx context.Context, req any) (any, error) {
				return srv.Refresh(ctx, req.(*auth.RefreshRequest))
			})
		})
	}
}

func (r *Handles) RegisterGRPCServer(engine *grpc.Server, srv *service.Auth) {
	auth.RegisterAuthServer(engine, srv)
	grpc_health_v1.RegisterHealthServer(engine, srv)
}
