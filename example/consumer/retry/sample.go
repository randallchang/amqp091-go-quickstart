package main

import (
	"context"
	"fmt"
	"time"

	mq "github.com/randallchang/amqp091-go-quickstart"
)

const (
	url                  = "amqp://guest:guest@localhost:/"
	retryLimit           = 3
	retryCountKey        = "x-retry-count"
	retryIntervalTimeKey = "x-delay"
	retryIntervalMillis  = 1000
	exchange             = "test-exchange"
)

var (
	consumer mq.Consumer
	producer mq.Producer
)

func main() {
	ctx := context.Background()

	consumer = mq.NewConsumer(
		url,
		mq.WithConsumerConfig(
			&mq.ConsumerConfig{Requeue: true},
		),
		mq.WithConsumerSessionConfig(&mq.SessionConfig{
			ConcurrencyLimit: 1 << 10,
			Heartbeat:        10 * time.Second,
			Locale:           "en_US",
		}),
	)
	producer = mq.NewProducer(url)
	_ = consumer.Connect()
	_ = producer.Connect()
	defer consumer.Close()
	defer producer.Close()

	_ = producer.Publish(
		ctx,
		exchange,
		"0",
		&mq.Message{
			CorrelationID: "1",
			Body:          []byte("hello, world"),
		})

	_ = consumer.Consume(ctx, "test-queue", handler)
}

func handler(ctx context.Context, msg *mq.Message) error {
	fmt.Printf("consume msg correlationId: [%v]", msg.CorrelationID)
	fmt.Printf("consume msg content: [%v]", string(msg.Body))
	_ = requeue(ctx, msg)
	return nil
}

func requeue(ctx context.Context, msg *mq.Message) error {
	if len(msg.Headers) == 0 {
		msg.Headers = make(map[string]interface{})
	}

	var retryCount int32
	if msg.Headers[retryCountKey] != nil {
		retryCount = msg.Headers[retryCountKey].(int32)
	}
	retryCount += 1
	msg.Headers[retryCountKey] = retryCount
	msg.Headers[retryIntervalTimeKey] = retryIntervalMillis

	routingKey := "0"
	if retryCount >= retryLimit {
		routingKey = "1"
	}

	if err := producer.Publish(ctx, exchange, routingKey, msg); err != nil {
		return err
	}

	return nil
}
