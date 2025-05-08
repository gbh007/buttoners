package client

import (
	"context"
	"errors"

	"github.com/gbh007/buttoners/services/log/internal/pb"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var errNoConnections = errors.New("no connection")

type Client struct {
	client pb.LogClient
	conn   *grpc.ClientConn
}

type UserInfo struct {
	ID int64
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
	c.client = pb.NewLogClient(conn)

	return c, nil
}

func (c *Client) Close() error {
	if c.conn == nil {
		return errNoConnections
	}

	return c.conn.Close()
}

func (c *Client) Activity(ctx context.Context, userID int64) (*LogData, error) {
	res, err := c.client.Activity(ctx, &pb.ActivityRequest{
		UserID: userID,
	})
	if err != nil {
		return nil, err
	}

	return &LogData{
		RequestCount: res.GetData().GetRequestCount(),
		LastRequest:  res.GetData().GetLastRequest().AsTime(),
	}, nil
}
