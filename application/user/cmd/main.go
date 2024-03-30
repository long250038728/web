package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/long250038728/web/application/user/ddd/domain"
	"github.com/long250038728/web/application/user/ddd/repository"
	"github.com/long250038728/web/application/user/ddd/service"
	user "github.com/long250038728/web/application/user/protoc"
	"github.com/long250038728/web/application/user/router"
	"github.com/long250038728/web/tool/app"
	"github.com/long250038728/web/tool/server/http"
	"github.com/long250038728/web/tool/server/rpc"
	"google.golang.org/grpc"
	"os"
	"path/filepath"
)

func main() {
	fmt.Println(Run())
}

func Run() error {
	wd, _ := os.Getwd()
	conf, err := app.NewAppConfig(filepath.Join(wd, "application", "user", "config"))
	if err != nil {
		return err
	}

	util, err := app.NewUtil(conf)
	if err != nil {
		return err
	}

	// 定义服务
	userService := service.NewUserService(
		service.UserDomain(
			domain.NewUserDomain(repository.NewUserRepository(util)),
		),
	)

	//启动服务
	application, err := app.NewApp(
		app.Register(util.Register()),                      //服务注册 && 发现
		app.Tracing(util.Exporter(), util.Info.ServerName), //链路
		app.Servers( // 服务
			http.NewHttp(util.Info.ServerName, util.Info.IP, util.Info.HttpPort, func(engine *gin.Engine) {
				router.RegisterUserServer(engine, userService, util)
			}),
			rpc.NewGrpc(util.Info.ServerName, util.Info.IP, util.Info.GrpcPort, func(engine *grpc.Server) {
				user.RegisterUserServer(engine, userService)
			}),
		),
	)

	defer application.Stop()
	if err != nil {
		return err
	}

	//程序运行
	return application.Start()
}
