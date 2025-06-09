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

	resultOK   = "ok"
	resultFail = "fail"
)

var (
	DefaultRegistry = prometheus.NewRegistry()

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
		return resultOK
	}

	return resultFail
}

func LogRequest(action string, ok bool, d time.Duration) {
	requestTime.WithLabelValues(action, ConvertOk(ok)).Observe(d.Seconds())
}
