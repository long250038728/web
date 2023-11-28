package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/long250038728/web/application/user/domain"
	user "github.com/long250038728/web/application/user/protoc"
	"github.com/long250038728/web/application/user/repository"
	"github.com/long250038728/web/application/user/router"
	"github.com/long250038728/web/application/user/service"
	"github.com/long250038728/web/tool/app"
	"github.com/long250038728/web/tool/server/http"
	"github.com/long250038728/web/tool/server/rpc"
	"github.com/long250038728/web/tool/tracing/opentelemetry"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println(Run())
}

func Run() error {
	//获取app配置信息
	config, err := app.NewConfig()
	if err != nil {
		return err
	}

	////创建consul客户端
	//register, err := consul.NewConsulRegister(config.RegisterAddr)
	//if err != nil {
	//	return err
	//}

	//创建链路
	exporter, err := opentelemetry.NewJaegerExporter(config.TracingUrl)
	if err != nil {
		return err
	}

	// 定义服务
	userService := service.NewUserService(
		service.UserDomain(domain.NewUserDomain(repository.NewUserRepository())),
	)

	//启动服务
	application, err := app.NewApp(
		// 服务
		app.Servers(
			http.NewHttp(config.ServerName, config.IP, config.HttpPort, func(engine *gin.Engine) {
				router.RegisterUserServerServer(engine, userService)
			}),
			rpc.NewGrpc(config.ServerName, config.IP, config.GrpcPort, func(engine *grpc.Server) {
				user.RegisterUserServerServer(engine, userService)
			}),
		),

		//服务注册 && 发现
		//app.Register(register),

		//链路
		app.Tracing(exporter, config.ServerName),
	)
	defer application.Stop()
	if err != nil {
		return err
	}

	//程序运行
	return application.Start()
}
