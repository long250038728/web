package kafka

import (
	"context"
	"github.com/long250038728/web/tool/mq"
	"github.com/segmentio/kafka-go"
)

// Send 发送消息
func (m *Mq) Send(ctx context.Context, topic string, key string, message *mq.Message) error {
	return m.BulkSend(ctx, topic, key, []*mq.Message{message})
}

// BulkSend 批量发送消息
func (m *Mq) BulkSend(ctx context.Context, topic string, key string, messages []*mq.Message) error {
	list := make([]kafka.Message, 0, len(messages))

	for _, message := range messages {
		headers := make([]kafka.Header, 0, len(message.Headers))
		for _, header := range message.Headers {
			headers = append(headers, kafka.Header{Key: header.Key, Value: header.Value})
		}
		msg := kafka.Message{
			Topic:   topic,
			Key:     []byte(key),
			Headers: headers,
			Value:   message.Data,
		}
		list = append(list, msg)
	}

	w := &kafka.Writer{
		Addr:                   kafka.TCP(m.address...),
		BatchSize:              len(messages),
		RequiredAcks:           1,
		AllowAutoTopicCreation: false,
	}
	defer w.Close()
	return w.WriteMessages(ctx, list...)
}
