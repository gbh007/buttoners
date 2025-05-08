package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/rabbitmq/amqp091-go"
	"go.opentelemetry.io/otel"
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
