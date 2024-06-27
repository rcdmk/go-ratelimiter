// Package rediscache provides a cache service that stores values in Redis, implementing the [cache.GetterSetter] interface.
package rediscache

import (
	"context"
	"time"

	"github.com/rcdmk/go-ratelimiter/cache"
	"github.com/redis/go-redis/v9"
)

// Redis represents a cache service that stores values in Redis.
type Redis struct {
	client *redis.Client
}

// New creates a new ready to use InMemory cache.
func New(client *redis.Client) *Redis {
	return &Redis{
		client: client,
	}
}

// Get retrieves a value from the cache.
// error is ErrCacheMiss if key is not present in the cache.
func (c *Redis) Get(key string) (int, error) {
	val, err := c.client.Get(context.Background(), key).Int()
	if err == redis.Nil {
		return 0, cache.ErrCacheMiss
	}

	if err != nil {
		return 0, err
	}

	return val, nil
}

// Set stores a value in the cache without expiration time.
func (c *Redis) Set(key string, value int) error {
	return c.SetWithExpiration(key, value, 0)
}

// SetWithExpiration stores a value in the cache with a given expiration time.
// If expiration is 0, the value never expires.
func (c *Redis) SetWithExpiration(key string, value int, expiration time.Duration) error {
	return c.client.Set(context.Background(), key, value, expiration).Err()
}
