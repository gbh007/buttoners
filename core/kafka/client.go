package kafka

import (
	"github.com/segmentio/kafka-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type Client struct {
	kafkaConn *kafka.Conn

	tracer trace.Tracer

	reader *kafka.Reader
	writer *kafka.Writer

	topic         string
	groupID       string
	addr          string
	numPartitions int
}

func New(addr, topic, groupID string, numPartitions int) *Client {
	return &Client{
		topic:         topic,
		groupID:       groupID,
		addr:          addr,
		numPartitions: numPartitions,
		tracer:        newTracer(otel.GetTracerProvider()),
	}
}
