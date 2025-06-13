//go:build wireinject
// +build wireinject

// 仅当编译时启用 wireinject 标签时，该文件才会被编译。
// 旧版本语法（Go 1.17 之前），功能与 //go:build 相同。
package wrie

// go get github.com/google/wire/cmd/wire@v0.6.0
// go install github.com/google/wire/cmd/wire
// cd .../cmd && wire     //生成wire_gen.go文件

import (
	"github.com/google/wire"
	"github.com/long250038728/web/application/agent/internal/domain"
	"github.com/long250038728/web/application/agent/internal/repository"
	"github.com/long250038728/web/application/agent/internal/service"
	"github.com/long250038728/web/tool/app"
)

// var ProviderSet = wire.NewSet(xxxNew()) 把New方法放入Set列表中，通过Build方法组织成有关系的链。
// 在xxxNew中如果有参数是无法直接初始化而需要通过传递的方式，那么传入到该方法中（如InitServer方法）
// 如果某个xxxNew() 返回了error或其他的即没有其他的xxxNew()能接收，那么就在init方法返回（如InitServer方法）

func InitServer(util *app.Util) *service.AgentService {
	panic(wire.Build(service.ProviderSet, domain.ProviderSet, repository.ProviderSet))
}
