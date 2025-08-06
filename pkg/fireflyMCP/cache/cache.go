// Package cache provides caching functionality for the Firefly III MCP server.
package cache

import (
	"context"
	"errors"
	"time"
)

// ErrCacheMiss is returned when a key is not found in the cache.
var ErrCacheMiss = errors.New("cache miss")

// Cache defines the interface for cache implementations.
type Cache interface {
	// Get retrieves a value from the cache.
	Get(ctx context.Context, key string) (interface{}, error)
	
	// Set stores a value in the cache with the given TTL.
	// If TTL is 0, the value never expires.
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	
	// Delete removes a value from the cache.
	Delete(ctx context.Context, key string) error
	
	// Clear removes all values from the cache.
	Clear(ctx context.Context) error
	
	// Exists checks if a key exists in the cache.
	Exists(ctx context.Context, key string) (bool, error)
	
	// TTL returns the remaining TTL for a key.
	// Returns 0 if the key doesn't exist or has no expiration.
	TTL(ctx context.Context, key string) (time.Duration, error)
	
	// Stats returns cache statistics.
	Stats() Stats
}

// Stats contains cache statistics.
type Stats struct {
	Hits       int64 `json:"hits"`
	Misses     int64 `json:"misses"`
	Sets       int64 `json:"sets"`
	Deletes    int64 `json:"deletes"`
	Evictions  int64 `json:"evictions"`
	Size       int64 `json:"size"`
	MaxSize    int64 `json:"max_size,omitempty"`
}

// HitRate calculates the cache hit rate.
func (s Stats) HitRate() float64 {
	total := s.Hits + s.Misses
	if total == 0 {
		return 0
	}
	return float64(s.Hits) / float64(total)
}

// KeyBuilder helps construct consistent cache keys.
type KeyBuilder struct {
	prefix string
	parts  []string
}

// NewKeyBuilder creates a new KeyBuilder with an optional prefix.
func NewKeyBuilder(prefix string) *KeyBuilder {
	return &KeyBuilder{
		prefix: prefix,
		parts:  make([]string, 0),
	}
}

// Add adds a part to the key.
func (kb *KeyBuilder) Add(part string) *KeyBuilder {
	kb.parts = append(kb.parts, part)
	return kb
}

// Build constructs the final cache key.
func (kb *KeyBuilder) Build() string {
	if kb.prefix == "" && len(kb.parts) == 0 {
		return ""
	}
	
	result := kb.prefix
	for _, part := range kb.parts {
		if result != "" {
			result += ":"
		}
		result += part
	}
	return result
}

// Option is a function that configures a cache.
type Option func(interface{})

// WithMaxSize sets the maximum cache size (for memory cache).
func WithMaxSize(size int64) Option {
	return func(c interface{}) {
		if mc, ok := c.(*MemoryCache); ok {
			mc.maxSize = size
		}
	}
}

// WithCleanupInterval sets the cleanup interval for expired entries (for memory cache).
func WithCleanupInterval(interval time.Duration) Option {
	return func(c interface{}) {
		if mc, ok := c.(*MemoryCache); ok {
			mc.cleanupInterval = interval
		}
	}
}

// WithDefaultTTL sets the default TTL for cache entries.
func WithDefaultTTL(ttl time.Duration) Option {
	return func(c interface{}) {
		switch cache := c.(type) {
		case *MemoryCache:
			cache.defaultTTL = ttl
		case *RedisCache:
			cache.defaultTTL = ttl
		}
	}
}