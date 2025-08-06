package cache

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestKeyBuilder(t *testing.T) {
	tests := []struct {
		name     string
		prefix   string
		parts    []string
		expected string
	}{
		{
			name:     "with prefix and parts",
			prefix:   "firefly",
			parts:    []string{"accounts", "123"},
			expected: "firefly:accounts:123",
		},
		{
			name:     "without prefix",
			prefix:   "",
			parts:    []string{"transactions", "456"},
			expected: "transactions:456",
		},
		{
			name:     "only prefix",
			prefix:   "cache",
			parts:    []string{},
			expected: "cache",
		},
		{
			name:     "empty",
			prefix:   "",
			parts:    []string{},
			expected: "",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kb := NewKeyBuilder(tt.prefix)
			for _, part := range tt.parts {
				kb.Add(part)
			}
			assert.Equal(t, tt.expected, kb.Build())
		})
	}
}

func TestStats_HitRate(t *testing.T) {
	tests := []struct {
		name     string
		stats    Stats
		expected float64
	}{
		{
			name: "50% hit rate",
			stats: Stats{
				Hits:   50,
				Misses: 50,
			},
			expected: 0.5,
		},
		{
			name: "100% hit rate",
			stats: Stats{
				Hits:   100,
				Misses: 0,
			},
			expected: 1.0,
		},
		{
			name: "0% hit rate",
			stats: Stats{
				Hits:   0,
				Misses: 100,
			},
			expected: 0.0,
		},
		{
			name: "no data",
			stats: Stats{
				Hits:   0,
				Misses: 0,
			},
			expected: 0.0,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.stats.HitRate())
		})
	}
}

func TestMemoryCache_BasicOperations(t *testing.T) {
	ctx := context.Background()
	cache := NewMemoryCache()
	defer cache.Close()
	
	// Test Set and Get
	err := cache.Set(ctx, "key1", "value1", 1*time.Hour)
	require.NoError(t, err)
	
	val, err := cache.Get(ctx, "key1")
	require.NoError(t, err)
	assert.Equal(t, "value1", val)
	
	// Test Get non-existent key
	_, err = cache.Get(ctx, "nonexistent")
	assert.ErrorIs(t, err, ErrCacheMiss)
	
	// Test Exists
	exists, err := cache.Exists(ctx, "key1")
	require.NoError(t, err)
	assert.True(t, exists)
	
	exists, err = cache.Exists(ctx, "nonexistent")
	require.NoError(t, err)
	assert.False(t, exists)
	
	// Test Delete
	err = cache.Delete(ctx, "key1")
	require.NoError(t, err)
	
	exists, err = cache.Exists(ctx, "key1")
	require.NoError(t, err)
	assert.False(t, exists)
	
	// Test Clear
	err = cache.Set(ctx, "key2", "value2", 1*time.Hour)
	require.NoError(t, err)
	err = cache.Set(ctx, "key3", "value3", 1*time.Hour)
	require.NoError(t, err)
	
	err = cache.Clear(ctx)
	require.NoError(t, err)
	
	exists, err = cache.Exists(ctx, "key2")
	require.NoError(t, err)
	assert.False(t, exists)
	
	exists, err = cache.Exists(ctx, "key3")
	require.NoError(t, err)
	assert.False(t, exists)
}

func TestMemoryCache_TTL(t *testing.T) {
	ctx := context.Background()
	cache := NewMemoryCache()
	defer cache.Close()
	
	// Set with TTL
	err := cache.Set(ctx, "key1", "value1", 2*time.Second)
	require.NoError(t, err)
	
	// Check TTL
	ttl, err := cache.TTL(ctx, "key1")
	require.NoError(t, err)
	assert.Greater(t, ttl, time.Duration(0))
	assert.LessOrEqual(t, ttl, 2*time.Second)
	
	// Set without TTL (no expiration)
	err = cache.Set(ctx, "key2", "value2", 0)
	require.NoError(t, err)
	
	ttl, err = cache.TTL(ctx, "key2")
	require.NoError(t, err)
	assert.Equal(t, time.Duration(0), ttl)
	
	// Non-existent key
	ttl, err = cache.TTL(ctx, "nonexistent")
	require.NoError(t, err)
	assert.Equal(t, time.Duration(0), ttl)
}

