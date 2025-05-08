package server

type KafkaConfig struct {
	TaskTopic     string
	LogTopic      string
	GroupID       string
	Addr          string
	NumPartitions int
}

type Config struct {
	SelfAddress         string
	AuthAddress         string
	LogAddress          string
	NotificationAddress string
	RedisAddress        string
	PrometheusAddress   string
	Kafka               KafkaConfig
}
