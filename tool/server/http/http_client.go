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

const (
	StatusCreated = http.StatusCreated
)

type Client struct {
	timeout            time.Duration
	isTracing          bool
	username, password string
	contentType        string
	client             *http.Client
}

type Opt func(c *Client)

// SetTimeout The request timeout ,Not client timeout
// The lifecycle is within one request, not throughout the entire client
func SetTimeout(timeout time.Duration) Opt {
	return func(c *Client) {
		c.timeout = timeout
	}
}

func SetIsTracing(isTracing bool) Opt {
	return func(c *Client) {
		c.isTracing = isTracing
	}
}

func SetBasicAuth(username, password string) Opt {
	return func(c *Client) {
		c.username = username
		c.password = password
	}
}

func SetContentType(contentType string) Opt {
	return func(c *Client) {
		c.contentType = contentType
	}
}

func SetHttpClient(client *http.Client) Opt {
	return func(c *Client) {
		c.client = client
	}
}

func NewClient(opts ...Opt) *Client {
	client := &Client{
		timeout:     time.Second * 3,       //默认3s超时(单个请求超时，并非整个client)
		isTracing:   true,                  //默认记录链路
		contentType: "application/json",    //默认json
		client:      NewCustomHttpClient(), //默认http client
	}
	for _, opt := range opts {
		opt(client)
	}
	return client
}

func (c *Client) Put(ctx context.Context, address string, data map[string]any) ([]byte, int, error) {
	jsonBody, err := json.Marshal(data)
	if err != nil {
		return nil, 0, err
	}
	return c.do(ctx, http.MethodPut, address, jsonBody)
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
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	request, err := http.NewRequestWithContext(ctx, method, address, strings.NewReader(string(data)))
	if err != nil {
		return nil, 0, err
	}
	request.Header.Set("Content-Type", c.contentType)
	if len(c.username) > 0 && len(c.password) > 0 {
		request.SetBasicAuth(c.username, c.password)
	}

	// 设置 cookies
	cookies := []*http.Cookie{
		{
			Name:  "_ga",
			Value: "GA1.2.2002000608.1686045235",
		},
		{
			Name:  "Hm_lvt_275b44e62899a0fb992899d723470c22",
			Value: "1713424940",
		},
		{
			Name:  "PHPSESSID",
			Value: "392hu7jqtadho0bokij446k1i5",
		},
		{
			Name:  "merchant_code",
			Value: "SC583267",
		},
		{
			Name:  "token",
			Value: "Pd9A2A8e8fbD1se81728436284",
		},
		{
			Name:  "_gid",
			Value: "GA1.2.558930651.1728357879",
		},
	}

	// 将 cookies 添加到请求
	for _, cookie := range cookies {
		request.AddCookie(cookie)
	}

	if !c.isTracing {
		return c.request(request)
	}
	//======================================= 加入链路 ===============================================

	//新增一个span
	span := opentelemetry.NewSpan(ctx, "HTTP Client")
	defer span.Close()

	//请求地址 	//请求参数  //把链路信息放到request中
	span.AddEvent(address)
	if len(data) > 0 {
		span.AddEvent(string(data))
	}
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
	res, err := c.client.Do(request)
	if err != nil {
		return nil, 0, err
	}
	defer func() {
		_ = res.Body.Close()
	}()
	b, err := io.ReadAll(res.Body)

	//scanner := bufio.NewScanner(res.Body)
	//for scanner.Scan() {
	//	fmt.Println(scanner.Text())
	//}
	//if err := scanner.Err(); err != nil {
	//	fmt.Println(err)
	//}
	return b, res.StatusCode, err
}
