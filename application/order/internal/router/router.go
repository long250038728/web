package router

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/long250038728/web/application/order/internal/service"
	"github.com/long250038728/web/protoc/order"
	"github.com/long250038728/web/tool/app"
	"github.com/long250038728/web/tool/server/http/gateway"
	"github.com/long250038728/web/tool/server/http/gateway/middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func RegisterHTTPServer(engine *gin.Engine, srv *service.OrderService) {
	authorized, err := app.NewUtil().Auth()
	if err != nil {
		panic(err)
	}

	orderGroup := engine.Group("/order/order/").Use(middleware.BaseHandle(authorized))
	{
		orderGroup.GET("detail", func(c *gin.Context) {
			gateway.Json(c, &order.OrderDetailRequest{}).Use(middleware.Limit()).Handle(func(ctx context.Context, req any) (any, error) {
				return srv.OrderDetail(ctx, req.(*order.OrderDetailRequest))
			})
		})
	}
}

func RegisterGRPCServer(engine *grpc.Server, srv *service.OrderService) {
	order.RegisterOrderServer(engine, srv)
	grpc_health_v1.RegisterHealthServer(engine, srv)
}
