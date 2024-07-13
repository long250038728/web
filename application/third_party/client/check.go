package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/long250038728/web/tool/app"
	"github.com/long250038728/web/tool/configurator"
	"github.com/long250038728/web/tool/git"
	"github.com/long250038728/web/tool/jenkins"
	"github.com/long250038728/web/tool/persistence/orm"
	"github.com/long250038728/web/tool/qy_hook"
	"github.com/long250038728/web/tool/ssh"
	"gorm.io/gorm"
	"os"
	"os/exec"
	"strings"
	"time"
)

type Check struct {
	outPath     string
	outFileName string
	hook        string

	services *Svc

	git     git.Git
	jenkins *jenkins.Client
	orm     *orm.Gorm
	ssh     *ssh.SSH
}

type OptsCheck func(o *Check)

func SetOutPathCheck(path string) OptsCheck {
	return func(o *Check) {
		o.outPath = path
	}
}

func SetFileNameCheck(fileName string) OptsCheck {
	return func(o *Check) {
		o.outFileName = fileName
	}
}

//==========================================

func SetQyHookCheck(hook string) OptsCheck {
	return func(o *Check) {
		o.hook = hook
	}
}

func SetGitCheck(git git.Git) OptsCheck {
	return func(o *Check) {
		o.git = git
	}
}

func SetJenkinsCheck(jenkins *jenkins.Client) OptsCheck {
	return func(o *Check) {
		o.jenkins = jenkins
	}
}

func SetOrmCheck(orm *orm.Gorm) OptsCheck {
	return func(o *Check) {
		o.orm = orm
	}
}

func SetRemoteShellCheck(ssh *ssh.SSH) OptsCheck {
	return func(o *Check) {
		o.ssh = ssh
	}
}

//==========================================

func NewCheckClient(opts ...OptsCheck) *Check {
	o := &Check{
		outPath:     "./",
		outFileName: "json.json",
		services:    &Svc{Kobe: make([]string, 0, 0), Marx: make([]string, 0, 0)},
	}

	for _, opt := range opts {
		opt(o)
	}
	return o
}

func (o *Check) Build(ctx context.Context, source, target, svcPath string) error {
	if len(svcPath) > 0 {
		if err := configurator.NewYaml().Load(svcPath, &o.services); err != nil {
			return err
		}
	}
	var list []*requestInfo
	var b []byte
	var err error
	if list, err = o.list(ctx, source, target); err != nil {
		o.hookSend(ctx, "生成失败: \n"+err.Error())
		return err
	}
	if b, err = json.MarshalIndent(list, "", "	"); err != nil {
		o.hookSend(ctx, "生成失败: \n"+err.Error())
		return err
	}
	if err = os.WriteFile(o.outPath+o.outFileName, b, os.ModePerm); err != nil {
		o.hookSend(ctx, "生成失败: \n"+err.Error())
		return err
	}

	projectNames := make([]string, 0, len(list))
	for index, val := range list {
		projectNames = append(projectNames, fmt.Sprintf("%d.%s", index+1, val.Project))
	}
	o.hookSend(ctx, "发布项目: \n"+strings.Join(projectNames, "\n\n"))
	return nil
}

func (o *Check) Request(ctx context.Context) error {
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
		if key == OnlineTypeGit && app.IsNil(o.git) {
			return errors.New("git is null")
		}
		if key == OnlineTypeJenkins && app.IsNil(o.jenkins) {
			return errors.New("jenkins is null")
		}
		if key == OnlineTypeSql && app.IsNil(o.orm) {
			return errors.New("orm is null")
		}
		if key == OnlineTypeRemoteShell && app.IsNil(o.ssh) {
			return errors.New("remote ssh is null")
		}
	}

	for _, request := range requestList {
		var err error
		var other string

		switch request.Type {
		case OnlineTypeGit: //合并
			err := o.git.Merge(ctx, request.Project, request.Num)
			if err != nil {
				err = errors.New(fmt.Sprintf("%s %s %s", request.Project, "pr merge", err))
			}
		case OnlineTypeShell: //shell
			project, ok := request.Params["project"].(string)
			if !ok {
				err = errors.New(fmt.Sprintf("%s %s", request.Project, "get project name is err"))
				break
			}
			err := exec.Command("sh", request.Project, project).Run()
			if err != nil {
				err = errors.New(fmt.Sprintf("%s %s %s", request.Project, "executing command", err))
				break
			}
		case OnlineTypeJenkins:
			//jenkins
			// jenkins 可能会构建失败，所以重试 3次重试还不行就报错
			isSuccess := false
			for i := 0; i < 3; i++ {
				err := o.jenkins.BlockBuild(ctx, request.Project, request.Params)
				if err == nil {
					isSuccess = true
					break
				}
				time.Sleep(time.Second * 2)
			}
			if !isSuccess {
				err = errors.New(fmt.Sprintf("%s %s %s", request.Project, "block build", err))
			}
		case OnlineTypeSql: //sql
			sqls := request.Params["sql"].([]string)
			err = o.orm.Transaction(func(tx *gorm.DB) error {
				for _, sql := range sqls {
					return o.orm.Exec(sql).Error
				}
				return nil
			})
			if err != nil {
				err = errors.New(fmt.Sprintf("%s %s %s", request.Project, "sql", err))
				break
			}
		case OnlineTypeRemoteShell: //remote shell
			success, err := o.ssh.Run(request.Project)
			if err != nil {
				err = errors.New(fmt.Sprintf("%s %s %s", request.Project, "remote shell", err))
				break
			}
			other = success
		default:
			err = errors.New("type is err")
		}

		if err != nil {
			o.hookSend(ctx, "action:\nproject: "+request.Project+"\nerr: "+err.Error())
			return err
		}
		o.hookSend(ctx, "action:\nproject: "+request.Project+"\nok\nother: "+other)
	}

	return nil
}

//============================================================================================

func (o *Check) list(ctx context.Context, source, target string) ([]*requestInfo, error) {
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
		address = append(address, &requestInfo{Type: OnlineTypeSql, Project: "sql", Params: map[string]any{"sql": sqls}})
	}

	if len(o.services.Shell) > 0 {
		address = append(address, &requestInfo{Type: OnlineTypeRemoteShell, Project: fmt.Sprintf("/soft/scripts/menu_script/run.sh 2024/%s/menu* 2024/%s/group* check", o.services.Shell, o.services.Shell)})
	}

	for _, addr := range productList {
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
		address = append(address, &requestInfo{Type: OnlineTypeGit, Project: addr, Num: list[0].Number})
	}
	return address, nil
}

func (o *Check) hookSend(ctx context.Context, text string) {
	if client, err := qy_hook.NewQyHookClient(&qy_hook.Config{Token: o.hook}); err == nil && len(text) > 0 {
		_ = client.SendHook(ctx, text, []string{})
	}
}
