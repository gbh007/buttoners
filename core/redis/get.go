package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/go-redis/redis"
)

func (c *Client[T]) Get(ctx context.Context, key string) (*T, error) {
	if c.client == nil {
		return nil, fmt.Errorf("%w: Get: %w", ErrRedisClient, ErrClientNotInitialized)
	}

	raw, err := c.client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, fmt.Errorf("%w: Get: %w", ErrRedisClient, ErrNotExists)
		}

		return nil, fmt.Errorf("%w: Get: %w", ErrRedisClient, err)
	}

	value := new(T)

	err = json.Unmarshal([]byte(raw), value)
	if err != nil {
		return nil, fmt.Errorf("%w: Get: %w", ErrRedisClient, err)
	}

	return value, nil
}
