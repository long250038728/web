package router

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/long250038728/web/application/order/internal/service"
	"github.com/long250038728/web/protoc/order"
	"github.com/long250038728/web/tool/app"
	"github.com/long250038728/web/tool/auth"
	"github.com/long250038728/web/tool/limiter"
	"github.com/long250038728/web/tool/server/http/tool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
	_ "net/http/pprof"
	"time"
)

//go func() {
//	log.Println(http.ListenAndServe("localhost:6060", nil))  //"net/http/pprof"
//}()

func RegisterHTTPServer(engine *gin.Engine, srv *service.OrderService) {
	opts := []tool.MiddlewareOpt{
		tool.Limiter( //设置限流
			limiter.NewCacheLimiter(
				app.NewUtil().Cache(),
				limiter.SetExpiration(time.Second), limiter.SetTimes(10),
			),
		),

		tool.Auth( //设置权限（权限信息可从数据库获取文件获取）
			auth.NewAuth(
				app.NewUtil().Cache(),
				auth.WhiteList(auth.NewLocalWhite([]string{"/", "/user/", "/user/hello", "/user/hello2", "/user/hello3"}, []string{})),
			),
		),

		tool.Error( //设置错误（错误信息可从数据库获取文件获取）
			[]*tool.MiddleErr{}, //可以通过数据库处理
		),
	}
	middleware := tool.NewMiddlewarePool(opts...)

	engine.GET("/", func(gin *gin.Context) {
		var request order.OrderDetailRequest
		middleware.JSON(gin, &request, func(ctx context.Context) (interface{}, error) {
			return srv.OrderDetail(ctx, &request)
		})
	})

	engine.POST("/order/detail", func(gin *gin.Context) {
		var request order.OrderDetailRequest
		middleware.File(gin, &request, func(ctx context.Context) (interface{}, error) {
			return srv.OrderDetail(ctx, &request)
		})
	})
}

func RegisterGRPCServer(engine *grpc.Server, srv *service.OrderService) {
	order.RegisterOrderServer(engine, srv)
	grpc_health_v1.RegisterHealthServer(engine, srv)
}
