package mq

import (
	"context"
	config2 "github.com/long250038728/web/tool/config"
	"testing"
	"time"
)

var topic = "bonus_message_queue_kafka"
var ctx = context.Background()
var consumerGroup = "hume_2"
var client *Kafka

func init() {
	var conf Config
	yaml := config2.Yaml{}
	_ = yaml.Load("", &conf)
	client = NewKafkaMq(&conf)
}

func TestMqCreateTopic(t *testing.T) {
	err := client.CreateTopic(ctx, "aaa", 1, 1)
	if err != nil {
		t.Log(err)
		return
	}
	t.Log("success")
}

func TestMqDelTopic(t *testing.T) {
	err := client.DeleteTopic(ctx, "aaa")
	if err != nil {
		t.Log(err)
		return
	}
	t.Log("success")
}

func TestMqSend(t *testing.T) {
	message := &Message{
		Data: []byte("hello1"),
	}
	err := client.Send(ctx, topic, "", message)
	if err != nil {
		t.Log(err)
	}
	t.Log("success")
}

func TestMqBulkSend(t *testing.T) {
	message := &Message{
		Data: []byte("hello2"),
	}
	err := client.BulkSend(ctx, topic, "", []*Message{message})
	if err != nil {
		t.Log(err)
	}
	t.Log("success")
}

func TestMqSubscribe(t *testing.T) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	client.Subscribe(ctx, topic, consumerGroup, func(message *Message, err error) error {
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
