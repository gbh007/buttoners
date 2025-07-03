package server

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/gbh007/buttoners/core/clients/authclient"
	"github.com/gbh007/buttoners/core/dto"
	"github.com/gbh007/buttoners/core/logger"
)

const cacheTTL = time.Minute * 5

var errUnauthorized = errors.New("unauthorized")

func (s *pbServer) authInfo(ctx context.Context) (*authclient.InfoResponse, error) {
	info, ok := ctx.Value(userInfoKey).(*authclient.InfoResponse)
	if !ok {
		return nil, errUnauthorized
	}

	return info, nil
}

func (s *pbServer) authInfoRaw(ctx context.Context, token string) (*authclient.InfoResponse, error) {
	redisStart := time.Now()

	redisData, err := s.redis.Get(ctx, token)

	redisFinish := time.Now()
	registerCacheHandle("redis", redisFinish.Sub(redisStart))

	if err != nil {
		// Ошибка отсутствия значения также логируется для отладки
		logger.LogWithMeta(s.logger, ctx, slog.LevelError, "get session from cache", "error", err.Error())
	} else {
		return &authclient.InfoResponse{
			UserID: redisData.ID,
		}, nil
	}

	authStart := time.Now()

	info, err := s.auth.Info(ctx, token)

	authFinish := time.Now()
	registerCacheHandle("auth", authFinish.Sub(authStart))

	if err != nil {
		return nil, err
	}

	// В данном случае кешер сеттится специально здесь, а не в сервисе авторизации
	err = s.redis.Set(ctx, token, dto.UserInfo{ID: info.UserID}, cacheTTL)
	if err != nil {
		logger.LogWithMeta(s.logger, ctx, slog.LevelError, "set session to cache", "error", err.Error())
	}

	return &info, nil
}
