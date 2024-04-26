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
	_ = configLoad.Load("/Users/linlong/Desktop/web/configurator/gitee.yaml", &gitConfig)
	_ = configLoad.Load("/Users/linlong/Desktop/web/configurator/jenkins.yaml", &jenkinsConfig)

	var giteeClinet, _ = git.NewGiteeClient(&gitConfig)
	var jenkinsClient, _ = jenkins.NewJenkinsClient(&jenkinsConfig)

	if err := NewOnlineClient(context.Background(), giteeClinet, jenkinsClient).Build("release/v3.5.57", "master", "/Users/linlong/Desktop/online/linl.yaml"); err != nil {
		t.Errorf("Build() error = %v ", err)
	}
	t.Log("ok")
}
