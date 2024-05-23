package ratelimiter

import (
	"sync"
	"time"
)

// RateLimiter represents a rate limiter that limits the rate of events, implemented using a token bucket algorithm.
type RateLimiter struct {
	maxRatePerSecond int        // The maximum rate of events allowed per second.
	maxBurst         int        // The maximum number of events that can be bursted.
	bucket           int        // The current number of events in the bucket.
	lastFill         time.Time  // The time when the bucket was last filled.
	mu               sync.Mutex // Mutex to synchronize access to the rate limiter.
}

// fillBucket fills the bucket with tokens based on the elapsed time since the last fill.
func (rl *RateLimiter) fillBucket() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	elapsed := time.Since(rl.lastFill).Seconds()
	rl.bucket += int(elapsed * float64(rl.maxRatePerSecond))
	if rl.bucket > rl.maxBurst {
		rl.bucket = rl.maxBurst
	}
	rl.lastFill = time.Now()
}

// Allow checks if the rate wasn't exhausted to allow or not an event to be executed.
func (rl *RateLimiter) Allow() bool {
	rl.fillBucket()

	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.bucket > 0 {
		rl.bucket--
		return true
	}

	return false
}

// Options represents the options for configuring a RateLimiter.
type Options struct {
	MaxRatePerSecond int // The maximum rate of events allowed per second.
	MaxBurst         int // The maximum number of events that can be bursted.
}

// New creates a new ready to use RateLimiter with the specified options.
func New(options Options) *RateLimiter {
	return &RateLimiter{
		maxRatePerSecond: options.MaxRatePerSecond,
		maxBurst:         options.MaxBurst,
	}
}
