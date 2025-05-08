package server

type DBConfig struct {
	Username, Password, Addr, DatabaseName string
}

type RabbitMQConfig struct {
	Username, Password, Addr, QueueName string
}

type Config struct {
	ServiceName         string
	NotificationAddress string
	PrometheusAddress   string
	DB                  DBConfig
	RabbitMQ            RabbitMQConfig
	RunnerCount         int
}
