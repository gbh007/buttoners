package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

func (c *Client[T]) Set(ctx context.Context, key string, value T, ttl time.Duration) error {
	if c.client == nil {
		return fmt.Errorf("%w: Set: %w", ErrRedisClient, ErrClientNotInitialized)
	}

	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("%w: Set: %w", ErrRedisClient, err)
	}

	err = c.client.Set(ctx, key, string(data), ttl).Err()
	if err != nil {
		return fmt.Errorf("%w: Set: %w", ErrRedisClient, err)
	}

	return nil
}
