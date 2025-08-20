package mq

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestRocket_NORMALSend(t *testing.T) {
	client := NewRocketMq(&RocketMqConfig{
		Endpoint: "rmq-18xpbagp4.rocketmq.gz.public.tencenttdmq.com:8080", AccessKey: "ak18xpbagp483ece0d1cc80", SecretKey: "sk48297705d3c10873",
	})
	t.Run("send", func(t *testing.T) {
		err := client.Send(context.Background(), "NORMAL", "first", &Message{
			Data:    []byte(fmt.Sprintf("%s", "hello")),
			Headers: NewRocketHeader(NewMessageHeaderNORMAL()),
		})
		if err != nil {
			t.Error(err)
			return
		}
		t.Log("ok")
	})

	t.Run("bulk_send", func(t *testing.T) {
		var messages []*Message
		for i := 0; i <= 1000; i++ {
			messages = append(messages, &Message{
				Data:    []byte(fmt.Sprintf("%s%d", "NORMAL", i)),
				Headers: NewRocketHeader(NewMessageHeaderNORMAL()),
			})
		}
		t.Log(client.BulkSend(context.Background(), "NORMAL", "first", messages))
	})
}

func TestRocket_FIFOSend(t *testing.T) {
	client := NewRocketMq(&RocketMqConfig{
		Endpoint: "rmq-18xpbagp4.rocketmq.gz.public.tencenttdmq.com:8080", AccessKey: "ak18xpbagp483ece0d1cc80", SecretKey: "sk48297705d3c10873",
	})

	t.Run("send", func(t *testing.T) {
		err := client.Send(context.Background(), "FIFO", "first", &Message{
			Data:    []byte(fmt.Sprintf("%s", "FIFO")),
			Headers: NewRocketHeader(NewMessageHeaderFIFO("FIFO")),
		})
		if err != nil {
			t.Error(err)
			return
		}
		t.Log("ok")
	})

	t.Run("bulk_send", func(t *testing.T) {
		var messages []*Message
		for i := 0; i <= 1000; i++ {
			messages = append(messages, &Message{
				Data:    []byte(fmt.Sprintf("%s%d", "FIFO", i)),
				Headers: NewRocketHeader(NewMessageHeaderFIFO("FIFO")),
			})
		}
		t.Log(client.BulkSend(context.Background(), "FIFO", "first", messages))
	})
}

func TestRocket_DELAYSend(t *testing.T) {
	client := NewRocketMq(&RocketMqConfig{
		Endpoint: "rmq-18xpbagp4.rocketmq.gz.public.tencenttdmq.com:8080", AccessKey: "ak18xpbagp483ece0d1cc80", SecretKey: "sk48297705d3c10873",
	})

	t.Run("send", func(t *testing.T) {
		err := client.Send(context.Background(), "DELAY", "first", &Message{
			Data:    []byte(fmt.Sprintf("%s", "DELAY")),
			Headers: NewRocketHeader(NewMessageHeaderDELAY(time.Now().Add(10 * time.Second))),
		})
		if err != nil {
			t.Error(err)
			return
		}
		t.Log("ok")
	})
}

//func TestRocket_BulkSend(t *testing.T) {
//	client := NewRocketMq(&RocketMqConfig{
//		Endpoint: "192.168.0.15:8080", AccessKey: "", SecretKey: "",
//	})
//	var messages []*Message
//	for i := 0; i <= 1000; i++ {
//		messages = append(messages, &Message{
//			Data:    []byte(fmt.Sprintf("%s%d", "hello", i)),
//			Headers: NewRocketHeader(NewMessageHeaderNORMAL()),
//		})
//	}
//	t.Log(client.BulkSend(context.Background(), "NORMAL", "first", messages))
//}

func TestRocket_Subscribe(t *testing.T) {
	client := NewRocketMq(&RocketMqConfig{
		Endpoint: "rmq-18xpbagp4.rocketmq.gz.public.tencenttdmq.com:8080", AccessKey: "ak18xpbagp483ece0d1cc80", SecretKey: "sk48297705d3c10873",
	})
	err := client.Subscribe(context.Background(), "DELAY", "consumer_group", func(ctx context.Context, c *Message, err error) error {
		if err != nil {
			t.Error(err)
			return nil
		}
		t.Log(string(c.Data))
		return nil
	})
	t.Log(err)
}
