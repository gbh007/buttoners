package authclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func marshal[T any](v T) (io.Reader, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(data), nil
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

	data, err := marshal(reqV)
	if err != nil {
		return empty, fmt.Errorf("%w: marshal: %w", ErrProcess, err)
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, c.addr+path, data)
	if err != nil {
		return empty, fmt.Errorf("%w: make request: %w", ErrProcess, err)
	}

	request.Header.Set("Authorization", c.token)
	request.Header.Set("X-Client-Name", c.name)
	request.Header.Set("Content-Type", ContentType)

	resp, err := c.client.Do(request)
	if err != nil {
		return empty, fmt.Errorf("%w: request: %w", ErrProcess, err)
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent {
		return empty, nil
	}

	if resp.StatusCode == http.StatusOK {
		v, err := unmarshal[RP](resp.Body)
		if err != nil {
			return empty, fmt.Errorf("%w: unmarshal: %w", ErrProcess, err)
		}

		return v, nil
	}

	v, err := unmarshal[ErrorResponse](resp.Body)
	if err != nil {
		return empty, fmt.Errorf("%w: unmarshal: %w", ErrProcess, err)
	}

	switch {
	case resp.StatusCode == http.StatusUnauthorized:
		err = ErrUnauthorized
	case resp.StatusCode == http.StatusForbidden:
		err = ErrForbidden
	case resp.StatusCode == http.StatusNotFound:
		err = ErrNotFound
	case resp.StatusCode == http.StatusConflict:
		err = ErrConflict
	case resp.StatusCode < 500:
		err = ErrBadRequest
	default:
		err = ErrInternal
	}

	return empty, fmt.Errorf("%w: code=%s details=%s", err, v.Code, v.Details)
}
