package mq

import (
	"context"
	"encoding/json"
	"io"
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

type Operate interface {
	CreateTopic(ctx context.Context, topic string, numPartitions int, replicationFactor int) error
	DeleteTopic(ctx context.Context, topic string) error
}

type Producer interface {
	Send(ctx context.Context, topic string, key string, message *Message) error
	BulkSend(ctx context.Context, topic string, key string, message []*Message) error
	io.Closer
}

type Consumer interface {
	Subscribe(subscribeCtx context.Context, callback func(ctx context.Context, c *Message, err error) error) error
	io.Closer
}

//type Transaction interface {
//	Send(ctx context.Context, topic string, key string, m *Message, handle func() bool) error
//}
