package http

import (
	"context"
	"github.com/long250038728/web/tool/configurator"
	"github.com/long250038728/web/tool/tracing/opentelemetry"
	"net/http"
	"os"
	"sync"
	"testing"
	"time"
)

var ctx = context.Background()

var exporter opentelemetry.SpanExporter

func init() {
	var err error
	var traceConfig opentelemetry.Config
	configurator.NewYaml().MustLoadConfigPath("tracing.yaml", &traceConfig)
	if exporter, err = opentelemetry.NewJaegerExporter(&traceConfig); err != nil {
		panic(err)
	}
}

func initTracing() (*opentelemetry.Trace, error) {
	return opentelemetry.NewTrace(ctx, exporter, "ServiceA")
}

func TestClient_Post(t *testing.T) {
	trace, err := initTracing()
	defer func() {
		if trace != nil {
			_ = trace.Close(ctx)
		}
	}()
	if err != nil {
		t.Error(err)
		return
	}

	httpClient := NewClient(SetTimeout(time.Second))
	data := map[string]any{
		"a": "Login",
		"m": "Admin",
		"p": "1",
		"r": "{\"merchant_code\":\"ab190735\",\"user_name\":\"yzt\",\"password\":\"123456\",\"last_admin_id\":\"\",\"last_admin_name\":\"\",\"shift_status\":\"1\"}",
		"t": "00000",
		"v": "2.4.4",
	}
	res, code, err := httpClient.Post(ctx, "http://test.zhubaoe.cn:9999/", data)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(string(res), code)
}

func TestClient_Get(t *testing.T) {
	trace, err := initTracing()
	defer func() {
		if trace != nil {
			_ = trace.Close(ctx)
		}
	}()
	if err != nil {
		t.Error(err)
		return
	}

	c := NewCustomHttpClient(Logger(os.Stdout), Handle(func(req *http.Request, requestBytes, responseBytes []byte, err error) {
		t.Log(req, requestBytes, responseBytes, err)
	}))

	httpClient := NewClient(SetTimeout(time.Second), SetContentType("application/JSON"), SetHttpClient(c))
	data := map[string]any{
		"merchant_id":      394,
		"merchant_shop_id": 1150,
		"start_date":       "2023-12-01",
		"end_date":         "2023-12-01",
		"field":            "goods_type_id",
		"client_name":      "app",
	}
	res, code, err := httpClient.Get(ctx, "http://test.zhubaoe.cn:8888/report/sale_report/inventory", data)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(string(res), code)
}

func TestClient_GoGet(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(50)
	for i := 0; i < 50; i++ {
		go func() {
			httpClient := NewClient(SetTimeout(time.Second))
			res, code, err := httpClient.Get(ctx, "http://192.168.0.30:8001/", nil)
			wg.Done()
			if err != nil {
				t.Error(err)
				return
			}
			t.Log(string(res), code)
		}()
	}
	wg.Wait()
}
