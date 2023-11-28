package kafka

import (
	"context"
	"errors"
	"github.com/long250038728/web/tool/mq"
	"github.com/segmentio/kafka-go"
)

// go get github.com/segmentio/kafka-go

type Mq struct {
	address []string
}

func NewKafkaMq(address ...string) *Mq {
	return &Mq{
		address: address,
	}
}

func (m *Mq) CreateTopic(context context.Context, topic string, numPartitions int, replicationFactor int) error {
	//如果外部关闭了就不退出循环
	select {
	case <-context.Done():
		return context.Err()
	default:
	}

	if len(m.address) == 0 || m.address[0] == "" {
		return errors.New("address is null")
	}

	conn, err := kafka.Dial("tcp", m.address[0]) // 未测试多主机地址
	if err != nil {
		return err
	}
	defer conn.Close()

	topicConfigs := []kafka.TopicConfig{
		{
			Topic:             topic,
			NumPartitions:     numPartitions,
			ReplicationFactor: replicationFactor,
		},
	}
	return conn.CreateTopics(topicConfigs...)
}

func (m *Mq) DeleteTopic(context context.Context, topic string) error {
	//如果外部关闭了就不退出循环
	select {
	case <-context.Done():
		return context.Err()
	default:
	}

	if len(m.address) == 0 || m.address[0] == "" {
		return errors.New("address is null")
	}

	conn, err := kafka.Dial("tcp", m.address[0]) // 未测试多主机地址
	if err != nil {
		return err
	}
	defer conn.Close()
	return conn.DeleteTopics(topic)
}

func (m *Mq) Send(context context.Context, topic string, key string, message *mq.Message) error {
	return m.BulkSend(context, topic, key, []*mq.Message{message})
}

func (m *Mq) BulkSend(context context.Context, topic string, key string, messages []*mq.Message) error {
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
	return w.WriteMessages(context, list...)
}

func (m *Mq) Subscribe(context context.Context, topic string, consumerGroup string, callback func(c *mq.Message, err error) error) {
	// 设置Kafka消费者配置
	config := kafka.ReaderConfig{
		Brokers: m.address,     // Kafka broker地址
		Topic:   topic,         // 消费的主题
		GroupID: consumerGroup, // 消费者组
	}

	// 创建Kafka消费者
	reader := kafka.NewReader(config)

	// 循环读取消息
	for {
		//如果外部关闭了就不退出循环
		select {
		case <-context.Done():
			_ = callback(nil, context.Err())
			return
		default:
		}

		// 读取消息
		kafkaMessage, err := reader.FetchMessage(context)
		if err != nil {
			_ = callback(nil, err)
			continue
		}

		// header头处理
		headers := make([]mq.Header, 0, len(kafkaMessage.Headers))
		for _, header := range kafkaMessage.Headers {
			headers = append(headers, mq.Header{Key: header.Key, Value: header.Value})
		}
		message := &mq.Message{Data: kafkaMessage.Value, Headers: headers}
		err = callback(message, nil)

		// 成功才提交
		if err == nil {
			MaxRetryNum := 3
			ErrRetryNum := 0

			//提交失败重试次数
			for {
				commitErr := reader.CommitMessages(context, kafkaMessage)
				if commitErr == nil || ErrRetryNum >= MaxRetryNum {
					break
				}
				ErrRetryNum++
			}
		}
	}
}
