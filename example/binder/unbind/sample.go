package main

import (
	mq "github.com/randallchang/amqp091-go-quickstart"
)

const (
	url          = "amqp://guest:guest@localhost:/"
	exchange     = "test-exchange"
	exchangeType = mq.DelayedMessage
)

func main() {
	binder := mq.NewBinder(url)
	_ = binder.Connect()
	defer binder.Close()

	_ = binder.Unbind(
		exchange,
		"test-queue",
		"0")

	_ = binder.Unbind(
		exchange,
		"test-queue.manual",
		"1")

	_ = binder.DeleteExchange(exchange)
}
