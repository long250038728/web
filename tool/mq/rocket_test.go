package mq

import (
	"context"
	"fmt"
	"testing"
)

func TestRocket_Send(t *testing.T) {
	for i := 0; i < 1000; i++ {
		client := NewRocketMq(&RocketMqConfig{
			Endpoint: "192.168.0.15:8080", AccessKey: "", SecretKey: "",
		})
		err := client.Send(context.Background(), "my_topic", "first", &Message{
			Data:    []byte(fmt.Sprintf("%s%d", "hello", i)),
			Headers: NewRocketHeader(NewMessageHeaderNORMAL()),
		})
		if err != nil {
			t.Error(err)
			return
		}
		t.Log("ok")
	}
}

func TestRocket_BulkSend(t *testing.T) {
	client := NewRocketMq(&RocketMqConfig{
		Endpoint: "192.168.0.15:8080", AccessKey: "", SecretKey: "",
	})
	var messages []*Message
	for i := 0; i <= 1000; i++ {
		messages = append(messages, &Message{
			Data:    []byte(fmt.Sprintf("%s%d", "hello", i)),
			Headers: NewRocketHeader(NewMessageHeaderNORMAL()),
		})
	}
	t.Log(client.BulkSend(context.Background(), "my_topic", "first", messages))
}

func TestRocket_Subscribe(t *testing.T) {
	client := NewRocketMq(&RocketMqConfig{
		Endpoint: "192.168.0.15:8080", AccessKey: "", SecretKey: "",
	})
	err := client.Subscribe(context.Background(), "my_topic", "consumer_group", func(ctx context.Context, c *Message, err error) error {
		if err != nil {
			t.Error(err)
			return nil
		}
		t.Log(string(c.Data))
		return nil
	})
	t.Log(err)
}
