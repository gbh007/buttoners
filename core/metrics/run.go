package metrics

import (
	"log/slog"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus/push"
)

type Config struct {
	Addr string
	Job  string
}

var InstanceName = "unknown"

func Run(l *slog.Logger, cfg Config) {
	l.Info("metrics start")

	host, _ := os.Hostname()

	if cfg.Job == "" {
		cfg.Job = "service"
	}

	pusher := push.New(cfg.Addr, cfg.Job).
		Collector(DefaultRegistry).
		Grouping(hostLabelName, host).
		Grouping(instanceLabelName, InstanceName)

	for range time.NewTicker(time.Second).C {
		if err := pusher.Push(); err != nil {
			l.Error("push metrics error", "error", err.Error())
		}
	}
}
