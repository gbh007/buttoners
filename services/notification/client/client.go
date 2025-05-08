package client

import (
	"context"
	"errors"

	"github.com/gbh007/buttoners/services/notification/internal/pb"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var errNoConnections = errors.New("no connection")

type Client struct {
	client pb.NotificationClient
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
	c.client = pb.NewNotificationClient(conn)

	return c, nil
}

func (c *Client) Close() error {
	if c.conn == nil {
		return errNoConnections
	}

	return c.conn.Close()
}

func (c *Client) New(ctx context.Context, userID int64, n *Notification) error {
	_, err := c.client.New(ctx, &pb.NewRequest{
		UserID: userID,
		Data: &pb.NotificationData{
			Kind:    n.Kind,
			Level:   n.Level,
			Title:   n.Title,
			Body:    n.Body,
			Created: timestamppb.New(n.Created),
		},
	})
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) List(ctx context.Context, userID int64) ([]*Notification, error) {
	res, err := c.client.List(ctx, &pb.ListRequest{UserID: userID})
	if err != nil {
		return nil, err
	}

	notifications := make([]*Notification, len(res.GetList()))

	for index, raw := range res.GetList() {
		notifications[index] = &Notification{
			ID:      raw.GetId(),
			Kind:    raw.GetKind(),
			Level:   raw.GetLevel(),
			Title:   raw.GetTitle(),
			Body:    raw.GetBody(),
			Created: raw.GetCreated().AsTime(),
		}
	}

	return notifications, nil
}

func (c *Client) Read(ctx context.Context, id int64) error {
	_, err := c.client.Read(ctx, &pb.ReadRequest{Id: id})
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) ReadAll(ctx context.Context, userID int64) error {
	_, err := c.client.ReadAll(ctx, &pb.ReadAllRequest{UserID: userID})
	if err != nil {
		return err
	}

	return nil
}
