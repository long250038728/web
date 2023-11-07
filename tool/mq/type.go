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

type Mq interface {
	CreateTopic(context context.Context, topic string, numPartitions int, replicationFactor int) error
	DeleteTopic(context context.Context, topic string) error

	Send(context context.Context, topic string, key string, message *Message) error
	BulkSend(context context.Context, topic string, key string, message []*Message) error

	Subscribe(context context.Context, topic string, consumerGroup string)
}
