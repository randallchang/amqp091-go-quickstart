package amqp091quickstart

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	_ Producer = (*producer)(nil)

	defaultProducerConfig = &ProducerConfig{
		ContentType:  "text/plain",
		DeliveryMode: amqp.Persistent,
	}
)

type Producer interface {
	Connect() error
	Close() error
	Publish(ctx context.Context, exchange, routingKey string, msg *Message) error
}

type producer struct {
	*session
	producerConfig *ProducerConfig
}

func NewProducer(url string, options ...ProducerOption) Producer {
	p := &producer{
		session:        NewSession(url),
		producerConfig: defaultProducerConfig,
	}
	for _, option := range options {
		option(p)
	}

	return p
}

func (p *producer) Publish(
	ctx context.Context,
	exchange string,
	routingKey string,
	msg *Message,
) error {
	return p.channel.PublishWithContext(
		ctx,
		exchange,
		routingKey,
		false,
		false,
		amqp.Publishing{
			CorrelationId:   msg.CorrelationID,
			Headers:         msg.Headers,
			Body:            msg.Body,
			ContentType:     p.producerConfig.ContentType,
			ContentEncoding: p.producerConfig.ContentEncoding,
			DeliveryMode:    p.producerConfig.DeliveryMode,
			Priority:        p.producerConfig.Priority,
			ReplyTo:         p.producerConfig.ReplyTo,
			Expiration:      p.producerConfig.Expiration,
			MessageId:       p.producerConfig.MessageId,
			Timestamp:       p.producerConfig.Timestamp,
			Type:            p.producerConfig.Type,
			UserId:          p.producerConfig.UserId,
			AppId:           p.producerConfig.AppId,
		})
}
