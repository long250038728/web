package client

import (
	"context"
	"errors"
	"fmt"
	"github.com/long250038728/web/tool/server/http"
	http2 "net/http"
)

type JenkinsClient struct {
	address string
	client  *http.Client
}

func NewJenkinsClient(address, username, password string) *JenkinsClient {
	return &JenkinsClient{
		address: address,
		client:  http.NewClient(http.SetBasicAuth(username, password)),
	}
}

// Build
//
// curl -X POST https://www.jenkins.cn/job/xxxxx/buildWithParameters \
// --user xxx:xxx \
// --data-urlencode "PARAMES=XXX" \
//
// curl -X POST https://www.jenkins.cn/job/xxxxx/build \
// --user xxx:xxx
func (j *JenkinsClient) Build(ctx context.Context, job string, params map[string]any) error {
	url := fmt.Sprintf("%s/job/%s/buildWithParameters", j.address, job)
	if params == nil {
		url = fmt.Sprintf("%s/job/%s/build", j.address, job)
	}
	_, code, err := j.client.Post(ctx, url, params)
	if err != nil {
		return err
	}
	if code != http2.StatusCreated {
		return errors.New("request code is not 201")
	}
	return nil
}

// Queue
// curl -X GET https://www.jenkins.cn/queue/api/json \
// --user xxx:xxx
func (j *JenkinsClient) Queue(ctx context.Context, job string, params map[string]any) {
	resp, code, err := j.client.Get(ctx, fmt.Sprintf("%s/queue/api/json", j.address), nil)
	if err != nil {
		return
	}
	fmt.Println(string(resp), code)
}
