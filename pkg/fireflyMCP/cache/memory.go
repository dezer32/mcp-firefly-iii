package cache

import (
	"context"
	"sync"
	"sync/atomic"
	"time"
)

// cacheEntry represents a single cache entry.
type cacheEntry struct {
	value      interface{}
	expiration time.Time
	size       int64
}

// isExpired checks if the entry has expired.
func (e *cacheEntry) isExpired() bool {
	if e.expiration.IsZero() {
		return false
	}
	return time.Now().After(e.expiration)
}

// MemoryCache is an in-memory cache implementation.
type MemoryCache struct {
	mu              sync.RWMutex
	data            map[string]*cacheEntry
	stats           Stats
	maxSize         int64
	currentSize     int64
	defaultTTL      time.Duration
	cleanupInterval time.Duration
	stopCleanup     chan struct{}
}

// NewMemoryCache creates a new in-memory cache.
func NewMemoryCache(opts ...Option) *MemoryCache {
	mc := &MemoryCache{
		data:            make(map[string]*cacheEntry),
		maxSize:         0, // 0 means unlimited
		cleanupInterval: 1 * time.Minute,
		stopCleanup:     make(chan struct{}),
	}
	
	for _, opt := range opts {
		opt(mc)
	}
	
	// Start cleanup goroutine
	go mc.cleanupExpired()
	
	return mc
}

// Get retrieves a value from the cache.
func (mc *MemoryCache) Get(ctx context.Context, key string) (interface{}, error) {
	mc.mu.RLock()
	entry, exists := mc.data[key]
	mc.mu.RUnlock()
	
	if !exists {
		atomic.AddInt64(&mc.stats.Misses, 1)
		return nil, ErrCacheMiss
	}
	
	if entry.isExpired() {
		mc.Delete(ctx, key)
		atomic.AddInt64(&mc.stats.Misses, 1)
		return nil, ErrCacheMiss
	}
	
	atomic.AddInt64(&mc.stats.Hits, 1)
	return entry.value, nil
}

// Set stores a value in the cache with the given TTL.
func (mc *MemoryCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	if ttl == 0 {
		ttl = mc.defaultTTL
	}
	
	var expiration time.Time
	if ttl > 0 {
		expiration = time.Now().Add(ttl)
	}
	
	entry := &cacheEntry{
		value:      value,
		expiration: expiration,
		size:       1, // Simple size calculation; could be improved
	}
	
	mc.mu.Lock()
	defer mc.mu.Unlock()
	
	// Check if we need to evict entries to make room
	if mc.maxSize > 0 && mc.currentSize+entry.size > mc.maxSize {
		mc.evictLRU()
	}
	
	// Update or add entry
	if oldEntry, exists := mc.data[key]; exists {
		mc.currentSize -= oldEntry.size
	}
	
	mc.data[key] = entry
	mc.currentSize += entry.size
	atomic.AddInt64(&mc.stats.Sets, 1)
	atomic.StoreInt64(&mc.stats.Size, int64(len(mc.data)))
	
	return nil
}

// Delete removes a value from the cache.
func (mc *MemoryCache) Delete(ctx context.Context, key string) error {
	mc.mu.Lock()
	defer mc.mu.Unlock()
	
	if entry, exists := mc.data[key]; exists {
		mc.currentSize -= entry.size
		delete(mc.data, key)
		atomic.AddInt64(&mc.stats.Deletes, 1)
		atomic.StoreInt64(&mc.stats.Size, int64(len(mc.data)))
	}
	
	return nil
}

// Clear removes all values from the cache.
func (mc *MemoryCache) Clear(ctx context.Context) error {
	mc.mu.Lock()
	defer mc.mu.Unlock()
	
	mc.data = make(map[string]*cacheEntry)
	mc.currentSize = 0
	atomic.StoreInt64(&mc.stats.Size, 0)
	
	return nil
}

// Exists checks if a key exists in the cache.
func (mc *MemoryCache) Exists(ctx context.Context, key string) (bool, error) {
	mc.mu.RLock()
	entry, exists := mc.data[key]
	mc.mu.RUnlock()
	
	if !exists {
		return false, nil
	}
	
	if entry.isExpired() {
		mc.Delete(ctx, key)
		return false, nil
	}
	
	return true, nil
}

// TTL returns the remaining TTL for a key.
func (mc *MemoryCache) TTL(ctx context.Context, key string) (time.Duration, error) {
	mc.mu.RLock()
	entry, exists := mc.data[key]
	mc.mu.RUnlock()
	
	if !exists {
		return 0, nil
	}
	
	if entry.expiration.IsZero() {
		return 0, nil // No expiration
	}
	
	ttl := time.Until(entry.expiration)
	if ttl < 0 {
		return 0, nil
	}
	
	return ttl, nil
}

// Stats returns cache statistics.
func (mc *MemoryCache) Stats() Stats {
	stats := mc.stats
	stats.MaxSize = mc.maxSize
	return stats
}

// Close stops the cleanup goroutine.
func (mc *MemoryCache) Close() error {
	close(mc.stopCleanup)
	return nil
}

// cleanupExpired periodically removes expired entries.
func (mc *MemoryCache) cleanupExpired() {
	ticker := time.NewTicker(mc.cleanupInterval)
	defer ticker.Stop()
	
	for {
		select {
		case <-ticker.C:
			mc.removeExpired()
		case <-mc.stopCleanup:
			return
		}
	}
}

// removeExpired removes all expired entries.
func (mc *MemoryCache) removeExpired() {
	mc.mu.Lock()
	defer mc.mu.Unlock()
	
	for key, entry := range mc.data {
		if entry.isExpired() {
			mc.currentSize -= entry.size
			delete(mc.data, key)
			atomic.AddInt64(&mc.stats.Evictions, 1)
		}
	}
	
	atomic.StoreInt64(&mc.stats.Size, int64(len(mc.data)))
}

// evictLRU evicts the least recently used entries.
// This is a simple implementation; could be improved with an actual LRU algorithm.
func (mc *MemoryCache) evictLRU() {
	// For simplicity, just remove the first expired entry or the first entry
	for key, entry := range mc.data {
		if entry.isExpired() || mc.currentSize > mc.maxSize {
			mc.currentSize -= entry.size
			delete(mc.data, key)
			atomic.AddInt64(&mc.stats.Evictions, 1)
			if mc.currentSize <= mc.maxSize {
				break
			}
		}
	}
}