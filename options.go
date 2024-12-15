package amqp091quickstart

type BinderOption func(*binder)
type ConsumerOption func(*consumer)
type ProducerOption func(*producer)

func WithBinderSessionConfig(sessionConfig *SessionConfig) BinderOption {
	return func(s *binder) {
		s.sessionConfig = sessionConfig
	}
}

func WithExchangeConfig(exchangeConfig *ExchangeConfig) BinderOption {
	return func(s *binder) {
		s.exchangeConfig = exchangeConfig
	}
}

func WithQueueConfig(queueConfig *QueueConfig) BinderOption {
	return func(s *binder) {
		s.queueConfig = queueConfig
	}
}

func WithQueueBindConfig(queueBindConfig *QueueBindConfig) BinderOption {
	return func(s *binder) {
		s.queueBindConfig = queueBindConfig
	}
}

func WithConsumerSessionConfig(sessionConfig *SessionConfig) ConsumerOption {
	return func(s *consumer) {
		s.sessionConfig = sessionConfig
	}
}

func WithConsumerConfig(consumerConfig *ConsumerConfig) ConsumerOption {
	return func(c *consumer) {
		c.consumerConfig = consumerConfig
	}
}

func WithProducerSessionConfig(sessionConfig *SessionConfig) ProducerOption {
	return func(s *producer) {
		s.sessionConfig = sessionConfig
	}
}

func WithProducerConfig(producerConfig *ProducerConfig) ProducerOption {
	return func(p *producer) {
		p.producerConfig = producerConfig
	}
}
