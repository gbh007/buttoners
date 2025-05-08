package client

import (
	"context"
	"errors"

	"github.com/gbh007/buttoners/services/auth/internal/pb"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var errNoConnections = errors.New("no connection")

type Client struct {
	client pb.AuthClient
	conn   *grpc.ClientConn
}

type UserInfo struct {
	ID    int64
	Token string
}

func New(addr string) (*Client, error) {
	c := new(Client)

	conn, err := grpc.Dial(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
	)
	if err != nil {
		return nil, err
	}

	c.conn = conn
	c.client = pb.NewAuthClient(conn)

	return c, nil
}

func (c *Client) Close() error {
	if c.conn == nil {
		return errNoConnections
	}

	return c.conn.Close()
}

func (c *Client) Login(ctx context.Context, login, pass string) (string, error) {
	res, err := c.client.Login(ctx, &pb.LoginRequest{
		Login:    login,
		Password: pass,
	})
	if err != nil {
		return "", err
	}

	return res.GetToken(), nil
}

func (c *Client) Register(ctx context.Context, login, pass string) error {
	_, err := c.client.Register(ctx, &pb.RegisterRequest{
		Login:    login,
		Password: pass,
	})
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) Info(ctx context.Context, token string) (*UserInfo, error) {
	res, err := c.client.Info(ctx, &pb.InfoRequest{
		Token: token,
	})
	if err != nil {
		return nil, err
	}

	return &UserInfo{
		ID:    res.GetUserID(),
		Token: token,
	}, nil
}
