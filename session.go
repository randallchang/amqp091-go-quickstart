package amqp091quickstart

import (
	"context"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	defaultSessionConfig = &SessionConfig{
		Heartbeat: 10 * time.Second,
		Locale:    "en_US",
	}
	reconnectIntervalTime time.Duration = 5 * time.Second
)

type session struct {
	sessionConfig *SessionConfig

	url                  string
	connection           *amqp.Connection
	channel              *amqp.Channel
	connectionNotifyChan chan *amqp.Error
	channelNotifyChan    chan *amqp.Error
	Exit                 chan struct{}
	Reconnect            chan struct{}
}

func NewSession(url string) *session {
	return &session{
		url:           url,
		sessionConfig: defaultSessionConfig,
		Exit:          make(chan struct{}),
		Reconnect:     make(chan struct{}),
	}
}

func (s *session) Connect() error {
	err := s.connect()
	if err != nil {
		return err
	}
	go s.handleReconnect()

	return nil
}

func (s *session) connect() error {
	newConnection, err := amqp.DialConfig(
		s.url,
		amqp.Config{
			SASL:            s.sessionConfig.SASL,
			Vhost:           s.sessionConfig.Vhost,
			ChannelMax:      s.sessionConfig.ChannelMax,
			FrameSize:       s.sessionConfig.FrameSize,
			Heartbeat:       s.sessionConfig.Heartbeat,
			TLSClientConfig: s.sessionConfig.TLSClientConfig,
			Properties:      amqp.Table(s.sessionConfig.Properties),
			Locale:          s.sessionConfig.Locale,
		})
	if err != nil {
		return err
	}
	newChannel, err := newConnection.Channel()
	if err != nil {
		return err
	}
	s.connection = newConnection
	s.channel = newChannel
	s.connectionNotifyChan = s.connection.NotifyClose(make(chan *amqp.Error))
	s.channelNotifyChan = s.channel.NotifyClose(make(chan *amqp.Error))

	go s.handleReconnect()

	return nil
}

func (s *session) Close() error {
	close(s.Exit)
	s.Exit = make(chan struct{})
	if err := s.channel.Close(); err != nil {
		return err
	}
	if err := s.connection.Close(); err != nil {
		return err
	}
	return nil
}

func (s *session) handleReconnect() {
	for {
		select {
		case <-s.channelNotifyChan:
		case <-s.connectionNotifyChan:
		case <-s.Exit:
			return
		}

		if !s.connection.IsClosed() {
			_ = s.channel.Close()
			_ = s.connection.Close()
		}

		// in case of channel block
		for range s.channelNotifyChan {
		}
		for range s.connectionNotifyChan {
		}

	connect:
		for {
			select {
			case <-s.Exit:
				return
			default:
				err := s.connect()
				if err != nil {
					ctx, cancel := context.WithTimeout(context.Background(), reconnectIntervalTime)
					<-ctx.Done()
					cancel()
					continue
				}

				break connect
			}
		}
		close(s.Reconnect)
		s.Reconnect = make(chan struct{})
	}
}
