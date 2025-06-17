package metrics

import (
	"fmt"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type GRPCClientMetrics struct {
	handle *prometheus.HistogramVec
	active *prometheus.GaugeVec
}

func NewGRPCClientMetrics(
	reg prometheus.Registerer,
	handleBuckets []float64,
) (*GRPCClientMetrics, error) {
	m := &GRPCClientMetrics{
		handle: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: MetricsNamespace,
				Subsystem: "grpc_client",
				Name:      "handle_seconds",
				Buckets:   handleBuckets,
			},
			[]string{"target_host", "target_path", "status"},
		),
		active: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: MetricsNamespace,
				Subsystem: "grpc_client",
				Name:      "active_handlers_total",
			},
			[]string{"target_host", "target_path"},
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

func (m *GRPCClientMetrics) AddHandle(host, path string, status string, d time.Duration) {
	m.handle.WithLabelValues(host, path, status).Observe(d.Seconds())
}

func (m *GRPCClientMetrics) IncActive(host, path string) {
	m.active.WithLabelValues(host, path).Inc()
}

func (m *GRPCClientMetrics) DecActive(host, path string) {
	m.active.WithLabelValues(host, path).Dec()
}

type GRPCServerMetrics struct {
	handle *prometheus.HistogramVec
	active *prometheus.GaugeVec
}

func NewGRPCServerMetrics(
	reg prometheus.Registerer,
	handleBuckets []float64,
) (*GRPCServerMetrics, error) {
	m := &GRPCServerMetrics{
		handle: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: MetricsNamespace,
				Subsystem: "grpc_server",
				Name:      "handle_seconds",
				Buckets:   handleBuckets,
			},
			[]string{"server_addr", "server_path", "status"},
		),
		active: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: MetricsNamespace,
				Subsystem: "grpc_server",
				Name:      "active_handlers_total",
			},
			[]string{"server_addr", "server_path"},
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

func (m *GRPCServerMetrics) AddHandle(addr, path string, status string, d time.Duration) {
	m.handle.WithLabelValues(addr, path, status).Observe(d.Seconds())
}

func (m *GRPCServerMetrics) IncActive(addr, path string) {
	m.active.WithLabelValues(addr, path).Inc()
}

func (m *GRPCServerMetrics) DecActive(addr, path string) {
	m.active.WithLabelValues(addr, path).Dec()
}
