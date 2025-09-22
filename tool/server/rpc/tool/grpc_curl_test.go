package tool

import (
	"context"
	"testing"
)

func TestCurl(t *testing.T) {
	ctx := context.Background()
	path := []string{"/Users/linlong/Desktop/web/protoc/order/", "/Users/linlong/Desktop/web/protoc/proto"}
	file := "order.proto"

	address := "192.168.1.101:19002"

	newC := NewCurl(path, file)
	serverMethods, err := newC.GetServerMethods()
	if err != nil {
		t.Error(err)
		return
	}

	for _, serverMethod := range serverMethods {
		resp, err := newC.Curl(ctx, address, serverMethod, map[string]any{"hello": "world"}, map[string]any{"id": 256253})
		if err != nil {
			t.Error(err)
			return
		}
		t.Log(resp)
	}
}
