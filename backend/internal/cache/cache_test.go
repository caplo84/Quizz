package cache

import (
    "context"
    "testing"
    "time"

    "github.com/alicebob/miniredis/v2"
    "github.com/go-redis/redis/v8"
    "github.com/stretchr/testify/assert"
)

func setupTestRedis(t *testing.T) *RedisCache {
    // Start miniredis server
    s, err := miniredis.Run()
    if err != nil {
        t.Fatalf("Failed to start miniredis: %v", err)
    }

    // Create Redis client
    client := redis.NewClient(&redis.Options{
        Addr: s.Addr(),
    })

    return &RedisCache{client: client}
}

func TestRedisCache_SetAndGet(t *testing.T) {
    cache := setupTestRedis(t)
    ctx := context.Background()

    key := "test:key"
    value := "test_value"
    ttl := time.Minute

    // Test Set
    err := cache.Set(ctx, key, value, ttl)
    assert.NoError(t, err)

    // Test Get
    result, err := cache.Get(ctx, key)
    assert.NoError(t, err)
    assert.Equal(t, value, result)
}

func TestRedisCache_GetNonExistent(t *testing.T) {
    cache := setupTestRedis(t)
    ctx := context.Background()

    // Test Get for non-existent key
    result, err := cache.Get(ctx, "non:existent")
    assert.Error(t, err)
    assert.Equal(t, "", result)
    assert.Equal(t, redis.Nil, err)
}

func TestRedisCache_Delete(t *testing.T) {
    cache := setupTestRedis(t)
    ctx := context.Background()

    key := "test:delete"
    value := "delete_me"

    // Set a value
    err := cache.Set(ctx, key, value, time.Minute)
    assert.NoError(t, err)

    // Verify it exists
    result, err := cache.Get(ctx, key)
    assert.NoError(t, err)
    assert.Equal(t, value, result)

    // Delete it
    err = cache.Delete(ctx, key)
    assert.NoError(t, err)

    // Verify it's gone
    _, err = cache.Get(ctx, key)
    assert.Error(t, err)
    assert.Equal(t, redis.Nil, err)
}