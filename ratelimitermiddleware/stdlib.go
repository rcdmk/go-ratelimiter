package ratelimitermiddleware

import (
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/rcdmk/go-ratelimiter"
	"github.com/rcdmk/go-ratelimiter/cache"
)

// Options represents the options for configuring rate limiter middleware.
type Options struct {
	MaxRatePerSecond int                // The maximum rate of events allowed per second.
	MaxBurst         int                // The maximum number of events that can be bursted.
	SourceHeaderKey  string             // The key in the request header to use as the rate limiting key.
	Cache            cache.GetterSetter // The cache to use for storing rate limiting data.
	CacheTTL         time.Duration      // The time-to-live for rate limiting data in the cache.
}

// StdLib wraps a standard lib handler in a rate limiter middleware.
// It returns an http.Handler that applies rate limiting to incoming requests, compatible with standard lib and frameworks that accept the same interface.
func StdLib(next http.Handler, options Options) http.Handler {
	limiter := ratelimiter.New(ratelimiter.Options{
		MaxRatePerSecond: options.MaxRatePerSecond,
		MaxBurst:         options.MaxBurst,
		Cache:            options.Cache,
		CacheTTL:         options.CacheTTL,
	})

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := r.Header.Get(options.SourceHeaderKey)

		burstResetSeconds := strconv.FormatFloat(math.Floor(float64(options.MaxBurst)/float64(options.MaxRatePerSecond)), 'f', 0, 64)

		w.Header().Add("RateLimit-Limit", strconv.Itoa(options.MaxRatePerSecond))
		w.Header().Add("RateLimit-Reset", burstResetSeconds)
		w.Header().Add("RateLimit-Remaining", strconv.Itoa(limiter.Remaining(key)))

		if !limiter.Allow(key) {
			w.Header().Add("Retry-After", burstResetSeconds)
			w.Header().Add("RateLimit-Remaining", "0")
			http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
