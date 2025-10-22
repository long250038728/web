package mq

import (
	"context"
	"errors"
	"fmt"
	rmq "github.com/apache/rocketmq-clients/golang/v5"
	"github.com/apache/rocketmq-clients/golang/v5/credentials"
	v2 "github.com/apache/rocketmq-clients/golang/v5/protocol/v2"
	"os"
	"time"
)

type RocketMqConfig struct {
	Endpoint  string `json:"endpoint" yaml:"endpoint"`
	NameSpace string `json:"name_space" yaml:"name_space"`
	AccessKey string `json:"access_key" yaml:"access_key"`
	SecretKey string `json:"secret_key" yaml:"secret_key"`
	Env       string `json:"env" yaml:"env"`
}

type Rocket struct {
	config *RocketMqConfig
}

func init() {
	_ = os.Setenv("mq.consoleAppender.enabled", "true")
	rmq.ResetLogger()
}

func NewRocketMq(config *RocketMqConfig) Mq {
	return &Rocket{config: config}
}

func (mq *Rocket) Send(ctx context.Context, topic string, key string, message *Message) error {
	return mq.BulkSend(ctx, topic, key, []*Message{message})
}

func (mq *Rocket) BulkSend(ctx context.Context, topic string, key string, message []*Message) error {
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

	producer, err := rmq.NewProducer(conf, opts...)
	if err != nil {
		return err
	}
	if err = producer.Start(); err != nil {
		return err
	}
	defer func() {
		_ = producer.GracefulStop()
	}()

	for _, m := range message {
		head, err := parseHeader(m.Headers)
		if err != nil {
			return err
		}

		//fifo 顺序队列
		if head.MsgType == RocketTypeFIFO && len(head.MessageGroup) == 0 {
			return errors.New("message group is empty")
		}

		//delay 延迟队列
		if head.MsgType == RocketTypeDELAY && head.DelayTimestamp.IsZero() {
			return errors.New("delay timestamp is empty")
		}

		//transaction 事务消息
		if head.MsgType == RocketTypeTRANSACTION {
			return errors.New("this method not support transaction")
		}
	}

	for _, m := range message {
		head, _ := parseHeader(m.Headers)
		msg := &rmq.Message{Topic: topic, Body: m.Data}

		if key != "" {
			msg.SetKeys(key) //一种为消息设置的唯一标识符或分类信息
		}
		if head.Tag != "" {
			msg.SetTag(head.Tag) //消费者则可以通过订阅特定的 tag 来过滤并消费消息(指定)
		} else {
			msg.SetTag(mq.config.Env) //消费者则可以通过订阅特定的 tag 来过滤并消费消息(如果不存在使用配置的env环境变量)
		}

		//fifo 顺序队列
		if head.MsgType == RocketTypeFIFO {
			msg.SetMessageGroup(head.MessageGroup)
		}

		//delay 延迟队列
		if head.MsgType == RocketTypeDELAY {
			msg.SetDelayTimestamp(head.DelayTimestamp)
		}

		// * 注意创建topic时需要指定类型（顺序，延迟，普通等），如果类型不一致会发生失败
		if head.IsAsync {
			producer.SendAsync(ctx, msg, func(ctx context.Context, receipts []*rmq.SendReceipt, err error) {
				fmt.Printf("============%#v\n%v====================", receipts, err)
			})
			continue
		}
		if _, err = producer.Send(ctx, msg); err != nil {
			return err
		}
	}

	return nil
}

//   ./mqadmin updateTopic -n <nameserver_address> -t <topic_name> -c <cluster_name> -a +message.type=<message_type>
//   ./mqadmin updateTopic  -n rmqnamesrv:9876 -t d_topic  -c DefaultCluster  -o true -a  +message.type=DELAY
//其中<message_type>可以替换为NORMAL、FIFO、DELAY或TRANSACTION

func (mq *Rocket) Subscribe(subscribeCtx context.Context, topic, consumerGroup string, callback func(ctx context.Context, c *Message, err error) error) error {
	var (
		awaitDuration     = time.Second * 5  // maximum waiting time for receive func
		maxMessageNum     = 16               // maximum number of messages received at one time
		invisibleDuration = time.Second * 20 // invisibleDuration should > 20s
	)

	//默认监听所有，但环境变量有值时监听环境变量的tag
	filter := rmq.SUB_ALL
	if len(mq.config.Env) > 0 {
		filter = rmq.NewFilterExpression(mq.config.Env)
	}

	ctx := context.Background()
	simpleConsumer, err := rmq.NewSimpleConsumer(&rmq.Config{
		Endpoint:      mq.config.Endpoint,
		ConsumerGroup: consumerGroup,
		NameSpace:     mq.config.NameSpace,
		Credentials: &credentials.SessionCredentials{
			AccessKey:    mq.config.AccessKey,
			AccessSecret: mq.config.SecretKey,
		},
	},
		rmq.WithSimpleAwaitDuration(awaitDuration),
		rmq.WithSimpleSubscriptionExpressions(map[string]*rmq.FilterExpression{
			topic: filter,
		}),
	)
	if err != nil {
		return err
	}
	if err = simpleConsumer.Start(); err != nil {
		return err
	}
	defer func() {
		_ = simpleConsumer.GracefulStop()
	}()

	for {
		select {
		case <-subscribeCtx.Done():
			return subscribeCtx.Err()
		default:
		}

		mvs, err := simpleConsumer.Receive(subscribeCtx, int32(maxMessageNum), invisibleDuration)
		if err != nil {
			var e *rmq.ErrRpcStatus
			if errors.As(err, &e) && e.GetCode() == int32(v2.Code_MESSAGE_NOT_FOUND) {
				continue
			}
			_ = callback(ctx, nil, err)
			continue
		}

		for _, mv := range mvs {
			//返回err则不提交偏移
			if err = callback(ctx, &Message{Data: mv.GetBody(), Headers: nil}, nil); err != nil {
				continue
			}
			//3次重试
			for retry := 0; retry < 3; retry++ {
				if commitErr := simpleConsumer.Ack(subscribeCtx, mv); commitErr == nil {
					break
				}
			}
		}
	}
}
