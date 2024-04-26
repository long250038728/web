package jenkins

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/long250038728/web/tool/server/http"
	"time"
)

type Client struct {
	address string
	client  *http.Client
}

var timeLayout = "2006-01-02 15:04:05"

type Config struct {
	Address  string `json:"address" yaml:"address"`
	Username string `json:"username" yaml:"username"`
	Password string `json:"password" yaml:"password"`
}

func NewJenkinsClient(config *Config) (*Client, error) {
	if len(config.Address) <= 0 {
		return nil, errors.New("address is empty")
	}

	var opts []http.Opt
	if len(config.Username) > 0 && len(config.Password) > 0 {
		opts = append(opts, http.SetBasicAuth(config.Username, config.Password))
	}

	return &Client{
		address: config.Address,
		client:  http.NewClient(opts...),
	}, nil
}

// Build 构建
func (j *Client) Build(ctx context.Context, job string, params map[string]any) error {
	url := fmt.Sprintf("%s/job/%s/buildWithParameters", j.address, job)
	if params == nil {
		url = fmt.Sprintf("%s/job/%s/build", j.address, job)
	}
	_, code, err := j.client.Post(ctx, url, params)

	if err != nil {
		return err
	}
	if code != http.StatusCreated {
		return errors.New("request code is not 201")
	}
	return nil
}

// GetLastNumber 获取最后一个id
func (j *Client) GetLastNumber(ctx context.Context, job string) (int32, error) {
	var resp []byte
	var err error
	if resp, _, err = j.client.Get(ctx, fmt.Sprintf("%s/job/%s/lastBuild/api/json", j.address, job), nil); err != nil {
		return 0, err
	}

	type lastBuild struct {
		Number int32 `json:"number"`
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
func (j *Client) Block(ctx context.Context, job string) error {
	type queueBuild struct {
		Result string `json:"result"`
	}

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

		if resp, _, err = j.client.Get(ctx, fmt.Sprintf("%s/job/%s/%d/api/json?tree=result,building,displayName,duration", j.address, job, number), nil); err != nil {
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
func (j *Client) BlockBuild(ctx context.Context, job string, params map[string]any) error {
	startTime := time.Now()
	fmt.Println("============== ", job, " ===============")
	fmt.Println("start time:", startTime.Format(timeLayout))
	defer func() {
		endTime := time.Now()
		fmt.Println("end time:", endTime.Format(timeLayout))
		fmt.Println("total:", endTime.Sub(startTime))
	}()

	if buildErr := j.Build(ctx, job, params); buildErr != nil {
		return buildErr
	}

	// 等待jenkins列表
	time.Sleep(time.Second * 10)
	fmt.Println("query start:", time.Now().Format(timeLayout))

	if queryErr := j.Block(ctx, job); queryErr != nil {
		return queryErr
	}
	return nil
}
