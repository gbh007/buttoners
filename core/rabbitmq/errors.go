package rabbitmq

import "errors"

var (
	ErrRabbitMQClient = errors.New("RabbitMQ client")

	ErrChannelNotInitialized = errors.New("channel not initialized")
)
