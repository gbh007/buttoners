package logclient

import (
	"context"
	"log/slog"
	"strings"

	"github.com/gbh007/buttoners/core/metrics"
	"github.com/gbh007/buttoners/core/observability"
	"github.com/imroc/req/v3"
	"go.opentelemetry.io/otel/trace"
)

type Client struct {
	client *req.Client
	addr   string
	token  string
	name   string
}

func New(logger *slog.Logger, tracer trace.Tracer, metrics *metrics.HTTPClientMetrics, addr, token, name string) (*Client, error) {
	c := &Client{
		addr:  strings.TrimRight(addr, "/"),
		token: token,
		name:  name,
	}

	c.client = req.C().
		SetBaseURL(c.addr).
		WrapRoundTrip(func(rt req.RoundTripper) req.RoundTripper {
			return observability.NewImroqReqRT(logger, metrics, tracer, rt, "log")
		}).
		SetCommonHeader("Authorization", c.token).
		SetCommonHeader("X-Client-Name", c.name).
		SetCommonHeader("Content-Type", ContentType)

	return c, nil
}

func (c *Client) Close() error {
	c.client.CloseIdleConnections()

	return nil
}

func (c *Client) Activity(ctx context.Context, userID int64) (ActivityResponse, error) {
	v, err := request[ActivityRequest, ActivityResponse](
		ctx,
		c,
		ActivityPath,
		ActivityRequest{
			UserID: userID,
		},
	)
	if err != nil {
		return ActivityResponse{}, err
	}

	return v, nil
}
