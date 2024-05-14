package client

import (
	"context"
	"encoding/json"
	"github.com/long250038728/web/tool/configurator"
	"github.com/long250038728/web/tool/git"
	"github.com/long250038728/web/tool/jenkins"
	"github.com/long250038728/web/tool/persistence/orm"
	"os"
	"testing"
)

func TestOnlineBuild(t *testing.T) {
	configLoad := configurator.NewYaml()

	var gitConfig git.Config
	var jenkinsConfig jenkins.Config
	var ormConfig orm.Config

	configLoad.MustLoad("/Users/linlong/Desktop/web/config/gitee.yaml", &gitConfig)
	configLoad.MustLoad("/Users/linlong/Desktop/web/config/jenkins.yaml", &jenkinsConfig)
	configLoad.MustLoad("/Users/linlong/Desktop/web/config/db2.yaml", &ormConfig)

	giteeClient, err := git.NewGiteeClient(&gitConfig)
	if err != nil {
		t.Error(err)
		return
	}
	jenkinsClient, err := jenkins.NewJenkinsClient(&jenkinsConfig)
	if err != nil {
		t.Error(err)
		return
	}

	ormClient, err := orm.NewGorm(&ormConfig)
	if err != nil {
		t.Error(err)
		return
	}

	if err := NewOnlineClient(context.Background(), giteeClient, jenkinsClient, ormClient).Build(
		"release/v3.5.63",
		"master",
		"/Users/linlong/Desktop/web/application/third_party/client/svc.yaml",
		"/Users/linlong/Desktop/web/application/third_party/client/sql.sql",
	); err != nil {
		t.Errorf("Build() error = %v ", err)
	}
	t.Log("ok")
}

func TestOnlineRequest(t *testing.T) {
	configLoad := configurator.NewYaml()

	var gitConfig git.Config
	var jenkinsConfig jenkins.Config
	var ormConfig orm.Config

	configLoad.MustLoad("/Users/linlong/Desktop/web/config/gitee.yaml", &gitConfig)
	configLoad.MustLoad("/Users/linlong/Desktop/web/config/jenkins.yaml", &jenkinsConfig)
	configLoad.MustLoad("/Users/linlong/Desktop/web/config/db.yaml", &ormConfig)

	giteeClient, err := git.NewGiteeClient(&gitConfig)
	if err != nil {
		t.Error(err)
		return
	}
	jenkinsClient, err := jenkins.NewJenkinsClient(&jenkinsConfig)
	if err != nil {
		t.Error(err)
		return
	}

	ormClient, err := orm.NewGorm(&ormConfig)
	if err != nil {
		t.Error(err)
		return
	}

	b, err := os.ReadFile("/Users/linlong/Desktop/web/application/third_party/client/project_list.md")
	if err != nil {
		t.Error(err)
	}

	var list []*requestInfo
	if err = json.Unmarshal(b, &list); err != nil {
		t.Error(err)
	}

	if err := NewOnlineClient(context.Background(), giteeClient, jenkinsClient, ormClient).Request(list); err != nil {
		t.Errorf("Build() error = %v ", err)
	}
	t.Log("ok")
}
