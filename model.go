package amqp091quickstart

import (
	"crypto/tls"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Table map[string]interface{}

type ExchangeType string

const (
	Direct         ExchangeType = "direct"
	Fanout         ExchangeType = "fanout"
	Topic          ExchangeType = "topic"
	Headers        ExchangeType = "headers"
	ConsistentHash ExchangeType = "x-consistent-hash"
	DelayedMessage ExchangeType = "x-delayed-message"
)

type Message struct {
	CorrelationID string
	Headers       map[string]interface{}
	Body          []byte
}

type SessionConfig struct {
	ConcurrencyLimit int
	SASL             []amqp.Authentication
	Vhost            string
	ChannelMax       int
	FrameSize        int
	Heartbeat        time.Duration
	TLSClientConfig  *tls.Config
	Properties       Table
	Locale           string
}

type ExchangeConfig struct {
	Durable    bool
	AutoDelete bool
	Internal   bool
	NoWait     bool
	Args       Table
}

type QueueConfig struct {
	Durable    bool
	AutoDelete bool
	Exclusive  bool
	NoWait     bool
	Args       Table
}

type QueueBindConfig struct {
	NoWait bool
	Args   Table
}

type ProducerConfig struct {
	ContentType     string
	ContentEncoding string
	DeliveryMode    uint8
	Priority        uint8
	ReplyTo         string
	Expiration      string
	MessageId       string
	Timestamp       time.Time
	Type            string
	UserId          string
	AppId           string
}

type ConsumerConfig struct {
	Consumer  string
	AutoAck   bool
	Exclusive bool
	NoLocal   bool
	NoWait    bool
	Args      Table
	Requeue   bool
}
