package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/long250038728/web/tool/server/http"
	http2 "net/http"
)

type Git interface {
	CreateFeature(ctx context.Context, repos, source, target string) error
	CreatePR(ctx context.Context, repos, source, target string) (*Info, error)
	GetPR(ctx context.Context, repos, source, target string) ([]*Info, error)
	Merge(ctx context.Context, repos string, num int32) error
}

type Info struct {
	HtmlUrl string `json:"html_url"`
	Url     string `json:"url"`
	Number  int32  `json:"number"`
}

type GiteeClient struct {
	client         *http.Client
	address, token string
}

type GiteeClientOpt func(c *GiteeClient)

func SetGiteeAddress(address string) GiteeClientOpt {
	return func(c *GiteeClient) {
		c.address = address
	}
}

func NewGiteeClinet(token string, opts ...GiteeClientOpt) Git {
	client := &GiteeClient{
		address: "https://gitee.com",
		token:   token,
		client:  http.NewClient(),
	}
	for _, opt := range opts {
		opt(client)
	}
	return client
}

func (g *GiteeClient) CreateFeature(ctx context.Context, repos, source, target string) error {
	data := map[string]any{
		"access_token": g.token,
		"refs":         source,
		"branch_name":  target,
	}
	_, code, err := g.client.Post(ctx, fmt.Sprintf("%s/api/v5/repos/%s/branches", g.address, repos), data)
	if err != nil {
		return errors.New(fmt.Sprintf("%s %s", repos, err.Error()))
	}

	if code != http2.StatusCreated {
		return errors.New("request code is not 201")
	}
	return nil
}

func (g *GiteeClient) CreatePR(ctx context.Context, repos, source, target string) (*Info, error) {
	data := map[string]any{
		"access_token": g.token,
		"title":        fmt.Sprintf("Merge branch %s into %s", source, target),
		"head":         source,
		"base":         target,
	}
	b, _, err := g.client.Post(ctx, fmt.Sprintf("%s/api/v5/repos/%s/pulls", g.address, repos), data)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("%s %s", repos, err.Error()))
	}
	//获取地址
	var item *Info
	err = json.Unmarshal(b, &item)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("%s %s", repos, err.Error()))
	}
	return item, nil
}

func (g *GiteeClient) GetPR(ctx context.Context, repos, source, target string) ([]*Info, error) {
	url := fmt.Sprintf("%s/api/v5/repos/%s/pulls", g.address, repos)
	data := map[string]any{
		"access_token": g.token,
		"state":        "all",
		"head":         source,
		"base":         target,
		"sort":         "created",
		"direction":    "desc",
		"page":         "1",
		"per_page":     "20",
	}

	b, _, err := g.client.Get(ctx, url, data)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("%s: %s", repos, err.Error()))
	}
	var list []*Info
	err = json.Unmarshal(b, &list)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("%s: %s", repos, err.Error()))
	}
	return list, nil
}

func (g *GiteeClient) Merge(ctx context.Context, repos string, num int32) error {
	url := fmt.Sprintf("%s/api/v5/repos/%s/pulls/%d/merge", g.address, repos, num)
	data := map[string]any{
		"access_token": g.token,
		"merge_method": "merge",
	}

	res, _, err := g.client.Get(ctx, url, data)
	if err != nil {
		return errors.New(fmt.Sprintf("%s: %s", repos, err.Error()))
	}
	println(string(res))
	return nil
}
