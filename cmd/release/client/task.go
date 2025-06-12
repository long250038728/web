package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/long250038728/web/tool/app"
	"github.com/long250038728/web/tool/configurator"
	"github.com/long250038728/web/tool/git"
	"github.com/long250038728/web/tool/hook"
	"github.com/long250038728/web/tool/jenkins"
	"github.com/long250038728/web/tool/persistence/orm"
	"github.com/long250038728/web/tool/ssh"
	"gorm.io/gorm"
	"os"
	"strings"
	"time"
)

type Task struct {
	outPath     string
	outFileName string
	tels        []string

	services *Svc

	git     git.Git
	jenkins *jenkins.Client
	orm     *orm.Gorm
	ssh     ssh.SSH
	hook    hook.Hook
}

type Opts func(o *Task)

func SetOutPath(path string) Opts {
	return func(o *Task) {
		o.outPath = path
	}
}

func SetFileName(fileName string) Opts {
	return func(o *Task) {
		o.outFileName = fileName
	}
}

//==========================================

func SetQyHook(hook hook.Hook, tels []string) Opts {
	return func(o *Task) {
		o.hook = hook
		o.tels = tels
	}
}

func SetGit(git git.Git) Opts {
	return func(o *Task) {
		o.git = git
	}
}

func SetJenkins(jenkins *jenkins.Client) Opts {
	return func(o *Task) {
		o.jenkins = jenkins
	}
}

func SetOrm(orm *orm.Gorm) Opts {
	return func(o *Task) {
		o.orm = orm
	}
}

func SetRemoteShell(ssh ssh.SSH) Opts {
	return func(o *Task) {
		o.ssh = ssh
	}
}

//==========================================

func NewTaskClient(opts ...Opts) *Task {
	o := &Task{
		outPath:     "./",
		outFileName: "json.json",
		services:    &Svc{Kobe: make([]string, 0, 0), Marx: make([]string, 0, 0)},
	}
	for _, opt := range opts {
		opt(o)
	}
	return o
}

//============================================================================================

func (o *Task) Build(ctx context.Context, source, target, svcPath string) error {
	if len(svcPath) > 0 {
		if err := configurator.NewYaml().Load(svcPath, &o.services); err != nil {
			return err
		}
	}
	var list []*requestInfo
	var err error
	if list, err = o.list(ctx, source, target); err != nil {
		_ = o.HookSend(ctx, "生成失败: \n"+err.Error())
		return err
	}
	if err = o.save(ctx, list); err != nil {
		_ = o.HookSend(ctx, "生成失败: \n"+err.Error())
		return err
	}

	var sqlBytes []byte
	var sqlErr error

	projectNames := make([]string, 0, len(list))
	for index, val := range list {
		projectNames = append(projectNames, fmt.Sprintf("%d.%s: %s", index+1, taskHashMap[val.Type], val.Project))
		if val.Type == TaskTypeSql {
			sqlBytes, sqlErr = json.MarshalIndent(val.Params["sql"], "", "	")
			if sqlErr != nil {
				_ = o.HookSend(ctx, sqlErr.Error())
			} else {
				_ = o.HookSend(ctx, string(sqlBytes))
			}
		}
	}
	_ = o.HookSend(ctx, "发布项目: \n"+strings.Join(projectNames, "\n\n"))
	return nil
}

func (o *Task) list(ctx context.Context, source, target string) ([]*requestInfo, error) {
	var address = make([]*requestInfo, 0, 100)

	if o.git == nil {
		return nil, errors.New("git client is null")
	}

	if len(o.services.SQL) > 0 {
		if o.orm == nil {
			return nil, errors.New("orm client is null")
		}
		sqls, err := o.orm.Parser(o.services.SQL)
		if err != nil {
			return nil, errors.New("sql parser is err: " + err.Error())
		}
		address = append(address, &requestInfo{Type: TaskTypeSql, Project: "sql", Params: map[string]any{"sql": sqls}})
	}

	if len(o.services.Shell) > 0 {
		//bash /soft/scripts/menu_script/run.sh 2024/sm1205/menu* 2024/sm1205/group* check
		address = append(address, &requestInfo{Type: TaskTypeRemoteShell, Project: fmt.Sprintf("bash /soft/scripts/menu_script/run.sh %s/menu* %s/group* prod", o.services.Shell, o.services.Shell)})
	}

	for _, addr := range ProductList {
		list, err := o.git.GetPR(ctx, addr, source, target)
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
		address = append(address, &requestInfo{Type: TaskTypeGit, Project: addr, Params: map[string]any{"num": list[0].Number}})

		//每个服务有两台服务器
		if addr == "zhubaoe-go/kobe" {
			address = append(address, &requestInfo{Type: TaskTypeRemoteShell, Project: "bash /data/project/tag.sh /data/project/kobe"})
			for _, svc := range o.services.Kobe {
				address = append(address, &requestInfo{Type: TaskTypeJenkins, Project: svc, Params: map[string]any{"BRANCH": "origin/master", "SYSTEM": "root@172.16.0.34"}}) //kobe

				if !strings.Contains(svc, "scheduler") {
					address = append(address, &requestInfo{Type: TaskTypeJenkins, Project: svc, Params: map[string]any{"BRANCH": "origin/master", "SYSTEM": "root@172.16.0.9"}})  //kobe-back
					address = append(address, &requestInfo{Type: TaskTypeJenkins, Project: svc, Params: map[string]any{"BRANCH": "origin/master", "SYSTEM": "root@172.16.0.17"}}) //kobe-back-2
				}
			}
		}

		//每个服务有一台服务器
		if addr == "zhubaoe/marx" {
			address = append(address, &requestInfo{Type: TaskTypeRemoteShell, Project: "bash /data/project/tag.sh /data/project/marx"})
			for _, svc := range o.services.Marx {
				address = append(address, &requestInfo{Type: TaskTypeJenkins, Project: svc})
			}
		}

		//有一个服务
		if addr == "zhubaoe/plato" {
			address = append(address, &requestInfo{Type: TaskTypeJenkins, Project: "plato-prod", Params: map[string]any{"BRANCH": "origin/master"}})
		}

		//有三个服务
		if addr == "zhubaoe/locke" {
			address = append(address, &requestInfo{Type: TaskTypeJenkins, Project: "locke-prod_32"})
			address = append(address, &requestInfo{Type: TaskTypeJenkins, Project: "locke-prod_64"})
			address = append(address, &requestInfo{Type: TaskTypeJenkins, Project: "locke-hot-prod-64"})
		}
	}
	return address, nil
}

