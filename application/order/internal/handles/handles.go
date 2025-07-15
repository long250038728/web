package handles

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

type Handles struct {
	util *app.Util
}

func NewHandles(util *app.Util) *Handles {
	return &Handles{util: util}
}

func (r *Handles) RegisterHTTPServer(engine *gin.Engine, srv *service.Order) {
	authorized, err := r.util.Auth()
	if err != nil {
		panic(err)
	}

	// 服务/领域/接口
	orderGroup := engine.Group("/order/order/").Use(middleware.BaseHandle(authorized))
	{
		orderGroup.GET("detail", func(c *gin.Context) {
			gateway.Json(c, &order.OrderDetailRequest{}).Use().Handle(func(ctx context.Context, req any) (any, error) {
				return srv.OrderDetail(ctx, req.(*order.OrderDetailRequest))
			})
		})
	}
}

func (r *Handles) RegisterGRPCServer(engine *grpc.Server, srv *service.Order) {
	order.RegisterOrderServer(engine, srv)
	grpc_health_v1.RegisterHealthServer(engine, srv)
}
