package metrics

import (
	"fmt"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type QueueWriterMetrics struct {
	handle *prometheus.HistogramVec
	active *prometheus.GaugeVec
}

func NewQueueWriterMetrics(
	reg prometheus.Registerer,
	handleBuckets []float64,
) (*QueueWriterMetrics, error) {
	m := &QueueWriterMetrics{
		handle: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: MetricsNamespace,
				Subsystem: "queue_writer",
				Name:      "handle_seconds",
				Buckets:   handleBuckets,
			},
			[]string{"target_host", "target_queue", "status"},
		),
		active: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: MetricsNamespace,
				Subsystem: "queue_writer",
				Name:      "active_handlers_total",
			},
			[]string{"target_host", "target_queue"},
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

func (m *QueueWriterMetrics) AddHandle(host, queue, status string, d time.Duration) {
	m.handle.WithLabelValues(host, queue, status).Observe(d.Seconds())
}

func (m *QueueWriterMetrics) IncActive(host, queue string) {
	m.active.WithLabelValues(host, queue).Inc()
}

func (m *QueueWriterMetrics) DecActive(host, queue string) {
	m.active.WithLabelValues(host, queue).Dec()
}

type QueueReaderMetrics struct {
	handle *prometheus.HistogramVec
	active *prometheus.GaugeVec
}

func NewQueueReaderMetrics(
	reg prometheus.Registerer,
	handleBuckets []float64,
) (*QueueReaderMetrics, error) {
	m := &QueueReaderMetrics{
		handle: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: MetricsNamespace,
				Subsystem: "queue_reader",
				Name:      "handle_seconds",
				Buckets:   handleBuckets,
			},
			[]string{"server_addr", "server_queue", "server_group", "status"},
		),
		active: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: MetricsNamespace,
				Subsystem: "queue_reader",
				Name:      "active_handlers_total",
			},
			[]string{"server_addr", "server_queue", "server_group"},
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

func (m *QueueReaderMetrics) AddHandle(addr, queue, group, status string, d time.Duration) {
	m.handle.WithLabelValues(addr, queue, group, status).Observe(d.Seconds())
}

func (m *QueueReaderMetrics) IncActive(addr, queue, group string) {
	m.active.WithLabelValues(addr, queue, group).Inc()
}

func (m *QueueReaderMetrics) DecActive(addr, queue, group string) {
	m.active.WithLabelValues(addr, queue, group).Dec()
}
