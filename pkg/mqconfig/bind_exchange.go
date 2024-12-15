package mqconfig

import (
	mq "github.com/randallchang/amqp091-go-quickstart"
	"github.com/randallchang/amqp091-go-quickstart/pkg/queue"
)

func QuickBindExchange(url string, mqConfigs ...*queue.MQConfig) error {
	for _, mqConfig := range mqConfigs {
		binder := mq.NewBinder(
			url,
			mq.WithExchangeConfig(&mq.ExchangeConfig{
				Durable: true,
				Args:    mq.Table{"x-delayed-type": "direct"},
			}),
		)

		err := binder.Connect()
		if err != nil {
			return err
		}
		defer binder.Close()

		err = binder.Bind(
			mqConfig.Exchange,
			mqConfig.ExchangeType,
			mqConfig.Queue,
			mqConfig.RoutingKey,
		)
		if err != nil {
			return err
		}
		if mqConfig.WithDLX {
			dlxMQConfig := mqConfig.DLXMQConfig()
			err = binder.Bind(
				dlxMQConfig.Exchange,
				dlxMQConfig.ExchangeType,
				dlxMQConfig.Queue,
				dlxMQConfig.RoutingKey,
			)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
