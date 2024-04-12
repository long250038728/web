package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/long250038728/web/tool/server/http"
)

type GiteeInfo struct {
	HtmlUrl string `json:"html_url"`
}

type GiteeClient struct {
	client *http.Client
	token  string
}

func NewGiteeClinet(token string) *GiteeClient {
	return &GiteeClient{
		token:  token,
		client: http.NewClient(),
	}
}

func (g *GiteeClient) CreateFeature(ctx context.Context, addr, source, target string) error {
	data := map[string]any{
		"access_token": g.token,
		"refs":         source,
		"branch_name":  target,
	}
	_, _, err := g.client.Post(ctx, fmt.Sprintf("https://gitee.com/api/v5/repos/%s/branches", addr), data)
	return err
}

func (g *GiteeClient) CreatePR(ctx context.Context, addr, source, target string) (*GiteeInfo, error) {
	data := map[string]any{
		"access_token": g.token,
		"title":        fmt.Sprintf("Merge branch %s into %s", source, target),
		"head":         source,
		"base":         target,
	}
	b, _, err := g.client.Post(ctx, fmt.Sprintf("https://gitee.com/api/v5/repos/%s/pulls", addr), data)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("%s %s", addr, err.Error()))
	}
	//获取地址
	var item *GiteeInfo
	err = json.Unmarshal(b, &item)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("%s %s", addr, err.Error()))
	}
	return item, nil
}

func (g *GiteeClient) GetPR(ctx context.Context, addr, source, target string) ([]*GiteeInfo, error) {
	url := fmt.Sprintf("https://gitee.com/api/v5/repos/%s/pulls?access_token=%s&state=open&head=%s&base=%s&sort=created&direction=desc&page=1&per_page=20",
		addr, g.token, source, target)

	b, _, err := g.client.Get(ctx, url, nil)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("%s: %s", addr, err.Error()))
	}

	var list []*GiteeInfo
	err = json.Unmarshal(b, &list)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("%s: %s", addr, err.Error()))
	}
	return list, nil
}
