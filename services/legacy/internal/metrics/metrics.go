package metrics

import (
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	namespace = "buttoners"
	subsystem = "example_baton_nagimator"
)

var (
	httpDuration = promauto.NewSummaryVec(prometheus.SummaryOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "http_seconds",
	}, []string{"path", "method", "code"})
	buttonLoversTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "button_lovers_total",
	}, []string{"name"})
)

func RecordHTTPRequest(path, method string, code int, d time.Duration) {
	httpDuration.WithLabelValues(path, method, strconv.Itoa(code)).Observe(d.Seconds())
}

func RecordButtonPress(name string) {
	buttonLoversTotal.WithLabelValues(name).Inc()
}
