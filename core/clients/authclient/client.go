package authclient

import (
	"context"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/gbh007/buttoners/core/metrics"
	"github.com/gbh007/buttoners/core/observability"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type Client struct {
	client *http.Client
	addr   string
	token  string
	name   string
}

func New(logger *slog.Logger, metrics *metrics.HTTPClientMetrics, addr, token, name string) (*Client, error) {
	c := &Client{
		client: &http.Client{
			Transport: observability.NewHTTPTransport(logger, metrics, otelhttp.NewTransport(http.DefaultTransport), "auth"),
			Timeout:   time.Second,
		},
		addr:  strings.TrimRight(addr, "/"),
		token: token,
	}

	return c, nil
}

func (c *Client) Close() error {
	c.client.CloseIdleConnections()

	return nil
}

func (c *Client) Login(ctx context.Context, login, pass string) (LoginResponse, error) {
	v, err := request[LoginRequest, LoginResponse](
		ctx,
		c,
		LoginPath,
		LoginRequest{
			Login:    login,
			Password: pass,
		},
	)
	if err != nil {
		return LoginResponse{}, err
	}

	return v, nil
}

func (c *Client) Logout(ctx context.Context, token string) error {
	_, err := request[LogoutRequest, struct{}](
		ctx,
		c,
		LogoutPath,
		LogoutRequest{
			Token: token,
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) Register(ctx context.Context, login, pass string) error {
	_, err := request[RegisterRequest, struct{}](
		ctx,
		c,
		RegisterPath,
		RegisterRequest{
			Login:    login,
			Password: pass,
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) Info(ctx context.Context, token string) (InfoResponse, error) {
	v, err := request[InfoRequest, InfoResponse](
		ctx,
		c,
		InfoPath,
		InfoRequest{
			Token: token,
		},
	)
	if err != nil {
		return InfoResponse{}, err
	}

	return v, nil
}
