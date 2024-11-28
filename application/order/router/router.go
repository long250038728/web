package router

import (
	"github.com/gin-gonic/gin"
	"github.com/long250038728/web/application/order/internal/service"
	"github.com/long250038728/web/protoc/order"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
)

//go func() {
//	log.Println(http.ListenAndServe("localhost:6060", nil))  //"net/http/pprof"
//}()

func RegisterHTTPServer(engine *gin.Engine, srv *service.OrderService) {
	//middleware := tool.NewHttpTools()
	//
	//engine.GET("/", func(gin *gin.Context) {
	//	var request order.OrderDetailRequest
	//	middleware.JSON(gin, &request, func(ctx context.Context) (interface{}, error) {
	//		return srv.OrderDetail(ctx, &request)
	//	})
	//})
	//
	//engine.POST("/order/detail", func(gin *gin.Context) {
	//	var request order.OrderDetailRequest
	//	middleware.File(gin, &request, func(ctx context.Context) (interface{}, error) {
	//		return srv.OrderDetail(ctx, &request)
	//	})
	//})
}

func RegisterGRPCServer(engine *grpc.Server, srv *service.OrderService) {
	order.RegisterOrderServer(engine, srv)
	grpc_health_v1.RegisterHealthServer(engine, srv)
}
