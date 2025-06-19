package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
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

	msg := amqp.Publishing{
		ContentType: contentTypeJSON,
		Body:        data,
		Headers:     fromMapCarrier(carrier),
	}

	err = c.ch.PublishWithContext(
		ctx,
		"",
		c.queue.Name,
		false,
		false,
		msg,
	)
	if err != nil {
		registerWriteHandleTime(false, time.Since(startTime))

		return fmt.Errorf("%w: Write: %w", ErrRabbitMQClient, err)
	}

	requestLog := []any{
		slog.String("msg_id", msg.MessageId),
		slog.String("queue", c.queueName),
	}

	if len(msg.Headers) > 0 {
		headers := make(map[string]string)

		for k, v := range msg.Headers {
			switch typedV := v.(type) {
			case string:
				headers[k] = typedV
			default:
				headers[k] = fmt.Sprint(typedV)
			}
		}

		requestLog = append(
			requestLog,
			slog.Any("headers", headers),
		)
	}

	if len(msg.Body) > 0 {
		requestLog = append(requestLog, slog.String("body", string(msg.Body)))
	}

	c.logger.InfoContext(
		ctx, "rabbit mq write",
		slog.String("trace_id", trace.SpanContextFromContext(ctx).TraceID().String()),
		slog.Group("request", requestLog...),
	)

	registerWriteHandleTime(true, time.Since(startTime))

	return nil
}
