package mq

import (
	"context"
	"github.com/long250038728/web/tool/configurator"
	"testing"
)

var topic = "bonus_message_queue_kafka"
var ctx = context.Background()
var consumerGroup = "hume_2"
var kafkaConf Config

func init() {
	configurator.NewYaml().MustLoadConfigPath("mq.yaml", &kafkaConf)
}

func TestOperate(t *testing.T) {
	client := NewKafkaOperate(&kafkaConf)
	ctx = context.Background()

	t.Run("create", func(t *testing.T) {
		err := client.CreateTopic(ctx, "hello", 2, 1)
		t.Error(err)
	})
	t.Run("delete", func(t *testing.T) {
		err := client.DeleteTopic(ctx, "hello")
		t.Error(err)
	})
}

func TestNewKafkaProducer(t *testing.T) {
	client := NewKafkaProducer(&kafkaConf)
	ctx = context.Background()
	message, err := NewMessage("hello1")
	if err != nil {
		t.Error(err)
		return
	}
	t.Run("send", func(t *testing.T) {
		err := client.Send(ctx, "hello", "", message)
		t.Error(err)
	})
	t.Run("bulkSend", func(t *testing.T) {
		err := client.BulkSend(ctx, "hello", "", []*Message{message})
		t.Error(err)
	})
	t.Error(client.Close())
}

func TestMqSubscribe(t *testing.T) {
	client := NewKafkaConsumer(&kafkaConf, "hello", consumerGroup)
	ctx = context.Background()

	_ = client.Subscribe(ctx, func(ctx context.Context, message *Message, err error) error {
		// 是否错误 （程序退出 或 reader报错）
		if err != nil {
			t.Log(err)
			return nil
		}
		//处理业务
		t.Log(string(message.Data))
		return nil
	})

	t.Error(client.Close())
}
