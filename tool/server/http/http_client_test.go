package http

import (
	"context"
	"fmt"
	"github.com/long250038728/web/tool/tracing/opentelemetry"
	"testing"
	"time"
)

func TestClient_Post(t *testing.T) {
	ctx := context.Background()
	exporter, err := opentelemetry.NewJaegerExporter("http://link.zhubaoe.cn:14268/api/traces")
	if err != nil {
		t.Error(err)
	}

	trace, err := opentelemetry.NewTrace(ctx, exporter, "ServiceA")
	if err != nil {
		t.Error(err)
	}
	defer func() {
		t.Error(trace.Close(ctx))
	}()

	c := NewClient(SetTimeout(time.Second * 5))

	data := map[string]any{
		"a": "Login",
		"m": "Admin",
		"p": "1",
		"r": "{\"merchant_code\":\"ab190735\",\"user_name\":\"yzt\",\"password\":\"123456\",\"last_admin_id\":\"\",\"last_admin_name\":\"\",\"shift_status\":\"1\"}",
		"t": "00000",
		"v": "2.4.4",
	}

	res, code, err := c.Post(context.Background(), "http://test.zhubaoe.cn:9999/", data)
	fmt.Println(string(res))
	fmt.Println(code)
	fmt.Println(err)

}

func TestClient_Get(t *testing.T) {
	ctx := context.Background()
	exporter, err := opentelemetry.NewJaegerExporter("http://link.zhubaoe.cn:14268/api/traces")
	if err != nil {
		t.Error(err)
	}

	trace, err := opentelemetry.NewTrace(ctx, exporter, "ServiceA")
	if err != nil {
		t.Error(err)
	}
	defer func() {
		t.Error(trace.Close(ctx))
	}()

	c := NewClient(SetTimeout(time.Second * 5))

	data := map[string]any{
		"merchant_id":      394,
		"merchant_shop_id": 1150,
		"start_date":       "2023-12-01",
		"end_date":         "2023-12-01",
		"field":            "goods_type_id",
		"client_name":      "app",
	}

	res, code, err := c.Get(ctx, "http://test.zhubaoe.cn:8888/report/sale_report/inventory", data)
	fmt.Println(string(res))
	fmt.Println(code)
	fmt.Println(err)

}
