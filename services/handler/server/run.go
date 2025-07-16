package server

import (
	"context"

	"github.com/gbh007/buttoners/core/metrics"
)

func (s *Server) Run(ctx context.Context) error {
	go metrics.Run(s.logger, metrics.Config{Addr: s.cfg.PrometheusAddress})

	defer s.Close(ctx)

	return s.kafkaClient.Start(ctx)
}
