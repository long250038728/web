package client

import (
	"context"
	"testing"
)

func TestClient_GetAccessToken(t *testing.T) {
	cli := &Client{}
	t.Log(cli.GetAccessToken(context.Background()))
}
