package amqp091quickstart

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"errors"

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
	doneCh         chan struct{}
	allowCh        chan struct{}
	exitCh         chan struct{}
	producerAckCh  chan *amqp.DeferredConfirmation
}

func NewProducer(url string, options ...ProducerOption) Producer {
	p := &producer{
		session:        NewSession(url),
		producerConfig: defaultProducerConfig,
		doneCh:         make(chan struct{}),
		allowCh:        make(chan struct{}, 1),
		exitCh:         make(chan struct{}),
		producerAckCh:  make(chan *amqp.DeferredConfirmation),
	}
	for _, option := range options {
		option(p)
	}

	return p
}

func (p *producer) Connect() error {
	p.concurrencyDaemon()
	return p.session.Connect()
}

func (p *producer) concurrencyDaemon() {
	p.notifyClose()

	go func() {
		producerDeferredConfirmations := make(map[uint64]*amqp.DeferredConfirmation)

		for {
			select {
			case <-p.exitCh:
				exitConcurrencyDaemon(
					p.doneCh,
					p.allowCh,
					p.producerAckCh,
					producerDeferredConfirmations)
				return
			default:
			}

			if len(producerDeferredConfirmations) < p.sessionConfig.ConcurrencyLimit {
				allowExecution(p.allowCh)
			}

			select {
			case confirmation := <-p.producerAckCh:
				producerDeferredConfirmations[confirmation.DeliveryTag] = confirmation
			case <-p.exitCh:
				exitConcurrencyDaemon(
					p.doneCh,
					p.allowCh,
					p.producerAckCh,
					producerDeferredConfirmations)
				return
			}

			for {
				checkProducerDeferredConfirmations(producerDeferredConfirmations)
				if len(producerDeferredConfirmations) != p.sessionConfig.ConcurrencyLimit {
					break
				}
				ctx, cancel := context.WithTimeout(context.Background(), produceIntervalTime)
				<-ctx.Done()
				cancel()
			}
		}
	}()
}

func (p *producer) notifyClose() {
	ch := make(chan os.Signal, 2)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-ch
		p.exitCh <- struct{}{}
	}()
}

func allowExecution(allowCh chan struct{}) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	select {
	case allowCh <- struct{}{}:
	case <-ctx.Done():
	}
}

func exitConcurrencyDaemon(
	doneCh chan struct{},
	allowCh chan struct{},
	producerAckCh chan *amqp.DeferredConfirmation,
	producerDeferredConfirmations map[uint64]*amqp.DeferredConfirmation) {
	for {
		checkProducerDeferredConfirmations(producerDeferredConfirmations)
		if len(producerDeferredConfirmations) == 0 {
			break
		}
		ctx, cancel := context.WithTimeout(context.Background(), produceIntervalTime)
		<-ctx.Done()
		cancel()
	}
	close(doneCh)
	close(allowCh)
	close(producerAckCh)
}

func checkProducerDeferredConfirmations(producerDeferredConfirmations map[uint64]*amqp.DeferredConfirmation) {
	for k, v := range producerDeferredConfirmations {
		if v.Acked() {
			delete(producerDeferredConfirmations, k)
		}
	}
}

func (p *producer) Publish(
	ctx context.Context,
	exchange string,
	routingKey string,
	msg *Message,
) error {
	if err := p.channel.Confirm(false); err != nil {
		return err
	}

	isAllow := false
	for {
		ctx, cancel := context.WithTimeout(context.Background(), produceIntervalTime)
		select {
		case <-p.exitCh:
			cancel()
			return errors.New("[rabbitmq-producer] interrupted")
		case <-p.allowCh:
			cancel()
			isAllow = true
		case <-ctx.Done():
			continue
		}
		if isAllow {
			break
		}
	}

	deferredConfirmation, err := p.channel.PublishWithDeferredConfirmWithContext(
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
	if err != nil {
		return err
	}

	select {
	case <-p.exitCh:
		return errors.New("[rabbitmq-producer] interrupted")
	case p.producerAckCh <- deferredConfirmation:
	}

	return nil
}
