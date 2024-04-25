package client

import (
	"context"
	"testing"
)

func TestJenkinsClient_Build(t *testing.T) {
	client := NewJenkinsClient("https://jenkins.zhubaoe.cn", "admin", "11739a99e314641a8f7c039db95458f6e1")
	if err := client.Build(context.Background(), "kobe-service-common", map[string]any{
		"BRANCH":      "feature/0413",
		"SERVICENAME": "order",
		"ENV":         "dev",
	}); err != nil {
		t.Error(err)
		return
	}
	t.Log("ok")
}

func TestJenkinsClient_Queue(t *testing.T) {
	client := NewJenkinsClient("https://jenkins.zhubaoe.cn", "admin", "11739a99e314641a8f7c039db95458f6e1")
	client.Block(context.Background(), "kobe-service-common", nil)
	t.Log("ok")
}

func TestJenkinsClient_BuildBlock(t *testing.T) {
	client := NewJenkinsClient("https://jenkins.zhubaoe.cn", "admin", "11739a99e314641a8f7c039db95458f6e1")
	data := map[string]any{
		"BRANCH":      "feature/0413",
		"SERVICENAME": "order",
		"ENV":         "dev",
	}
	if err := client.BlockBuild(context.Background(), "kobe-service-common", data); err != nil {
		t.Error(err)
		return
	}
	t.Log("ok")
}
