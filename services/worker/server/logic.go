package server

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"runtime"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

var (
	errFailed           = errors.New("failed")
	errInvalidInputData = errors.New("invalid")
)

func max(a, b int64) int64 {
	if a > b {
		return a
	}

	return b
}

func (s *Server) someBusinessLogic(ctx context.Context, duration, failChance int64) (int64, string, error) {
	_, span := s.tracer.Start(ctx, "business-logic")
	defer span.End()

	span.SetAttributes(
		attribute.Int64("duration", duration),
		attribute.Int64("fail_chance", failChance),
	)

	if duration < 1 || duration > 60 {
		span.SetStatus(codes.Error, "invalid duration")

		return 0, "", fmt.Errorf("%w duration %d", errInvalidInputData, duration)
	}

	if failChance < 0 || failChance > 100 {
		span.SetStatus(codes.Error, "invalid fail chance")

		return 0, "", fmt.Errorf("%w fail chance %d", errInvalidInputData, failChance)
	}

	totalSleep := rand.Int63n(max(duration*80/100, 1)) + 1 + duration*60/100

	runtime.Gosched()
	// Имитация выполнения
	time.Sleep(time.Second * time.Duration(totalSleep))

	// Результат выполнения [1, 100]
	result := rand.Int63n(100) + 1

	span.SetAttributes(
		attribute.Int64("work_time", totalSleep),
		attribute.Int64("result_chance", result),
	)

	if failChance >= result {
		span.SetStatus(codes.Error, "unsuccess result")

		return result, "", fmt.Errorf("%w: result %d ", errFailed, result)
	}

	return result, fmt.Sprintf(
		"Успешно обработано\nПродолжительность: %d секунд\nШанс успеха: %d%%",
		totalSleep, 100-failChance,
	), nil
}
