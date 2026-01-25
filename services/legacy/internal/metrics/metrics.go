package metrics

import (
	"strconv"
	"time"

	"github.com/gbh007/buttoners/core/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	namespace = "buttoners"
	subsystem = "example_baton_nagimator" // прим. оставлено специально
)

var (
	httpDuration = prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "http_seconds",
	}, []string{"path", "method", "code"})
	buttonLoversTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "button_lovers_total",
	}, []string{"name"})
)

func init() {
	metrics.DefaultRegistry.MustRegister(
		httpDuration,
		buttonLoversTotal,
	)
}

// FIXME: заюзать
func RecordHTTPRequest(path, method string, code int, d time.Duration) {
	httpDuration.WithLabelValues(path, method, strconv.Itoa(code)).Observe(d.Seconds())
}

func RecordButtonPress(name string) {
	buttonLoversTotal.WithLabelValues(name).Inc()
}
