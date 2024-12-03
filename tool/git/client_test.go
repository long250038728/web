package git

import (
	"context"
	"github.com/long250038728/web/tool/configurator"
	"testing"
)

var gitClient Git

func init() {
	var err error
	var giteeConfig Config
	configurator.NewYaml().MustLoadConfigPath("gitee.yaml", &giteeConfig)
	if gitClient, err = NewGiteeClient(&giteeConfig); err != nil {
		panic(err)
	}
}

func TestClient_Merge(t *testing.T) {
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
	if err := gitClient.MergePR(context.Background(), "zhubaoe/socrates", "release/v3.5.59", "master"); err != nil {
		t.Error(err)
		return
	}
	t.Log("ok")
}

func TestClient_GetPR(t *testing.T) {
	list, err := gitClient.GetPR(context.Background(), "zhubaoe/socrates", "release/v3.6.19.0", "check")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(list)
}
