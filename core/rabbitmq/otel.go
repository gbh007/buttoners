package rabbitmq

import (
	"fmt"

	"github.com/rabbitmq/amqp091-go"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

const tracerName = "buttoners/rabbitmq"

func toMapCarrier(in amqp091.Table) propagation.MapCarrier {
	out := make(map[string]string, len(in))

	for k, v := range in {
		switch typedV := v.(type) {
		case string:
			out[k] = typedV
		default:
			out[k] = fmt.Sprint(typedV)
		}
	}

	return out
}

func fromMapCarrier(in propagation.MapCarrier) amqp091.Table {
	out := make(amqp091.Table, len(in))

	for k, v := range in {
		out[k] = v
	}

	return out
}

func newTracer(tp trace.TracerProvider) trace.Tracer {
	return tp.Tracer(tracerName)
}
