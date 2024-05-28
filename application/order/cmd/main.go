package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/long250038728/web/application/order/internal/domain"
	"github.com/long250038728/web/application/order/internal/repository"
	"github.com/long250038728/web/application/order/internal/service"
	"github.com/long250038728/web/application/order/router"
	"github.com/long250038728/web/protoc"
	"github.com/long250038728/web/tool/app"
	"github.com/long250038728/web/tool/server/http"
	"github.com/long250038728/web/tool/server/rpc"
	"google.golang.org/grpc"
)

func main() {
	path := flag.String("config", "", "config path")
	flag.Parse()
	app.InitCenterInfo(*path, protoc.OrderService)

	fmt.Println(Run())
}

func Run() error {
	util := app.NewUtil()
	orderService := service.NewService(
		service.SetOrderDomain(
			domain.NewOrderDomain(repository.NewOrderRepository(util)),
		),
	)

	opts := []app.Option{
		app.Servers( // 服务
			http.NewHttp(util.Info.ServerName, util.Info.IP, util.Info.HttpPort, func(engine *gin.Engine) {
				router.RegisterHTTPServer(engine, orderService)
			}),
			rpc.NewGrpc(util.Info.ServerName, util.Info.IP, util.Info.GrpcPort, func(engine *grpc.Server) {
				router.RegisterGRPCServer(engine, orderService)
			}),
		),
		app.Register(util.Register()),                      //服务注册 && 发现
		app.Tracing(util.Exporter(), util.Info.ServerName), //链路
	}

	//启动服务
	application, err := app.NewApp(opts...)
	defer application.Stop()
	if err != nil {
		return err
	}

	//程序运行
	return application.Start()
}