//============================================================================================

func (o *Task) Request(ctx context.Context) error {
	b, err := os.ReadFile(o.outPath + o.outFileName)
	if err != nil {
		return err
	}

	var requestList []*requestInfo
	if err = json.Unmarshal(b, &requestList); err != nil {
		return err
	}

	// 查询有什么类型
	for _, val := range requestList {
		key := val.Type
		if key == TaskTypeGit && app.IsNil(o.git) {
			return errors.New("git is null")
		}
		if key == TaskTypeJenkins && app.IsNil(o.jenkins) {
			return errors.New("jenkins is null")
		}
		if key == TaskTypeSql && app.IsNil(o.orm) {
			return errors.New("orm is null")
		}
		if key == TaskTypeRemoteShell && app.IsNil(o.ssh) {
			return errors.New("remote ssh is null")
		}
	}
	for index, request := range requestList {
		//已经成功的就不再处理
		if request.Success {
			continue
		}

		startTime := time.Now().Local()
		var err error
		var other = "empty"

		switch request.Type {
		case TaskTypeGit: //合并
			err = o.git.Merge(ctx, request.Project, int32(request.Params["num"].(float64)))
		case TaskTypeShell: //shell
			project, ok := request.Params["project"].(string)
			if !ok {
				err = errors.New("shell script is error")
				break
			}
			other, err = ssh.NewLocalSSH().Run(fmt.Sprintf("%s %s", request.Project, project))
		case TaskTypeJenkins:
			// jenkins 可能会构建失败，所以重试 3次重试还不行就报错
			isSuccess := false
			if requestParams, jsonErr := json.Marshal(request.Params); jsonErr == nil {
				other = string(requestParams)
			}
			for i := 0; i < 3; i++ {
				err := o.jenkins.BlockBuild(ctx, request.Project, request.Params)
				if err == nil {
					isSuccess = true
					break
				}
				time.Sleep(time.Second * 2)
			}
			if !isSuccess {
				err = errors.New("jenkins build is failure")
			}
		case TaskTypeSql: //sql
			sql := request.Params["sql"].([]interface{})
			sqlList := make([]string, 0, len(sql))
			for _, s := range sql {
				str, ok := s.(string)
				if !ok {
					err = errors.New("sql is failure")
					break
				}
				sqlList = append(sqlList, str)
			}

			err = o.orm.Transaction(func(tx *gorm.DB) error {
				for _, sql := range sqlList {
					if err = tx.Exec(sql).Error; err != nil {
						return err
					}
				}
				return nil
			})
		case TaskTypeRemoteShell: //remote shell
			other, err = o.ssh.Run(request.Project)
		default:
			err = errors.New("type is err")
		}

		//============================================================================
		endTime := time.Now().Local()
		if err != nil {
			_ = o.HookSend(ctx, fmt.Sprintf("project: %s \nstatus: %s \nstart: %s   end: %s   sub: %s \nother: \n%s", request.Project, "failure", startTime.Format(time.TimeOnly), endTime.Format(time.TimeOnly), endTime.Sub(startTime).String(), err.Error()))
			return err
		}

		_ = o.HookSend(ctx, fmt.Sprintf("project: %s \nstatus: %s \nstart: %s   end: %s   sub: %s \nother: \n%s", request.Project, "success", startTime.Format(time.TimeOnly), endTime.Format(time.TimeOnly), endTime.Sub(startTime).String(), other))
		requestList[index].Success = true
		_ = o.save(ctx, requestList)
	}

	_ = o.HookSend(ctx, "The entire pipeline has been processed")

	return nil
}

//============================================================================================

func (o *Task) HookSend(ctx context.Context, text string) error {
	if len(text) == 0 {
		return errors.New("text is empty")
	}
	return o.hook.SendHook(ctx, text, o.tels)
}

func (o *Task) save(ctx context.Context, list []*requestInfo) error {
	b, err := json.MarshalIndent(list, "", "	")
	if err != nil {
		_ = o.HookSend(ctx, "生成失败: \n"+err.Error())
		return err
	}
	if err := os.WriteFile(o.outPath+o.outFileName, b, os.ModePerm); err != nil {
		_ = o.HookSend(ctx, "生成失败: \n"+err.Error())
		return err
	}
	return nil
}
