package cache

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"
)

// CacheMiddleware provides caching functionality for MCP handlers.
type CacheMiddleware struct {
	cache      Cache
	ttl        time.Duration
	keyPrefix  string
	skipCache  func(toolName string, args interface{}) bool
}

// NewCacheMiddleware creates a new cache middleware.
func NewCacheMiddleware(cache Cache, ttl time.Duration, keyPrefix string) *CacheMiddleware {
	return &CacheMiddleware{
		cache:     cache,
		ttl:       ttl,
		keyPrefix: keyPrefix,
		skipCache: func(toolName string, args interface{}) bool {
			// By default, don't skip any tools
			return false
		},
	}
}

// WithSkipCache sets a function to determine whether to skip caching for specific tools.
func (cm *CacheMiddleware) WithSkipCache(skip func(toolName string, args interface{}) bool) *CacheMiddleware {
	cm.skipCache = skip
	return cm
}

// generateCacheKey generates a cache key based on tool name and arguments.
func (cm *CacheMiddleware) generateCacheKey(toolName string, args interface{}) (string, error) {
	// Serialize arguments to JSON for consistent key generation
	argsJSON, err := json.Marshal(args)
	if err != nil {
		return "", err
	}
	
	// Create MD5 hash of arguments for shorter keys
	hash := md5.Sum(argsJSON)
	hashStr := hex.EncodeToString(hash[:])
	
	// Build cache key
	kb := NewKeyBuilder(cm.keyPrefix)
	kb.Add(toolName).Add(hashStr)
	
	return kb.Build(), nil
}

// Wrap wraps a handler function with caching logic.
func (cm *CacheMiddleware) Wrap(toolName string, handler func(ctx context.Context, args interface{}) (interface{}, error)) func(context.Context, interface{}) (interface{}, error) {
	return func(ctx context.Context, args interface{}) (interface{}, error) {
		// Check if we should skip caching for this tool
		if cm.skipCache(toolName, args) {
			return handler(ctx, args)
		}
		
		// Generate cache key
		cacheKey, err := cm.generateCacheKey(toolName, args)
		if err != nil {
			// If we can't generate a key, proceed without caching
			return handler(ctx, args)
		}
		
		// Try to get from cache
		cachedResult, err := cm.cache.Get(ctx, cacheKey)
		if err == nil && cachedResult != nil {
			// Cache hit
			return cachedResult, nil
		}
		
		// Cache miss - execute handler
		result, err := handler(ctx, args)
		if err != nil {
			return nil, err
		}
		
		// Store in cache (ignore errors - caching is best-effort)
		_ = cm.cache.Set(ctx, cacheKey, result, cm.ttl)
		
		return result, nil
	}
}

// InvalidatePattern invalidates cache entries matching a pattern.
// This is useful when data is updated and cached entries need to be cleared.
func (cm *CacheMiddleware) InvalidatePattern(ctx context.Context, pattern string) error {
	// This would require a more sophisticated cache implementation
	// that supports pattern-based deletion (like Redis SCAN + DEL)
	// For now, we can only clear all cache
	return cm.cache.Clear(ctx)
}

// CacheConfig provides configuration for cache middleware.
type CacheConfig struct {
	Enabled        bool          `yaml:"enabled" json:"enabled"`
	Type           string        `yaml:"type" json:"type"`                     // "memory" or "redis"
	TTL            time.Duration `yaml:"ttl" json:"ttl"`                       // Default TTL
	MaxSize        int64         `yaml:"max_size" json:"max_size"`             // For memory cache
	RedisAddr      string        `yaml:"redis_addr" json:"redis_addr"`         // For Redis cache
	RedisPassword  string        `yaml:"redis_password" json:"redis_password"` // For Redis cache
	RedisDB        int           `yaml:"redis_db" json:"redis_db"`             // For Redis cache
	KeyPrefix      string        `yaml:"key_prefix" json:"key_prefix"`         // Cache key prefix
}

// DefaultCacheConfig returns default cache configuration.
func DefaultCacheConfig() *CacheConfig {
	return &CacheConfig{
		Enabled:   false,
		Type:      "memory",
		TTL:       5 * time.Minute,
		MaxSize:   100,
		KeyPrefix: "firefly",
	}
}

// BuildCache creates a cache instance based on configuration.
func BuildCache(config *CacheConfig) (Cache, error) {
	if !config.Enabled {
		return nil, nil
	}
	
	switch config.Type {
	case "memory":
		return NewMemoryCache(
			WithMaxSize(config.MaxSize),
			WithDefaultTTL(config.TTL),
		), nil
		
	case "redis":
		// In a real implementation, you would create a real Redis client here
		// For now, we'll return an error
		return nil, fmt.Errorf("redis cache not yet implemented with real client")
		
	default:
		return nil, fmt.Errorf("unknown cache type: %s", config.Type)
	}
}