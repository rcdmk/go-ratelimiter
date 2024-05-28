// Description: This package provides an interface and default implementation for caching that can be used with the rate limiter.
package cache

import (
	"errors"
	"time"
)

// ErrCacheMiss represents an error when a cache miss occurs.
// It is provided so clients can differentiate between a cache miss and a cache hit with a zero value.
var ErrCacheMiss = errors.New("cache miss")

// GetterSetter represents an interface for getting and setting values in a cache.
type GetterSetter interface {
	Get(key string) (int, error)
	Set(key string, value int) error
	SetWithExpiration(key string, value int, expiration time.Duration) error
}
