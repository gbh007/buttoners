package server

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/gbh007/buttoners/core/clients/logclient"
	"github.com/gbh007/buttoners/core/logger"
	"github.com/gbh007/buttoners/core/metrics"
	"github.com/gofiber/contrib/otelfiber/v2"
	"github.com/gofiber/fiber/v2"
)

func (s *Server) Run(ctx context.Context) error {
	go metrics.Run(s.l, metrics.Config{Addr: s.cfg.PrometheusAddress})

	fb := fiber.New(fiber.Config{DisableStartupMessage: true})
	otelHandler := otelfiber.Middleware(
		otelfiber.WithoutMetrics(true),
	)

	fb.Post(logclient.ActivityPath, otelHandler, s.observabilityHandler, s.authHandler, s.Activity)

	go func() {
		<-ctx.Done()
		sCtx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		err := fb.ShutdownWithContext(sCtx)
		if err != nil {
			logger.LogWithMeta(s.l, ctx, slog.LevelWarn, "unsuccess shutdown web server", "error", err.Error())
		}
	}()

	wg := new(sync.WaitGroup)
	wg.Add(2)

	go func() {
		defer wg.Done()
		defer s.kafka.Close()

		err := s.kafka.Start(ctx)
		if err != nil {
			logger.LogWithMeta(s.l, ctx, slog.LevelWarn, "unsuccess handler result", "error", err.Error())
		}
	}()

	go func() {
		defer wg.Done()

		err := fb.Listen(s.cfg.SelfAddress)
		if err != nil {
			logger.LogWithMeta(s.l, ctx, slog.LevelWarn, "unsuccess server result", "error", err.Error())
		}
	}()

	wg.Wait()

	return nil
}
