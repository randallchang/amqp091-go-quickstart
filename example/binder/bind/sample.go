package main

import (
	"time"

	mq "github.com/randallchang/amqp091-go-quickstart"
)

const (
	url          = "amqp://guest:guest@localhost:/"
	exchange     = "test-exchange"
	exchangeType = mq.DelayedMessage
)

func main() {
	binder := mq.NewBinder(
		url,
		mq.WithExchangeConfig(&mq.ExchangeConfig{
			Durable: true,
			Args:    mq.Table{"x-delayed-type": "direct"},
		}),
		mq.WithBinderSessionConfig(&mq.SessionConfig{
			ConcurrencyLimit: 1 << 10,
			Heartbeat:        10 * time.Second,
			Locale:           "en_US",
		}),
	)

	_ = binder.Connect()
	defer binder.Close()

	_ = binder.Bind(
		exchange,
		exchangeType,
		"test-queue",
		"0")

	_ = binder.Bind(
		exchange,
		exchangeType,
		"test-queue.manual",
		"1")
}
