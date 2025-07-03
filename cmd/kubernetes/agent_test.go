package kubernetes

import (
	"context"
	"testing"
)

func getClient() *Agent {
	client, err := NewAgent("/Users/linlong/Downloads/cls-09eyrddg-config")
	if err != nil {
		panic(err)
	}
	return client
}

func TestAgent_GetLogs(t *testing.T) {
	logs, err := getClient().GetLogs(context.Background(), "dev", "aristotle-56d79d59c5-f9dvq", "aristotle")
	t.Log(logs, err)
}

func TestAgent_GetPodEvents(t *testing.T) {
	logs, err := getClient().GetPodEvents(context.Background(), "Pod", "dev")
	t.Log(logs, err)
}

func TestAgent_ListResource(t *testing.T) {
	logs, err := getClient().ListResource(context.Background(), "Pod", "dev")
	t.Log(logs, err)
}

func TestAgent_DeleteResource(t *testing.T) {
	err := getClient().DeleteResource(context.Background(), "Pod", "dev", "aristotle-56d79d59c5-lzgfm")
	t.Log(err)
}
