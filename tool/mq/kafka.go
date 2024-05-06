package mq

import (
	"context"
	"errors"
	"github.com/segmentio/kafka-go"
)

// go get github.com/segmentio/kafka-go

type Config struct {
	Address []string `json:"address" yaml:"address"`
}

type Kafka struct {
	config *Config
}

func NewKafkaMq(config *Config) *Kafka {
	return &Kafka{
		config: config,
	}
}

//======================================================================================================================

// CreateTopic 创建主题
func (m *Kafka) CreateTopic(ctx context.Context, topic string, numPartitions int, replicationFactor int) error {
	//如果外部关闭了就不退出循环
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	if len(m.config.Address) == 0 || m.config.Address[0] == "" {
		return errors.New("IP / Project Not Find")
	}

	conn, err := kafka.Dial("tcp", m.config.Address[0]) // 未测试多主机地址
	if err != nil {
		return err
	}
	defer func() {
		_ = conn.Close()
	}()

	topicConfigs := []kafka.TopicConfig{
		{
			Topic:             topic,
			NumPartitions:     numPartitions,
			ReplicationFactor: replicationFactor,
		},
	}
	return conn.CreateTopics(topicConfigs...)
}

// DeleteTopic 删除主题
func (m *Kafka) DeleteTopic(ctx context.Context, topic string) error {
	//如果外部关闭了就不退出循环
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	if len(m.config.Address) == 0 || m.config.Address[0] == "" {
		return errors.New("IP / Project Not Find")
	}

	conn, err := kafka.Dial("tcp", m.config.Address[0]) // 未测试多主机地址
	if err != nil {
		return err
	}
	defer func() {
		_ = conn.Close()
	}()
	return conn.DeleteTopics(topic)
}

//======================================================================================================================

// Send 发送消息
func (m *Kafka) Send(ctx context.Context, topic string, key string, message *Message) error {
	return m.BulkSend(ctx, topic, key, []*Message{message})
}

// BulkSend 批量发送消息
func (m *Kafka) BulkSend(ctx context.Context, topic string, key string, messages []*Message) error {
	list := make([]kafka.Message, 0, len(messages))

	// 通过自定义的message 转换 为 kafka内部的message
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
		Addr:                   kafka.TCP(m.config.Address...),
		BatchSize:              len(messages),
		RequiredAcks:           1,     //0:无需主节点写入成功  1:需要主节点写入成功  -1:所有的ISR节点写入成功
		AllowAutoTopicCreation: false, //主题不存在不自动创建主题
	}
	defer func() {
		_ = w.Close()
	}()
	return w.WriteMessages(ctx, list...)
}

//======================================================================================================================

// Subscribe 消费者
func (m *Kafka) Subscribe(ctx context.Context, topic string, consumerGroup string, callback func(c *Message, err error) error) {
	// 设置Kafka消费者配置
	config := kafka.ReaderConfig{
		Brokers: m.config.Address, // Kafka broker地址
		Topic:   topic,            // 消费的主题
		GroupID: consumerGroup,    // 消费者组
	}

	// 创建Kafka消费者
	reader := kafka.NewReader(config)
	defer func() {
		_ = reader.Close()
	}()

	// 循环读取消息
	for {
		//如果外部关闭了就不退出循环

		// 读取消息
		kafkaMessage, err := reader.FetchMessage(ctx)
		if err != nil {
			_ = callback(nil, err)
			continue
		}

		// 通过kafka内部的message 转换为 自定义的message
		// header头处理
		headers := make([]Header, 0, len(kafkaMessage.Headers))
		for _, header := range kafkaMessage.Headers {
			headers = append(headers, Header{Key: header.Key, Value: header.Value})
		}
		message := &Message{Data: kafkaMessage.Value, Headers: headers}
		err = callback(message, nil)

		// 成功才提交
		if err == nil {
			MaxRetryNum := 3
			ErrRetryNum := 0

			//提交失败重试次数
			for {
				commitErr := reader.CommitMessages(ctx, kafkaMessage)
				if commitErr == nil || ErrRetryNum >= MaxRetryNum {
					break
				}
				ErrRetryNum++
			}
		}
	}
}
