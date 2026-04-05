package cache

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func newTestMemoryCache() Cache {
	return NewMemoryCache()
}

func TestMemoryCache_SetAndGet(t *testing.T) {
	c := newTestMemoryCache()
	ctx := context.Background()

	key := "test:key"
	value := []byte("hello")

	err := c.Set(ctx, key, value, time.Minute)
	assert.NoError(t, err)

	got, err := c.Get(ctx, key)
	assert.NoError(t, err)
	assert.Equal(t, value, got)
}

func TestMemoryCache_GetMiss(t *testing.T) {
	c := newTestMemoryCache()
	ctx := context.Background()

	_, err := c.Get(ctx, "nonexistent")
	assert.ErrorIs(t, err, ErrCacheMiss)
}

func TestMemoryCache_Expiration(t *testing.T) {
	c := newTestMemoryCache()
	ctx := context.Background()

	key := "expire:key"
	err := c.Set(ctx, key, []byte("v"), 50*time.Millisecond)
	assert.NoError(t, err)

	// Should be present before expiry
	got, err := c.Get(ctx, key)
	assert.NoError(t, err)
	assert.Equal(t, []byte("v"), got)

	// Wait for expiry
	time.Sleep(100 * time.Millisecond)

	_, err = c.Get(ctx, key)
	assert.ErrorIs(t, err, ErrCacheMiss)
}

func TestMemoryCache_Delete(t *testing.T) {
	c := newTestMemoryCache()
	ctx := context.Background()

	key := "del:key"
	c.Set(ctx, key, []byte("data"), time.Minute)

	err := c.Delete(ctx, key)
	assert.NoError(t, err)

	_, err = c.Get(ctx, key)
	assert.ErrorIs(t, err, ErrCacheMiss)
}

func TestMemoryCache_Delete_NonExistent(t *testing.T) {
	c := newTestMemoryCache()
	ctx := context.Background()

	// Deleting a non-existent key should not error
	err := c.Delete(ctx, "no:such:key")
	assert.NoError(t, err)
}

func TestMemoryCache_Exists_Present(t *testing.T) {
	c := newTestMemoryCache()
	ctx := context.Background()

	c.Set(ctx, "exists:key", []byte("1"), time.Minute)

	exists, err := c.Exists(ctx, "exists:key")
	assert.NoError(t, err)
	assert.True(t, exists)
}

func TestMemoryCache_Exists_Missing(t *testing.T) {
	c := newTestMemoryCache()
	ctx := context.Background()

	exists, err := c.Exists(ctx, "missing:key")
	assert.NoError(t, err)
	assert.False(t, exists)
}

func TestMemoryCache_Exists_Expired(t *testing.T) {
	c := newTestMemoryCache()
	ctx := context.Background()

	c.Set(ctx, "exp:key", []byte("v"), 50*time.Millisecond)
	time.Sleep(100 * time.Millisecond)

	exists, err := c.Exists(ctx, "exp:key")
	assert.NoError(t, err)
	assert.False(t, exists)
}

func TestMemoryCache_FlushAll(t *testing.T) {
	c := newTestMemoryCache()
	ctx := context.Background()

	c.Set(ctx, "k1", []byte("v1"), time.Minute)
	c.Set(ctx, "k2", []byte("v2"), time.Minute)

	err := c.FlushAll(ctx)
	assert.NoError(t, err)

	_, err = c.Get(ctx, "k1")
	assert.ErrorIs(t, err, ErrCacheMiss)

	_, err = c.Get(ctx, "k2")
	assert.ErrorIs(t, err, ErrCacheMiss)
}

func TestMemoryCache_Overwrite(t *testing.T) {
	c := newTestMemoryCache()
	ctx := context.Background()

	c.Set(ctx, "key", []byte("original"), time.Minute)
	c.Set(ctx, "key", []byte("updated"), time.Minute)

	got, err := c.Get(ctx, "key")
	assert.NoError(t, err)
	assert.Equal(t, []byte("updated"), got)
}
