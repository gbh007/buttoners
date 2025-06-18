package server

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gbh007/buttoners/core/clients/logclient"
	"github.com/gbh007/buttoners/core/kafka"
	"github.com/gbh007/buttoners/core/metrics"
	"github.com/gbh007/buttoners/core/observability"
	"github.com/gbh007/buttoners/services/log/internal/storage"
	"github.com/gofiber/contrib/otelfiber/v2"
	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel"
)

func Run(ctx context.Context, cfg Config) error {
	go metrics.Run(metrics.Config{Addr: cfg.PrometheusAddress})

	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))
	logger = logger.With("service_name", metrics.InstanceName)

	httpServerMetrics, err := metrics.NewHTTPServerMetrics(metrics.DefaultRegistry, metrics.DefaultTimeBuckets)
	if err != nil {
		return err
	}

	db, err := storage.Init(ctx, cfg.DB.Username, cfg.DB.Password, cfg.DB.Addr, cfg.DB.DatabaseName)
	if err != nil {
		return err
	}

	kafkaClient := kafka.New(logger, cfg.Kafka.Addr, cfg.Kafka.Topic, cfg.Kafka.GroupID, cfg.Kafka.NumPartitions)

	err = kafkaClient.Connect(cfg.Kafka.NumPartitions > 0)
	if err != nil {
		return err
	}

	defer kafkaClient.Close()

	handler := &handler{
		kafka:  kafkaClient,
		db:     db,
		tracer: otel.GetTracerProvider().Tracer(cfg.ServiceName),
	}

	server := &pbServer{
		db: db,
	}

	fb := fiber.New()
	otelHandler := otelfiber.Middleware(
		otelfiber.WithoutMetrics(true),
	)
	observabilityHandler := func(ctx *fiber.Ctx) error {
		tStart := time.Now()

		httpServerMetrics.IncActive(string(ctx.Request().Host()), string(ctx.Request().URI().Path()), string(ctx.Request().Header.Method()))
		defer httpServerMetrics.DecActive(string(ctx.Request().Host()), string(ctx.Request().URI().Path()), string(ctx.Request().Header.Method()))

		defer func() {
			httpServerMetrics.AddHandle(string(ctx.Request().Host()), string(ctx.Request().URI().Path()), string(ctx.Request().Header.Method()), ctx.Response().StatusCode(), time.Since(tStart))
		}()

		defer observability.LogFastHTTPData(ctx.UserContext(), logger, "log server request", ctx.Request(), ctx.Response())

		return ctx.Next()
	}
	authHandler := func(ctx *fiber.Ctx) error {
		token := string(ctx.Request().Header.Peek("Authorization"))

		if token == "" {
			ctx.Set(fiber.HeaderContentType, logclient.ContentType)

			return ctx.
				Status(http.StatusUnauthorized).
				JSON(logclient.ErrorResponse{
					Code:    "unauthorized",
					Details: "empty token",
				})
		}

		if token != cfg.SelfToken {
			ctx.Set(fiber.HeaderContentType, logclient.ContentType)

			return ctx.
				Status(http.StatusForbidden).
				JSON(logclient.ErrorResponse{
					Code:    "forbidden",
					Details: "invalid token",
				})
		}

		return ctx.Next()
	}

	fb.Post(logclient.ActivityPath, otelHandler, observabilityHandler, authHandler, server.Activity)

	go func() {
		<-ctx.Done()
		sCtx, _ := context.WithTimeout(context.Background(), time.Second*10)
		fb.ShutdownWithContext(sCtx)
	}()

	wg := new(sync.WaitGroup)
	wg.Add(2)

	go func() {
		defer wg.Done()

		err := handler.Run(ctx)
		if err != nil {
			log.Println(err)
		}
	}()

	go func() {
		defer wg.Done()

		err := fb.Listen(cfg.SelfAddress)
		if err != nil {
			log.Println(err)
		}
	}()

	wg.Wait()

	return nil
}
