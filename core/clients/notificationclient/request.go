package notificationclient

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gbh007/buttoners/core/observability"
	"github.com/valyala/fasthttp"
	"go.opentelemetry.io/otel"
)

func marshal[T any](w io.Writer, v T) error {
	err := json.NewEncoder(w).Encode(v)
	if err != nil {
		return err
	}

	return nil
}

func unmarshal[T any](data []byte) (T, error) {
	var v T

	err := json.Unmarshal(data, &v)
	if err != nil {
		return v, err
	}

	return v, nil
}

func request[RQ, RP any](ctx context.Context, c *Client, path string, reqV RQ) (RP, error) {
	tStart := time.Now()

	ctx, span := c.tracer.Start(ctx, "notification:"+path)
	defer span.End()

	var empty RP

	request := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(request)

	request.SetRequestURI(c.addr + path)
	request.Header.SetMethod(fasthttp.MethodPost)
	request.Header.SetContentType(ContentType)
	request.Header.Set("Authorization", c.token)
	request.Header.Set("X-Client-Name", c.name)

	otel.GetTextMapPropagator().Inject(ctx, &headerWraper{
		raw: &request.Header,
	})

	// TODO: метрики в RoundTripper?
	c.metrics.IncActive(string(request.Host()), string(request.URI().Path()), string(request.Header.Method()))
	defer c.metrics.DecActive(string(request.Host()), string(request.URI().Path()), string(request.Header.Method()))

	// FIXME: использовать более быстрые библиотеки для json

	err := marshal(request.BodyWriter(), reqV)
	if err != nil {
		return empty, fmt.Errorf("%w: marshal: %w", ErrProcess, err)
	}

	resp := fasthttp.AcquireResponse()

	err = c.client.DoTimeout(request, resp, time.Second)

	defer fasthttp.ReleaseResponse(resp)

	defer func() {
		c.metrics.AddHandle(string(request.Host()), string(request.URI().Path()), string(request.Header.Method()), resp.StatusCode(), time.Since(tStart))
	}()

	defer observability.LogFastHTTPData(ctx, c.logger, "notification client request", request, resp)

	if err != nil {
		return empty, fmt.Errorf("%w: request: %w", ErrProcess, err)
	}

	if resp.StatusCode() == http.StatusNoContent {
		return empty, nil
	}

	if resp.StatusCode() == http.StatusOK {
		v, err := unmarshal[RP](resp.Body())
		if err != nil {
			return empty, fmt.Errorf("%w: unmarshal: %w", ErrProcess, err)
		}

		return v, nil
	}

	v, err := unmarshal[ErrorResponse](resp.Body())
	if err != nil {
		return empty, fmt.Errorf("%w: unmarshal: %w", ErrProcess, err)
	}

	switch {
	case resp.StatusCode() == http.StatusUnauthorized:
		err = ErrUnauthorized
	case resp.StatusCode() == http.StatusForbidden:
		err = ErrForbidden
	case resp.StatusCode() == http.StatusNotFound:
		err = ErrNotFound
	case resp.StatusCode() < 500:
		err = ErrBadRequest
	default:
		err = ErrInternal
	}

	return empty, fmt.Errorf("%w: code=%s details=%s", err, v.Code, v.Details)
}
