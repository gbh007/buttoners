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

type HTTPMiddleware struct {
	logger     *slog.Logger
	metrics    *metrics.HTTPServerMetrics
	next       http.Handler
	serverName string
}

func NewHTTPMiddleware(
	logger *slog.Logger,
	metrics *metrics.HTTPServerMetrics,
	serverName string,
	next http.Handler,
) *HTTPMiddleware {
	return &HTTPMiddleware{
		logger:     logger,
		metrics:    metrics,
		next:       next,
		serverName: serverName,
	}
}

func (h *HTTPMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	tStart := time.Now()

	h.metrics.IncActive(r.URL.Host, r.URL.Path, r.Method)
	defer h.metrics.DecActive(r.URL.Host, r.URL.Path, r.Method)

	var (
		requestBody             []byte
		err                     error
		requestLog, responseLog []any
	)

	requestLog = append(
		requestLog,
		slog.String("host", r.URL.Host),
		slog.String("path", r.URL.Path),
	)

	if len(r.Header) > 0 {
		headers := make(map[string]string)

		for k, v := range r.Header {
			headers[k] = strings.Join(v, ";")
		}

		requestLog = append(
			requestLog,
			slog.Any("headers", headers),
		)
	}

	if r.GetBody != nil {
		body, err := r.GetBody()
		if err != nil {
			h.logger.ErrorContext(
				r.Context(), h.serverName+" get request body",
				slog.String("error", err.Error()),
			)

			return
		}

		requestBody, err = io.ReadAll(body)
		if err != nil {
			h.logger.ErrorContext(
				r.Context(), h.serverName+" read request body",
				slog.String("error", err.Error()),
			)

			return
		}
	} else if r.Body != nil {
		requestBody, err = io.ReadAll(r.Body)
		if err != nil {
			h.logger.ErrorContext(
				r.Context(), h.serverName+" read request body",
				slog.String("error", err.Error()),
			)

			return
		}

		err = r.Body.Close()
		if err != nil {
			h.logger.ErrorContext(
				r.Context(), h.serverName+" close request body",
				slog.String("error", err.Error()),
			)

			return
		}

		r.Body = io.NopCloser(bytes.NewReader(requestBody))
		r.GetBody = func() (io.ReadCloser, error) {
			return io.NopCloser(bytes.NewReader(requestBody)), nil
		}
	}

	if len(requestBody) > 0 {
		requestLog = append(
			requestLog,
			slog.String("body", string(requestBody)),
		)
	}

	wReplace := &httpResponseWriter{
		old:  w,
		body: &bytes.Buffer{},
	}

	w = wReplace

	defer func() {
		h.metrics.AddHandle(r.URL.Host, r.URL.Path, r.Method, wReplace.code, time.Since(tStart))
	}()

	h.next.ServeHTTP(w, r)

	responseLog = append(
		responseLog,
		slog.String("status", strconv.Itoa(wReplace.code)),
	)

	if len(wReplace.Header()) > 0 {
		headers := make(map[string]string)

		for k, v := range wReplace.Header() {
			headers[k] = strings.Join(v, ";")
		}

		responseLog = append(
			responseLog,
			slog.Any("headers", headers),
		)
	}

	if wReplace.body != nil && wReplace.body.Len() > 0 {
		responseLog = append(
			responseLog,
			slog.String("body", wReplace.body.String()),
		)
	}

	h.logger.InfoContext(
		r.Context(), h.serverName+" server request",
		slog.String("trace_id", trace.SpanContextFromContext(r.Context()).TraceID().String()),
		slog.Group("request", requestLog...),
		slog.Group("response", responseLog...),
	)
}

type httpResponseWriter struct {
	old  http.ResponseWriter
	code int
	body *bytes.Buffer
}

func (rw *httpResponseWriter) Header() http.Header {
	return rw.old.Header()
}

func (rw *httpResponseWriter) Write(data []byte) (int, error) {
	_, _ = rw.body.Write(data)

	return rw.old.Write(data)
}

func (rw *httpResponseWriter) WriteHeader(statusCode int) {
	rw.code = statusCode
	rw.old.WriteHeader(statusCode)
}
