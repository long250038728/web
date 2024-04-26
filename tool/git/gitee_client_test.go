package git

import (
	"context"
	"github.com/long250038728/web/tool/configurator"
	"testing"
)

var gitToken = "5f8aaa1e024cad5e24e86fda85c57f49"

func TestClient_Merge(t *testing.T) {
	var cfg Config
	err := configurator.NewYaml().Load("/Users/linlong/Desktop/web/configurator/gitee.yaml", &cfg)
	if err != nil {
		t.Error(err)
		return
	}

	gitClient, err := NewGiteeClient(&cfg)
	if err != nil {
		t.Error(err)
		return
	}
	list, err := gitClient.GetPR(context.Background(), "zhubaoe/socrates", "release/v3.5.59", "master")
	if err != nil {
		t.Error(err)
		return
	}
	if len(list) != 1 {
		t.Error("pr list is not one")
	}

	err = gitClient.Merge(context.Background(), "zhubaoe/socrates", list[0].Number)
	if err != nil {
		t.Error(err)
		return
	}

	t.Log("ok")
}

func TestClient_MergePR(t *testing.T) {
	gitClient, err := NewGiteeClient(&Config{Token: gitToken})
	if err != nil {
		t.Error(err)
		return
	}
	if err = gitClient.MergePR(context.Background(), "zhubaoe/socrates", "release/v3.5.59", "master"); err != nil {
		t.Error(err)
		return
	}
	t.Log("ok")
}
