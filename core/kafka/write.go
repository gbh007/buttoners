package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
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

	err = c.writer.WriteMessages(ctx, kafka.Message{
		Key:     []byte(k),
		Value:   data,
		Headers: fromMapCarrier(carrier),
	})
	if err != nil {
		registerWriteHandleTime(false, time.Since(startTime))

		return fmt.Errorf("%w: Write: %w", ErrKafkaClient, err)
	}

	registerWriteHandleTime(true, time.Since(startTime))

	return nil
}
