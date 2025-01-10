package router

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/long250038728/web/application/agent/internal/service"
	"github.com/long250038728/web/protoc/agent"
	"github.com/long250038728/web/tool/app"
	"github.com/long250038728/web/tool/server/http/gateway"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func RegisterHTTPServer(engine *gin.Engine, srv *service.Service) {
	cache, _ := app.NewUtil().Cache()
	infoGroup := engine.Group("/agent/info/").Use(gateway.BaseHandle(cache))
	{
		infoGroup.GET("logs", func(c *gin.Context) {
			gateway.Json(c, &agent.LogsRequest{}).Use(gateway.Limit(), gateway.Cache()).Handle(func(ctx context.Context, req any) (any, error) {
				return srv.Logs(ctx, req.(*agent.LogsRequest))
			})
		})

		infoGroup.GET("events", func(c *gin.Context) {
			gateway.Json(c, &agent.EventsRequest{}).Use(gateway.Limit(), gateway.Cache()).Handle(func(ctx context.Context, req any) (any, error) {
				return srv.Events(ctx, req.(*agent.EventsRequest))
			})
		})

		infoGroup.GET("resources", func(c *gin.Context) {
			gateway.Json(c, &agent.ResourcesRequest{}).Use(gateway.Limit(), gateway.Cache()).Handle(func(ctx context.Context, req any) (any, error) {
				return srv.Resources(ctx, req.(*agent.ResourcesRequest))
			})
		})
	}
}

func RegisterGRPCServer(engine *grpc.Server, srv *service.Service) {
	agent.RegisterAgentServer(engine, srv)
	grpc_health_v1.RegisterHealthServer(engine, srv)
}
