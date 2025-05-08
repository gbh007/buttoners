package server

import (
	"context"

	"github.com/gbh007/buttoners/services/log/internal/pb"
	"github.com/gbh007/buttoners/services/log/internal/storage"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type pbServer struct {
	pb.UnimplementedLogServer

	db *storage.Database
}

func (s *pbServer) Activity(ctx context.Context, req *pb.ActivityRequest) (*pb.ActivityResponse, error) {
	count, last, err := s.db.SelectCompressedUserLogByUserID(ctx, req.GetUserID())
	if err != nil {
		return nil, err
	}

	return &pb.ActivityResponse{
		Data: &pb.LogData{
			RequestCount: count,
			LastRequest:  timestamppb.New(last),
		},
	}, nil
}
