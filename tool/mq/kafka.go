package mq

import (
	"context"
	"errors"
	"github.com/segmentio/kafka-go"
	"time"
)

// go get github.com/segmentio/kafka-go

const envKey = "env"

type Config struct {
	Address      string        `json:"address" yaml:"address"`
	Env          string        `json:"env" yaml:"env"`
	BatchSize    int           `json:"batch_size" yaml:"batch_size"`
	BatchBytes   int64         `json:"batch_bytes" yaml:"batchBytes"`
	BatchTimeout time.Duration `json:"batch_timeout" yaml:"batchTimeout"`
}

//======================================================================================================================

type operate struct {
	config *Config
}

func NewKafkaOperate(config *Config) Operate {
	return &operate{config: config}
}

// CreateTopic 创建主题
func (m *operate) CreateTopic(ctx context.Context, topic string, numPartitions int, replicationFactor int) error {
	//如果外部关闭了就不退出循环
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	if m.config.Address == "" {
		return errors.New("IP / Project Not Find")
	}

	conn, err := kafka.Dial("tcp", m.config.Address) // 未测试多主机地址
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
func (m *operate) DeleteTopic(ctx context.Context, topic string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	if m.config.Address == "" {
		return errors.New("IP / Project Not Find")
	}

	conn, err := kafka.Dial("tcp", m.config.Address) // 未测试多主机地址
	if err != nil {
		return err
	}
	defer func() {
		_ = conn.Close()
	}()
	return conn.DeleteTopics(topic)
}

//======================================================================================================================

type producer struct {
	writer *kafka.Writer
	config *Config
}

func NewKafkaProducer(config *Config) Producer {
	return &producer{
		writer: &kafka.Writer{
			Addr:                   kafka.TCP(config.Address),
			BatchSize:              config.BatchSize,    // 批数量（默认100）
			BatchBytes:             config.BatchBytes,   // 批大小(默认1048576)
			BatchTimeout:           config.BatchTimeout, // 批时间(默认1s)
			RequiredAcks:           kafka.RequireOne,    // 0:无需主节点写入成功  1:需要主节点写入成功  -1:所有的ISR节点写入成功
			AllowAutoTopicCreation: false,               // 主题不存在不自动创建主题
		},
		config: config,
	}
}

// Send 发送消息
func (m *producer) Send(ctx context.Context, topic string, key string, message *Message) error {
	return m.BulkSend(ctx, topic, key, []*Message{message})
}

// BulkSend 批量发送消息
func (m *producer) BulkSend(ctx context.Context, topic string, key string, messages []*Message) error {
	list := make([]kafka.Message, 0, len(messages))

	// 通过自定义的message 转换 为 kafka内部的message
	for _, message := range messages {
		headers := make([]kafka.Header, 0, len(message.Headers)+1)

		// config中如果带有环境变量，那么就把环境变量的值写入到header中
		if len(m.config.Env) > 0 {
			headers = append(headers, kafka.Header{Key: envKey, Value: []byte(m.config.Env)})
		}

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
	return m.writer.WriteMessages(ctx, list...)
}

func (m *producer) Close() error {
	return m.writer.Close()
}

//======================================================================================================================

type consumer struct {
	reader *kafka.Reader
	config *Config
}

func NewKafkaConsumer(config *Config, topic, consumerGroup string) Consumer {
	return &consumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: []string{config.Address}, // Kafka broker地址
			Topic:   topic,                    // 消费的主题
			GroupID: consumerGroup,            // 消费者组
		}),
		config: config,
	}
}

// Subscribe 消费者
func (m *consumer) Subscribe(subscribeCtx context.Context, callback func(ctx context.Context, c *Message, err error) error) error {
	// 循环读取消息
	for {
		ctx := context.Background()

		//如果外部关闭退出循环
		select {
		case <-subscribeCtx.Done():
			return subscribeCtx.Err()
		default:
		}

		// 读取消息
		kafkaMessage, err := m.reader.FetchMessage(subscribeCtx)
		if err != nil {
			_ = callback(ctx, nil, err)
			continue
		}

		// 通过kafka内部的message 转换为 自定义的message
		env := true
		headers := make([]Header, 0, len(kafkaMessage.Headers))
		for _, header := range kafkaMessage.Headers {
			if len(m.config.Env) > 0 && header.Key == envKey && m.config.Env != string(header.Value) { // config中如果带有环境变量，那么判断环境变量的值与消息中的环境变量的值是否一致，如果不一致就提交不处理
				env = false
				break
			}
			headers = append(headers, Header{Key: header.Key, Value: header.Value})
		}

		//环境不同，不处理消息(直接提交)
		if !env {
			_ = m.reader.CommitMessages(subscribeCtx, kafkaMessage)
			continue
		}

		// 成功才提交
		if err = callback(ctx, &Message{Data: kafkaMessage.Value, Headers: headers}, nil); err != nil {
			continue
		}
		//3次重试
		for retry := 0; retry < 3; retry++ {
			if commitErr := m.reader.CommitMessages(subscribeCtx, kafkaMessage); commitErr == nil {
				break
			}
		}
	}
}

func (m *consumer) Close() error {
	return m.reader.Close()
}
