package router

import (
	"github.com/gin-gonic/gin"
	"github.com/long250038728/web/application/order/internal/service"
	"github.com/long250038728/web/protoc/order"
	"github.com/long250038728/web/tool/app"
	"github.com/long250038728/web/tool/server/http"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func RegisterHTTPServer(engine *gin.Engine, srv *service.OrderService) {
	cache, _ := app.NewUtil().Cache()
	orderGroup := engine.Group("/order/order/").Use(http.BaseHandle(cache), http.LimitHandle(cache))
	{
		orderGroup.GET("detail", func(c *gin.Context) {
			var request order.OrderDetailRequest
			http.NewGateway().JSON(c, &request, func() (interface{}, error) {
				return srv.OrderDetail(c.Request.Context(), &request)
			})
		})
	}
}

func RegisterGRPCServer(engine *grpc.Server, srv *service.OrderService) {
	order.RegisterOrderServer(engine, srv)
	grpc_health_v1.RegisterHealthServer(engine, srv)
}
