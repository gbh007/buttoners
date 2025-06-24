package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	endpointLabelName = "endpoint"
	resultLabelName   = "result"
	hostLabelName     = "host"
	instanceLabelName = "instance"

	MetricsNamespace = "buttoners"
	subsystemName    = "protobuf"

	ResultOK    = "ok"
	ResultError = "err"
	ResultFail  = "fail"
)

var (
	DefaultRegistry = prometheus.NewRegistry()

	DefaultTimeBuckets = []float64{0.0001, 0.0005, 0.001, 0.005, 0.01, 0.05, 0.1, 0.2, 0.5, 1, 1.5, 2}

	requestTime = promauto.With(DefaultRegistry).NewHistogramVec(prometheus.HistogramOpts{
		Namespace: MetricsNamespace,
		Subsystem: subsystemName,
		Name:      "request_duration",
		Help:      "Распределение времени запроса",
		Buckets:   prometheus.DefBuckets,
	}, []string{endpointLabelName, resultLabelName})
)

func init() {
	DefaultRegistry.MustRegister(
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{
			Namespace: MetricsNamespace,
		}),
		collectors.NewGoCollector(),
	)
}

func ConvertOk(ok bool) string {
	if ok {
		return ResultOK
	}

	return ResultFail
}

func LogRequest(action string, ok bool, d time.Duration) {
	requestTime.WithLabelValues(action, ConvertOk(ok)).Observe(d.Seconds())
}
