package cache

import (
	"context"
	"time"
)

type Cache interface {
	Get(ctx context.Context, key string) ([]byte, error)                        // Return []byte, not string
	Set(ctx context.Context, key string, value []byte, ttl time.Duration) error // Accept []byte, not string
	Delete(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) (bool, error)
	FlushAll(ctx context.Context) error
}
