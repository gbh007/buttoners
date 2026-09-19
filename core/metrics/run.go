package metrics

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/client_golang/prometheus/push"
)

type Config struct {
	Addr     string
	Job      string
	PushMode bool
}

var InstanceName = "unknown"

func Run(parentCtx context.Context, l *slog.Logger, cfg Config) {
	l.Info("metrics start")

	host, _ := os.Hostname()

	if cfg.Job == "" {
		cfg.Job = "service"
	}

	if cfg.PushMode {

		pusher := push.New(cfg.Addr, cfg.Job).
			Collector(DefaultRegistry).
			Grouping(hostLabelName, host).
			Grouping(instanceLabelName, InstanceName)

		for range time.NewTicker(time.Second).C {
			if err := pusher.Push(); err != nil {
				l.Error("push metrics error", "error", err.Error())
			}
		}
	} else {
		mux := http.NewServeMux()

		mux.Handle("/metrics", promhttp.HandlerFor(DefaultRegistry, promhttp.HandlerOpts{}))

		server := &http.Server{ //nolint:gosec // будет исправлено позднее
			Handler: mux,
			Addr:    cfg.Addr,
		}

		go func() {
			err := server.ListenAndServe()
			if err != nil && !errors.Is(err, http.ErrServerClosed) {
				l.Error("start metrics server", "error", err.Error())
			}
		}()

		go func() {
			<-parentCtx.Done()

			shutdownCtx, cancel := context.WithTimeout(context.WithoutCancel(parentCtx), time.Second*10)
			defer cancel()

			err := server.Shutdown(shutdownCtx)
			if err != nil {
				l.Error("stop metrics server", "error", err.Error())
			}
		}()
	}
}
