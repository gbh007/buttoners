package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/url"
	"time"

	"github.com/gbh007/buttoners/core/metrics"
	"github.com/rabbitmq/amqp091-go"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type Reader[T any] struct {
	tracer        trace.Tracer
	logger        *slog.Logger
	readerMetrics *metrics.QueueReaderMetrics

	handler func(context.Context, string, T) error

	conn *amqp.Connection
	ch   *amqp.Channel

	user, pass, addr, queueName string
}

func NewReader[T any](
	logger *slog.Logger,
	user, pass, addr, queueName string,
	readerMetrics *metrics.QueueReaderMetrics,
	handler func(context.Context, string, T) error,
) (*Reader[T], error) {
	r := &Reader[T]{
		logger:        logger,
		user:          user,
		pass:          pass,
		addr:          addr,
		queueName:     queueName,
		tracer:        newTracer(otel.GetTracerProvider()),
		readerMetrics: readerMetrics,
		handler:       handler,
	}

	var err error

	u := url.URL{
		Scheme: "amqp",
		User:   url.UserPassword(r.user, r.pass),
		Host:   r.addr,
	}

	r.conn, err = amqp.Dial(u.String())
	if err != nil {
		return nil, fmt.Errorf("%w: dial: %w", ErrRabbitMQClient, err)
	}

	r.ch, err = r.conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("%w: create channel: %w", ErrRabbitMQClient, err)
	}

	return r, nil
}

func (r *Reader[T]) Close() error {
	var errs []error

	if r.ch != nil {
		err := r.ch.Close()
		if err != nil {
			errs = append(errs, err)
		}
	}

	if r.conn != nil {
		err := r.conn.Close()
		if err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) == 0 {
		return nil
	}

	err := fmt.Errorf("%w: close", ErrRabbitMQClient)

	for _, e := range errs {
		err = fmt.Errorf("%w: %w", err, e)
	}

	return err
}

func (c *Reader[T]) Start(ctx context.Context) error {
	messages, err := c.ch.Consume(
		c.queueName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("%w: start read: %w", ErrRabbitMQClient, err)
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		case msg := <-messages:
			_ = c.handleMsg(ctx, msg)
		}
	}
}

func (c *Reader[T]) handleMsg(ctx context.Context, msg amqp091.Delivery) (returnedErr error) {
	startTime := time.Now()

	c.readerMetrics.IncActive(c.addr, c.queueName, "")
	defer c.readerMetrics.DecActive(c.addr, c.queueName, "")

	defer func() {
		status := metrics.ResultOK
		if returnedErr != nil {
			status = metrics.ResultError
		}

		c.readerMetrics.AddHandle(c.addr, c.queueName, "", status, time.Since(startTime))
	}()

	// Распространение трассировки
	ctx = otel.GetTextMapPropagator().Extract(ctx, toMapCarrier(msg.Headers))

	ctx, span := c.tracer.Start(ctx, "rabbitmq-read")
	defer span.End()

	requestLog := []any{
		slog.String("message_key", msg.MessageId),
		slog.String("queue", c.queueName),
		slog.String("addr", c.addr),
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
			"rabbit mq read",
			args...,
		)
	}()

	var v T

	err := json.Unmarshal(msg.Body, &v)
	if err != nil {
		return fmt.Errorf("%w: unmarshal: %w", ErrRabbitMQClient, err)
	}

	err = c.handler(ctx, msg.MessageId, v)
	if err != nil {
		return fmt.Errorf("%w: handle: %w", ErrRabbitMQClient, err)
	}

	return nil
}
