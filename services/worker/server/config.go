package server

import "github.com/gbh007/buttoners/core/config"

type DBConfig struct {
	Username, Password, Addr, DatabaseName string
}

type RabbitMQConfig struct {
	Username, Password, Addr, QueueName string
}

type Config struct {
	ServiceName         string
	NotificationService config.Service
	PrometheusAddress   string
	DB                  DBConfig
	RabbitMQ            RabbitMQConfig
	RunnerCount         int
}
