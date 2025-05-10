package notificationclient

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/valyala/fasthttp"
)

func marshal[T any](w io.Writer, v T) error {
	err := json.NewEncoder(w).Encode(v)
	if err != nil {
		return err
	}

	return nil
}

func unmarshal[T any](r io.Reader) (T, error) {
	var v T

	err := json.NewDecoder(r).Decode(&v)
	if err != nil {
		return v, err
	}

	return v, nil
}

func request[RQ, RP any](ctx context.Context, c *Client, path string, reqV RQ) (RP, error) {
	var empty RP

	request := fasthttp.AcquireRequest()
	request.SetRequestURI(c.addr + path)
	request.Header.SetMethod(fasthttp.MethodPost)
	request.Header.SetContentType(ContentType)
	request.Header.Set("Authorization", c.token)
	request.Header.Set("X-Client-Name", c.name)

	// FIXME: поддержать телеметрию, использовать более быстрые библиотеки для json

	err := marshal(request.BodyWriter(), reqV)
	if err != nil {
		return empty, fmt.Errorf("%w: marshal: %w", ErrProcess, err)
	}

	resp := fasthttp.AcquireResponse()

	err = c.client.DoTimeout(request, resp, time.Second)

	fasthttp.ReleaseRequest(request)
	defer fasthttp.ReleaseResponse(resp)

	if err != nil {
		return empty, fmt.Errorf("%w: request: %w", ErrProcess, err)
	}

	if resp.StatusCode() == http.StatusNoContent {
		return empty, nil
	}

	if resp.StatusCode() == http.StatusOK {
		v, err := unmarshal[RP](resp.BodyStream())
		if err != nil {
			return empty, fmt.Errorf("%w: unmarshal: %w", ErrProcess, err)
		}

		return v, nil
	}

	v, err := unmarshal[ErrorResponse](resp.BodyStream())
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
