package observability

import (
	"bytes"
	"io"
	"log/slog"
	"strconv"
	"strings"
	"time"

	"github.com/gbh007/buttoners/core/metrics"
	imrocReq "github.com/imroc/req/v3"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

type ImroqReqRT struct {
	logger     *slog.Logger
	metrics    *metrics.HTTPClientMetrics
	tracer     trace.Tracer
	next       imrocReq.RoundTripper
	clientName string
}

func NewImroqReqRT(
	logger *slog.Logger,
	metrics *metrics.HTTPClientMetrics,
	tracer trace.Tracer,
	next imrocReq.RoundTripper,
	clientName string,
) *ImroqReqRT {
	return &ImroqReqRT{
		logger:     logger,
		metrics:    metrics,
		tracer:     tracer,
		next:       next,
		clientName: clientName,
	}
}

func (h *ImroqReqRT) RoundTrip(req *imrocReq.Request) (*imrocReq.Response, error) {
	tStart := time.Now()

	_, span := h.tracer.Start(req.Context(), h.clientName+":"+req.URL.Path)
	defer span.End()

	otel.GetTextMapPropagator().Inject(req.Context(), propagation.HeaderCarrier(req.Headers))

	h.metrics.IncActive(req.URL.Host, req.URL.Path, req.Method)
	defer h.metrics.DecActive(req.URL.Host, req.URL.Path, req.Method)

	var (
		responseBody            []byte
		err                     error
		requestLog, responseLog []any
	)

	requestLog = append(
		requestLog,
		slog.String("host", req.URL.Host),
		slog.String("path", req.URL.Path),
	)

	if len(req.Headers) > 0 {
		headers := make(map[string]string)

		for k, v := range req.Headers {
			headers[k] = strings.Join(v, ";")
		}

		requestLog = append(
			requestLog,
			slog.Any("headers", headers),
		)
	}

	if len(req.Body) > 0 {
		requestLog = append(
			requestLog,
			slog.String("body", string(req.Body)),
		)
	}

	resp, err := h.next.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	defer func() {
		h.metrics.AddHandle(req.URL.Host, req.URL.Path, req.Method, resp.StatusCode, time.Since(tStart))
	}()

	responseLog = append(
		responseLog,
		slog.String("status", strconv.Itoa(resp.StatusCode)),
	)

	if len(resp.Header) > 0 {
		headers := make(map[string]string)

		for k, v := range resp.Header {
			headers[k] = strings.Join(v, ";")
		}

		responseLog = append(
			responseLog,
			slog.Any("headers", headers),
		)
	}

	if resp.Body != nil {
		responseBody, err = io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		err = resp.Body.Close()
		if err != nil {
			return nil, err
		}

		resp.Body = io.NopCloser(bytes.NewReader(responseBody))
	}

	if len(responseBody) > 0 {
		responseLog = append(
			responseLog,
			slog.String("body", string(responseBody)),
		)
	}

	h.logger.InfoContext(
		req.Context(), h.clientName+" client request",
		slog.String("trace_id", trace.SpanContextFromContext(req.Context()).TraceID().String()),
		slog.Group("request", requestLog...),
		slog.Group("response", responseLog...),
	)

	return resp, nil
}
