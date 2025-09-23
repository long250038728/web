package http

import (
	"bytes"
	"context"
	"fmt"
	"github.com/long250038728/web/tool/tracing/opentelemetry"
	"io"
	"net"
	"net/http"
	"os"
	"time"
)

type TransportOpt func(c *CustomTransport)

func Name(name string) TransportOpt {
	return func(c *CustomTransport) {
		c.name = name
	}
}

func Logger(logger io.Writer) TransportOpt {
	return func(c *CustomTransport) {
		c.logger = logger
	}
}

func Handle(handle func(req *http.Request, requestBytes, responseBytes []byte, err error)) TransportOpt {
	return func(c *CustomTransport) {
		c.handle = handle
	}
}

func NewCustomHttpClient(opts ...TransportOpt) *http.Client {
	return &http.Client{Transport: NewCustomTransport(opts...)}
}

func NewCustomTransport(opts ...TransportOpt) http.RoundTripper {
	//Transport: 控制 HTTP 请求的连接管理、代理设置、超时等行为。
	//	参数Proxy:
	//
	//	参数DialContext: 配置连接超时 (Timeout) 和保持活动 (KeepAlive) 时间
	//
	//	参数ForceAttemptHTTP2: 强制尝试使用 HTTP/2
	//
	//	参数MaxIdleConns: 全局空闲连接池的最大连接数
	//	参数IdleConnTimeout: 空闲连接在被关闭前的最大空闲时间。
	//
	//	参数TLSHandshakeTimeout: TLS 握手的超时时间，防止卡在 SSL/TLS 握手阶段
	//	参数ExpectContinueTimeout: 设置发送带 Expect: 100-Continue 标头的请求时，等待服务器响应的超时时间
	//
	transport := &CustomTransport{
		Transport: http.Transport{ // 控制 HTTP 请求的连接管理、代理设置、超时等行为。
			Proxy: http.ProxyFromEnvironment, // 设置代理服务器，通过读取环境变量 HTTP_PROXY 或 HTTPS_PROXY 决定是否使用代理
			DialContext: defaultTransportDialContext(&net.Dialer{
				Timeout:   time.Second * 2,  // 配置连接超时拨号超时设置为 2~5 秒 以容忍更高网络延迟。
				KeepAlive: 30 * time.Second, // 保持活动 (KeepAlive) 时间
			}),
			ForceAttemptHTTP2:     true,             // 强制尝试使用 HTTP/2
			MaxIdleConns:          100,              // "全局"空闲连接池的最大连接数
			IdleConnTimeout:       90 * time.Second, // "全局"空闲连接在被关闭前的最大空闲时间。
			TLSHandshakeTimeout:   10 * time.Second, // TLS 握手的超时时间，防止卡在 SSL/TLS 握手阶段
			ExpectContinueTimeout: 1 * time.Second,  // 置发送带 Expect: 100-Continue 标头的请求时，等待服务器响应的超时时间

			// 调整写入和读取缓冲区大小。
			//WriteBufferSize: 32 * 1024, // 32KB
			//ReadBufferSize:  32 * 1024, // 32K

			// PerHost 设置
			//MaxConnsPerHost: 100,		// "每个主机" 最大连接数，包括正在使用的连接和空闲连接。
			//MaxIdleConnsPerHost: 10, 	// "每个主机" 保持的最大空闲连接数

			// 等待服务器响应头的超时时间
			//ResponseHeaderTimeout: 5 * time.Second,
		},
		logger: os.Stdout,
		name:   "HTTP",
	}
	for _, opt := range opts {
		opt(transport)
	}
	return transport
}

type CustomTransport struct {
	http.Transport
	name   string
	logger io.Writer
	handle func(req *http.Request, requestBytes, responseBytes []byte, err error)
}

func (c *CustomTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	var requestBytes []byte
	var responseBytes []byte
	var requestReader io.ReadCloser
	var responseReader io.ReadCloser
	var err error

	url := req.URL.Scheme + "://" + req.URL.Host + req.URL.Path
	if len(req.URL.RawQuery) > 0 {
		url = url + "?" + req.URL.RawQuery
	}

	var span *opentelemetry.Span
	if req.Method != http.MethodHead {
		if s, err := opentelemetry.SpanFromContext(req.Context()); s != nil && err == nil {
			span = opentelemetry.NewSpan(req.Context(), fmt.Sprintf("%s %s", c.name, req.URL.Host))
		}
	}
	// 链路信息写入http header
	opentelemetry.InjectHttp(req.Context(), req)

	defer func() {
		if c.handle != nil {
			c.handle(req, requestBytes, responseBytes, err)
		}

		if req.Method != http.MethodHead {
			c.writeLog("url: " + url)
			c.writeLog("method: " + req.Method)
			c.writeLog("request: " + string(requestBytes))
			c.writeLog("response: " + string(responseBytes))
			if err != nil {
				c.writeLog("err: " + err.Error())
			}
			c.writeLog("=================================================================")
		}

		if span != nil {
			span.AddEvent(fmt.Sprintf("%s: %s", req.Method, url))
			if len(requestBytes) > 0 {
				span.AddEvent(string(requestBytes))
			}
			span.AddEvent(string(responseBytes))
			span.Close()
		}
	}()

	requestBytes, requestReader, err = readBody(req.Body)
	if err != nil {
		return nil, err
	}
	req.Body = requestReader

	res, err := c.Transport.RoundTrip(req)
	if err != nil {
		return nil, err
	}
	responseBytes, responseReader, err = readBody(res.Body)
	res.Body = responseReader
	return res, err
}

func readBody(body io.ReadCloser) ([]byte, io.ReadCloser, error) {
	if body == nil {
		return nil, nil, nil
	}

	b, err := io.ReadAll(body)
	if err != nil {
		return nil, nil, err
	}
	_ = body.Close()
	return b, io.NopCloser(bytes.NewReader(b)), nil
}

func (c *CustomTransport) writeLog(str string) {
	_, _ = c.logger.Write([]byte(str + "\n"))
}

func defaultTransportDialContext(dialer *net.Dialer) func(context.Context, string, string) (net.Conn, error) {
	return dialer.DialContext
}
