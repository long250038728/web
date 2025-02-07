package hook

import (
	"context"
	"errors"
	"github.com/long250038728/web/tool/server/http"
)

type Config struct {
	Token string `json:"token" yaml:"token"`
}

type Client struct {
	client *http.Client
	token  string
}

func NewQyHookClient(config *Config) (Hook, error) {
	if len(config.Token) <= 0 {
		return nil, errors.New("token is empty")
	}

	return &Client{
		token:  config.Token,
		client: http.NewClient(),
	}, nil
}

func (c *Client) SendHook(ctx context.Context, content string, mobileList []string) error {
	data := map[string]any{
		"msgtype": "text",
		"text": map[string]any{
			"content":               content,
			"mentioned_mobile_list": mobileList,
		},
	}
	_, _, err := http.NewClient().Post(ctx, "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key="+c.token, data)
	return err
}
