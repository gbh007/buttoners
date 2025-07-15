package server

import (
	"context"
	"fmt"

	"github.com/gbh007/buttoners/core/clients/gateclient/gen/pb"
	"github.com/gbh007/buttoners/core/dto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Server) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	resp, err := s.auth.Login(ctx, req.GetLogin(), req.GetPassword())
	if err != nil {
		return nil, err
	}

	return &pb.LoginResponse{
		Token: resp.Token,
	}, nil
}

func (s *Server) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	err := s.auth.Register(ctx, req.GetLogin(), req.GetPassword())
	if err != nil {
		return nil, err
	}

	return new(pb.RegisterResponse), nil
}

func (s *Server) Button(ctx context.Context, req *pb.ButtonRequest) (*pb.ButtonResponse, error) {
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

func (s *Server) List(ctx context.Context, _ *pb.NotificationListRequest) (*pb.NotificationListResponse, error) {
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

func (s *Server) Read(ctx context.Context, req *pb.NotificationReadRequest) (*pb.NotificationReadResponse, error) {
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

func (s *Server) Activity(ctx context.Context, _ *pb.ActivityRequest) (*pb.ActivityResponse, error) {
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
