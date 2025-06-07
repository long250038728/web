package client

import (
	"context"
	"github.com/long250038728/web/tool/configurator"
	"github.com/long250038728/web/tool/git"
	"github.com/long250038728/web/tool/hook"
	"github.com/long250038728/web/tool/jenkins"
	"github.com/long250038728/web/tool/persistence/orm"
	"github.com/long250038728/web/tool/ssh"
	"testing"
)

var gitClient git.Git
var jenkinsClient *jenkins.Client
var ormClient *orm.Gorm
var sshClient ssh.SSH
var hookClient hook.Hook

var tels = []string{"18575538087"}

func init() {
	var gitConfig git.Config
	var jenkinsConfig jenkins.Config
	var sshConfig ssh.Config
	var hookToken = "bb3f6f61-04b8-4b46-a167-08a2c91d408d"

	var err error
	configLoad := configurator.NewYaml()
	configLoad.MustLoadConfigPath("other/gitee.yaml", &gitConfig)
	configLoad.MustLoadConfigPath("other/jenkins.yaml", &jenkinsConfig)
	configLoad.MustLoadConfigPath("other/ssh.yaml", &sshConfig)

	if gitClient, err = git.NewGiteeClient(&gitConfig); err != nil {
		panic(err)
	}
	if jenkinsClient, err = jenkins.NewJenkinsClient(&jenkinsConfig); err != nil {
		panic(err)
	}

	if sshClient, err = ssh.NewRemoteSSH(&sshConfig); err != nil {
		panic(err)
	}
	if hookClient, err = hook.NewQyHookClient(&hook.Config{Token: hookToken}); err != nil {
		panic(err)
	}
}

func TestOnlineBuild(t *testing.T) {
	return
	var err error
	var ormConfig orm.Config
	configLoad := configurator.NewYaml()
	configLoad.MustLoadConfigPath("online/db.yaml", &ormConfig)
	if ormClient, err = orm.NewMySQLGorm(&ormConfig); err != nil {
		panic(err)
	}

	if err := NewTaskClient(
		SetGit(gitClient),
		SetJenkins(jenkinsClient),
		SetOrm(ormClient),
		SetRemoteShell(sshClient),
		SetQyHook(hookClient, tels),
	).Build(
		context.Background(),
		"release/v3.5.96",
		"master",
		"./script/svc.yaml",
	); err != nil {
		t.Errorf("Build() error = %v ", err)
	}
	t.Log("ok")
}

func TestRequest(t *testing.T) {
	if err := NewTaskClient(
		SetGit(gitClient),
		SetJenkins(jenkinsClient),
		SetOrm(ormClient),
		SetRemoteShell(sshClient),
		SetQyHook(hookClient, tels),
	).Request(context.Background()); err != nil {
		t.Errorf("Build() error = %v ", err)
	}
	t.Log("ok")
}
