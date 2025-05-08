package server

import (
	"context"
	"log"
	"strings"

	"github.com/gbh007/buttoners/core/redis"
	"github.com/gbh007/buttoners/services/auth/internal/pb"
	"github.com/gbh007/buttoners/services/auth/internal/storage"
	"github.com/gbh007/buttoners/services/gate/dto"
)

type authServer struct {
	pb.UnimplementedAuthServer

	db    *storage.Database
	redis *redis.Client[dto.UserInfo]
}

func (s *authServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	login := strings.ToLower(req.GetLogin())
	pass := req.GetPassword()

	token, err := s.createSession(ctx, login, pass)
	if err != nil {
		return nil, err
	}

	// Кеш в редисе мог сеттится в этом месте

	return &pb.LoginResponse{
		Token: token,
	}, nil
}

func (s *authServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	login := strings.ToLower(req.GetLogin())
	pass := req.GetPassword()

	_, err := s.createUser(ctx, login, pass)
	if err != nil {
		return nil, err
	}

	return new(pb.RegisterResponse), nil
}

func (s *authServer) Logout(ctx context.Context, req *pb.LogoutRequest) (*pb.LogoutResponse, error) {
	err := s.deleteSession(ctx, req.GetToken())
	if err != nil {
		return nil, err
	}

	// Инвалидация кеша
	err = s.redis.Del(ctx, req.GetToken())
	if err != nil {
		log.Println(err)
	}

	return new(pb.LogoutResponse), nil
}

func (s *authServer) Info(ctx context.Context, req *pb.InfoRequest) (*pb.InfoResponse, error) {
	user, err := s.getUser(ctx, req.GetToken())
	if err != nil {
		return nil, err
	}

	return &pb.InfoResponse{
		UserID: user.ID,
	}, nil
}
