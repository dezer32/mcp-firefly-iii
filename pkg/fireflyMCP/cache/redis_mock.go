package cache

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
	"time"
)

// MockRedisClient is a mock implementation of RedisClient for testing.
type MockRedisClient struct {
	mu     sync.RWMutex
	data   map[string]string
	expiry map[string]time.Time
}

// NewMockRedisClient creates a new mock Redis client.
func NewMockRedisClient() *MockRedisClient {
	return &MockRedisClient{
		data:   make(map[string]string),
		expiry: make(map[string]time.Time),
	}
}

// Get retrieves a value from the mock Redis.
func (m *MockRedisClient) Get(ctx context.Context, key string) (string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	val, exists := m.data[key]
	if !exists {
		return "", errors.New("redis: nil")
	}
	
	// Check expiration
	if exp, hasExp := m.expiry[key]; hasExp && time.Now().After(exp) {
		return "", errors.New("redis: nil")
	}
	
	return val, nil
}

// Set stores a value in the mock Redis.
func (m *MockRedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	// Convert value to string
	var strVal string
	switch v := value.(type) {
	case string:
		strVal = v
	case []byte:
		strVal = string(v)
	default:
		data, err := json.Marshal(v)
		if err != nil {
			return err
		}
		strVal = string(data)
	}
	
	m.data[key] = strVal
	
	if expiration > 0 {
		m.expiry[key] = time.Now().Add(expiration)
	} else {
		delete(m.expiry, key)
	}
	
	return nil
}

// Del deletes keys from the mock Redis.
func (m *MockRedisClient) Del(ctx context.Context, keys ...string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	for _, key := range keys {
		delete(m.data, key)
		delete(m.expiry, key)
	}
	
	return nil
}

// Exists checks if keys exist in the mock Redis.
func (m *MockRedisClient) Exists(ctx context.Context, keys ...string) (int64, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	count := int64(0)
	for _, key := range keys {
		if _, exists := m.data[key]; exists {
			// Check expiration
			if exp, hasExp := m.expiry[key]; !hasExp || time.Now().Before(exp) {
				count++
			}
		}
	}
	
	return count, nil
}

// TTL returns the TTL of a key in the mock Redis.
func (m *MockRedisClient) TTL(ctx context.Context, key string) (time.Duration, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	if _, exists := m.data[key]; !exists {
		return -2 * time.Second, nil // Key doesn't exist
	}
	
	exp, hasExp := m.expiry[key]
	if !hasExp {
		return -1 * time.Second, nil // Key exists but has no expiration
	}
	
	ttl := time.Until(exp)
	if ttl < 0 {
		return -2 * time.Second, nil // Key expired
	}
	
	return ttl, nil
}

// FlushAll clears all data from the mock Redis.
func (m *MockRedisClient) FlushAll(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	m.data = make(map[string]string)
	m.expiry = make(map[string]time.Time)
	
	return nil
}

// Close closes the mock Redis client.
func (m *MockRedisClient) Close() error {
	return nil
}