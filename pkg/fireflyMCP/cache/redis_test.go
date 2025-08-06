package cache

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRedisCache_BasicOperations(t *testing.T) {
	ctx := context.Background()
	mockClient := NewMockRedisClient()
	cache := NewRedisCache(mockClient, "test")
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
}

func TestRedisCache_WithPrefix(t *testing.T) {
	ctx := context.Background()
	mockClient := NewMockRedisClient()
	cache := NewRedisCache(mockClient, "firefly")
	
	// Set a value
	err := cache.Set(ctx, "accounts:123", map[string]string{"id": "123"}, 1*time.Hour)
	require.NoError(t, err)
	
	// Value should be stored with prefix
	rawVal, err := mockClient.Get(ctx, "firefly:accounts:123")
	require.NoError(t, err)
	assert.NotEmpty(t, rawVal)
	
	// Get should work with unprefixed key
	val, err := cache.Get(ctx, "accounts:123")
	require.NoError(t, err)
	assert.NotNil(t, val)
}

func TestRedisCache_TTL(t *testing.T) {
	ctx := context.Background()
	mockClient := NewMockRedisClient()
	cache := NewRedisCache(mockClient, "")
	
	// Set with TTL
	err := cache.Set(ctx, "key1", "value1", 2*time.Second)
	require.NoError(t, err)
	
	// Check TTL
	ttl, err := cache.TTL(ctx, "key1")
	require.NoError(t, err)
	assert.Greater(t, ttl, time.Duration(0))
	assert.LessOrEqual(t, ttl, 2*time.Second)
	
	// Non-existent key
	ttl, err = cache.TTL(ctx, "nonexistent")
	require.NoError(t, err)
	assert.Equal(t, time.Duration(0), ttl)
}

func TestRedisCache_ComplexTypes(t *testing.T) {
	ctx := context.Background()
	mockClient := NewMockRedisClient()
	cache := NewRedisCache(mockClient, "test")
	
	// Test struct
	type TestStruct struct {
		ID   string `json:"id"`
		Name string `json:"name"`
		Age  int    `json:"age"`
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
	
	// The value will be returned as a map since we marshal/unmarshal through JSON
	valMap, ok := val.(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, "123", valMap["id"])
	assert.Equal(t, "Test", valMap["name"])
	assert.Equal(t, float64(30), valMap["age"]) // JSON numbers are float64
	
	// Test slice
	slice := []string{"a", "b", "c"}
	err = cache.Set(ctx, "slice", slice, 1*time.Hour)
	require.NoError(t, err)
	
	val, err = cache.Get(ctx, "slice")
	require.NoError(t, err)
	
	valSlice, ok := val.([]interface{})
	require.True(t, ok)
	assert.Len(t, valSlice, 3)
	assert.Equal(t, "a", valSlice[0])
	assert.Equal(t, "b", valSlice[1])
	assert.Equal(t, "c", valSlice[2])
}

func TestRedisCache_Stats(t *testing.T) {
	ctx := context.Background()
	mockClient := NewMockRedisClient()
	cache := NewRedisCache(mockClient, "test")
	
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
	
	_, _ = cache.Get(ctx, "key1")
	stats = cache.Stats()
	assert.Equal(t, int64(1), stats.Hits)
	
	_, _ = cache.Get(ctx, "nonexistent")
	stats = cache.Stats()
	assert.Equal(t, int64(1), stats.Misses)
	
	_ = cache.Delete(ctx, "key1")
	stats = cache.Stats()
	assert.Equal(t, int64(1), stats.Deletes)
}

func TestRedisCache_DefaultTTL(t *testing.T) {
	ctx := context.Background()
	mockClient := NewMockRedisClient()
	cache := NewRedisCache(mockClient, "test", WithDefaultTTL(100*time.Millisecond))
	
	// Set without explicit TTL (should use default)
	err := cache.Set(ctx, "key1", "value1", 0)
	require.NoError(t, err)
	
	// Value should exist initially
	val, err := cache.Get(ctx, "key1")
	require.NoError(t, err)
	assert.Equal(t, "value1", val)
	
	// Wait for default TTL to expire
	time.Sleep(150 * time.Millisecond)
	
	// Value should be expired
	_, err = mockClient.Get(ctx, "test:key1")
	assert.Error(t, err)
}

func TestRedisCache_Clear(t *testing.T) {
	ctx := context.Background()
	mockClient := NewMockRedisClient()
	
	// Test without prefix (flushes all)
	cache := NewRedisCache(mockClient, "")
	
	err := cache.Set(ctx, "key1", "value1", 1*time.Hour)
	require.NoError(t, err)
	err = cache.Set(ctx, "key2", "value2", 1*time.Hour)
	require.NoError(t, err)
	
	err = cache.Clear(ctx)
	require.NoError(t, err)
	
	exists, err := cache.Exists(ctx, "key1")
	require.NoError(t, err)
	assert.False(t, exists)
	
	exists, err = cache.Exists(ctx, "key2")
	require.NoError(t, err)
	assert.False(t, exists)
	
	// Test with prefix (not yet fully implemented)
	cache2 := NewRedisCache(mockClient, "prefix")
	err = cache2.Clear(ctx)
	assert.Error(t, err) // Should return "not yet implemented" error
}