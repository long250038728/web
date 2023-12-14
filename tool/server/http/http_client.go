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
	timeout time.Duration
}

type Otp func(c *Client)

func SetTimeout(timeout time.Duration) Otp {
	return func(c *Client) {
		c.timeout = timeout
	}
}

func NewClient(opts ...Otp) *Client {
	client := &Client{}
	for _, opt := range opts {
		opt(client)
	}
	return client
}

func (c *Client) Post(ctx context.Context, address string, data map[string]any) ([]byte, int, error) {
	return c.do(ctx, http.MethodPost, address, data)
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

func (c *Client) do(ctx context.Context, method string, address string, data map[string]any) ([]byte, int, error) {
	client := &http.Client{
		Timeout: c.timeout,
	}
	span := opentelemetry.NewSpan(ctx, "HTTP Client")
	defer span.Close()
	span.AddEvent(address)

	var body io.Reader
	if data != nil {
		jsonBody, err := json.Marshal(data)
		if err != nil {
			return nil, 0, err
		}
		span.AddEvent(string(jsonBody))
		body = strings.NewReader(string(jsonBody))
	}

	request, err := http.NewRequest(method, address, body)
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
