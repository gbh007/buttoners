package redis

import "errors"

var (
	ErrRedisClient = errors.New("redis client")

	ErrClientNotInitialized = errors.New("client not initialized")
	ErrNotExists            = errors.New("not exists")
)
