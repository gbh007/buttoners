package server

import "github.com/gbh007/buttoners/core/config"

type KafkaConfig struct {
	TaskTopic string
	LogTopic  string
	Addr      string
}

type Config struct {
	SelfAddress         string
	AuthService         config.Service
	LogService          config.Service
	NotificationService config.Service
	RedisAddress        string
	PrometheusAddress   string
	Kafka               KafkaConfig
}
