package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync/atomic"
	"time"
)

// RedisClient defines the interface for Redis operations.
// This allows for easy mocking in tests.
type RedisClient interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Del(ctx context.Context, keys ...string) error
	Exists(ctx context.Context, keys ...string) (int64, error)
	TTL(ctx context.Context, key string) (time.Duration, error)
	FlushAll(ctx context.Context) error
	Close() error
}

// RedisCache is a Redis-based cache implementation.
type RedisCache struct {
	client     RedisClient
	stats      Stats
	prefix     string
	defaultTTL time.Duration
}

// NewRedisCache creates a new Redis cache.
func NewRedisCache(client RedisClient, prefix string, opts ...Option) *RedisCache {
	rc := &RedisCache{
		client: client,
		prefix: prefix,
	}
	
	for _, opt := range opts {
		opt(rc)
	}
	
	return rc
}

// prefixKey adds the prefix to a key.
func (rc *RedisCache) prefixKey(key string) string {
	if rc.prefix == "" {
		return key
	}
	return fmt.Sprintf("%s:%s", rc.prefix, key)
}

// Get retrieves a value from the cache.
func (rc *RedisCache) Get(ctx context.Context, key string) (interface{}, error) {
	prefixedKey := rc.prefixKey(key)
	
	val, err := rc.client.Get(ctx, prefixedKey)
	if err != nil {
		// Check if the error message contains "redis: nil" which indicates key not found
		if err.Error() == "redis: nil" {
			atomic.AddInt64(&rc.stats.Misses, 1)
			return nil, ErrCacheMiss
		}
		return nil, fmt.Errorf("redis get error: %w", err)
	}
	
	// Try to unmarshal as JSON
	var result interface{}
	if err := json.Unmarshal([]byte(val), &result); err != nil {
		// If unmarshal fails, return as string
		result = val
	}
	
	atomic.AddInt64(&rc.stats.Hits, 1)
	return result, nil
}

// Set stores a value in the cache with the given TTL.
func (rc *RedisCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	prefixedKey := rc.prefixKey(key)
	
	if ttl == 0 {
		ttl = rc.defaultTTL
	}
	
	// Marshal value to JSON
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}
	
	err = rc.client.Set(ctx, prefixedKey, data, ttl)
	if err != nil {
		return fmt.Errorf("redis set error: %w", err)
	}
	
	atomic.AddInt64(&rc.stats.Sets, 1)
	return nil
}

// Delete removes a value from the cache.
func (rc *RedisCache) Delete(ctx context.Context, key string) error {
	prefixedKey := rc.prefixKey(key)
	
	err := rc.client.Del(ctx, prefixedKey)
	if err != nil {
		return fmt.Errorf("redis delete error: %w", err)
	}
	
	atomic.AddInt64(&rc.stats.Deletes, 1)
	return nil
}

// Clear removes all values from the cache.
// WARNING: This will clear ALL keys in Redis if no prefix is set!
func (rc *RedisCache) Clear(ctx context.Context) error {
	// If we have a prefix, we should scan and delete keys with that prefix
	// For simplicity, we'll just flush the entire database
	// In production, you'd want to implement a scan-and-delete approach
	if rc.prefix == "" {
		return rc.client.FlushAll(ctx)
	}
	
	// TODO: Implement prefix-based clearing using SCAN
	return errors.New("prefix-based clearing not yet implemented")
}

// Exists checks if a key exists in the cache.
func (rc *RedisCache) Exists(ctx context.Context, key string) (bool, error) {
	prefixedKey := rc.prefixKey(key)
	
	count, err := rc.client.Exists(ctx, prefixedKey)
	if err != nil {
		return false, fmt.Errorf("redis exists error: %w", err)
	}
	
	return count > 0, nil
}

// TTL returns the remaining TTL for a key.
func (rc *RedisCache) TTL(ctx context.Context, key string) (time.Duration, error) {
	prefixedKey := rc.prefixKey(key)
	
	ttl, err := rc.client.TTL(ctx, prefixedKey)
	if err != nil {
		return 0, fmt.Errorf("redis ttl error: %w", err)
	}
	
	if ttl < 0 {
		return 0, nil // Key doesn't exist or has no expiration
	}
	
	return ttl, nil
}

// Stats returns cache statistics.
func (rc *RedisCache) Stats() Stats {
	return rc.stats
}

// Close closes the Redis connection.
func (rc *RedisCache) Close() error {
	return rc.client.Close()
}