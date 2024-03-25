package kafka

import (
	"context"
	"github.com/long250038728/web/tool/mq"
	"github.com/segmentio/kafka-go"
)

func (m *Mq) Subscribe(ctx context.Context, topic string, consumerGroup string, callback func(c *mq.Message, err error) error) {
	// 设置Kafka消费者配置
	config := kafka.ReaderConfig{
		Brokers: m.address,     // Kafka broker地址
		Topic:   topic,         // 消费的主题
		GroupID: consumerGroup, // 消费者组
	}

	// 创建Kafka消费者
	reader := kafka.NewReader(config)
	defer reader.Close()

	// 循环读取消息
	for {
		//如果外部关闭了就不退出循环

		// 读取消息
		kafkaMessage, err := reader.FetchMessage(ctx)
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
				commitErr := reader.CommitMessages(ctx, kafkaMessage)
				if commitErr == nil || ErrRetryNum >= MaxRetryNum {
					break
				}
				ErrRetryNum++
			}
		}
	}
}
