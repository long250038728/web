package mq

import (
	"context"
	"encoding/json"
)

type Header struct {
	Key   string
	Value []byte
}

type Message struct {
	Headers []Header
	Data    []byte
}

func NewMessage(data interface{}) (msg *Message, err error) {
	var bytes []byte

	switch d := data.(type) {
	case []byte:
		bytes = d
	case string:
		bytes = []byte(d)
	default:
		bytes, err = json.Marshal(d)
		if err != nil {
			return nil, err
		}
	}
	return &Message{
		Headers: make([]Header, 0, 10),
		Data:    bytes,
	}, nil
}

type KafkaMq interface {
	Operate
	Send
	Subscribe
}

type Mq interface {
	Send
	Subscribe
}

type Operate interface {
	CreateTopic(ctx context.Context, topic string, numPartitions int, replicationFactor int) error
	DeleteTopic(ctx context.Context, topic string) error
}

type Send interface {
	Send(ctx context.Context, topic string, key string, message *Message) error
	BulkSend(ctx context.Context, topic string, key string, message []*Message) error
}

type Subscribe interface {
	Subscribe(ctx context.Context, topic, consumerGroup string, callback func(ctx context.Context, c *Message, err error) error) error
}

type Transaction interface {
	Send(ctx context.Context, topic string, key string, m *Message, handle func() bool) error
}
