package server

import (
	"context"
	"sync"

	"github.com/gbh007/buttoners/core/metrics"
)

func (s *Server) Run(ctx context.Context) error {
	go metrics.Run(s.logger, metrics.Config{Addr: s.cfg.PrometheusAddress})

	runnerCtx, runnerCnl := context.WithCancel(context.TODO())
	runnerWg := new(sync.WaitGroup)

	for i := 0; i < s.cfg.RunnerCount; i++ {
		runnerWg.Add(1)

		go func() {
			defer runnerWg.Done()

			rnErr := s.rabbitClient.Start(runnerCtx)
			if rnErr != nil {
				s.logger.Error("runner end unsuccessful", "error", rnErr.Error())
			}
		}()
	}

	<-ctx.Done()

	runnerCnl()

	runnerWg.Wait()

	return nil
}
