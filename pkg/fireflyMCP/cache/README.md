# Cache Package

The cache package provides a flexible caching layer for the Firefly III MCP server, supporting both in-memory and Redis-based caching with metrics tracking.

## Features

- **Multiple Implementations**: In-memory and Redis cache implementations
- **TTL Support**: Time-to-live for automatic expiration
- **Metrics Tracking**: Hit/miss rates, evictions, and size tracking
- **Middleware Integration**: Easy integration with MCP handlers
- **Type-Safe**: Works with any Go type through interface{}
- **Thread-Safe**: All operations are concurrent-safe

## Architecture

### Core Interface

```go
type Cache interface {
    Get(ctx context.Context, key string) (interface{}, error)
    Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
    Delete(ctx context.Context, key string) error
    Clear(ctx context.Context) error
    Exists(ctx context.Context, key string) (bool, error)
    TTL(ctx context.Context, key string) (time.Duration, error)
    Stats() Stats
}
```

### Implementations

#### Memory Cache
- Fast, in-process caching
- LRU eviction when max size is reached
- Automatic cleanup of expired entries
- Configurable cleanup interval

#### Redis Cache
- Distributed caching across multiple servers
- Persistent storage option
- Built-in TTL support
- JSON serialization for complex types

## Usage

### Basic Usage

```go
// Create an in-memory cache
cache := NewMemoryCache(
    WithMaxSize(1000),
    WithDefaultTTL(5 * time.Minute),
    WithCleanupInterval(1 * time.Minute),
)
defer cache.Close()

// Set a value
err := cache.Set(ctx, "user:123", userData, 10*time.Minute)

// Get a value
value, err := cache.Get(ctx, "user:123")
if err == ErrCacheMiss {
    // Handle cache miss
}

// Check existence
exists, err := cache.Exists(ctx, "user:123")

// Delete a value
err = cache.Delete(ctx, "user:123")

// Get statistics
stats := cache.Stats()
fmt.Printf("Hit rate: %.2f%%\n", stats.HitRate()*100)
```

### With Middleware

```go
// Create cache middleware
middleware := NewCacheMiddleware(cache, 5*time.Minute, "firefly")

// Configure which tools to skip
middleware.WithSkipCache(func(toolName string, args interface{}) bool {
    // Skip caching for write operations
    return strings.HasPrefix(toolName, "store_") || 
           strings.HasPrefix(toolName, "update_")
})

// Wrap a handler
cachedHandler := middleware.Wrap("list_accounts", originalHandler)
```

### Key Builder

```go
// Build consistent cache keys
kb := NewKeyBuilder("firefly")
key := kb.Add("accounts").Add("user").Add("123").Build()
// Result: "firefly:accounts:user:123"
```

## Configuration

### YAML Configuration

```yaml
cache:
  enabled: true
  type: memory          # or "redis"
  ttl: 5m
  max_size: 1000       # for memory cache
  key_prefix: firefly
  
  # Redis-specific settings
  redis_addr: localhost:6379
  redis_password: ""
  redis_db: 0
```

### Programmatic Configuration

```go
config := &CacheConfig{
    Enabled:   true,
    Type:      "memory",
    TTL:       5 * time.Minute,
    MaxSize:   1000,
    KeyPrefix: "firefly",
}

cache, err := BuildCache(config)
```

## Performance Considerations

### Memory Cache
- **Pros**: 
  - Very fast (nanosecond access)
  - No network overhead
  - Simple deployment
- **Cons**:
  - Limited to single process
  - Data lost on restart
  - Memory constraints

### Redis Cache
- **Pros**:
  - Distributed across servers
  - Persistent storage
  - Large capacity
- **Cons**:
  - Network latency
  - Serialization overhead
  - Additional infrastructure

## Metrics

The cache tracks the following metrics:

- **Hits**: Successful cache retrievals
- **Misses**: Failed cache retrievals
- **Sets**: Cache write operations
- **Deletes**: Cache removal operations
- **Evictions**: Automatic removals (LRU or expiration)
- **Size**: Current number of cached items
- **MaxSize**: Maximum allowed items (memory cache)

### Hit Rate Calculation

```go
stats := cache.Stats()
hitRate := stats.HitRate() // Returns 0.0 to 1.0
fmt.Printf("Cache hit rate: %.1f%%\n", hitRate * 100)
```

## Best Practices

1. **Choose appropriate TTLs**: Balance between freshness and performance
2. **Use consistent key patterns**: Leverage KeyBuilder for maintainability
3. **Monitor metrics**: Track hit rates and adjust cache size accordingly
4. **Handle cache misses gracefully**: Always have a fallback to the source
5. **Consider data volatility**: Don't cache frequently changing data
6. **Size appropriately**: Set max size based on available memory

## Integration with Handlers

### Example: Caching Account Lists

```go
func (h *AccountHandler) ListAccounts(ctx context.Context, args interface{}) (interface{}, error) {
    // Generate cache key based on arguments
    kb := NewKeyBuilder("firefly")
    key := kb.Add("accounts").Add(fmt.Sprintf("%v", args)).Build()
    
    // Try cache first
    if cached, err := h.cache.Get(ctx, key); err == nil {
        return cached, nil
    }
    
    // Cache miss - fetch from API
    result, err := h.client.ListAccounts(args)
    if err != nil {
        return nil, err
    }
    
    // Store in cache
    _ = h.cache.Set(ctx, key, result, 5*time.Minute)
    
    return result, nil
}
```

## Testing

The package includes comprehensive tests:

```bash
# Run all tests
go test ./pkg/fireflyMCP/cache -v

# Run with coverage
go test ./pkg/fireflyMCP/cache -cover

# Run benchmarks
go test ./pkg/fireflyMCP/cache -bench=.
```

## Thread Safety

All cache implementations are thread-safe and can be safely used from multiple goroutines concurrently.

## Error Handling

- `ErrCacheMiss`: Key not found in cache
- Network errors: For Redis implementation
- Serialization errors: When storing complex types

## Future Enhancements

- [ ] Pattern-based invalidation
- [ ] Compression for large values
- [ ] Warm-up/pre-loading capability
- [ ] Circuit breaker for Redis
- [ ] Metrics export (Prometheus)
- [ ] Cache warming strategies
- [ ] Two-tier caching (L1/L2)