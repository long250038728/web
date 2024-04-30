package client

import (
	"context"
	"github.com/long250038728/web/tool/configurator"
	"github.com/long250038728/web/tool/git"
	"github.com/long250038728/web/tool/jenkins"
	"testing"
)

func TestOnlineBuild(t *testing.T) {
	configLoad := configurator.NewYaml()

	var gitConfig git.Config
	var jenkinsConfig jenkins.Config
	configLoad.MustLoad("/Users/linlong/Desktop/web/config/gitee.yaml", &gitConfig)
	configLoad.MustLoad("/Users/linlong/Desktop/web/config/jenkins.yaml", &jenkinsConfig)

	var giteeClient, _ = git.NewGiteeClient(&gitConfig)
	var jenkinsClient, _ = jenkins.NewJenkinsClient(&jenkinsConfig)

	if err := NewOnlineClient(context.Background(), giteeClient, jenkinsClient).Build("hotfix/staff_20240428", "master", "/Users/linlong/Desktop/online/linl.yaml"); err != nil {
		t.Errorf("Build() error = %v ", err)
	}
	t.Log("ok")
}
