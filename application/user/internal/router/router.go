package router

import (
	"github.com/gin-gonic/gin"
	"github.com/long250038728/web/application/user/internal/service"
	"github.com/long250038728/web/protoc/user"
	"github.com/long250038728/web/tool/app"
	"github.com/long250038728/web/tool/server/http"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
)

//go func() {
//	log.Println(http.ListenAndServe("localhost:6060", nil))  //"net/http/pprof"
//}()

func RegisterHTTPServer(engine *gin.Engine, srv *service.UserService) {
	cache, _ := app.NewUtil().Cache()
	gateway := http.NewGateway()

	userGroup := engine.Group("/user/user/").Use(http.BaseHandle(cache), http.LimitHandle(cache))
	{
		userGroup.GET("say_hello", func(c *gin.Context) {
			var request user.RequestHello
			gateway.JSON(c, &request, func() (interface{}, error) {
				return srv.SayHello(c.Request.Context(), &request)
			})
		})

		//userGroup.GET("say", func(c *gin.Context) {
		//	// 添加拦截器链并处理请求
		//	var request user.RequestHello
		//	resp, err := g.NewAAA(c, &request).Handle(func(ctx context.Context, req map[string]any) (any, error) {
		//		return srv.SayHello(c.Request.Context(), &request)
		//	})
		//	fmt.Println(resp, err)
		//})

		//userGroup.GET("say_hello", func(c *gin.Context) {
		//	var request user.RequestHello
		//	gateway.Run(http.Json(), c, &request, func() (interface{}, error) {
		//		return srv.SayHello(c.Request.Context(), &request)
		//	})
		//})
		//userGroup.GET("xls", func(c *gin.Context) {
		//	var request user.RequestHello
		//	gateway.Run(http.File(), c, &request, func() (interface{}, error) {
		//		return srv.SayHello(c.Request.Context(), &request)
		//	})
		//})
		//userGroup.GET("sse", func(c *gin.Context) {
		//	var request user.RequestHello
		//	gateway.Run(http.File(), c, &request, func() (interface{}, error) {
		//		return srv.SendSSE(c.Request.Context(), &request)
		//	})
		//})
	}
}

func RegisterGRPCServer(engine *grpc.Server, srv *service.UserService) {
	user.RegisterUserServer(engine, srv)
	grpc_health_v1.RegisterHealthServer(engine, srv)
}
