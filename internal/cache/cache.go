package cache

import (
	"context"
	"time"
)

// Cache defines the caching contract for the BFF layer.
// Implementations should be safe for concurrent use.
type Cache interface {
	// Get retrieves a cached value by key. Returns nil, nil if the key does not exist.
	Get(ctx context.Context, key string) ([]byte, error)

	// Set stores a value in the cache with the given TTL.
	Set(ctx context.Context, key string, value []byte, ttl time.Duration) error

	// Delete removes one or more specific keys from the cache.
	Delete(ctx context.Context, keys ...string) error

	// DeleteByPattern removes all keys matching the given glob pattern (e.g. "wallet:abc:*").
	DeleteByPattern(ctx context.Context, pattern string) error

	// Close gracefully shuts down the cache connection.
	Close() error
}
