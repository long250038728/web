package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/long250038728/web/application/order/internal/domain"
	"github.com/long250038728/web/application/order/internal/handles"
	"github.com/long250038728/web/application/order/internal/repository"
	"github.com/long250038728/web/application/order/internal/service"
	"github.com/long250038728/web/tool/app"
	"github.com/long250038728/web/tool/server/http"
	"github.com/long250038728/web/tool/server/rpc"
	"google.golang.org/grpc"
)

func Run(serverName string) error {
	util := app.NewUtil()
	if !util.CheckPort(serverName) {
		return fmt.Errorf("server %s is not bind port", serverName)
	}

	// 定义服务
	svc := service.NewService(
		service.SetDomain(domain.NewOrderDomain(repository.NewOrderRepository(util))),
	)

	handle := handles.NewHandles(util)

	opts := []app.Option{
		app.Servers( // 服务
			http.NewHttp(serverName, util.Info.IP, util.Port(serverName).HttpPort, func(engine *gin.Engine) {
				handle.RegisterHTTPServer(engine, svc)
			}),
			rpc.NewGrpc(serverName, util.Info.IP, util.Port(serverName).GrpcPort, func(engine *grpc.Server) {
				handle.RegisterGRPCServer(engine, svc)
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
