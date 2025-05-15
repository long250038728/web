package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/long250038728/web/application/order/internal/domain"
	"github.com/long250038728/web/application/order/internal/repository"
	"github.com/long250038728/web/application/order/internal/router"
	"github.com/long250038728/web/application/order/internal/service"
	"github.com/long250038728/web/protoc"
	"github.com/long250038728/web/tool/app"
	"github.com/long250038728/web/tool/server/http"
	"github.com/long250038728/web/tool/server/rpc"
	"google.golang.org/grpc"
)

func main() {
	path := flag.String("path", "", "root path")
	flag.Parse()
	app.InitPathInfo(path)
	fmt.Println(Run(protoc.OrderService))
}

func Run(serverName string) error {
	util := app.NewUtil()
	port, ok := util.Port(serverName)
	if !ok {
		return fmt.Errorf("server %s is not bind port", serverName)
	}

	// 定义服务
	svc := service.NewService(
		service.SetDomain(domain.NewDomain(repository.NewRepository(util))),
	)

	r := router.NewRouter(util)

	opts := []app.Option{
		app.Servers( // 服务
			http.NewHttp(serverName, util.Info.IP, port.HttpPort, func(engine *gin.Engine) {
				r.RegisterHTTPServer(engine, svc)
			}),
			rpc.NewGrpc(serverName, util.Info.IP, port.GrpcPort, func(engine *grpc.Server) {
				r.RegisterGRPCServer(engine, svc)
			}),
		),
	}
	if register, err := util.Register(); err == nil {
		opts = append(opts, app.Register(register)) //服务注册 && 发现
	}
	if exporter, err := util.Exporter(); err == nil {
		opts = append(opts, app.Tracing(exporter, serverName)) //服务注册 && 发现
	}

	//启动服务
	application, err := app.NewApp(opts...)
	defer func() {
		application.Stop()
		util.Close()
	}()
	if err != nil {
		return err
	}

	//程序运行
	return application.Start()
}
