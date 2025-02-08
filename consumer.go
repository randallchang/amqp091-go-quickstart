package amqp091quickstart

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	_                     Consumer = (*consumer)(nil)
	defaultConsumerConfig          = &ConsumerConfig{}
)

type Consumer interface {
	Connect() error
	Close() error
	Consume(ctx context.Context, queue string, handler func(ctx context.Context, msg *Message) error) error
}

type consumer struct {
	*session
	consumerConfig *ConsumerConfig
}

func NewConsumer(url string, options ...ConsumerOption) Consumer {
	c := &consumer{
		session:        NewSession(url),
		consumerConfig: defaultConsumerConfig,
	}
	for _, option := range options {
		option(c)
	}

	return c
}

func (c *consumer) Consume(
	ctx context.Context,
	queue string,
	handler func(ctx context.Context, msg *Message) error,
) error {
	err := c.channel.Qos(1, 0, false)
	if err != nil {
		return err
	}

	deliveryChan, err := c.channel.Consume(
		queue,
		c.consumerConfig.Consumer,
		c.consumerConfig.AutoAck,
		c.consumerConfig.Exclusive,
		c.consumerConfig.NoLocal,
		c.consumerConfig.NoWait,
		amqp.Table(c.consumerConfig.Args),
	)
	if err != nil {
		return err
	}

	go c.handleReconnect(&deliveryChan, queue)

	for {
		err := c.consumeDelivery(ctx, deliveryChan, handler)
		if err != nil {
			return err
		}
	}
}

func (c *consumer) consumeDelivery(
	ctx context.Context,
	deliveryChan <-chan amqp.Delivery,
	handler func(ctx context.Context, msg *Message) error,
) error {
	for delivery := range deliveryChan {
		err := handler(
			ctx,
			&Message{
				CorrelationID: delivery.CorrelationId,
				Headers:       delivery.Headers,
				Body:          delivery.Body,
			})
		if err != nil {
			_ = delivery.Nack(false, c.consumerConfig.Requeue)
			return err
		}
		if !c.consumerConfig.AutoAck {
			_ = delivery.Ack(false)
		}
	}

	return nil
}

func (c *consumer) handleReconnect(deliveryChan *<-chan amqp.Delivery, queue string) {
	for {
		select {
		case <-c.Exit:
			return
		case <-c.Reconnect:
			*deliveryChan, _ = c.channel.Consume(
				queue,
				c.consumerConfig.Consumer,
				c.consumerConfig.AutoAck,
				c.consumerConfig.Exclusive,
				c.consumerConfig.NoLocal,
				c.consumerConfig.NoWait,
				amqp.Table(c.consumerConfig.Args),
			)
		}
	}
}
