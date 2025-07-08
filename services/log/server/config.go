package server

type KafkaConfig struct {
	Topic   string
	GroupID string
	Addr    string
}

type DBConfig struct {
	Username, Password, Addr, DatabaseName string
}

type Config struct {
	ServiceName       string
	SelfAddress       string
	SelfToken         string
	PrometheusAddress string
	Kafka             KafkaConfig
	DB                DBConfig
}
