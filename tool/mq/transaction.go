package mq

import (
	"context"
	"errors"
	"fmt"
	rmq "github.com/apache/rocketmq-clients/golang/v5"
	"github.com/apache/rocketmq-clients/golang/v5/credentials"
	"os"
)

type CheckHandle func(msg *rmq.MessageView) rmq.TransactionResolution

type RocketTransaction struct {
	config      *RocketMqConfig
	checkHandle CheckHandle
}

func init() {
	_ = os.Setenv("mq.consoleAppender.enabled", "true")
	rmq.ResetLogger()
}

func NewRocketTransactionMq(config *RocketMqConfig, checkHandle CheckHandle) Transaction {
	return &RocketTransaction{config: config, checkHandle: checkHandle}
}

func (mq *RocketTransaction) Send(ctx context.Context, topic string, key string, m *Message, handle func() bool) error {
	conf := &rmq.Config{
		Endpoint:  mq.config.Endpoint,
		NameSpace: mq.config.NameSpace,
		Credentials: &credentials.SessionCredentials{
			AccessKey:    mq.config.AccessKey,
			AccessSecret: mq.config.SecretKey,
		},
	}
	opts := []rmq.ProducerOption{
		rmq.WithTopics(topic),
	}

	if mq.checkHandle != nil {
		opts = append(opts, rmq.WithTransactionChecker(&rmq.TransactionChecker{
			Check: func(msg *rmq.MessageView) rmq.TransactionResolution {
				return mq.checkHandle(msg)
			},
		}))
	}

	producer, err := rmq.NewProducer(conf, opts...)
	if err != nil {
		return err
	}
	if err = producer.Start(); err != nil {
		return err
	}
	defer producer.GracefulStop()

	head, err := parseHeader(m.Headers)
	if err != nil {
		return err
	}

	if head.MsgType != RocketTypeTRANSACTION {
		return errors.New("this method is only support transaction")
	}

	msg := &rmq.Message{Topic: topic, Body: m.Data}

	if key != "" {
		msg.SetKeys(key) //一种为消息设置的唯一标识符或分类信息
	}
	if head.Tag != "" {
		msg.SetTag(head.Tag) //消费者则可以通过订阅特定的 tag 来过滤并消费消息(指定)
	} else {
		msg.SetTag(mq.config.Env) //消费者则可以通过订阅特定的 tag 来过滤并消费消息(如果不存在使用配置的env环境变量)
	}

	transaction := producer.BeginTransaction()
	resp, err := producer.SendWithTransaction(context.TODO(), msg, transaction)
	if err != nil {
		return err
	}
	for i := 0; i < len(resp); i++ {
		fmt.Printf("%#v\n", resp[i])
	}

	ok := handle()

	if ok {
		return transaction.Commit()
	} else {
		return transaction.RollBack()
	}
}
