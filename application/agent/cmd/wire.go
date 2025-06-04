//go:build wireinject
// +build wireinject

// 仅当编译时启用 wireinject 标签时，该文件才会被编译。
// 旧版本语法（Go 1.17 之前），功能与 //go:build 相同。
package main

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

func InitServer(util *app.Util) *service.AgentService {
	panic(wire.Build(service.ProviderSet, domain.ProviderSet, repository.ProviderSet))
}
