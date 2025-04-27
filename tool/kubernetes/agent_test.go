package kubernetes

import (
	"context"
	"testing"
)

func TestAgent_GetLogs(t *testing.T) {
	client, err := NewAgent("/Users/linlong/Downloads/cls-09eyrddg-config")
	if err != nil {
		t.Error(err)
		return
	}
	ctx := context.Background()

	////http://192.168.1.136:8011/agent/info/logs?ns=dev&name=aristotle-6fc64f487b-c76js&container=aristotle
	logs, err := client.GetLogs(ctx, "dev", "aristotle-56d79d59c5-f9dvq", "aristotle")
	t.Log(logs, err)
}

func TestAgent_GetPodEvents(t *testing.T) {
	client, err := NewAgent("/Users/linlong/Downloads/cls-09eyrddg-config")
	if err != nil {
		t.Error(err)
		return
	}
	ctx := context.Background()

	//http://192.168.1.136:8011/agent/info/events?ns=dev&resource=Pod
	logs, err := client.GetPodEvents(ctx, "Pod", "dev")
	t.Log(logs, err)
}

func TestAgent_ListResource(t *testing.T) {
	client, err := NewAgent("/Users/linlong/Downloads/cls-09eyrddg-config")
	if err != nil {
		t.Error(err)
		return
	}
	ctx := context.Background()

	//http://192.168.1.136:8011/agent/info/events?ns=dev&resource=Pod
	logs, err := client.ListResource(ctx, "Pod", "dev")
	t.Log(logs, err)
}

func TestAgent_DeleteResource(t *testing.T) {
	client, err := NewAgent("/Users/linlong/Downloads/cls-09eyrddg-config")
	if err != nil {
		t.Error(err)
		return
	}
	ctx := context.Background()

	err = client.DeleteResource(ctx, "Pod", "dev", "aristotle-56d79d59c5-lzgfm")
	t.Log(err)
}
