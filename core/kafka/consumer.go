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
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type Consumer[T any] struct {
	tracer        trace.Tracer
	logger        *slog.Logger
	readerMetrics *metrics.QueueReaderMetrics

	reader  *kafka.Reader
	handler func(context.Context, string, T) error

	topic   string
	groupID string
	addr    string
}

func NewConsumer[T any](
	logger *slog.Logger,
	addr, topic, groupID string,
	readerMetrics *metrics.QueueReaderMetrics,
	handler func(context.Context, string, T) error,
) *Consumer[T] {
	c := &Consumer[T]{
		logger:        logger,
		topic:         topic,
		groupID:       groupID,
		addr:          addr,
		tracer:        newTracer(otel.GetTracerProvider()),
		readerMetrics: readerMetrics,
		handler:       handler,
	}

	c.reader = kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{c.addr},
		GroupID: c.groupID,
		Topic:   c.topic,
	})

	return c
}

func (c *Consumer[T]) Close() error {
	err := c.reader.Close()
	if err != nil {
		return err
	}

	return nil
}

func (c *Consumer[T]) Start(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		_ = c.read(ctx)
	}
}

func (c *Consumer[T]) read(ctx context.Context) (returnedErr error) {
	realStartTime := time.Now()

	var v T

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
			"kafka consume",
			args...,
		)
	}()

	c.readerMetrics.IncActive(c.addr, c.topic, c.groupID)
	defer c.readerMetrics.DecActive(c.addr, c.topic, c.groupID)

	msg, err := c.reader.ReadMessage(ctx)
	if err != nil {
		return fmt.Errorf("%w: read: %w", ErrKafkaClient, err)
	}

	startTime := time.Now()

	defer func() {
		status := metrics.ResultOK
		if returnedErr != nil {
			status = metrics.ResultError
		}

		c.readerMetrics.AddHandle(c.addr, c.topic, c.groupID, status, time.Since(startTime))
	}()

	// Распространение трассировки
	ctx = otel.GetTextMapPropagator().Extract(ctx, toMapCarrier(msg.Headers))

	ctx, span := c.tracer.Start(ctx, "kafka-read", trace.WithTimestamp(startTime))
	defer span.End()

	span.SetAttributes(attribute.String("wait_time", startTime.Sub(realStartTime).String()))

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

	err = json.Unmarshal(msg.Value, &v)
	if err != nil {
		return fmt.Errorf("%w: unmarshal: %w", ErrKafkaClient, err)
	}

	err = c.handler(ctx, string(msg.Key), v)
	if err != nil {
		return fmt.Errorf("%w: handle: %w", ErrKafkaClient, err)
	}

	return nil
}
