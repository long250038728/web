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

func init() {
	var err error
	configLoad := configurator.NewYaml()
	configLoad.MustLoad("/Users/linlong/Desktop/web/config/gitee.yaml", &gitConfig)
	configLoad.MustLoad("/Users/linlong/Desktop/web/config/check/jenkins.yaml", &jenkinsConfig)
	configLoad.MustLoad("/Users/linlong/Desktop/web/config/check/db.yaml", &ormConfig)
	configLoad.MustLoad("/Users/linlong/Desktop/web/config/ssh.yaml", &sshConfig)

	if gitClient, err = git.NewGiteeClient(&gitConfig); err != nil {
		panic(err)
	}
	if jenkinsClient, err = jenkins.NewJenkinsClient(&jenkinsConfig); err != nil {
		panic(err)
	}
	if ormClient, err = orm.NewGorm(&ormConfig); err != nil {
		panic(err)
	}
	if sshClient, err = ssh.NewSSH(&sshConfig); err != nil {
		panic(err)
	}
}

func TestCheckBuild(t *testing.T) {
	if err := NewCheckClient(
		SetGitCheck(gitClient),
		SetJenkinsCheck(jenkinsClient),
		SetOrmCheck(ormClient),
		SetRemoteShellCheck(sshClient),
		SetQyHookCheck(hookToken),
	).BuildCheck(
		context.Background(),
		"release/v3.5.80",
		"check",
		"./script/svc.yaml",
	); err != nil {
		t.Errorf("Build() error = %v ", err)
	}
	t.Log("ok")
}

func TestCheckRequest(t *testing.T) {
	if err := NewCheckClient(
		SetGitCheck(gitClient),
		SetJenkinsCheck(jenkinsClient),
		SetOrmCheck(ormClient),
		SetRemoteShellCheck(sshClient),
		SetQyHookCheck(hookToken),
	).RequestCheck(context.Background()); err != nil {
		t.Errorf("Build() error = %v ", err)
	}
	t.Log("ok")
}
