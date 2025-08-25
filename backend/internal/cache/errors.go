package cache

import "errors"

var (
	ErrCacheMiss = errors.New("cache miss")
	ErrCacheSet  = errors.New("failed to set cache")
)
