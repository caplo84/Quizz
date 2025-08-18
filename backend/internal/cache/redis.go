package cache

import (
    "context"
    "time"

    "github.com/go-redis/redis/v8"
)

// redisCache implements Cache interface using Redis
type redisCache struct {
    client *redis.Client
}

// NewRedisCache creates a new Redis cache instance
func NewRedisCache(client *redis.Client) Cache {
    return &redisCache{
        client: client,
    }
}

func (r *redisCache) Get(ctx context.Context, key string) ([]byte, error) {
    val, err := r.client.Get(ctx, key).Result()
    if err != nil {
        return nil, err
    }
    return []byte(val), nil
}

func (r *redisCache) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
    return r.client.Set(ctx, key, string(value), ttl).Err()
}

func (r *redisCache) Delete(ctx context.Context, key string) error {
    return r.client.Del(ctx, key).Err()
}

func (r *redisCache) Exists(ctx context.Context, key string) (bool, error) {
    count, err := r.client.Exists(ctx, key).Result()
    return count > 0, err
}

func (r *redisCache) FlushAll(ctx context.Context) error {
    return r.client.FlushAll(ctx).Err()
}