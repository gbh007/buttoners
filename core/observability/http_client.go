package observability

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gbh007/buttoners/core/metrics"
	"go.opentelemetry.io/otel/trace"
)

type HTTPTransport struct {
	logger     *slog.Logger
	metrics    *metrics.HTTPClientMetrics
	next       http.RoundTripper
	clientName string
}

func NewHTTPTransport(
	logger *slog.Logger,
	metrics *metrics.HTTPClientMetrics,
	next http.RoundTripper,
	clientName string,
) *HTTPTransport {
	return &HTTPTransport{
		logger:     logger,
		metrics:    metrics,
		next:       next,
		clientName: clientName,
	}
}

func (h *HTTPTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req = req.Clone(req.Context())
	tStart := time.Now()

	h.metrics.IncActive(req.URL.Host, req.URL.Path, req.Method)
	defer h.metrics.DecActive(req.URL.Host, req.URL.Path, req.Method)

	var (
		requestBody, responseBody []byte
		err                       error
		requestLog, responseLog   []any
	)

	requestLog = append(
		requestLog,
		slog.String("host", req.URL.Host),
		slog.String("path", req.URL.Path),
	)

	if len(req.Header) > 0 {
		headers := make(map[string]string)

		for k, v := range req.Header {
			headers[k] = strings.Join(v, ";")
		}

		requestLog = append(
			requestLog,
			slog.Any("headers", headers),
		)
	}

	if req.GetBody != nil {
		body, err := req.GetBody()
		if err != nil {
			return nil, err
		}

		requestBody, err = io.ReadAll(body)
		if err != nil {
			return nil, err
		}
	} else if req.Body != nil {
		requestBody, err = io.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}

		err = req.Body.Close()
		if err != nil {
			return nil, err
		}

		req.Body = io.NopCloser(bytes.NewReader(requestBody))
		req.GetBody = func() (io.ReadCloser, error) {
			return io.NopCloser(bytes.NewReader(requestBody)), nil
		}
	}

	if len(requestBody) > 0 {
		requestLog = append(
			requestLog,
			slog.String("body", string(requestBody)),
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
