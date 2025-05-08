package kafka

import "errors"

var (
	ErrKafkaClient = errors.New("kafka client")

	ErrFailToCreateTopic        = errors.New("fail to create topic")
	ErrConnectionNotInitialized = errors.New("connection not initialized")
)
