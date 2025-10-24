package main

import (
	"flag"
	"fmt"
	"github.com/long250038728/web/application/user/cmd/server"
	"github.com/long250038728/web/protoc"
	"github.com/long250038728/web/tool/app"
)

// main 服务运行
// config配置文件路径 离项目的启动方式越近越优先（1.启动指定路径   2.当前目录下   3.环境变量）
func main() {
	path := flag.String("config", "", "config path")
	flag.Parse()

	app.InitPathInfo(path)
	fmt.Println(server.Run(protoc.UserService))
}
