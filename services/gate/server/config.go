package server

import "github.com/gbh007/buttoners/core/config"

type KafkaConfig struct {
	TaskTopic     string
	LogTopic      string
	GroupID       string
	Addr          string
	NumPartitions int
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
