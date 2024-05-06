package client

import (
	"context"
	"errors"
	"fmt"
	"github.com/long250038728/web/tool/configurator"
	"github.com/long250038728/web/tool/git"
	"github.com/long250038728/web/tool/jenkins"
	"os/exec"
)

type requestInfo struct {
	Type    int32
	Project string
	Params  map[string]any
}

type Online struct {
	git      git.Git
	jenkins  *jenkins.Client
	ctx      context.Context
	services *Svc
}

const (
	OnlineTypeGit     int32 = 1 //git
	OnlineTypeJenkins int32 = 2 //构建
	OnlineTypeShell   int32 = 3 //脚本
)

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
	if len(svcPath) > 0 {
		if err := configurator.NewYaml().Load(svcPath, &o.services); err != nil {
			return err
		}
	}

	list, err := o.list(source, target)
	if err != nil {
		return err
	}
	fmt.Println(list)
	return nil
	//return o.request(list)
}

func (o *Online) list(source, target string) ([]*requestInfo, error) {
	var address = make([]*requestInfo, 0, 100)
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
		address = append(address, &requestInfo{Type: OnlineTypeGit, Project: addr, Params: map[string]any{"num": list[0].Number}})

		//两台服务器
		if addr == "zhubaoe-go/kobe" {
			address = append(address, &requestInfo{Type: OnlineTypeShell, Project: "./change_tag.sh kobe", Params: nil})
			for _, svc := range o.services.Kobe {
				address = append(address, &requestInfo{Type: OnlineTypeJenkins, Project: svc, Params: map[string]any{"BRANCH": "origin/master", "SYSTEM": "root@172.16.0.34"}})
				address = append(address, &requestInfo{Type: OnlineTypeJenkins, Project: svc, Params: map[string]any{"BRANCH": "origin/master", "SYSTEM": "root@172.16.0.9"}})
			}
		}

		// 一台服务器
		if addr == "zhubaoe/marx" {
			address = append(address, &requestInfo{Type: OnlineTypeShell, Project: "./change_tag.sh marx", Params: nil})
			for _, svc := range o.services.Marx {
				address = append(address, &requestInfo{Type: OnlineTypeJenkins, Project: svc})
			}
		}

		if addr == "zhubaoe/plato" {
			address = append(address, &requestInfo{Type: OnlineTypeJenkins, Project: "plato-prod", Params: map[string]any{"BRANCH": "origin/master"}})
		}

		// 三个服务
		if addr == "zhubaoe-go/locke" {
			address = append(address, &requestInfo{Type: OnlineTypeJenkins, Project: "locke-prod_32"})
			address = append(address, &requestInfo{Type: OnlineTypeJenkins, Project: "locke-prod_64"})
			address = append(address, &requestInfo{Type: OnlineTypeJenkins, Project: "locke-hot-prod-64"})
		}
	}
	return address, nil
}

func (o *Online) request(requestList []*requestInfo) error {
	for _, request := range requestList {

		fmt.Printf("=================== %s ===============", request.Project)

		switch request.Type {
		case OnlineTypeGit: //合并
			err := o.git.Merge(o.ctx, request.Project, request.Params["num"].(int32))
			if err != nil {
				fmt.Printf("=================== %s  err ===============", err)
				return errors.New(fmt.Sprintf("%s %s %s", request.Project, "pr merge", err))
			}
		case OnlineTypeShell: //shell
			out, err := exec.Command("sh", "-c", request.Project).Output()
			if err != nil {
				fmt.Printf("=================== %s  err ===============", err)
				return errors.New(fmt.Sprintf("%s %s %s", request.Project, "executing command", err))
			}
			// 输出命令执行结果
			fmt.Println("Command output:", string(out))
		case OnlineTypeJenkins: //jenkins
			err := o.jenkins.BlockBuild(o.ctx, request.Project, request.Params)
			if err != nil {
				fmt.Printf("=================== %s  err ===============", err)
				return errors.New(fmt.Sprintf("%s %s %s", request.Project, "block build", err))
			}
		default:
			return errors.New("type is err")
		}

		fmt.Printf("=================== %s  ok ===============", request.Project)
	}

	return nil
}
