package client

import (
	"context"
	"github.com/long250038728/web/tool/configurator"
	"github.com/long250038728/web/tool/git"
	"github.com/long250038728/web/tool/jenkins"
	"github.com/long250038728/web/tool/persistence/orm"
	"github.com/long250038728/web/tool/ssh"
	"os"
	"path/filepath"
	"testing"
)

var gitConfig git.Config
var jenkinsConfig jenkins.Config
var ormConfig orm.Config
var sshConfig ssh.Config

var gitClient git.Git
var jenkinsClient *jenkins.Client
var ormClient *orm.Gorm
var sshClient ssh.SSH

var hookToken = "bb3f6f61-04b8-4b46-a167-08a2c91d408d"

func init() {
	path := os.Getenv("WEB")
	if len(path) == 0 {
		path = "/Users/linlong/Desktop/web"
	}

	var err error
	configLoad := configurator.NewYaml()
	configLoad.MustLoad(filepath.Join(path, "config", "gitee.yaml"), &gitConfig)
	configLoad.MustLoad(filepath.Join(path, "config", "jenkins.yaml"), &jenkinsConfig)
	configLoad.MustLoad(filepath.Join(path, "config", "online/db.yaml"), &ormConfig)
	configLoad.MustLoad(filepath.Join(path, "config", "ssh.yaml"), &sshConfig)

	if gitClient, err = git.NewGiteeClient(&gitConfig); err != nil {
		panic(err)
	}
	if jenkinsClient, err = jenkins.NewJenkinsClient(&jenkinsConfig); err != nil {
		panic(err)
	}
	if ormClient, err = orm.NewGorm(&ormConfig); err != nil {
		panic(err)
	}
	if sshClient, err = ssh.NewRemoteSSH(&sshConfig); err != nil {
		panic(err)
	}
}

func TestCheckBuild(t *testing.T) {
	if err := NewTaskClient(
		SetGit(gitClient),
		SetJenkins(jenkinsClient),
		SetOrm(ormClient),
		SetRemoteShell(sshClient),
		SetQyHook(hookToken),
	).BuildCheck(
		context.Background(),
		"release/v3.5.85",
		"check",
		"./script/svc.yaml",
	); err != nil {
		t.Errorf("Build() error = %v ", err)
	}
	t.Log("ok")
}

func TestOnlineBuild(t *testing.T) {
	if err := NewTaskClient(
		SetGit(gitClient),
		SetJenkins(jenkinsClient),
		SetOrm(ormClient),
		SetRemoteShell(sshClient),
		SetQyHook(hookToken),
	).Build(
		context.Background(),
		"release/v3.5.85",
		"master",
		"./script/svc.yaml",
	); err != nil {
		t.Errorf("Build() error = %v ", err)
	}
	t.Log("ok")
}

func TestOnlineRequest(t *testing.T) {
	if err := NewTaskClient(
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
