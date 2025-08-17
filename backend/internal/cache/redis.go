package cache

import (
    "context"
    "github.com/go-redis/redis/v8"
)

type RedisCache struct {
    client *redis.Client
}

func NewRedisCache(addr, password string, db int) *RedisCache {
    rdb := redis.NewClient(&redis.Options{
        Addr:     addr,
        Password: password,
        DB:       db,
    })
    return &RedisCache{client: rdb}
}

func (r *RedisCache) Ping(ctx context.Context) error {
    return r.client.Ping(ctx).Err()
}

func (r *RedisCache) Close() error {
    return r.client.Close()
}