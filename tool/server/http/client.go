package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	StatusCreated = http.StatusCreated
)

type Client struct {
	timeout            time.Duration
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

	request, err := http.NewRequestWithContext(ctx, method, address, bytes.NewReader(data))
	if err != nil {
		return nil, 0, err
	}
	request.Header.Set("Content-Type", c.contentType)
	if len(c.username) > 0 && len(c.password) > 0 {
		request.SetBasicAuth(c.username, c.password)
	}
	return c.request(request)
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
