package ratelimitermiddleware

import (
	"net/http"
	"strconv"

	"github.com/rcdmk/go-ratelimiter"
)

// Options represents the options for configuring rate limiter middleware.
type Options struct {
	MaxRatePerSecond int    // The maximum rate of events allowed per second.
	MaxBurst         int    // The maximum number of events that can be bursted.
	SourceHeaderKey  string // The key in the request header to use as the rate limiting key.
}

// StdLib wraps a standard lib handler in a rate limiter middleware.
// It returns an http.Handler that applies rate limiting to incoming requests, compatible with standard lib and frameworks that accept the same interface.
func StdLib(next http.Handler, options Options) http.Handler {
	limiter := ratelimiter.New(ratelimiter.Options{
		MaxRatePerSecond: options.MaxRatePerSecond,
		MaxBurst:         options.MaxBurst,
	})

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := r.Header.Get(options.SourceHeaderKey)

		if !limiter.Allow(key) {
			http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
			w.Header().Add("Retry-After", "1")
			w.Header().Add("RateLimit-Limit", strconv.Itoa(options.MaxRatePerSecond))
			w.Header().Add("RateLimit-Remaining", "0")
			w.Header().Add("RateLimit-Reset", "1")
			return
		}

		next.ServeHTTP(w, r)
	})
}
