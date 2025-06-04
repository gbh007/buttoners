package gateclient

import (
	"context"
	"errors"
	"time"

	"github.com/gbh007/buttoners/core/clients/gateclient/gen/pb"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

var errNoConnections = errors.New("no connection")

type Client struct {
	gateClient         pb.GateClient
	notificationClient pb.NotificationClient
	logClient          pb.LogClient
	conn               *grpc.ClientConn
}

type UserInfo struct {
	ID int64
}

func New(addr string) (*Client, error) {
	c := new(Client)

	conn, err := grpc.NewClient(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
	)
	if err != nil {
		return nil, err
	}

	c.conn = conn
	c.gateClient = pb.NewGateClient(conn)
	c.notificationClient = pb.NewNotificationClient(conn)
	c.logClient = pb.NewLogClient(conn)

	return c, nil
}

func (c *Client) Close() error {
	if c.conn == nil {
		return errNoConnections
	}

	return c.conn.Close()
}

func (c *Client) Login(ctx context.Context, login, pass string) (string, error) {
	res, err := c.gateClient.Login(ctx, &pb.LoginRequest{
		Login:    login,
		Password: pass,
	})
	if err != nil {
		return "", err
	}

	return res.GetToken(), nil
}

func (c *Client) Register(ctx context.Context, login, pass string) error {
	_, err := c.gateClient.Register(ctx, &pb.RegisterRequest{
		Login:    login,
		Password: pass,
	})
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) ButtonClick(ctx context.Context, token string, duration, chance int64) error {
	ctx = metadata.AppendToOutgoingContext(ctx, SessionHeader, token)

	_, err := c.gateClient.Button(ctx, &pb.ButtonRequest{
		Duration: duration,
		Chance:   chance,
	})
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) Read(ctx context.Context, token string, all bool, id int64) error {
	ctx = metadata.AppendToOutgoingContext(ctx, SessionHeader, token)

	_, err := c.notificationClient.Read(ctx, &pb.NotificationReadRequest{
		Id:  id,
		All: all,
	})
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) List(ctx context.Context, token string) ([]NotificationData, error) {
	ctx = metadata.AppendToOutgoingContext(ctx, SessionHeader, token)

	res, err := c.notificationClient.List(ctx, new(pb.NotificationListRequest))
	if err != nil {
		return nil, err
	}

	notifications := make([]NotificationData, len(res.GetList()))

	for index, raw := range res.GetList() {
		notifications[index] = NotificationData{
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

func (c *Client) Activity(ctx context.Context, token string) (int64, time.Time, error) {
	ctx = metadata.AppendToOutgoingContext(ctx, SessionHeader, token)

	res, err := c.logClient.Activity(ctx, new(pb.ActivityRequest))
	if err != nil {
		return 0, time.Time{}, err
	}

	return res.GetRequestCount(), res.GetLastRequest().AsTime(), nil
}
