package handles

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/long250038728/web/application/agent/internal/service"
	"github.com/long250038728/web/protoc/agent"
	"github.com/long250038728/web/tool/app"
	"github.com/long250038728/web/tool/server/http/gateway"
	"github.com/long250038728/web/tool/server/http/gateway/middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
)

var ProviderSet = wire.NewSet(NewHandles)

type Handles struct {
	srv  *service.Agent
	util *app.Util
}

func NewHandles(srv *service.Agent, util *app.Util) *Handles {
	return &Handles{srv: srv, util: util}
}

func (r *Handles) RegisterHTTPServer(engine *gin.Engine) {
	authorized, err := r.util.Auth()
	if err != nil {
		panic(err)
	}

	// 服务/领域/接口
	infoGroup := engine.Group("/agent/info/").Use(middleware.BaseHandle(authorized))
	{
		infoGroup.GET("logs", func(c *gin.Context) {
			gateway.Json(c, &agent.LogsRequest{}).Use(middleware.Limit()).Handle(func(ctx context.Context, req any) (any, error) {
				return r.srv.Logs(ctx, req.(*agent.LogsRequest))
			})
		})

		infoGroup.GET("events", func(c *gin.Context) {
			gateway.Json(c, &agent.EventsRequest{}).Use(middleware.Limit()).Handle(func(ctx context.Context, req any) (any, error) {
				return r.srv.Events(ctx, req.(*agent.EventsRequest))
			})
		})

		infoGroup.GET("resources", func(c *gin.Context) {
			gateway.Json(c, &agent.ResourcesRequest{}).Use(middleware.Limit()).Handle(func(ctx context.Context, req any) (any, error) {
				return r.srv.Resources(ctx, req.(*agent.ResourcesRequest))
			})
		})
	}
}

func (r *Handles) RegisterGRPCServer(engine *grpc.Server) {
	agent.RegisterAgentServer(engine, r.srv)
	grpc_health_v1.RegisterHealthServer(engine, r.srv)
}
