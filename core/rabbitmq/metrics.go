package rabbitmq

import (
	"time"

	"github.com/gbh007/buttoners/core/metrics"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	resultLabelName = "result"
	subsystemName   = "rabbitmq"
)

var (
	writeHandleTime = promauto.With(metrics.DefaultRegistry).NewHistogramVec(prometheus.HistogramOpts{
		Namespace: metrics.MetricsNamespace,
		Subsystem: subsystemName,
		Name:      "write_handle_time",
		Help:      "Время обработки записи в rabbitmq",
		Buckets:   prometheus.DefBuckets,
	}, []string{resultLabelName})
	readHandleTime = promauto.With(metrics.DefaultRegistry).NewHistogramVec(prometheus.HistogramOpts{
		Namespace: metrics.MetricsNamespace,
		Subsystem: subsystemName,
		Name:      "read_handle_time",
		Help:      "Время обработки чтения из rabbitmq",
		Buckets:   prometheus.DefBuckets,
	}, []string{resultLabelName})
)

func registerWriteHandleTime(ok bool, d time.Duration) {
	writeHandleTime.WithLabelValues(metrics.ConvertOk(ok)).Observe(d.Seconds())
}

func registerReadHandleTime(ok bool, d time.Duration) {
	readHandleTime.WithLabelValues(metrics.ConvertOk(ok)).Observe(d.Seconds())
}
