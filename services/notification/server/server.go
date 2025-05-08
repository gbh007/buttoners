package server

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/gbh007/buttoners/services/notification/internal/pb"
	"github.com/gbh007/buttoners/services/notification/internal/storage"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

var errMissingInputData = errors.New("missing")

type pbServer struct {
	pb.UnimplementedNotificationServer

	db *storage.Database
}

func (s *pbServer) New(ctx context.Context, req *pb.NewRequest) (*pb.NewResponse, error) {
	userID := req.GetUserID()
	if userID == 0 {
		return nil, fmt.Errorf("%w: user id", errMissingInputData)
	}

	if req.GetData() == nil {
		return nil, fmt.Errorf("%w: notification", errMissingInputData)
	}

	err := s.db.CreateNotification(ctx, &storage.Notification{
		UserID: userID,
		Kind:   req.GetData().GetKind(),
		Level:  req.GetData().GetLevel(),
		Title:  req.GetData().GetTitle(),
		Body: sql.NullString{
			String: req.GetData().GetBody(),
			Valid:  req.GetData().GetBody() != "",
		},
		Created: req.GetData().GetCreated().AsTime(),
	})
	if err != nil {
		return nil, err
	}

	return &pb.NewResponse{}, nil
}

func (s *pbServer) List(ctx context.Context, req *pb.ListRequest) (*pb.ListResponse, error) {
	userID := req.GetUserID()
	if userID == 0 {
		return nil, fmt.Errorf("%w: user id", errMissingInputData)
	}

	rawNotifications, err := s.db.GetNotificationsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	notifications := make([]*pb.NotificationData, len(rawNotifications))

	for index, raw := range rawNotifications {
		notifications[index] = &pb.NotificationData{
			Kind:    raw.Kind,
			Level:   raw.Level,
			Title:   raw.Title,
			Body:    raw.Body.String,
			Id:      raw.ID,
			Created: timestamppb.New(raw.Created),
		}
	}

	return &pb.ListResponse{
		List: notifications,
	}, nil
}

func (s *pbServer) Read(ctx context.Context, req *pb.ReadRequest) (*pb.ReadResponse, error) {
	id := req.GetId()
	if id == 0 {
		return nil, fmt.Errorf("%w: id", errMissingInputData)
	}

	err := s.db.MarkReadByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &pb.ReadResponse{}, nil
}

func (s *pbServer) ReadAll(ctx context.Context, req *pb.ReadAllRequest) (*pb.ReadAllResponse, error) {
	userID := req.GetUserID()
	if userID == 0 {
		return nil, fmt.Errorf("%w: user id", errMissingInputData)
	}

	err := s.db.MarkReadByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &pb.ReadAllResponse{}, nil
}
