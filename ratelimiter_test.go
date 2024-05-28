package ratelimiter_test

import (
	"errors"
	"testing"
	"time"

	"github.com/rcdmk/go-ratelimiter"
)

func TestRateLimiter_Allow(t *testing.T) {
	sourceKey := "test"

	options := ratelimiter.Options{
		MaxRatePerSecond: 10,
		MaxBurst:         5,
	}
	limiter := ratelimiter.New(options)

	// Allow 5 events
	for i := 0; i < 5; i++ {
		if !limiter.Allow(sourceKey) {
			t.Errorf("Expected limiter to allow event, but it didn't")
		}
	}

	// Allow 5 more events after waiting for 1 second
	time.Sleep(time.Second)
	for i := 0; i < 5; i++ {
		if !limiter.Allow(sourceKey) {
			t.Errorf("Expected limiter to allow event, but it didn't")
		}
	}

	// Try to allow 1 more event, which should be rate-limited
	if limiter.Allow(sourceKey) {
		t.Errorf("Expected limiter to rate-limit event, but it didn't")
	}
}

func TestRateLimiter_Allow_ZeroRate(t *testing.T) {
	sourceKey := "test"

	// Create a new RateLimiter with zero rate
	options := ratelimiter.Options{
		MaxRatePerSecond: 0,
		MaxBurst:         5,
	}
	limiter := ratelimiter.New(options)

	// Try to allow an event, which should be rate-limited
	if limiter.Allow(sourceKey) {
		t.Errorf("Expected limiter to rate-limit event, but it didn't")
	}
}

func TestRateLimiter_Allow_ZeroBurst(t *testing.T) {
	limitKey := "test"

	options := ratelimiter.Options{
		MaxRatePerSecond: 10,
		MaxBurst:         0,
	}
	limiter := ratelimiter.New(options)

	// Try to allow an event, which should be rate-limited
	if limiter.Allow(limitKey) {
		t.Errorf("Expected limiter to rate-limit event, but it didn't")
	}
}

func TestRateLimiter_Allow_BurstHigherThanMaxRate(t *testing.T) {
	sourceKey := "test"

	options := ratelimiter.Options{
		MaxRatePerSecond: 10,
		MaxBurst:         15,
	}
	limiter := ratelimiter.New(options)

	// Allow 15 events
	for i := 0; i < 15; i++ {
		if !limiter.Allow(sourceKey) {
			t.Errorf("Expected limiter to allow event, but it didn't")
		}
	}

	// Try to allow 1 more event, which should be rate-limited
	if limiter.Allow(sourceKey) {
		t.Errorf("Expected limiter to rate-limit event, but it didn't")
	}

	// Allow 10 more events after waiting for 1 second
	time.Sleep(time.Second)
	for i := 0; i < 10; i++ {
		if !limiter.Allow(sourceKey) {
			t.Errorf("Expected limiter to allow event, but it didn't")
		}
	}

	// Try to allow 1 more event, which should be rate-limited
	if limiter.Allow(sourceKey) {
		t.Errorf("Expected limiter to rate-limit event, but it didn't")
	}
}

func TestRateLimiter_Allow_Multiple_Keys(t *testing.T) {
	sourceKey1 := "test1"
	sourceKey2 := "test2"

	options := ratelimiter.Options{
		MaxRatePerSecond: 10,
		MaxBurst:         5,
	}
	limiter := ratelimiter.New(options)

	// Allow 5 events for limitKey1
	for i := 0; i < 5; i++ {
		if !limiter.Allow(sourceKey1) {
			t.Errorf("Expected limiter to allow event, but it didn't")
		}
	}

	// Allow 5 events for limitKey2
	for i := 0; i < 5; i++ {
		if !limiter.Allow(sourceKey2) {
			t.Errorf("Expected limiter to allow event, but it didn't")
		}
	}

	// Try to allow 1 more event for limitKey1, which should be rate-limited
	if limiter.Allow(sourceKey1) {
		t.Errorf("Expected limiter to rate-limit event, but it didn't")
	}

	// Try to allow 1 more event for limitKey2, which should be rate-limited
	if limiter.Allow(sourceKey2) {
		t.Errorf("Expected limiter to rate-limit event, but it didn't")
	}
}

func TestRateLimiter_Allow_Always_If_Cache_Fails(t *testing.T) {
	sourceKey := "test"

	options := ratelimiter.Options{
		MaxRatePerSecond: 10,
		MaxBurst:         5,
		Cache:            &mockFailedCache{},
	}

	limiter := ratelimiter.New(options)

	// Allow 15 events
	for i := 0; i < 30; i++ {
		if !limiter.Allow(sourceKey) {
			t.Errorf("Expected limiter to allow event, but it didn't: %d", i)
		}
	}
}

// mockFailedCache is a mock implementation of the cache.GetterSetter interface that fails on all operations.
type mockFailedCache struct{}

func (c *mockFailedCache) Get(key string) (int, error) {
	return 0, errors.New("mock cache error: get")
}

func (c *mockFailedCache) Set(key string, value int) error {
	return errors.New("mock cache error: set")
}

func (c *mockFailedCache) SetWithExpiration(key string, value int, expiration time.Duration) error {
	return errors.New("mock cache error: set with expiration")
}
