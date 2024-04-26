package client

import (
	"context"
	"errors"
	"fmt"
	"github.com/long250038728/web/tool/configurator"
	"github.com/long250038728/web/tool/git"
	"github.com/long250038728/web/tool/jenkins"
)

type requestInfo struct {
	Type    int
	Project string
	Params  map[string]any
}

type Online struct {
	git      git.Git
	jenkins  *jenkins.Client
	ctx      context.Context
	services *Svc
}

func NewOnlineClient(ctx context.Context, git git.Git, jenkins *jenkins.Client) *Online {
	return &Online{
		ctx:      ctx,
		git:      git,
		jenkins:  jenkins,
		services: &Svc{Kobe: make([]string, 0, 0), Marx: make([]string, 0, 0)},
	}
}

var productList = []string{
	"zhubaoe/locke",
	"zhubaoe-go/kobe",
	"zhubaoe/hume",
	"zhubaoe/socrates",
	"zhubaoe/aristotle",
	"fissiongeek/h5-sales",
	"zhubaoe/plato",
	"zhubaoe/marx",
}

func (o *Online) Build(source, target, svcPath string) error {
	list, err := o.list(source, target, svcPath)
	if err != nil {
		return err
	}
	fmt.Println(list)
	return nil
	//return o.request(list)
}

func (o *Online) list(source, target, svcPath string) ([]*requestInfo, error) {
	var address = make([]*requestInfo, 0, 100)

	if len(svcPath) > 0 {
		if err := configurator.NewYaml().Load(svcPath, &o.services); err != nil {
			return address, err
		}
	}

	for _, addr := range productList {
		list, err := o.git.GetPR(o.ctx, addr, source, target)
		if err != nil || len(list) != 1 {
			continue
		}
		if addr == "zhubaoe-go/kobe" && len(o.services.Kobe) == 0 {
			return address, errors.New("有kobe项目，但是未添加服务")
		}
		if addr == "zhubaoe/marx" && len(o.services.Marx) == 0 {
			return address, errors.New("有marx项目，但是未添加服务")
		}

		//调用合并分支
		address = append(address, &requestInfo{Type: 1, Project: addr, Params: map[string]any{"num": list[0].Number}})

		//两台服务器
		if addr == "zhubaoe-go/kobe" {
			for _, svc := range o.services.Kobe {
				address = append(address, &requestInfo{Type: 2, Project: svc, Params: map[string]any{"BRANCH": "origin/master", "SYSTEM": "root@172.16.0.34"}})
				address = append(address, &requestInfo{Type: 2, Project: svc, Params: map[string]any{"BRANCH": "origin/master", "SYSTEM": "root@172.16.0.9"}})
			}
		}

		// 一台服务器
		if addr == "zhubaoe/marx" {
			for _, svc := range o.services.Marx {
				address = append(address, &requestInfo{Type: 2, Project: svc})
			}
		}

		if addr == "zhubaoe/plato" {
			address = append(address, &requestInfo{Type: 2, Project: "plato-prod", Params: map[string]any{"BRANCH": "origin/master"}})
		}

		// 三个服务
		if addr == "zhubaoe-go/locke" {
			address = append(address, &requestInfo{Type: 2, Project: "locke-prod_32"})
			address = append(address, &requestInfo{Type: 2, Project: "locke-prod_64"})
			address = append(address, &requestInfo{Type: 2, Project: "locke-hot-prod-64"})
		}
	}
	return address, nil
}

func (o *Online) request(requestList []*requestInfo) error {
	for _, request := range requestList {
		switch request.Type {
		case 1:
			err := o.git.Merge(o.ctx, request.Project, request.Params["num"].(int32))
			if err != nil {
				return errors.New(fmt.Sprintf("%s %s %s", request.Project, "pr merge", err))
			}
		case 2:
			err := o.jenkins.BlockBuild(o.ctx, request.Project, request.Params)
			if err != nil {
				return errors.New(fmt.Sprintf("%s %s %s", request.Project, "block build", err))
			}
		default:
			return errors.New("type is err")
		}
	}

	return nil
}
