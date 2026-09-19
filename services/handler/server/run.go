package server

import (
	"context"

	"github.com/gbh007/buttoners/core/metrics"
)

func (s *Server) Run(ctx context.Context) error {
	go metrics.Run(ctx, s.logger, metrics.Config{Addr: s.cfg.MetricAddr})

	defer s.Close(ctx)

	return s.kafkaClient.Start(ctx)
}
