package user

import (
	"context"
	"fmt"

	"github.com/gbh007/buttoners/core/clients/authclient"
	"github.com/gbh007/buttoners/services/legacy/internal/domain"
)

type Service struct {
	randomizer *domain.Randomizer
	authClient *authclient.Client
}

func New(authClient *authclient.Client) *Service {
	return &Service{
		authClient: authClient,
		randomizer: domain.NewRandomizer(),
	}
}

func (s *Service) CreateUser(ctx context.Context, login string, pass string) (domain.User, error) {
	u := domain.User{
		Name: login,
	}

	err := s.authClient.Register(ctx, login, pass)
	if err != nil {
		return domain.User{}, fmt.Errorf("register: %w", err)
	}

	session, err := s.authClient.Login(ctx, login, pass)
	if err != nil {
		return domain.User{}, fmt.Errorf("login: %w", err)
	}

	u.Token = session.Token

	info, err := s.authClient.Info(ctx, u.Token)
	if err != nil {
		return domain.User{}, fmt.Errorf("info: %w", err)
	}

	u.ID = int(info.UserID)

	return u, nil
}

func (s *Service) GetUser(ctx context.Context, token string) (domain.User, error) {
	u := domain.User{
		Token: token,
	}

	info, err := s.authClient.Info(ctx, u.Token)
	if err != nil {
		return domain.User{}, fmt.Errorf("info: %w", err)
	}

	u.ID = int(info.UserID)
	u.Name = info.Login

	return u, nil
}
