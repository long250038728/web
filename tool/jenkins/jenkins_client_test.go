package jenkins

import (
	"context"
	"testing"
)

var client, _ = NewJenkinsClient(&Config{Address: "https://jenkins.zhubaoe.cn", Username: "admin", Password: "11739a99e314641a8f7c039db95458f6e1"})

func TestJenkinsClient_Build(t *testing.T) {
	var job = "kobe-service-common"
	var data = map[string]any{
		"BRANCH":      "feature/0413",
		"SERVICENAME": "order",
		"ENV":         "dev",
	}
	if err := client.Build(context.Background(), job, data); err != nil {
		t.Error(err)
		return
	}
	t.Log("ok")
}

func TestJenkinsClient_BuildBlock(t *testing.T) {
	var job = "kobe-service-common"
	var data = map[string]any{
		"BRANCH":      "feature/0413",
		"SERVICENAME": "order",
		"ENV":         "dev",
	}
	if err := client.BlockBuild(context.Background(), job, data); err != nil {
		t.Error(err)
		return
	}
	t.Log("ok")
}
