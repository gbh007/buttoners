package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/segmentio/kafka-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

func (c *Client) Write(ctx context.Context, k string, v any) error {
	ctx, span := c.tracer.Start(ctx, "kafka-write")
	defer span.End()

	startTime := time.Now()

	if c.writer == nil {
		registerWriteHandleTime(false, time.Since(startTime))

		return fmt.Errorf("%w: Write: %w", ErrKafkaClient, ErrConnectionNotInitialized)
	}

	data, err := json.Marshal(v)
	if err != nil {
		registerWriteHandleTime(false, time.Since(startTime))

		return fmt.Errorf("%w: Write: %w", ErrKafkaClient, err)
	}

	// Распространение трассировки
	carrier := propagation.MapCarrier(make(map[string]string, 3))
	otel.GetTextMapPropagator().Inject(ctx, carrier)

	msg := kafka.Message{
		Key:     []byte(k),
		Value:   data,
		Headers: fromMapCarrier(carrier),
	}

	err = c.writer.WriteMessages(ctx, msg)
	if err != nil {
		registerWriteHandleTime(false, time.Since(startTime))

		return fmt.Errorf("%w: Write: %w", ErrKafkaClient, err)
	}

	registerWriteHandleTime(true, time.Since(startTime))

	requestLog := []any{
		slog.String("message_key", string(msg.Key)),
		slog.String("topic", c.topic),
	}

	if len(msg.Headers) > 0 {
		headers := make(map[string]string)

		for _, v := range msg.Headers {
			old := headers[v.Key]
			if old != "" {
				old += ";"
			}
			old += string(v.Value)
			headers[v.Key] = old
		}

		requestLog = append(
			requestLog,
			slog.Any("headers", headers),
		)
	}

	if len(msg.Value) > 0 {
		requestLog = append(requestLog, slog.String("body", string(msg.Value)))
	}

	c.logger.InfoContext(
		ctx, "kafka produce",
		slog.String("trace_id", trace.SpanContextFromContext(ctx).TraceID().String()),
		slog.Group("request", requestLog...),
	)

	return nil
}
