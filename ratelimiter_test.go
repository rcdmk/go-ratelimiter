package ratelimiter_test

import (
	"testing"
	"time"

	"github.com/rcdmk/go-ratelimiter"
)

func TestRateLimiter_Allow(t *testing.T) {
	// Create a new RateLimiter with options
	options := ratelimiter.Options{
		MaxRatePerSecond: 10,
		MaxBurst:         5,
	}
	limiter := ratelimiter.New(options)

	// Allow 5 events
	for i := 0; i < 5; i++ {
		if !limiter.Allow() {
			t.Errorf("Expected limiter to allow event, but it didn't")
		}
	}

	// Allow 5 more events after waiting for 1 second
	time.Sleep(time.Second)
	for i := 0; i < 5; i++ {
		if !limiter.Allow() {
			t.Errorf("Expected limiter to allow event, but it didn't")
		}
	}

	// Try to allow 1 more event, which should be rate-limited
	if limiter.Allow() {
		t.Errorf("Expected limiter to rate-limit event, but it didn't")
	}
}

func TestRateLimiter_Allow_ZeroRate(t *testing.T) {
	// Create a new RateLimiter with zero rate
	options := ratelimiter.Options{
		MaxRatePerSecond: 0,
		MaxBurst:         5,
	}
	limiter := ratelimiter.New(options)

	// Try to allow an event, which should be rate-limited
	if limiter.Allow() {
		t.Errorf("Expected limiter to rate-limit event, but it didn't")
	}
}

func TestRateLimiter_Allow_ZeroBurst(t *testing.T) {
	// Create a new RateLimiter with zero burst
	options := ratelimiter.Options{
		MaxRatePerSecond: 10,
		MaxBurst:         0,
	}
	limiter := ratelimiter.New(options)

	// Try to allow an event, which should be rate-limited
	if limiter.Allow() {
		t.Errorf("Expected limiter to rate-limit event, but it didn't")
	}
}

func TestRateLimiter_Allow_BurstHigherThanMaxRate(t *testing.T) {
	// Create a new RateLimiter with burst higher than max rate
	options := ratelimiter.Options{
		MaxRatePerSecond: 10,
		MaxBurst:         15,
	}
	limiter := ratelimiter.New(options)
	// Allow 15 events
	for i := 0; i < 15; i++ {
		if !limiter.Allow() {
			t.Errorf("Expected limiter to allow event, but it didn't")
		}
	}
	// Try to allow 1 more event, which should be rate-limited
	if limiter.Allow() {
		t.Errorf("Expected limiter to rate-limit event, but it didn't")
	}
	// Allow 10 more events after waiting for 1 second
	time.Sleep(time.Second)
	for i := 0; i < 10; i++ {
		if !limiter.Allow() {
			t.Errorf("Expected limiter to allow event, but it didn't")
		}
	}
	// Try to allow 1 more event, which should be rate-limited
	if limiter.Allow() {
		t.Errorf("Expected limiter to rate-limit event, but it didn't")
	}
}
