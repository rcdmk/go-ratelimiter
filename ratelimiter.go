package ratelimiter

import (
	"errors"
	"time"

	"github.com/rcdmk/go-ratelimiter/cache"
)

const (
	bucketKeyPrefix   = "rl:bucket:"
	lastFillKeyPrefix = "rl:fill:"
)

// RateLimiter represents a rate limiter that limits the rate of events, implemented using a token bucket algorithm.
// This implementation assumes cache operations are fast, reliable and concurrency-safe.
type RateLimiter struct {
	maxRatePerMillisecond float64            // The maximum rate of events allowed per millisecond.
	maxBurst              int                // The maximum number of events that can be bursted.
	cache                 cache.GetterSetter // Cache to store the bucket and lastFill values.
	cacheTTL              time.Duration      // The time-to-live for the cache entries.
}

func (rl *RateLimiter) getBucketKeyFor(sourceKey string) string {
	return bucketKeyPrefix + sourceKey
}

func (rl *RateLimiter) getLastFillKeyFor(sourceKey string) string {
	return lastFillKeyPrefix + sourceKey
}

// getBucketFor retrieves the current bucket value for a particular key.
// If cache operations fail, it will always return a full bucket.
func (rl *RateLimiter) getBucketFor(sourceKey string) int {
	bucketKey := rl.getBucketKeyFor(sourceKey)
	bucket, err := rl.cache.Get(bucketKey)
	if err != nil && !errors.Is(err, cache.ErrCacheMiss) {
		// if cache fails, bucket is always full. Allow the event to be executed
		return rl.maxBurst
	}

	return bucket
}

// getLastFillFor retrieves the last fill time for a particular key.
// If cache operations fail, it will always return the current time.
func (rl *RateLimiter) getLastFillFor(sourceKey string) int {
	lastFillKey := rl.getLastFillKeyFor(sourceKey)
	lastFill, err := rl.cache.Get(lastFillKey)
	if err != nil && !errors.Is(err, cache.ErrCacheMiss) {
		// if cache fails, the bucket is always full
		return int(time.Now().UnixMilli())
	}
	return lastFill
}

func (rl *RateLimiter) setBucketFor(sourceKey string, value int) {
	key := rl.getBucketKeyFor(sourceKey)
	_ = rl.cache.SetWithExpiration(key, value, rl.cacheTTL)
}

func (rl *RateLimiter) setLastFillFor(sourceKey string, value int) {
	key := rl.getLastFillKeyFor(sourceKey)
	_ = rl.cache.SetWithExpiration(key, value, rl.cacheTTL)
}

// fillBucket fills the bucket with tokens based on the elapsed time since the last fill.
func (rl *RateLimiter) fillBucket(sourceKey string) {
	now := int(time.Now().UnixMilli())
	lastFill := rl.getLastFillFor(sourceKey)
	elapsed := now - lastFill

	// important to use floating points for partial bucket filling, eg. 10 tokens per second = 1 token per 0.1 seconds
	newTokens := int(float64(elapsed) * rl.maxRatePerMillisecond)

	bucket := rl.getBucketFor(sourceKey) + newTokens

	if bucket > rl.maxBurst {
		bucket = rl.maxBurst
	}

	rl.setBucketFor(sourceKey, bucket)
	rl.setLastFillFor(sourceKey, now)
}

// Remaining returns the number of remaining requests for the given source key.
func (rl *RateLimiter) Remaining(sourceKey string) int {
	rl.fillBucket(sourceKey)

	return rl.getBucketFor(sourceKey)
}

// Allow checks if the rate wasn't exhausted for a particular key to allow or not an event to be executed.
// If cache operations fail, it will always return false.
func (rl *RateLimiter) Allow(sourceKey string) bool {
	rl.fillBucket(sourceKey)

	bucket := rl.getBucketFor(sourceKey)
	if bucket > 0 {
		bucket--
		rl.setBucketFor(sourceKey, bucket)
		return true
	}

	return false
}

// Options represents the options for configuring a RateLimiter.
type Options struct {
	MaxRatePerSecond int                // The maximum rate of events allowed per second.
	MaxBurst         int                // The maximum number of events that can be bursted.
	Cache            cache.GetterSetter // The cache to store the bucket and lastFill values. If not provided, an in-memory cache will be used.
	CacheTTL         time.Duration      // The time-to-live for the cache entries. Default is 10 seconds.
}

// New creates a new ready to use RateLimiter with the specified options.
func New(options Options) *RateLimiter {
	if options.Cache == nil {
		options.Cache = cache.NewInMemory()
	}

	if options.CacheTTL == 0 {
		options.CacheTTL = 10 * time.Second
	}

	return &RateLimiter{
		maxRatePerMillisecond: float64(options.MaxRatePerSecond) / 1000.0,
		maxBurst:              options.MaxBurst,
		cache:                 options.Cache,
		cacheTTL:              options.CacheTTL,
	}
}
