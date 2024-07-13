package client

import (
	"context"
	"github.com/long250038728/web/tool/configurator"
	"github.com/long250038728/web/tool/git"
	"github.com/long250038728/web/tool/jenkins"
	"github.com/long250038728/web/tool/persistence/orm"
	"github.com/long250038728/web/tool/ssh"
	"testing"
)

var gitConfig git.Config
var jenkinsConfig jenkins.Config
var ormConfig orm.Config

var gitClient git.Git
var jenkinsClient *jenkins.Client
var ormClient *orm.Gorm
var sshClient *ssh.SSH

var hookToken = "bb3f6f61-04b8-4b46-a167-08a2c91d408d"

func init() {
	var err error
	configLoad := configurator.NewYaml()
	configLoad.MustLoad("/Users/linlong/Desktop/web/config/gitee.yaml", &gitConfig)
	configLoad.MustLoad("/Users/linlong/Desktop/web/config/jenkins.yaml", &jenkinsConfig)
	configLoad.MustLoad("/Users/linlong/Desktop/web/config/online/db.yaml", &ormConfig)

	if gitClient, err = git.NewGiteeClient(&gitConfig); err != nil {
		panic(err)
	}
	if jenkinsClient, err = jenkins.NewJenkinsClient(&jenkinsConfig); err != nil {
		panic(err)
	}
	if ormClient, err = orm.NewGorm(&ormConfig); err != nil {
		panic(err)
	}
	if sshClient, err = ssh.NewSSH(&ssh.Config{Host: "42.193.172.210", Port: 22, User: "root", Password: "199481&&Shuai"}); err != nil {
		panic(err)
	}
}

func TestOnlineBuild(t *testing.T) {
	if err := NewOnlineClient(
		SetGit(gitClient),
		SetJenkins(jenkinsClient),
		SetOrm(ormClient),
		SetRemoteShell(sshClient),
		SetQyHook(hookToken),
	).Build(
		context.Background(),
		"release/v3.5.80",
		"master",
		"/Users/linlong/Desktop/online/linl.yaml",
	); err != nil {
		t.Errorf("Build() error = %v ", err)
	}
	t.Log("ok")
}

func TestOnlineRequest(t *testing.T) {
	if err := NewOnlineClient(
		SetGit(gitClient),
		SetJenkins(jenkinsClient),
		SetOrm(ormClient),
		SetRemoteShell(sshClient),
		SetQyHook(hookToken),
	).Request(context.Background()); err != nil {
		t.Errorf("Build() error = %v ", err)
	}
	t.Log("ok")
}
