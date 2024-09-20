package mq

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	rmq "github.com/apache/rocketmq-clients/golang"
	"github.com/apache/rocketmq-clients/golang/credentials"
)

const (
	Topic     = "xxxxxx"
	Endpoint  = "xxxxxx"
	AccessKey = "xxxxxx"
	SecretKey = "xxxxxx"
)

func main() {
	// log to console
	os.Setenv("mq.consoleAppender.enabled", "true")
	rmq.ResetLogger()
	// new producer instance
	producer, err := rmq.NewProducer(&rmq.Config{
		Endpoint: Endpoint,
		Credentials: &credentials.SessionCredentials{
			AccessKey:    AccessKey,
			AccessSecret: SecretKey,
		},
	},
		rmq.WithTransactionChecker(&rmq.TransactionChecker{
			Check: func(msg *rmq.MessageView) rmq.TransactionResolution {
				log.Printf("check transaction message: %v", msg)
				return rmq.COMMIT
			},
		}),
		rmq.WithTopics(Topic),
	)
	if err != nil {
		log.Fatal(err)
	}
	// start producer
	err = producer.Start()
	if err != nil {
		log.Fatal(err)
	}
	// graceful stop producer
	defer producer.GracefulStop()
	for i := 0; i < 10; i++ {
		// new a message
		msg := &rmq.Message{
			Topic: Topic,
			Body:  []byte("this is a message : " + strconv.Itoa(i)),
		}
		// set keys and tag
		msg.SetKeys("a", "b")
		msg.SetTag("ab")
		// send message in sync
		transaction := producer.BeginTransaction()
		resp, err := producer.SendWithTransaction(context.TODO(), msg, transaction)
		if err != nil {
			log.Fatal(err)
		}
		for i := 0; i < len(resp); i++ {
			fmt.Printf("%#v\n", resp[i])
		}
		// commit transaction message
		err = transaction.Commit()
		if err != nil {
			log.Fatal(err)
		}
		// wait a moment
		time.Sleep(time.Second * 1)
	}
}
