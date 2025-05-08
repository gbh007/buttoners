package server

import (
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/gbh007/buttoners/core/metrics"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	typeLabelName = "type"

	subsystemName = "gate"
)

var cacheTimeTotal = promauto.With(metrics.DefaultRegistry).NewHistogramVec(prometheus.HistogramOpts{
	Namespace: metrics.MetricsNamespace,
	Subsystem: subsystemName,
	Name:      "cache_time",
	Help:      "Суммарное время обращений по кешу",
	Buckets:   []float64{0.0001, 0.0005, 0.001, 0.005, 0.01, 0.05, 0.1, 0.5},
}, []string{typeLabelName})

func registerCacheHandle(t string, d time.Duration) {
	cacheTimeTotal.WithLabelValues(t).Observe(d.Seconds())
}

func randomSHA256String() string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(time.Now().String())))
}

func fistOf(in []string) string {
	if len(in) > 0 {
		return in[0]
	}

	return ""
}
