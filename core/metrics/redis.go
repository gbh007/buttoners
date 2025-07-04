package metrics

import (
	"fmt"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type RedisMetrics struct {
	handle *prometheus.HistogramVec
	active *prometheus.GaugeVec
}

func NewRedisMetrics(
	reg prometheus.Registerer,
	handleBuckets []float64,
) (*RedisMetrics, error) {
	m := &RedisMetrics{
		handle: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: MetricsNamespace,
				Subsystem: "redis",
				Name:      "handle_seconds",
				Buckets:   handleBuckets,
			},
			[]string{"target_host", "method", "status"},
		),
		active: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: MetricsNamespace,
				Subsystem: "redis",
				Name:      "active_handlers_total",
			},
			[]string{"target_host", "method"},
		),
	}

	err := reg.Register(m.handle)
	if err != nil {
		return nil, fmt.Errorf("register handle: %w", err)
	}

	err = reg.Register(m.active)
	if err != nil {
		return nil, fmt.Errorf("register active: %w", err)
	}

	return m, nil
}

func (m *RedisMetrics) AddHandle(host, method, status string, d time.Duration) {
	m.handle.WithLabelValues(host, method, status).Observe(d.Seconds())
}

func (m *RedisMetrics) IncActive(host, method string) {
	m.active.WithLabelValues(host, method).Inc()
}

func (m *RedisMetrics) DecActive(host, method string) {
	m.active.WithLabelValues(host, method).Dec()
}
