package router

import (
	"github.com/gin-gonic/gin"
	"github.com/long250038728/web/application/user/internal/service"
	"github.com/long250038728/web/protoc/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
)

//go func() {
//	log.Println(http.ListenAndServe("localhost:6060", nil))  //"net/http/pprof"
//}()

func RegisterHTTPServer(engine *gin.Engine, srv *service.UserService) {
	//middleware := tool.NewHttpTools()
	//
	//engine.GET("/", func(gin *gin.Context) {
	//	var request user.RequestHello
	//	middleware.JSON(gin, &request, func(ctx context.Context) (interface{}, error) {
	//		return srv.SayHello(ctx, &request)
	//	})
	//})
	//
	//engine.GET("/user/", func(gin *gin.Context) {
	//	var request user.RequestHello
	//	middleware.JSON(gin, &request, func(ctx context.Context) (interface{}, error) {
	//		return srv.SayHello(ctx, &request)
	//	})
	//})
	//
	//engine.POST("/user/xls", func(gin *gin.Context) {
	//	var request user.RequestHello
	//	middleware.File(gin, &request, func(ctx context.Context) (interface{}, error) {
	//		return srv.SayHello(ctx, &request)
	//	})
	//})
	//
	//engine.GET("/user/sse", func(gin *gin.Context) {
	//	var request user.RequestHello
	//	middleware.SSE(gin, &request, func(ctx context.Context) (interface{}, error) {
	//		return srv.SendSSE(ctx, &request)
	//	})
	//})
}

func RegisterGRPCServer(engine *grpc.Server, srv *service.UserService) {
	user.RegisterUserServer(engine, srv)
	grpc_health_v1.RegisterHealthServer(engine, srv)
}
