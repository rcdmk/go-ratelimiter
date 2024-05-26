package ratelimiter

import (
	"sync"
	"time"
)

// RateLimiter represents a rate limiter that limits the rate of events, implemented using a token bucket algorithm.
type RateLimiter struct {
	maxRatePerSecond int                  // The maximum rate of events allowed per second.
	maxBurst         int                  // The maximum number of events that can be bursted.
	bucket           map[string]int       // The current number of events in the bucket.
	lastFill         map[string]time.Time // The time when the bucket was last filled.
	mu               sync.Mutex           // Mutex to synchronize access to the rate limiter.
}

// fillBucket fills the bucket with tokens based on the elapsed time since the last fill.
func (rl *RateLimiter) fillBucket(sourceKey string) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	elapsed := time.Since(rl.lastFill[sourceKey]).Seconds()
	rl.bucket[sourceKey] += int(elapsed * float64(rl.maxRatePerSecond))
	if rl.bucket[sourceKey] > rl.maxBurst {
		rl.bucket[sourceKey] = rl.maxBurst
	}
	rl.lastFill[sourceKey] = time.Now()
}

// Allow checks if the rate wasn't exhausted for a particular key to allow or not an event to be executed.
func (rl *RateLimiter) Allow(sourceKey string) bool {
	rl.fillBucket(sourceKey)

	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.bucket[sourceKey] > 0 {
		rl.bucket[sourceKey]--
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
		bucket:           make(map[string]int),
		lastFill:         make(map[string]time.Time),
	}
}
