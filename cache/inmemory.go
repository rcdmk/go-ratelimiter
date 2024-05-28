package cache

import (
	"sync"
	"time"
)

// inMemoryEntry represents a cache entry that supports expiration of values.
type inMemoryEntry struct {
	value      int // The value stored in the cache
	expiration int // Unix time in milliseconds when the entry expires
}

// InMemory represents an in-memory cache that stores values in memory.
type InMemory struct {
	cache map[string]inMemoryEntry
	mu    sync.Mutex
}

// NewInMemory creates a new ready to use InMemory cache.
func NewInMemory() *InMemory {
	return &InMemory{
		cache: make(map[string]inMemoryEntry),
	}
}

// Get retrieves a value from the cache.
// error can only be nil or ErrCacheMiss for this implementation.
func (c *InMemory) Get(key string) (int, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if entry, ok := c.cache[key]; ok {
		if entry.expiration > 0 && entry.expiration <= int(time.Now().UnixMilli()) {
			delete(c.cache, key)
			return 0, ErrCacheMiss
		}
		return entry.value, nil
	}
	return 0, ErrCacheMiss
}

// Set stores a value in the cache without expiration time.
func (c *InMemory) Set(key string, value int) error {
	return c.SetWithExpiration(key, value, 0)
}

// Set stores a value in the cache with a given expiration time.
// If expiration is 0, the value never expires.
// error is always nil for this implementation.
func (c *InMemory) SetWithExpiration(key string, value int, expiration time.Duration) error {
	var expirationTime int
	if expiration > 0 {
		expirationTime = int(time.Now().Add(expiration).UnixMilli())
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache[key] = inMemoryEntry{value: value, expiration: expirationTime}
	return nil
}
