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

type Producer[T any] struct {
	tracer        trace.Tracer
	logger        *slog.Logger
	writerMetrics *metrics.QueueWriterMetrics

	writer *kafka.Writer

	topic string
	addr  string
}

func NewProducer[T any](
	logger *slog.Logger,
	addr, topic string,
	writerMetrics *metrics.QueueWriterMetrics,
) *Producer[T] {
	p := &Producer[T]{
		logger:        logger,
		topic:         topic,
		addr:          addr,
		tracer:        newTracer(otel.GetTracerProvider()),
		writerMetrics: writerMetrics,
	}

	p.writer = &kafka.Writer{
		Addr:            kafka.TCP(p.addr),
		Topic:           p.topic,
		Balancer:        kafka.Murmur2Balancer{},
		WriteBackoffMin: 10 * time.Millisecond,
		WriteBackoffMax: 50 * time.Millisecond,
		BatchTimeout:    100 * time.Millisecond,
	}

	return p
}

func (p *Producer[T]) Close() error {
	err := p.writer.Close()
	if err != nil {
		return err
	}

	return nil
}

func (p *Producer[T]) Write(ctx context.Context, k string, v T) (returnedErr error) {
	startTime := time.Now()

	p.writerMetrics.IncActive(p.addr, p.topic)
	defer p.writerMetrics.DecActive(p.addr, p.topic)

	ctx, span := p.tracer.Start(ctx, "kafka-write")
	defer span.End()

	requestLog := []any{
		slog.String("queue", p.topic),
		slog.String("addr", p.addr),
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

		p.logger.InfoContext(
			ctx,
			"kafka produce",
			args...,
		)
	}()

	if p.writer == nil {
		return fmt.Errorf("%w: Write: %w", ErrKafkaClient, ErrConnectionNotInitialized)
	}

	defer func() {
		status := metrics.ResultOK
		if returnedErr != nil {
			status = metrics.ResultError
		}

		p.writerMetrics.AddHandle(p.addr, p.topic, status, time.Since(startTime))
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

	err = p.writer.WriteMessages(ctx, msg)
	if err != nil {
		return fmt.Errorf("%w: Write: %w", ErrKafkaClient, err)
	}

	return nil
}
