package client

import (
	"context"
	"testing"
)

func TestJenkinsClient_Build(t *testing.T) {
	client := NewJenkinsClient("https://jenkins.zhubaoe.cn", "admin", "11739a99e314641a8f7c039db95458f6e1")
	if err := client.Build(context.Background(), "kobe-service-common", map[string]any{
		"BRANCH":      "check",
		"SERVICENAME": "order",
		"ENV":         "check",
	}); err != nil {
		t.Error(err)
		return
	}
	t.Log("ok")
}

func TestJenkinsClient_Queue(t *testing.T) {
	client := NewJenkinsClient("https://jenkins.zhubaoe.cn", "admin", "11739a99e314641a8f7c039db95458f6e1")
	client.Queue(context.Background(), "", nil)
	t.Log("ok")
}
