package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/gbh007/buttoners/core/metrics"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

const contentTypeJSON = "application/json"

func (c *Client[T]) Write(ctx context.Context, v T) (returnedErr error) {
	startTime := time.Now()

	c.writerMetrics.IncActive(c.addr, c.queueName)
	defer c.writerMetrics.DecActive(c.addr, c.queueName)

	defer func() {
		status := metrics.ResultOK
		if returnedErr != nil {
			status = metrics.ResultError
		}

		c.writerMetrics.AddHandle(c.addr, c.queueName, status, time.Since(startTime))
	}()

	ctx, span := c.tracer.Start(ctx, "rabbitmq-write")
	defer span.End()

	requestLog := []any{
		slog.String("queue", c.queueName),
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
			"rabbit mq write",
			args...,
		)
	}()

	if c.ch == nil {
		return fmt.Errorf("%w: Write: %w", ErrRabbitMQClient, ErrChannelNotInitialized)
	}

	data, err := json.Marshal(v)
	if err != nil {
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

	requestLog = append(requestLog, slog.String("message_key", msg.MessageId))

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

	err = c.ch.PublishWithContext(
		ctx,
		"",
		c.queue.Name,
		false,
		false,
		msg,
	)
	if err != nil {
		return fmt.Errorf("%w: Write: %w", ErrRabbitMQClient, err)
	}

	return nil
}
