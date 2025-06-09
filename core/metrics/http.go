package metrics

import (
	"fmt"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type HTTPClientMetrics struct {
	handle *prometheus.HistogramVec
	active *prometheus.GaugeVec
}

func NewHTTPClientMetrics(
	reg prometheus.Registerer,
	handleBuckets []float64,
) (*HTTPClientMetrics, error) {
	m := &HTTPClientMetrics{
		handle: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: MetricsNamespace,
				Subsystem: "http_client",
				Name:      "handle_seconds",
				Buckets:   handleBuckets,
			},
			[]string{"target_host", "target_path", "method", "status"},
		),
		active: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: MetricsNamespace,
				Subsystem: "http_client",
				Name:      "active_handlers_total",
			},
			[]string{"target_host", "target_path", "method"},
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

func (m *HTTPClientMetrics) AddHandle(host, path, method string, status int, d time.Duration) {
	m.handle.WithLabelValues(host, path, method, strconv.Itoa(status)).Observe(d.Seconds())
}

func (m *HTTPClientMetrics) IncActive(host, path, method string) {
	m.active.WithLabelValues(host, path, method).Inc()
}

func (m *HTTPClientMetrics) DecActive(host, path, method string) {
	m.active.WithLabelValues(host, path, method).Dec()
}

type HTTPServerMetrics struct {
	handle *prometheus.HistogramVec
	active *prometheus.GaugeVec
}

func NewHTTPServerMetrics(
	reg prometheus.Registerer,
	handleBuckets []float64,
) (*HTTPServerMetrics, error) {
	m := &HTTPServerMetrics{
		handle: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: MetricsNamespace,
				Subsystem: "http_server",
				Name:      "handle_seconds",
				Buckets:   handleBuckets,
			},
			[]string{"server_addr", "server_path", "method", "status"},
		),
		active: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: MetricsNamespace,
				Subsystem: "http_server",
				Name:      "active_handlers_total",
			},
			[]string{"server_addr", "server_path", "method"},
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

func (m *HTTPServerMetrics) AddHandle(addr, path, method string, status int, d time.Duration) {
	m.handle.WithLabelValues(addr, path, method, strconv.Itoa(status)).Observe(d.Seconds())
}

func (m *HTTPServerMetrics) IncActive(addr, path, method string) {
	m.active.WithLabelValues(addr, path, method).Inc()
}

func (m *HTTPServerMetrics) DecActive(addr, path, method string) {
	m.active.WithLabelValues(addr, path, method).Dec()
}
