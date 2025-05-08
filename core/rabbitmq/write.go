package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

const contentTypeJSON = "application/json"

func (c *Client[T]) Write(ctx context.Context, v T) error {
	ctx, span := c.tracer.Start(ctx, "rabbitmq-write")
	defer span.End()

	startTime := time.Now()
	if c.ch == nil {
		registerWriteHandleTime(false, time.Since(startTime))

		return fmt.Errorf("%w: Write: %w", ErrRabbitMQClient, ErrChannelNotInitialized)
	}

	data, err := json.Marshal(v)
	if err != nil {
		registerWriteHandleTime(false, time.Since(startTime))

		return fmt.Errorf("%w: Write: %w", ErrRabbitMQClient, err)
	}

	// Распространение трассировки
	carrier := propagation.MapCarrier(make(map[string]string, 3))
	otel.GetTextMapPropagator().Inject(ctx, carrier)

	err = c.ch.PublishWithContext(ctx,
		"",
		c.queue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: contentTypeJSON,
			Body:        data,
			Headers:     fromMapCarrier(carrier),
		})
	if err != nil {
		registerWriteHandleTime(false, time.Since(startTime))

		return fmt.Errorf("%w: Write: %w", ErrRabbitMQClient, err)
	}

	registerWriteHandleTime(true, time.Since(startTime))

	return nil
}
