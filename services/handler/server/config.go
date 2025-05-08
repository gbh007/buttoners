package server

type KafkaConfig struct {
	Topic         string
	GroupID       string
	Addr          string
	NumPartitions int
}

type RabbitMQConfig struct {
	Username, Password, Addr, QueueName string
}

type Config struct {
	ServiceName       string
	PrometheusAddress string
	Kafka             KafkaConfig
	RabbitMQ          RabbitMQConfig
}
