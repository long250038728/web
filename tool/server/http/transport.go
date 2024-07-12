package http

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"
)

type TransportOpt func(c *CustomTransport)

func NewCustomTransport(opts ...TransportOpt) http.RoundTripper {
	transport := &CustomTransport{
		Transport: http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: defaultTransportDialContext(&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}),
			ForceAttemptHTTP2:     true,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}
	for _, opt := range opts {
		opt(transport)
	}
	return transport
}

type CustomTransport struct {
	http.Transport
}

func (c *CustomTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	requestBytes, requestReader, err := readBody(req.Body)
	if err != nil {
		return nil, err
	}
	req.Body = requestReader

	res, err := c.Transport.RoundTrip(req)
	responseBytes, responseReader, err := readBody(res.Body)
	if err != nil {
		return nil, err
	}
	res.Body = responseReader

	fmt.Println("url: " + req.URL.Host + req.URL.Path)
	fmt.Println("request: " + string(requestBytes))
	fmt.Println("response: " + string(responseBytes))
	fmt.Println("=================================================================")
	return res, err
}

func readBody(body io.ReadCloser) ([]byte, io.ReadCloser, error) {
	b, err := io.ReadAll(body)
	if err != nil {
		return nil, nil, err
	}
	_ = body.Close()
	return b, io.NopCloser(bytes.NewReader(b)), nil
}

func defaultTransportDialContext(dialer *net.Dialer) func(context.Context, string, string) (net.Conn, error) {
	return dialer.DialContext
}
