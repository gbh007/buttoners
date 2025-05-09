package authclient

import (
	"context"
	"net/http"
	"strings"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type Client struct {
	client *http.Client
	addr   string
	token  string
	name   string
}

func New(addr, token, name string) (*Client, error) {
	c := &Client{
		client: &http.Client{
			Transport: otelhttp.NewTransport(http.DefaultTransport),
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
		"/api/v1/login",
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
		"/api/v1/logout",
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
		"/api/v1/register",
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
		"/api/v1/info",
		InfoRequest{
			Token: token,
		},
	)
	if err != nil {
		return InfoResponse{}, err
	}

	return v, nil
}
