package http

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/long250038728/web/tool/tracing/opentelemetry"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Client struct {
	timeout   time.Duration
	isTracing bool
}

type httpClientOpt func(c *Client)

func SetTimeout(timeout time.Duration) httpClientOpt {
	return func(c *Client) {
		c.timeout = timeout
	}
}

func SetIsTracing(isTracing bool) httpClientOpt {
	return func(c *Client) {
		c.isTracing = isTracing
	}
}

func NewClient(opts ...httpClientOpt) *Client {
	client := &Client{
		timeout:   time.Second * 3,
		isTracing: true,
	}
	for _, opt := range opts {
		opt(client)
	}
	return client
}

func (c *Client) Post(ctx context.Context, address string, data map[string]any) ([]byte, int, error) {
	jsonBody, err := json.Marshal(data)
	if err != nil {
		return nil, 0, err
	}
	return c.do(ctx, http.MethodPost, address, jsonBody)
}

func (c *Client) Get(ctx context.Context, address string, data map[string]any) ([]byte, int, error) {
	reqURL, err := url.Parse(address)
	if err != nil {
		return nil, 0, err
	}

	reqQuery := reqURL.Query()
	for k, v := range data {
		reqQuery.Set(k, fmt.Sprintf("%v", v))
	}
	reqURL.RawQuery = reqQuery.Encode()

	return c.do(ctx, http.MethodGet, reqURL.String(), nil)
}

func (c *Client) do(ctx context.Context, method string, address string, data []byte) ([]byte, int, error) {
	request, err := http.NewRequest(method, address, strings.NewReader(string(data)))
	if err != nil {
		return nil, 0, err
	}
	request.Header.Set("Content-Type", "application/json")

	if !c.isTracing {
		b, code, err := c.request(request)
		return b, code, err
	}
	//======================================= 加入链路 ===============================================

	//新增一个span
	span := opentelemetry.NewSpan(ctx, "HTTP Client")
	defer span.Close()

	//请求地址 	//请求参数  //把链路信息放到request中
	span.AddEvent(address)
	span.AddEvent(string(data))
	opentelemetry.InjectHttp(span.Context(), request) //把链路信息写到http header中

	//响应参数
	b, code, err := c.request(request)
	if err != nil || code != http.StatusOK {
		span.SetAttributes("err", true)
		span.AddEvent(fmt.Sprintf("err: %v , code :%d", err, code))
	}
	span.AddEvent(string(b))
	return b, code, err
}

func (c *Client) request(request *http.Request) ([]byte, int, error) {
	client := http.Client{Timeout: c.timeout}
	res, err := client.Do(request)
	if err != nil {
		return nil, 0, err
	}
	defer func() {
		_ = res.Body.Close()
	}()
	b, err := io.ReadAll(res.Body)
	return b, res.StatusCode, err
}
