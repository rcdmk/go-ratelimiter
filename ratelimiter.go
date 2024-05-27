package ratelimiter

import (
	"sync"
	"time"
)

// RateLimiter represents a rate limiter that limits the rate of events, implemented using a token bucket algorithm.
type RateLimiter struct {
	maxRatePerMillisecond float64        // The maximum rate of events allowed per millisecond.
	maxBurst              int            // The maximum number of events that can be bursted.
	bucket                map[string]int // The current number of events in the bucket.
	lastFill              map[string]int // The Unix time in milliseconds when the bucket was last filled.
	mu                    sync.Mutex     // Mutex to synchronize access to the rate limiter.
}

// fillBucket fills the bucket with tokens based on the elapsed time since the last fill.
func (rl *RateLimiter) fillBucket(sourceKey string) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := int(time.Now().UnixMilli())
	elapsed := now - rl.lastFill[sourceKey]

	rl.bucket[sourceKey] += int(float64(elapsed) * rl.maxRatePerMillisecond)
	if rl.bucket[sourceKey] > rl.maxBurst {
		rl.bucket[sourceKey] = rl.maxBurst
	}
	rl.lastFill[sourceKey] = now
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
		maxRatePerMillisecond: float64(options.MaxRatePerSecond) / 1000.0,
		maxBurst:              options.MaxBurst,
		bucket:                make(map[string]int),
		lastFill:              make(map[string]int),
	}
}
