package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/long250038728/web/tool/configurator"
	"github.com/long250038728/web/tool/git"
	"github.com/long250038728/web/tool/jenkins"
	"github.com/long250038728/web/tool/persistence/orm"
	"os"
	"os/exec"
	"strings"
)

type requestInfo struct {
	Type    int32          `json:"type"`
	Project string         `json:"project"`
	Params  map[string]any `json:"params"`
	Num     int32          `json:"num"`
}

type Online struct {
	git     git.Git
	jenkins *jenkins.Client
	orm     *orm.Gorm

	ctx      context.Context
	services *Svc
	sql      string
}

const (
	OnlineTypeGit     int32 = 1 //git
	OnlineTypeJenkins int32 = 2 //构建
	OnlineTypeShell   int32 = 3 //脚本
	OnlineTypeSql     int32 = 4 //数据库
)

func NewOnlineClient(ctx context.Context, git git.Git, jenkins *jenkins.Client, orm *orm.Gorm) *Online {
	return &Online{
		ctx:      ctx,
		git:      git,
		jenkins:  jenkins,
		orm:      orm,
		services: &Svc{Kobe: make([]string, 0, 0), Marx: make([]string, 0, 0)},
		sql:      "",
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

func (o *Online) Build(source, target, svcPath, sqlPath string) error {
	if len(svcPath) > 0 {
		if err := configurator.NewYaml().Load(svcPath, &o.services); err != nil {
			return err
		}
	}

	if len(sqlPath) > 0 {
		b, err := os.ReadFile(sqlPath)
		if err != nil {
			return err
		}
		o.sql = string(b)
	}

	list, err := o.list(source, target)
	if err != nil {
		return err
	}

	b, err := json.MarshalIndent(list, "", "	")
	if err != nil {
		return err
	}
	err = os.WriteFile("/Users/linlong/Desktop/web/application/third_party/client/project_list.md", b, os.ModePerm)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}

func (o *Online) list(source, target string) ([]*requestInfo, error) {
	var address = make([]*requestInfo, 0, 100)

	if len(o.sql) > 0 {
		for _, sql := range strings.Split(o.sql, ";") {
			address = append(address, &requestInfo{Type: OnlineTypeSql, Project: strings.Replace(sql, "\n", " ", -1)})
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
		address = append(address, &requestInfo{Type: OnlineTypeGit, Project: addr, Num: list[0].Number})

		//两台服务器
		if addr == "zhubaoe-go/kobe" {
			address = append(address, &requestInfo{Type: OnlineTypeShell, Project: "/Users/linlong/Desktop/web/application/third_party/client/change_tag.sh", Params: map[string]any{"project": "kobe"}})
			for _, svc := range o.services.Kobe {
				address = append(address, &requestInfo{Type: OnlineTypeJenkins, Project: svc, Params: map[string]any{"BRANCH": "origin/master", "SYSTEM": "root@172.16.0.34"}})
				address = append(address, &requestInfo{Type: OnlineTypeJenkins, Project: svc, Params: map[string]any{"BRANCH": "origin/master", "SYSTEM": "root@172.16.0.9"}})
			}
		}

		// 一台服务器
		if addr == "zhubaoe/marx" {
			address = append(address, &requestInfo{Type: OnlineTypeShell, Project: "/Users/linlong/Desktop/web/application/third_party/client/change_tag.sh", Params: map[string]any{"project": "marx"}})
			for _, svc := range o.services.Marx {
				address = append(address, &requestInfo{Type: OnlineTypeJenkins, Project: svc})
			}
		}

		if addr == "zhubaoe/plato" {
			address = append(address, &requestInfo{Type: OnlineTypeJenkins, Project: "plato-prod", Params: map[string]any{"BRANCH": "origin/master"}})
		}

		// 三个服务
		if addr == "zhubaoe/locke" {
			address = append(address, &requestInfo{Type: OnlineTypeJenkins, Project: "locke-prod_32"})
			address = append(address, &requestInfo{Type: OnlineTypeJenkins, Project: "locke-prod_64"})
			address = append(address, &requestInfo{Type: OnlineTypeJenkins, Project: "locke-hot-prod-64"})
		}
	}
	return address, nil
}

func (o *Online) Request(requestList []*requestInfo) error {
	for _, request := range requestList {
		switch request.Type {
		case OnlineTypeGit: //合并
			err := o.git.Merge(o.ctx, request.Project, request.Num)
			if err != nil {
				fmt.Printf("=================== %s  err ===============\n", err)
				return errors.New(fmt.Sprintf("%s %s %s", request.Project, "pr merge", err))
			}
		case OnlineTypeShell: //shell
			project, ok := request.Params["project"].(string)
			if !ok {
				return errors.New(fmt.Sprintf("%s %s", request.Project, "get project name is err"))
			}
			err := exec.Command("sh", request.Project, project).Run()
			if err != nil {
				fmt.Printf("=================== %s  err ===============\n", err)
				return errors.New(fmt.Sprintf("%s %s %s", request.Project, "executing command", err))
			}
		case OnlineTypeJenkins: //jenkins
			err := o.jenkins.BlockBuild(o.ctx, request.Project, request.Params)
			if err != nil {
				fmt.Printf("=================== %s  err ===============\n", err)
				return errors.New(fmt.Sprintf("%s %s %s", request.Project, "block build", err))
			}
		case OnlineTypeSql: //sql
			err := o.orm.Exec(request.Project).Error
			if err != nil {
				fmt.Printf("=================== %s  err ===============\n", err)
				return errors.New(fmt.Sprintf("%s %s %s", request.Project, "sql build", err))
			}
		default:
			return errors.New("type is err")
		}
	}

	return nil
}
