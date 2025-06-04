package logclient

import (
	"context"
	"fmt"
	"net/http"
)

func request[RQ, RP any](ctx context.Context, c *Client, path string, reqV RQ) (RP, error) {
	var (
		empty, successResponse RP
		errorResponse          ErrorResponse
	)

	resp, err := c.client.R().
		SetContext(ctx).
		SetBody(reqV).
		SetSuccessResult(&successResponse).
		SetErrorResult(&errorResponse).
		Post(path)

	if err != nil {
		return empty, fmt.Errorf("%w: request: %w", ErrProcess, err)
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent {
		return empty, nil
	}

	if resp.IsSuccessState() {
		return successResponse, nil
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

	return empty, fmt.Errorf("%w: code=%s details=%s", err, errorResponse.Code, errorResponse.Details)
}
