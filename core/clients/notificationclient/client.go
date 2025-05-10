package notificationclient

import (
	"context"
	"strings"
	"time"

	"github.com/valyala/fasthttp"
)

type Client struct {
	client *fasthttp.Client
	addr   string
	token  string
	name   string
}

func New(addr, token, name string) (*Client, error) {
	c := &Client{
		client: &fasthttp.Client{
			Transport:    fasthttp.DefaultTransport,
			ReadTimeout:  time.Second,
			WriteTimeout: time.Second,
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

func (c *Client) New(ctx context.Context, req NewRequest) error {
	_, err := request[NewRequest, struct{}](
		ctx,
		c,
		NewPath,
		req,
	)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) List(ctx context.Context, userID int64) (ListResponse, error) {
	v, err := request[ListRequest, ListResponse](
		ctx,
		c,
		ListPath,
		ListRequest{
			UserID: userID,
		},
	)
	if err != nil {
		return ListResponse{}, err
	}

	return v, nil
}

func (c *Client) Read(ctx context.Context, userID, notificationID int64) error {
	_, err := request[ReadRequest, struct{}](
		ctx,
		c,
		ReadPath,
		ReadRequest{
			UserID: userID,
			ID:     notificationID,
		},
	)
	if err != nil {
		return err
	}

	return nil
}
