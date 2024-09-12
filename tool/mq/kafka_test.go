package mq

import (
	"context"
	"github.com/long250038728/web/tool/configurator"
	"testing"
)

var topic = "bonus_message_queue_kafka"
var ctx = context.Background()
var consumerGroup = "hume_2"

var client Mq

func init() {
	var kafkaConf Config
	configurator.NewYaml().MustLoadConfigPath("kafka.yaml", &kafkaConf)
	client = NewKafkaMq(&kafkaConf)
}

//func TestMqCreateTopic(t *testing.T) {
//	err := client.CreateTopic(ctx, "aaa", 1, 1)
//	if err != nil {
//		t.Log(err)
//		return
//	}
//	t.Log("success")
//}

//func TestMqDelTopic(t *testing.T) {
//	err := client.DeleteTopic(ctx, "aaa")
//	if err != nil {
//		t.Log(err)
//		return
//	}
//	t.Log("success")
//}

func TestMqSend(t *testing.T) {
	message, err := NewMessage("hello1")
	if err != nil {
		t.Log(err)
	}
	err = client.Send(ctx, topic, "", message)
	if err != nil {
		t.Log(err)
	}
	t.Log("success")
}

func TestMqBulkSend(t *testing.T) {
	message, err := NewMessage([]byte("hello2"))
	if err != nil {
		t.Log(err)
	}
	err = client.BulkSend(ctx, topic, "", []*Message{message})
	if err != nil {
		t.Log(err)
	}
	t.Log("success")
}

func TestMqSubscribe(t *testing.T) {
	client.Subscribe(ctx, "canal", consumerGroup, func(ctx context.Context, message *Message, err error) error {
		// 是否错误 （程序退出 或 reader报错）
		if err != nil {
			t.Log(err)
			return nil
		}

		//处理业务
		t.Log(string(message.Data))
		return nil
	})
}
