package rabbitmq

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type Client[T any] struct {
	tracer trace.Tracer

	conn  *amqp.Connection
	ch    *amqp.Channel
	queue amqp.Queue

	out chan Read[T]

	user, pass, addr, queueName string
}

func New[T any](user, pass, addr, queueName string) *Client[T] {
	return &Client[T]{
		user:      user,
		pass:      pass,
		addr:      addr,
		queueName: queueName,
		tracer:    newTracer(otel.GetTracerProvider()),
	}
}
