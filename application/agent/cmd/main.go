package main

import (
	"flag"
	"fmt"
	"github.com/long250038728/web/application/agent/cmd/server"
	"github.com/long250038728/web/protoc"
	"github.com/long250038728/web/tool/app"
)

// main 1.默认读取命令行config配置信息，2.读取Config环境变量，3.获取当前路径下面的config文件
func main() {
	path := flag.String("config", "", "config path")
	flag.Parse()

	app.InitPathInfo(path)
	fmt.Println(server.Run(protoc.AgentService))
}
