package config

import "strconv"

type Addr struct {
	Host string `envconfig:"default=localhost"`
	Port int64  `envconfig:"default=50051"`
}

func (a Addr) Full() string {
	return a.Host + ":" + strconv.FormatInt(a.Port, 10)
}

type Database struct {
	User string
	Pass string
	Addr string
	Name string
}

type Kafka struct {
	Addr          string `envconfig:"default=kafka:9092"`
	TaskTopic     string `envconfig:"default=gate"`
	LogTopic      string `envconfig:"default=log"`
	GroupID       string `envconfig:"optional"`
	NumPartitions int    `envconfig:"optional"`
}

type RabbitMQ struct {
	User  string
	Pass  string
	Addr  string `envconfig:"default=rabbitmq:5672"`
	Queue string `envconfig:"default=task"`
}

type Jaeger struct {
	URL string `envconfig:"default=http://jaeger:14268/api/traces"`
}
