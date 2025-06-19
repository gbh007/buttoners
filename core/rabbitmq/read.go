package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/rabbitmq/amqp091-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type Read[T any] func(ctx context.Context) (context.Context, *T, error)

func (c *Client[T]) StartRead(ctx context.Context) (chan Read[T], error) {
	if c.out != nil {
		return c.out, nil
	}

	if c.ch == nil {
		return nil, fmt.Errorf("%w: StartRead: %w", ErrRabbitMQClient, ErrChannelNotInitialized)
	}

	messages, err := c.ch.Consume(
		c.queue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("%w: StartRead: %w", ErrRabbitMQClient, err)
	}

	// Создаем не буферизированный канал
	c.out = make(chan Read[T])

	go func() {
		for {
			select {
			case <-ctx.Done():
				return

			case msg := <-messages:
				c.out <- c.handleMsg(msg)
			}
		}
	}()

	return c.out, nil
}

func (c *Client[T]) handleMsg(msg amqp091.Delivery) Read[T] {
	return func(ctx context.Context) (context.Context, *T, error) {
		// Распространение трассировки
		ctx = otel.GetTextMapPropagator().Extract(ctx, toMapCarrier(msg.Headers))

		ctx, span := c.tracer.Start(ctx, "rabbitmq-read")
		defer span.End()

		startTime := time.Now()

		requestLog := []any{
			slog.String("routing_key", msg.RoutingKey),
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
			ctx, "rabbit mq read",
			slog.String("trace_id", trace.SpanContextFromContext(ctx).TraceID().String()),
			slog.Group("request", requestLog...),
		)

		v := new(T)

		err := json.Unmarshal(msg.Body, &v)
		if err != nil {
			registerReadHandleTime(false, time.Since(startTime))

			return ctx, nil, fmt.Errorf("%w: Read: %w", ErrRabbitMQClient, err)
		}

		// Находиться после занесения в очередь, по причине того,
		// что чтение рассматриваем как процесс перемещения задачи в раннер.
		registerReadHandleTime(true, time.Since(startTime))

		return ctx, v, nil
	}
}
