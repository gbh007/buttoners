package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/gbh007/buttoners/core/metrics"
	"github.com/segmentio/kafka-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

func (c *Client) Write(ctx context.Context, k string, v any) (returnedErr error) {
	startTime := time.Now()

	c.writerMetrics.IncActive(c.addr, c.topic)
	defer c.writerMetrics.DecActive(c.addr, c.topic)

	ctx, span := c.tracer.Start(ctx, "kafka-write")
	defer span.End()

	requestLog := []any{
		slog.String("queue", c.topic),
		slog.String("addr", c.addr),
	}

	defer func() {
		args := []any{
			slog.Bool("success", returnedErr == nil),
			slog.String("trace_id", trace.SpanContextFromContext(ctx).TraceID().String()),
			slog.Group("request", requestLog...),
		}

		if returnedErr != nil {
			args = append(args, slog.String("error", returnedErr.Error()))
		}

		c.logger.InfoContext(
			ctx,
			"kafka produce",
			args...,
		)
	}()

	if c.writer == nil {
		return fmt.Errorf("%w: Write: %w", ErrKafkaClient, ErrConnectionNotInitialized)
	}

	defer func() {
		status := metrics.ResultOK
		if returnedErr != nil {
			status = metrics.ResultError
		}

		c.writerMetrics.AddHandle(c.addr, c.topic, status, time.Since(startTime))
	}()

	data, err := json.Marshal(v)
	if err != nil {
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

	requestLog = append(requestLog, slog.String("message_key", string(msg.Key)))

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

	err = c.writer.WriteMessages(ctx, msg)
	if err != nil {
		return fmt.Errorf("%w: Write: %w", ErrKafkaClient, err)
	}

	return nil
}
