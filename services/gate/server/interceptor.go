package server

import (
	"context"
	"log/slog"
	"time"

	"github.com/gbh007/buttoners/core/clients/gateclient"
	"github.com/gbh007/buttoners/core/dto"
	"github.com/gbh007/buttoners/core/logger"
	"github.com/gbh007/buttoners/core/metrics"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

type contextKey struct {
	name string
}

var (
	requestIDKey = &contextKey{"requestIDKey"}
	userInfoKey  = &contextKey{"userInfoKey"}
)

func (s *Server) logInterceptor(
	ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler,
) (resp interface{}, err error) {
	requestID := randomSHA256String()
	ctx = context.WithValue(ctx, requestIDKey, requestID)

	addr := "unknown"
	routeName := "unknown"

	p, ok := peer.FromContext(ctx)
	if ok {
		addr = p.Addr.String()
	}

	if info != nil {
		routeName = info.FullMethod
	}

	kafkaData := dto.KafkaLogData{
		Action:      routeName,
		Addr:        addr,
		RequestTime: time.Now().UTC(),
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		kafkaData.RealIP = fistOf(md.Get("X-Real-IP"))
		kafkaData.ForwardedFor = md.Get("X-Forwarded-For")
		kafkaData.SessionToken = fistOf(md.Get(gateclient.SessionHeader))
	}

	requestStart := time.Now()

	// Пытаемся идентифицировать пользователя.
	// Время на идентификацию тоже считается частью запроса.
	if kafkaData.SessionToken != "" {
		userInfo, err := s.authInfoRaw(ctx, kafkaData.SessionToken)
		if err != nil {
			logger.LogWithMeta(s.logger, ctx, slog.LevelError, "get user info", "error", err.Error())
		} else {
			ctx = context.WithValue(ctx, userInfoKey, userInfo)
			kafkaData.UserID = userInfo.UserID
		}
	}

	// Выполняем сам запрос
	resp, err = handler(ctx, req)

	if err != nil {
		kafkaData.ErrorText = err.Error()
	}

	_ = s.kafkaLog.Write(ctx, requestID, kafkaData)

	metrics.LogRequest(routeName, err == nil, time.Since(requestStart))

	return
}
