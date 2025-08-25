package cache

import (
	"context"
	"sync"
	"time"
)

// memoryCache implements Cache interface using in-memory storage
type memoryCache struct {
	data map[string]cacheItem
	mu   sync.RWMutex
}

type cacheItem struct {
	value     []byte
	expiresAt time.Time
}

// NewMemoryCache creates a new in-memory cache instance
func NewMemoryCache() Cache {
	cache := &memoryCache{
		data: make(map[string]cacheItem),
	}

	// Start cleanup goroutine
	go cache.cleanup()

	return cache
}

func (m *memoryCache) Get(ctx context.Context, key string) ([]byte, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	item, exists := m.data[key]
	if !exists {
		return nil, ErrCacheMiss
	}

	if time.Now().After(item.expiresAt) {
		delete(m.data, key)
		return nil, ErrCacheMiss
	}

	return item.value, nil
}

func (m *memoryCache) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.data[key] = cacheItem{
		value:     value,
		expiresAt: time.Now().Add(ttl),
	}

	return nil
}

func (m *memoryCache) Delete(ctx context.Context, key string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.data, key)
	return nil
}

func (m *memoryCache) Exists(ctx context.Context, key string) (bool, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	item, exists := m.data[key]
	if !exists {
		return false, nil
	}

	if time.Now().After(item.expiresAt) {
		delete(m.data, key)
		return false, nil
	}

	return true, nil
}

func (m *memoryCache) FlushAll(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.data = make(map[string]cacheItem)
	return nil
}

// cleanup removes expired items periodically
func (m *memoryCache) cleanup() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			m.mu.Lock()
			now := time.Now()
			for key, item := range m.data {
				if now.After(item.expiresAt) {
					delete(m.data, key)
				}
			}
			m.mu.Unlock()
		}
	}
}
