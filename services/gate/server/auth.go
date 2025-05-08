package server

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/gbh007/buttoners/services/auth/client"
	"github.com/gbh007/buttoners/services/gate/dto"
)

const cacheTTL = time.Minute * 5

var errUnauthorized = errors.New("unauthorized")

func (s *pbServer) authInfo(ctx context.Context) (*client.UserInfo, error) {
	info, ok := ctx.Value(userInfoKey).(*client.UserInfo)
	if !ok {
		return nil, errUnauthorized
	}

	return info, nil
}

func (s *pbServer) authInfoRaw(ctx context.Context, token string) (*client.UserInfo, error) {
	redisStart := time.Now()

	redisData, err := s.redis.Get(ctx, token)

	redisFinish := time.Now()
	registerCacheHandle("redis", redisFinish.Sub(redisStart))

	if err != nil {
		// Ошибка отсутствия значения также логируется для отладки
		log.Printf("%s error from redis: %s\n", token, err.Error())
	} else {
		return &client.UserInfo{
			ID:    redisData.ID,
			Token: token,
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
	err = s.redis.Set(ctx, token, dto.UserInfo{ID: info.ID}, cacheTTL)
	if err != nil {
		log.Println(err)
	}

	return info, nil
}
