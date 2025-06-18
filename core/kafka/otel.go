package kafka

import (
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/protocol"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

const tracerName = "buttoners/kafka"

func toMapCarrier(in []kafka.Header) propagation.MapCarrier {
	out := make(map[string]string, len(in))

	for _, raw := range in {
		out[raw.Key] = string(raw.Value)
	}

	return out
}

func fromMapCarrier(in propagation.MapCarrier) []kafka.Header {
	out := make([]kafka.Header, 0, len(in))

	for k, v := range in {
		out = append(out, protocol.Header{
			Key:   k,
			Value: []byte(v),
		})
	}

	return out
}

func newTracer(tp trace.TracerProvider) trace.Tracer {
	return tp.Tracer(tracerName)
}
