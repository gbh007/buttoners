package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/url"
	"time"

	"github.com/gbh007/buttoners/core/metrics"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

const contentTypeJSON = "application/json"

type Writer[T any] struct {
	tracer        trace.Tracer
	logger        *slog.Logger
	writerMetrics *metrics.QueueWriterMetrics

	conn *amqp.Connection
	ch   *amqp.Channel

	user, pass, addr, queueName string
}

func NewWriter[T any](
	logger *slog.Logger,
	user, pass, addr, queueName string,
	writerMetrics *metrics.QueueWriterMetrics,
) (*Writer[T], error) {
	w := &Writer[T]{
		logger:        logger,
		user:          user,
		pass:          pass,
		addr:          addr,
		queueName:     queueName,
		tracer:        newTracer(otel.GetTracerProvider()),
		writerMetrics: writerMetrics,
	}

	var err error

	u := url.URL{
		Scheme: "amqp",
		User:   url.UserPassword(w.user, w.pass),
		Host:   w.addr,
	}

	w.conn, err = amqp.Dial(u.String())
	if err != nil {
		return nil, fmt.Errorf("%w: dial: %w", ErrRabbitMQClient, err)
	}

	w.ch, err = w.conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("%w: create channel: %w", ErrRabbitMQClient, err)
	}

	// FIXME: конфигурировать не в коде
	_, err = w.ch.QueueDeclare(
		w.queueName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("%w: create queue: %w", ErrRabbitMQClient, err)
	}

	return w, nil
}

func (w *Writer[T]) Close() error {
	var errs []error

	if w.ch != nil {
		err := w.ch.Close()
		if err != nil {
			errs = append(errs, err)
		}
	}

	if w.conn != nil {
		err := w.conn.Close()
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

func (w *Writer[T]) Write(ctx context.Context, k string, v T) (returnedErr error) {
	startTime := time.Now()

	w.writerMetrics.IncActive(w.addr, w.queueName)
	defer w.writerMetrics.DecActive(w.addr, w.queueName)

	defer func() {
		status := metrics.ResultOK
		if returnedErr != nil {
			status = metrics.ResultError
		}

		w.writerMetrics.AddHandle(w.addr, w.queueName, status, time.Since(startTime))
	}()

	ctx, span := w.tracer.Start(ctx, "rabbitmq-write")
	defer span.End()

	requestLog := []any{
		slog.String("queue", w.queueName),
		slog.String("addr", w.addr),
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

		w.logger.InfoContext(
			ctx,
			"rabbit mq write",
			args...,
		)
	}()

	data, err := json.Marshal(v)
	if err != nil {
		return fmt.Errorf("%w: marshal: %w", ErrRabbitMQClient, err)
	}

	// Распространение трассировки
	carrier := propagation.MapCarrier(make(map[string]string, 3))
	otel.GetTextMapPropagator().Inject(ctx, carrier)

	msg := amqp.Publishing{
		MessageId:   k,
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

	err = w.ch.PublishWithContext(
		ctx,
		"",
		w.queueName,
		false,
		false,
		msg,
	)
	if err != nil {
		return fmt.Errorf("%w: publish: %w", ErrRabbitMQClient, err)
	}

	return nil
}
