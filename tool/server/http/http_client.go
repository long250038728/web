package http

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/long250038728/web/tool/tracing/opentelemetry"
	"io"
	"net/http"
	url2 "net/url"
	"strings"
	"time"
)

type Client struct {
	timeout time.Duration
}

type otp func(c *Client)

func NewClient(timeout time.Duration, opts ...otp) *Client {
	client := &Client{
		timeout: timeout,
	}
	for _, opt := range opts {
		opt(client)
	}
	return client
}

func (c *Client) Post(ctx context.Context, url string, data map[string]any) ([]byte, int, error) {
	return c.do(ctx, http.MethodPost, url, data)
}

func (c *Client) Get(ctx context.Context, url string, data map[string]any) ([]byte, int, error) {
	reqURL, err := url2.Parse(url)
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

func (c *Client) do(ctx context.Context, method string, url string, data map[string]any) ([]byte, int, error) {
	client := &http.Client{
		Timeout: c.timeout,
	}
	span := opentelemetry.NewSpan(ctx, "HTTP Client")
	defer span.Close()
	span.AddEvent(url)

	var body io.Reader
	if data != nil {
		jsonBody, err := json.Marshal(data)
		if err != nil {
			return nil, 0, err
		}
		span.AddEvent(string(jsonBody))
		body = strings.NewReader(string(jsonBody))
	}

	request, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, 0, err
	}

	request.Header.Set("Content-Type", "application/json")
	opentelemetry.InjectHttp(span.Context(), request) //把链路信息写到http header中

	res, err := client.Do(request)
	if err != nil {
		return nil, 0, err
	}

	b, err := io.ReadAll(res.Body)
	defer res.Body.Close()

	if err != nil || res.StatusCode != http.StatusOK {
		span.AddEvent(fmt.Sprintf("err: %v , code :%d", err, res.StatusCode))
		span.SetAttributes("err", true)
	}
	span.AddEvent(string(b))

	return b, res.StatusCode, err
}
