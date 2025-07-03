package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/gbh007/buttoners/core/clients/authclient"
	"github.com/gbh007/buttoners/core/clients/gateclient/gen/pb"
	"github.com/gbh007/buttoners/core/clients/logclient"
	"github.com/gbh007/buttoners/core/clients/notificationclient"
	"github.com/gbh007/buttoners/core/dto"
	"github.com/gbh007/buttoners/core/kafka"
	"github.com/gbh007/buttoners/core/redis"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var errInvalidInputData = errors.New("invalid")

type pbServer struct {
	pb.UnimplementedGateServer
	pb.UnimplementedNotificationServer
	pb.UnimplementedLogServer

	logger *slog.Logger

	auth         *authclient.Client
	notification *notificationclient.Client
	log          *logclient.Client
	kafkaTask    *kafka.Client
	kafkaLog     *kafka.Client
	redis        *redis.Client[dto.UserInfo]
}

func (s *pbServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	resp, err := s.auth.Login(ctx, req.GetLogin(), req.GetPassword())
	if err != nil {
		return nil, err
	}

	return &pb.LoginResponse{
		Token: resp.Token,
	}, nil
}

func (s *pbServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	err := s.auth.Register(ctx, req.GetLogin(), req.GetPassword())
	if err != nil {
		return nil, err
	}

	return new(pb.RegisterResponse), nil
}

func (s *pbServer) Button(ctx context.Context, req *pb.ButtonRequest) (*pb.ButtonResponse, error) {
	requestID, _ := ctx.Value(requestIDKey).(string)

	if req.GetDuration() <= 0 {
		err := fmt.Errorf("%w duration %d", errInvalidInputData, req.GetDuration())

		return nil, err
	}

	info, err := s.authInfo(ctx)
	if err != nil {
		return nil, err
	}

	kafkaData := dto.KafkaTaskData{
		UserID:   info.UserID,
		Chance:   req.GetChance(),
		Duration: req.GetDuration(),
	}

	err = s.kafkaTask.Write(ctx, requestID, kafkaData)
	if err != nil {
		return nil, err
	}

	return new(pb.ButtonResponse), nil
}

func (s *pbServer) List(ctx context.Context, _ *pb.NotificationListRequest) (*pb.NotificationListResponse, error) {
	info, err := s.authInfo(ctx)
	if err != nil {
		return nil, err
	}

	rawNotifications, err := s.notification.List(ctx, info.UserID)
	if err != nil {
		return nil, err
	}

	notifications := make([]*pb.NotificationData, len(rawNotifications.Notifications))
	for index, raw := range rawNotifications.Notifications {
		notifications[index] = &pb.NotificationData{
			Kind:    raw.Kind,
			Level:   raw.Level,
			Title:   raw.Title,
			Body:    raw.Body,
			Id:      raw.ID,
			Created: timestamppb.New(raw.Created),
		}
	}

	return &pb.NotificationListResponse{
		List: notifications,
	}, nil
}

func (s *pbServer) Read(ctx context.Context, req *pb.NotificationReadRequest) (*pb.NotificationReadResponse, error) {
	info, err := s.authInfo(ctx)
	if err != nil {
		return nil, err
	}

	if req.GetAll() {
		err = s.notification.Read(ctx, info.UserID, 0)
	} else {
		err = s.notification.Read(ctx, info.UserID, req.GetId())
	}

	if err != nil {
		return nil, err
	}

	return new(pb.NotificationReadResponse), nil
}

func (s *pbServer) Activity(ctx context.Context, _ *pb.ActivityRequest) (*pb.ActivityResponse, error) {
	info, err := s.authInfo(ctx)
	if err != nil {
		return nil, err
	}

	data, err := s.log.Activity(ctx, info.UserID)
	if err != nil {
		return nil, err
	}

	return &pb.ActivityResponse{
		RequestCount: data.RequestCount,
		LastRequest:  timestamppb.New(data.LastRequest),
	}, nil
}
