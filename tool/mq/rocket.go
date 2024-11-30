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
	AccessKey string `json:"access_key" yaml:"access_key"`
	SecretKey string `json:"secret_key" yaml:"secret_key"`
	Env       string `json:"env" yaml:"env"`
}

type Rocket struct {
	config *RocketMqConfig
}

func NewRocketMq(config *RocketMqConfig) Mq {
	os.Setenv("mq.consoleAppender.enabled", "true")
	rmq.ResetLogger()
	return &Rocket{
		config: config,
	}
}

func (mq *Rocket) Send(ctx context.Context, topic string, key string, message *Message) error {
	return mq.BulkSend(ctx, topic, key, []*Message{message})
}

func (mq *Rocket) BulkSend(ctx context.Context, topic string, key string, message []*Message) error {
	producer, err := rmq.NewProducer(&rmq.Config{
		Endpoint: mq.config.Endpoint,
		Credentials: &credentials.SessionCredentials{
			AccessKey:    mq.config.AccessKey,
			AccessSecret: mq.config.SecretKey,
		},
	}, rmq.WithTopics(topic),
	)
	if err != nil {
		return err
	}
	if err = producer.Start(); err != nil {
		return err
	}
	defer producer.GracefulStop()

	for _, m := range message {
		head, err := parseHeader(m.Headers)
		if err != nil {
			return err
		}
		msg := &rmq.Message{
			Topic: topic,
			Body:  m.Data,
		}

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
			if len(head.MessageGroup) == 0 {
				return errors.New("message group is empty")
			}
			msg.SetMessageGroup(head.MessageGroup)
		}

		//delay 延迟队列
		if head.MsgType == RocketTypeDELAY {
			if head.DelayTimestamp.IsZero() {
				return errors.New("delay timestamp is empty")
			}
			//fmt.Println(head.DelayTimestamp.Format(time.DateTime))
			msg.SetDelayTimestamp(head.DelayTimestamp)
		}

		// * 注意创建topic时需要指定类型（顺序，延迟，普通等），如果类型不一致会发生失败
		if head.MsgType == RocketTypeNORMAL || head.MsgType == RocketTypeFIFO || head.MsgType == RocketTypeDELAY {
			if head.IsAsync {
				producer.SendAsync(ctx, msg, func(ctx context.Context, receipts []*rmq.SendReceipt, err error) {
					fmt.Printf("============%#v\n%v====================", receipts, err)
				})
			} else {
				if _, err = producer.Send(ctx, msg); err != nil {
					return err
				}
			}
			continue
		}

		// 事务
		if head.MsgType == RocketTypeTRANSACTION {
			//SendWithTransaction(context.Context, *Message, Transaction) ([]*SendReceipt, error) //带事务的
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
		Credentials: &credentials.SessionCredentials{
			AccessKey:    mq.config.AccessKey,
			AccessSecret: mq.config.SecretKey,
		},
	},
		rmq.WithAwaitDuration(awaitDuration),
		rmq.WithSubscriptionExpressions(map[string]*rmq.FilterExpression{
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
			if e, ok := err.(*rmq.ErrRpcStatus); ok && e.GetCode() == int32(v2.Code_MESSAGE_NOT_FOUND) {
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
