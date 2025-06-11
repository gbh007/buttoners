package logclient

import (
	"context"
	"log/slog"
	"net/http"
	"strings"

	"github.com/gbh007/buttoners/core/metrics"
	"github.com/gbh007/buttoners/core/observability"
	"github.com/imroc/req/v3"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type Client struct {
	client *req.Client
	addr   string
	token  string
	name   string
}

func New(logger *slog.Logger, metrics *metrics.HTTPClientMetrics, addr, token, name string) (*Client, error) {
	c := &Client{
		addr:  strings.TrimRight(addr, "/"),
		token: token,
		name:  name,
	}

	client := req.C()

	client.Transport = req.T().WrapRoundTrip(func(rt http.RoundTripper) http.RoundTripper {
		return observability.NewHTTPTransport(logger, metrics, otelhttp.NewTransport(rt), "log")
	})

	client = client.
		SetBaseURL(c.addr).
		SetCommonHeader("Authorization", c.token).
		SetCommonHeader("X-Client-Name", c.name).
		SetCommonHeader("Content-Type", ContentType)

	c.client = client

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
