package redis

import (
	"context"
	"fmt"

	"github.com/gbh007/buttoners/core/observability"
	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
)

func (c *Client[T]) Connect(ctx context.Context, rh *observability.RedisHook) (err error) {
	// Правильно сделать полную настройку, но в данном проекте это не требуется
	c.client = redis.NewClient(&redis.Options{
		Addr:     c.addr,
		Password: "",
		DB:       0,
	})

	err = c.client.Ping(ctx).Err()
	if err != nil {
		return fmt.Errorf("%w: Connect: %w", ErrRedisClient, err)
	}

	err = redisotel.InstrumentTracing(c.client)
	if err != nil {
		return fmt.Errorf("%w: Tracing: %w", ErrRedisClient, err)
	}

	c.client.AddHook(rh)

	return nil
}

func (c *Client[T]) Close() error {
	err := c.client.Close()
	if err != nil {
		return fmt.Errorf("%w: Close: %w", ErrRedisClient, err)
	}

	return nil
}
