package client

import (
	"context"
	"encoding/json"
	"github.com/long250038728/web/tool/configurator"
	"github.com/long250038728/web/tool/git"
	"github.com/long250038728/web/tool/jenkins"
	"os"
	"testing"
)

func TestOnlineBuild(t *testing.T) {
	configLoad := configurator.NewYaml()

	var gitConfig git.Config
	var jenkinsConfig jenkins.Config
	configLoad.MustLoad("/Users/linlong/Desktop/web/config/gitee.yaml", &gitConfig)
	configLoad.MustLoad("/Users/linlong/Desktop/web/config/jenkins.yaml", &jenkinsConfig)

	giteeClient, err := git.NewGiteeClient(&gitConfig)
	if err != nil {
		t.Error(err)
	}
	jenkinsClient, err := jenkins.NewJenkinsClient(&jenkinsConfig)
	if err != nil {
		t.Error(err)
	}

	if err := NewOnlineClient(context.Background(), giteeClient, jenkinsClient).Build("release/v3.5.62", "master", "/Users/linlong/Desktop/online/linl.yaml"); err != nil {
		t.Errorf("Build() error = %v ", err)
	}
	t.Log("ok")
}

func TestOnlineRequest(t *testing.T) {
	configLoad := configurator.NewYaml()

	var gitConfig git.Config
	var jenkinsConfig jenkins.Config
	configLoad.MustLoad("/Users/linlong/Desktop/web/config/gitee.yaml", &gitConfig)
	configLoad.MustLoad("/Users/linlong/Desktop/web/config/jenkins.yaml", &jenkinsConfig)

	giteeClient, err := git.NewGiteeClient(&gitConfig)
	if err != nil {
		t.Error(err)
	}
	jenkinsClient, err := jenkins.NewJenkinsClient(&jenkinsConfig)
	if err != nil {
		t.Error(err)
	}

	b, err := os.ReadFile("/Users/linlong/Desktop/web/application/third_party/client/project_list.md")
	if err != nil {
		t.Error(err)
	}

	var list []*requestInfo
	if err = json.Unmarshal(b, &list); err != nil {
		t.Error(err)
	}

	if err := NewOnlineClient(context.Background(), giteeClient, jenkinsClient).Request(list); err != nil {
		t.Errorf("Build() error = %v ", err)
	}
	t.Log("ok")
}
