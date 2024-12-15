package queue

import mq "github.com/randallchang/amqp091-go-quickstart"

type MQConfig struct {
	ExchangeType mq.ExchangeType
	Exchange     string
	RoutingKey   string
	Queue        string
	WithDLX      bool
}

func (c *MQConfig) DLXMQConfig() *MQConfig {
	return &MQConfig{
		Exchange:     c.Exchange,
		ExchangeType: c.ExchangeType,
		RoutingKey:   c.RoutingKey + ".dlx",
		Queue:        c.Queue + ".dlx",
		WithDLX:      false,
	}
}
