package kafka

import (
	"context"
	"github.com/long250038728/web/tool/server/http/tool"
	"github.com/segmentio/kafka-go"
)

// go get github.com/segmentio/kafka-go

type Mq struct {
	address []string
}

func NewKafkaMq(config *Config) *Mq {
	return &Mq{
		address: config.Address,
	}
}

// CreateTopic 创建主题
func (m *Mq) CreateTopic(ctx context.Context, topic string, numPartitions int, replicationFactor int) error {
	//如果外部关闭了就不退出循环
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	if len(m.address) == 0 || m.address[0] == "" {
		return tool.Address
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

// DeleteTopic 删除主题
func (m *Mq) DeleteTopic(ctx context.Context, topic string) error {
	//如果外部关闭了就不退出循环
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	if len(m.address) == 0 || m.address[0] == "" {
		return tool.Address
	}

	conn, err := kafka.Dial("tcp", m.address[0]) // 未测试多主机地址
	if err != nil {
		return err
	}
	defer conn.Close()
	return conn.DeleteTopics(topic)
}
