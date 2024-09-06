package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/long250038728/web/application/auth/internal/domain"
	"github.com/long250038728/web/application/auth/internal/repository"
	"github.com/long250038728/web/application/auth/internal/service"
	"github.com/long250038728/web/application/auth/router"
	"github.com/long250038728/web/protoc"
	"github.com/long250038728/web/tool/app"
	"github.com/long250038728/web/tool/server/http"
	"github.com/long250038728/web/tool/server/rpc"
	"google.golang.org/grpc"
)

func main() {
	path := flag.String("path", "", "root path")
	flag.Parse()
	app.InitPathInfo(*path, protoc.AuthService)

	fmt.Println(Run())
}

func Run() error {
	util := app.NewUtil()

	// 定义服务
	userService := service.NewService(
		service.SetDomain(domain.NewDomain(repository.NewRepository(util))),
	)

	opts := []app.Option{
		app.Servers( // 服务
			http.NewHttp(util.Info.ServerName, util.Info.IP, util.Info.HttpPort, func(engine *gin.Engine) {
				router.RegisterHTTPServer(engine, userService)
			}),
			rpc.NewGrpc(util.Info.ServerName, util.Info.IP, util.Info.GrpcPort, func(engine *grpc.Server) {
				router.RegisterGRPCServer(engine, userService)
			}),
		),
	}
	//if register, err := util.Register(); err == nil {
	//	opts = append(opts, app.Register(register)) //服务注册 && 发现
	//}
	if exporter, err := util.Exporter(); err == nil {
		opts = append(opts, app.Tracing(exporter, util.Info.ServerName)) //服务注册 && 发现
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
