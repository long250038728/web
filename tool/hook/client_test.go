package hook

import (
	"context"
	"testing"
)

func TestNewQyHookClient(t *testing.T) {
	client, err := NewQyHookClient(&Config{""})
	if err != nil {
		t.Error(err)
		return
	}
	err = client.SendHook(context.Background(), "this is test", nil)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("ok")
}
