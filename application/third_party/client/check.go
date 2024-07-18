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

func (o *Check) BuildCheck(ctx context.Context, source, target, svcPath string) error {
	if len(svcPath) > 0 {
		if err := configurator.NewYaml().Load(svcPath, &o.services); err != nil {
			return err
		}
	}
	var list []*requestInfo
	var err error
	if list, err = o.listCheck(ctx, source, target); err != nil {
		o.hookSendCheck(ctx, "生成失败: \n"+err.Error())
		return err
	}
	if err = o.saveCheck(ctx, list); err != nil {
		o.hookSendCheck(ctx, "生成失败: \n"+err.Error())
		return err
	}

	projectNames := make([]string, 0, len(list))
	for index, val := range list {
		projectNames = append(projectNames, fmt.Sprintf("%d.%s", index+1, val.Project))
	}
	o.hookSendCheck(ctx, "发布项目: \n"+strings.Join(projectNames, "\n\n"))
	return nil
}

func (o *Check) RequestCheck(ctx context.Context) error {
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

	for index, request := range requestList {
		//已经成功的就不再处理
		if request.Success {
			continue
		}

		startTime := time.Now().Local()
		var err error
		var other = "empty"

		switch request.Type {
		case OnlineTypeGit: //合并
			err = o.git.Merge(ctx, request.Project, request.Num)
		case OnlineTypeShell: //shell
			project, ok := request.Params["project"].(string)
			if !ok {
				err = errors.New("shell script is error")
				break
			}
			err = exec.Command("sh", request.Project, project).Run()
		case OnlineTypeJenkins:
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
				err = errors.New("jenkins build is failure")
			}
		case OnlineTypeSql: //sql
			sql := request.Params["sql"].([]interface{})
			sqls := make([]string, 0, len(sql))
			for _, s := range sql {
				str, ok := s.(string)
				if !ok {
					err = errors.New("sql is failure")
					break
				}
				sqls = append(sqls, str)
			}

			err = o.orm.Transaction(func(tx *gorm.DB) error {
				for _, sql := range sqls {
					if err = o.orm.Exec(sql).Error; err != nil {
						return err
					}
				}
				return nil
			})
		case OnlineTypeRemoteShell: //remote shell
			other, err = o.ssh.Run(request.Project)
		default:
			err = errors.New("type is err")
		}

		//============================================================================
		endTime := time.Now().Local()
		if err != nil {
			o.hookSendCheck(ctx, fmt.Sprintf("project: %s \nstatus: %s \nstart: %s   end: %s   sub: %s \nother: \n%s", request.Project, "failure", startTime.Format(time.TimeOnly), endTime.Format(time.TimeOnly), endTime.Sub(startTime).String(), err.Error()))
			return err
		}

		o.hookSendCheck(ctx, fmt.Sprintf("project: %s \nstatus: %s \nstart: %s   end: %s   sub: %s \nother: \n%s", request.Project, "success", startTime.Format(time.TimeOnly), endTime.Format(time.TimeOnly), endTime.Sub(startTime).String(), other))
		requestList[index].Success = true
		_ = o.saveCheck(ctx, requestList)
	}

	return nil
}

//============================================================================================

func (o *Check) listCheck(ctx context.Context, source, target string) ([]*requestInfo, error) {
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
		address = append(address, &requestInfo{Type: OnlineTypeRemoteShell, Project: fmt.Sprintf("/soft/scripts/menu_script/run.sh 2024/%s/menu* 2024/%s/group* prod", o.services.Shell, o.services.Shell)})
	}

	for _, addr := range productList {
		list, err := o.git.GetPR(ctx, addr, source, target)
		if err != nil || len(list) != 1 {
			continue
		}
		//调用合并分支
		address = append(address, &requestInfo{Type: OnlineTypeGit, Project: addr, Num: list[0].Number})
	}
	return address, nil
}

func (o *Check) hookSendCheck(ctx context.Context, text string) {
	if client, err := qy_hook.NewQyHookClient(&qy_hook.Config{Token: o.hook}); err == nil && len(text) > 0 {
		_ = client.SendHook(ctx, text, []string{})
	}
}

func (o *Check) saveCheck(ctx context.Context, list []*requestInfo) error {
	b, err := json.MarshalIndent(list, "", "	")
	if err != nil {
		o.hookSendCheck(ctx, "生成失败: \n"+err.Error())
		return err
	}
	if err := os.WriteFile(o.outPath+o.outFileName, b, os.ModePerm); err != nil {
		o.hookSendCheck(ctx, "生成失败: \n"+err.Error())
		return err
	}

	return nil
}
