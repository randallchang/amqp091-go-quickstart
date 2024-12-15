package main

import (
	"context"
	"strconv"

	mq "github.com/randallchang/amqp091-go-quickstart"
)

const (
	url      = "amqp://guest:guest@localhost:/"
	exchange = "test-exchange"
)

func main() {
	producer := mq.NewProducer(url)
	_ = producer.Connect()
	defer producer.Close()

	for i := 0; i < 1000; i++ {
		_ = producer.Publish(
			context.Background(),
			exchange,
			strconv.Itoa(i%2),
			getMsg(i))
	}
}

func getMsg(i int) *mq.Message {
	return &mq.Message{
		CorrelationID: strconv.Itoa(i),
		Body:          []byte("hello, world"),
	}
}
