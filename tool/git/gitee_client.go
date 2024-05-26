package git

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/long250038728/web/tool/server/http"
	"time"
)

type Client struct {
	client         *http.Client
	address, token string
}

type Config struct {
	Token string `json:"token" yaml:"token"`
}

type Opt func(c *Client)

func NewGiteeClient(config *Config) (Git, error) {
	if len(config.Token) <= 0 {
		return nil, errors.New("token is empty")
	}
	client := &Client{
		address: "https://gitee.com",
		token:   config.Token,
		client:  http.NewClient(http.SetTimeout(time.Second * 10)),
	}
	return client, nil
}

func (g *Client) CreateFeature(ctx context.Context, repos, source, target string) error {
	data := map[string]any{
		"access_token": g.token,
		"refs":         source,
		"branch_name":  target,
	}
	_, code, err := g.client.Post(ctx, fmt.Sprintf("%s/api/v5/repos/%s/branches", g.address, repos), data)
	if err != nil {
		return errors.New(fmt.Sprintf("%s %s", repos, err.Error()))
	}

	if code != http.StatusCreated {
		return errors.New("request code is not 201")
	}
	return nil
}

func (g *Client) CreatePR(ctx context.Context, repos, source, target string) (*Info, error) {
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
	defer fmt.Println(string(b))

	//获取地址
	var item *Info
	err = json.Unmarshal(b, &item)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("%s %s", repos, err.Error()))
	}
	return item, nil
}

func (g *Client) GetPR(ctx context.Context, repos, source, target string) ([]*Info, error) {
	url := fmt.Sprintf("%s/api/v5/repos/%s/pulls", g.address, repos)
	data := map[string]any{
		"access_token": g.token,
		"state":        "open",
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

func (g *Client) Merge(ctx context.Context, repos string, num int32) error {
	type response struct {
		Merged  bool   `json:"merged"`
		Message string `json:"message"`
	}

	url := fmt.Sprintf("%s/api/v5/repos/%s/pulls/%d/merge", g.address, repos, num)
	data := map[string]any{
		"access_token": g.token,
		"merge_method": "merge",
	}

	resp, _, err := g.client.Put(ctx, url, data)
	if err != nil {
		return errors.New(fmt.Sprintf("%s: %s", repos, err.Error()))
	}

	var res response
	if err := json.Unmarshal(resp, &res); err != nil {
		return nil
	}

	if !res.Merged {
		return errors.New(res.Message)
	}
	fmt.Println("")
	fmt.Println(string(resp))
	return nil
}

func (g *Client) MergePR(ctx context.Context, repos string, source, target string) error {
	list, err := g.GetPR(ctx, repos, source, target)
	if err != nil {
		return err
	}
	if len(list) != 1 {
		return errors.New("pr list is not one")
	}
	return g.Merge(context.Background(), repos, list[0].Number)
}
