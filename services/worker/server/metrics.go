package server

import (
	"time"

	"github.com/gbh007/buttoners/core/metrics"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	resultLabelName = "result"
	subsystemName   = "worker"
)

var (
	handleTimeTotal = promauto.With(metrics.DefaultRegistry).NewHistogram(prometheus.HistogramOpts{
		Namespace: metrics.MetricsNamespace,
		Subsystem: subsystemName,
		Name:      "handle_time",
		Help:      "Суммарное время обработки задачи в worker",
		Buckets:   prometheus.DefBuckets,
	})
	businessHandleTimeTotal = promauto.With(metrics.DefaultRegistry).NewHistogramVec(prometheus.HistogramOpts{
		Namespace: metrics.MetricsNamespace,
		Subsystem: subsystemName,
		Name:      "business_handle_time",
		Help:      "Бизнесовое время обработки задачи в worker",
		Buckets:   []float64{0.5, 0.75, 1, 1.5, 2, 3, 5, 7, 10},
	}, []string{resultLabelName})
	activeTaskTotal = promauto.With(metrics.DefaultRegistry).NewGauge(prometheus.GaugeOpts{
		Namespace: metrics.MetricsNamespace,
		Subsystem: subsystemName,
		Name:      "active_task",
		Help:      "Общее количество активных задач в worker",
	})
)

func registerHandleTime(d time.Duration) {
	handleTimeTotal.Observe(d.Seconds())
}

func registerBusinessHandleTime(ok bool, d time.Duration) {
	businessHandleTimeTotal.WithLabelValues(metrics.ConvertOk(ok)).Observe(d.Seconds())
}
