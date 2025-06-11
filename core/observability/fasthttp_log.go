package observability

import (
	"context"
	"log/slog"

	"github.com/valyala/fasthttp"
	"go.opentelemetry.io/otel/trace"
)

func LogFastHTTPData(ctx context.Context, logger *slog.Logger, msg string, request *fasthttp.Request, resp *fasthttp.Response) {
	rqLog := []any{}

	if request != nil {
		rqLog = append(
			rqLog,
			slog.String("host", string(request.Host())),
			slog.String("path", string(request.URI().Path())),
			slog.String("method", string(request.Header.Method())),
			slog.String("body", string(request.Body())),
		)

		headers := make(map[string]string)
		request.Header.VisitAll(func(key, value []byte) {
			old := headers[string(key)]
			if old != "" {
				old += ";"
			}
			old += string(value)
			headers[string(key)] = old
		})

		if len(headers) > 0 {
			rqLog = append(
				rqLog,
				slog.Any("headers", headers),
			)
		}
	}

	rpLog := []any{}

	if resp != nil {
		rpLog = append(
			rpLog,
			slog.Int("status", resp.StatusCode()),
			slog.String("body", string(resp.Body())),
		)

		headers := make(map[string]string)
		resp.Header.VisitAll(func(key, value []byte) {
			old := headers[string(key)]
			if old != "" {
				old += ";"
			}
			old += string(value)
			headers[string(key)] = old
		})

		if len(headers) > 0 {
			rpLog = append(
				rpLog,
				slog.Any("headers", headers),
			)
		}
	}

	logger.InfoContext(
		ctx, msg,
		slog.String("trace_id", trace.SpanContextFromContext(ctx).TraceID().String()),
		slog.Group("request", rqLog...),
		slog.Group("response", rpLog...),
	)
}
