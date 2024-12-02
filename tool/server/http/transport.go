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
	transport := &CustomTransport{
		Transport: http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: defaultTransportDialContext(&net.Dialer{
				Timeout:   time.Second,
				KeepAlive: 30 * time.Second,
			}),
			ForceAttemptHTTP2:     true,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
		logger: os.Stdout,
	}
	for _, opt := range opts {
		opt(transport)
	}
	return transport
}

type CustomTransport struct {
	http.Transport
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
		span = opentelemetry.NewSpan(req.Context(), fmt.Sprintf("HTTP %s", req.URL.Host))
	}

	defer func() {
		if c.handle != nil {
			c.handle(req, requestBytes, responseBytes, err)
		}

		c.writeLog("url: " + url)
		c.writeLog("method: " + req.Method)
		c.writeLog("request: " + string(requestBytes))
		c.writeLog("response: " + string(responseBytes))
		if err != nil {
			c.writeLog("err: " + err.Error())
		}
		c.writeLog("=================================================================")
		c.writeLog(fmt.Sprintf("%v,%v", span == nil, span))

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
