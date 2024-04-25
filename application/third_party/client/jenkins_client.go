package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/long250038728/web/tool/server/http"
	http2 "net/http"
	"time"
)

type JenkinsClient struct {
	address string
	client  *http.Client
}

type lastBuild struct {
	Number int32 `json:"number"`
}
type queueBuild struct {
	Result string `json:"result"`
}

func NewJenkinsClient(address, username, password string) *JenkinsClient {
	return &JenkinsClient{
		address: address,
		client:  http.NewClient(http.SetBasicAuth(username, password)),
	}
}

// Build 构建
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

// GetLastNumber 获取最后一个id
func (j *JenkinsClient) GetLastNumber(ctx context.Context, job string) (int32, error) {
	var resp []byte
	var err error
	if resp, _, err = j.client.Get(ctx, fmt.Sprintf("%s/job/%s/lastBuild/api/json", j.address, job), nil); err != nil {
		return 0, err
	}
	var b lastBuild
	if err := json.Unmarshal(resp, &b); err != nil {
		return 0, err
	}
	if b.Number <= 0 {
		return 0, errors.New("number is error")
	}
	return b.Number, nil
}

// Block 阻塞获取是否构建完成
func (j *JenkinsClient) Block(ctx context.Context, job string, params map[string]any) error {
	number, err := j.GetLastNumber(ctx, job)
	if err != nil {
		return err
	}

	index := 1
	for {
		//检查是否退出
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		var resp []byte
		var err error
		var q queueBuild

		if resp, _, err = j.client.Get(context.Background(), fmt.Sprintf("%s/job/%s/%d/api/json?tree=result,building,displayName,duration", j.address, job, number), nil); err != nil {
			return err
		}
		if err := json.Unmarshal(resp, &q); err != nil {
			return err
		}
		if q.Result == "SUCCESS" {
			return nil
		}
		fmt.Println("for check num:", index)
		index += 1
		time.Sleep(10 * time.Second)
	}
}

// BlockBuild  阻塞构建
func (j *JenkinsClient) BlockBuild(ctx context.Context, job string, params map[string]any) error {
	startTime := time.Now()
	fmt.Println("============== ", job, " ===============")
	fmt.Println("start time:", startTime.Format("2006-01-02 15:04:05"))
	defer func() {
		endTime := time.Now()
		fmt.Println("end time:", endTime.Format("2006-01-02 15:04:05"))
		fmt.Println("total:", endTime.Sub(startTime))
	}()

	if buildErr := j.Build(ctx, job, params); buildErr != nil {
		return buildErr
	}

	// 等待jenkins列表
	time.Sleep(time.Second * 10)
	fmt.Println("query start:", time.Now().Format("2006-01-02 15:04:05"))

	if queryErr := j.Block(ctx, job, params); queryErr != nil {
		return queryErr
	}
	return nil
}
