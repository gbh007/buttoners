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
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel/trace"
)

func NewEchoMiddleware(
	logger *slog.Logger,
	metrics *metrics.HTTPServerMetrics,
	serverName string,
) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			tStart := time.Now()
			r := c.Request()
			path := c.Path()
			host := r.Host

			metrics.IncActive(host, path, r.Method)
			defer metrics.DecActive(host, path, r.Method)

			var (
				requestBody             []byte
				err                     error
				requestLog, responseLog []any
			)

			requestLog = append(
				requestLog,
				slog.String("host", host),
				slog.String("path", path),
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

			if r.Body != nil {
				requestBody, err = io.ReadAll(r.Body)
				if err != nil {
					logger.ErrorContext(
						r.Context(), serverName+" read request body",
						slog.String("error", err.Error()),
					)

					return err
				}

				err = r.Body.Close()
				if err != nil {
					logger.ErrorContext(
						r.Context(), serverName+" close request body",
						slog.String("error", err.Error()),
					)

					return err
				}

				r.Body = io.NopCloser(bytes.NewReader(requestBody))
			}

			if len(requestBody) > 0 {
				requestLog = append(
					requestLog,
					slog.String("body", string(requestBody)),
				)
			}

			wReplace := &httpResponseWriter{
				old:  c.Response().Writer,
				body: &bytes.Buffer{},
			}

			c.Response().Writer = wReplace

			defer func() {
				metrics.AddHandle(host, path, r.Method, wReplace.code, time.Since(tStart))
			}()

			handleErr := next(c)

			if wReplace.code == 0 && handleErr != nil {
				wReplace.code = http.StatusInternalServerError
				responseLog = append(
					responseLog,
					slog.String("error", handleErr.Error()),
				)
			}

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

			logger.InfoContext(
				r.Context(), serverName+" server request",
				slog.String("trace_id", trace.SpanContextFromContext(r.Context()).TraceID().String()),
				slog.Group("request", requestLog...),
				slog.Group("response", responseLog...),
			)

			return handleErr
		}
	}
}