func TestMemoryCache_Expiration(t *testing.T) {
	ctx := context.Background()
	cache := NewMemoryCache()
	defer cache.Close()
	
	// Set with short TTL
	err := cache.Set(ctx, "key1", "value1", 100*time.Millisecond)
	require.NoError(t, err)
	
	// Should exist initially
	exists, err := cache.Exists(ctx, "key1")
	require.NoError(t, err)
	assert.True(t, exists)
	
	// Wait for expiration
	time.Sleep(150 * time.Millisecond)
	
	// Should not exist after expiration
	_, err = cache.Get(ctx, "key1")
	assert.ErrorIs(t, err, ErrCacheMiss)
	
	exists, err = cache.Exists(ctx, "key1")
	require.NoError(t, err)
	assert.False(t, exists)
}

func TestMemoryCache_Stats(t *testing.T) {
	ctx := context.Background()
	cache := NewMemoryCache()
	defer cache.Close()
	
	// Initial stats
	stats := cache.Stats()
	assert.Equal(t, int64(0), stats.Hits)
	assert.Equal(t, int64(0), stats.Misses)
	assert.Equal(t, int64(0), stats.Sets)
	assert.Equal(t, int64(0), stats.Deletes)
	
	// Perform operations
	_ = cache.Set(ctx, "key1", "value1", 1*time.Hour)
	stats = cache.Stats()
	assert.Equal(t, int64(1), stats.Sets)
	assert.Equal(t, int64(1), stats.Size)
	
	_, _ = cache.Get(ctx, "key1")
	stats = cache.Stats()
	assert.Equal(t, int64(1), stats.Hits)
	
	_, _ = cache.Get(ctx, "nonexistent")
	stats = cache.Stats()
	assert.Equal(t, int64(1), stats.Misses)
	
	_ = cache.Delete(ctx, "key1")
	stats = cache.Stats()
	assert.Equal(t, int64(1), stats.Deletes)
	assert.Equal(t, int64(0), stats.Size)
}

func TestMemoryCache_Options(t *testing.T) {
	ctx := context.Background()
	
	// Test with max size
	cache := NewMemoryCache(
		WithMaxSize(2),
		WithDefaultTTL(1*time.Hour),
		WithCleanupInterval(100*time.Millisecond),
	)
	defer cache.Close()
	
	// Add items up to max size
	err := cache.Set(ctx, "key1", "value1", 0)
	require.NoError(t, err)
	err = cache.Set(ctx, "key2", "value2", 0)
	require.NoError(t, err)
	
	// Adding another should trigger eviction
	err = cache.Set(ctx, "key3", "value3", 0)
	require.NoError(t, err)
	
	stats := cache.Stats()
	assert.Equal(t, int64(2), stats.MaxSize)
	
	// Test default TTL
	cache2 := NewMemoryCache(WithDefaultTTL(100 * time.Millisecond))
	defer cache2.Close()
	
	err = cache2.Set(ctx, "key1", "value1", 0) // Use default TTL
	require.NoError(t, err)
	
	time.Sleep(150 * time.Millisecond)
	
	_, err = cache2.Get(ctx, "key1")
	assert.ErrorIs(t, err, ErrCacheMiss)
}

func TestMemoryCache_ComplexTypes(t *testing.T) {
	ctx := context.Background()
	cache := NewMemoryCache()
	defer cache.Close()
	
	// Test struct
	type TestStruct struct {
		ID   string
		Name string
		Age  int
	}
	
	testData := TestStruct{
		ID:   "123",
		Name: "Test",
		Age:  30,
	}
	
	err := cache.Set(ctx, "struct", testData, 1*time.Hour)
	require.NoError(t, err)
	
	val, err := cache.Get(ctx, "struct")
	require.NoError(t, err)
	assert.Equal(t, testData, val)
	
	// Test slice
	slice := []string{"a", "b", "c"}
	err = cache.Set(ctx, "slice", slice, 1*time.Hour)
	require.NoError(t, err)
	
	val, err = cache.Get(ctx, "slice")
	require.NoError(t, err)
	assert.Equal(t, slice, val)
	
	// Test map
	m := map[string]int{"one": 1, "two": 2}
	err = cache.Set(ctx, "map", m, 1*time.Hour)
	require.NoError(t, err)
	
	val, err = cache.Get(ctx, "map")
	require.NoError(t, err)
	assert.Equal(t, m, val)
}