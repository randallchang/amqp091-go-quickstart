package amqp091quickstart

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	_ Binder = (*binder)(nil)

	defaultExchangeConfig  = &ExchangeConfig{Durable: true}
	defaultQueueConfig     = &QueueConfig{Durable: true}
	defaultQueueBindConfig = &QueueBindConfig{}
)

type Binder interface {
	Connect() error
	Close() error
	Bind(exchange string, exchangeType ExchangeType, queue string, routingKey string) error
	Unbind(exchange string, queue string, routingKey string) error
	DeleteExchange(exchange string) error
}

type binder struct {
	*session
	exchangeConfig  *ExchangeConfig
	queueConfig     *QueueConfig
	queueBindConfig *QueueBindConfig
}

func NewBinder(url string, options ...BinderOption) Binder {
	b := &binder{
		session:         NewSession(url),
		exchangeConfig:  defaultExchangeConfig,
		queueConfig:     defaultQueueConfig,
		queueBindConfig: defaultQueueBindConfig,
	}
	for _, option := range options {
		option(b)
	}

	return b
}

func (b *binder) Bind(
	exchange string,
	exchangeType ExchangeType,
	queue string,
	routingKey string,
) error {
	if err := b.channel.ExchangeDeclare(
		exchange,
		string(exchangeType),
		b.exchangeConfig.Durable,
		b.exchangeConfig.AutoDelete,
		b.exchangeConfig.Internal,
		b.exchangeConfig.NoWait,
		amqp.Table(b.exchangeConfig.Args)); err != nil {

		return err
	}
	if _, err := b.channel.QueueDeclare(
		queue,
		b.queueConfig.Durable,
		b.queueConfig.AutoDelete,
		b.queueConfig.Exclusive,
		b.queueConfig.NoWait,
		amqp.Table(b.queueConfig.Args)); err != nil {

		return err
	}
	if err := b.channel.QueueBind(
		queue,
		routingKey,
		exchange,
		b.queueBindConfig.NoWait,
		amqp.Table(b.queueBindConfig.Args)); err != nil {
		return err
	}

	return nil
}

func (b *binder) Unbind(
	exchange string,
	queue string,
	routingKey string,
) error {
	if err := b.channel.QueueUnbind(
		queue,
		routingKey,
		exchange,
		amqp.Table(b.queueConfig.Args)); err != nil {
		return err
	}

	return nil
}

func (b *binder) DeleteExchange(exchange string) error {
	if err := b.channel.ExchangeDelete(exchange, true, false); err != nil {
		return err
	}
	return nil
}
