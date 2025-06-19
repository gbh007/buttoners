package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func (c *Client) Read(ctx context.Context, v any) (context.Context, string, error) {
	realStartTime := time.Now()

	if c.reader == nil {
		registerReadHandleTime(false, time.Since(realStartTime))

		return ctx, "", fmt.Errorf("%w: Read: %w", ErrKafkaClient, ErrConnectionNotInitialized)
	}

	msg, err := c.reader.ReadMessage(ctx)
	if err != nil {
		registerReadHandleTime(false, time.Since(realStartTime))

		return ctx, "", fmt.Errorf("%w: Read: %w", ErrKafkaClient, err)
	}

	startTime := time.Now()

	// Распространение трассировки
	ctx = otel.GetTextMapPropagator().Extract(ctx, toMapCarrier(msg.Headers))

	ctx, span := c.tracer.Start(ctx, "kafka-read", trace.WithTimestamp(startTime))
	defer span.End()

	span.SetAttributes(attribute.String("wait_time", startTime.Sub(realStartTime).String()))

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
		ctx, "kafka consume",
		slog.String("trace_id", trace.SpanContextFromContext(ctx).TraceID().String()),
		slog.Group("request", requestLog...),
	)

	err = json.Unmarshal(msg.Value, &v)
	if err != nil {
		registerReadHandleTime(false, time.Since(startTime))

		return ctx, "", fmt.Errorf("%w: Read: %w", ErrKafkaClient, err)
	}

	registerReadHandleTime(true, time.Since(startTime))

	return ctx, string(msg.Key), nil
}
