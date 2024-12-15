package main

import (
	"context"
	"fmt"
	"time"

	mq "github.com/randallchang/amqp091-go-quickstart"
)

const (
	url = "amqp://guest:guest@localhost:/"
)

func main() {
	ctx := context.Background()
	consumer := mq.NewConsumer(
		url,
		mq.WithConsumerConfig(&mq.ConsumerConfig{Requeue: true}),
		mq.WithConsumerSessionConfig(&mq.SessionConfig{
			ConcurrencyLimit: 1 << 10,
			Heartbeat:        10 * time.Second,
			Locale:           "en_US",
		}),
	)
	_ = consumer.Connect()
	defer consumer.Close()

	_ = consumer.Consume(ctx, "test-queue", handler)
	_ = consumer.Consume(ctx, "test-queue.manual", handler)
}

func handler(ctx context.Context, msg *mq.Message) error {
	fmt.Printf("consume msg correlationId: [%v], ", msg.CorrelationID)
	fmt.Printf("consume msg content: [%v]\n", string(msg.Body))
	return nil
}
