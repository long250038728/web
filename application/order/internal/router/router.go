package router

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/long250038728/web/application/order/internal/service"
	"github.com/long250038728/web/protoc/order"
	"github.com/long250038728/web/tool/app"
	"github.com/long250038728/web/tool/server/http/gateway"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func RegisterHTTPServer(engine *gin.Engine, srv *service.OrderService) {
	cache, _ := app.NewUtil().Cache()
	orderGroup := engine.Group("/order/order/").Use(gateway.BaseHandle(cache))
	{
		orderGroup.GET("detail", func(c *gin.Context) {
			gateway.Json(c, &order.OrderDetailRequest{}).Use(gateway.Limit(), gateway.Cache()).Handle(func(ctx context.Context, req any) (any, error) {
				return srv.OrderDetail(ctx, req.(*order.OrderDetailRequest))
			})
		})
	}
}

func RegisterGRPCServer(engine *grpc.Server, srv *service.OrderService) {
	order.RegisterOrderServer(engine, srv)
	grpc_health_v1.RegisterHealthServer(engine, srv)
}
