package client

import (
	"context"
	"testing"
)

func TestClient_GetAccessToken(t *testing.T) {
	t.Log(cli.GetAccessToken(context.Background()))
}
