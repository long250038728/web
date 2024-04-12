package client

import (
	"context"
	"github.com/long250038728/web/tool/server/http"
)

type QyHookClient struct {
	client *http.Client
	token  string
}

func NewQyHookClient(token string) *QyHookClient {
	return &QyHookClient{
		token:  token,
		client: http.NewClient(),
	}
}

func (q *QyHookClient) sendHook(ctx context.Context, content string, mobileList []string) error {
	data := map[string]any{
		"msgtype": "text",
		"text": map[string]any{
			"content":               content,
			"mentioned_mobile_list": mobileList,
		},
	}
	_, _, err := http.NewClient().Post(ctx, "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key="+q.token, data)
	return err
}
