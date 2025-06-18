package kafka

import (
	"log/slog"

	"github.com/segmentio/kafka-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type Client struct {
	kafkaConn *kafka.Conn

	tracer trace.Tracer
	logger *slog.Logger

	reader *kafka.Reader
	writer *kafka.Writer

	topic         string
	groupID       string
	addr          string
	numPartitions int
}

func New(
	logger *slog.Logger, addr, topic, groupID string, numPartitions int) *Client {
	return &Client{
		logger:        logger,
		topic:         topic,
		groupID:       groupID,
		addr:          addr,
		numPartitions: numPartitions,
		tracer:        newTracer(otel.GetTracerProvider()),
	}
}
