package ratelimitermiddleware

import (
	"net/http"

	"github.com/rcdmk/go-ratelimiter"
)

// Options represents the options for configuring rate limiter middleware.
type Options struct {
	MaxRatePerSecond int    // The maximum rate of events allowed per second.
	MaxBurst         int    // The maximum number of events that can be bursted.
	SourceHeaderKey  string // The key in the request header to use as the rate limiting key.
}

// stdLib is a middleware that implements rate limiting using the standard library.
type stdLib struct {
	handler         http.Handler
	limiter         *ratelimiter.RateLimiter
	sourceHeaderKey string
}

// ServeHTTP handles the HTTP request by checking if the rate limit is exceeded.
// If the rate limit is exceeded, it returns a "Too Many Requests" error.
// Otherwise, it passes the request to the next handler in the chain.
func (m stdLib) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	key := r.Header.Get(m.sourceHeaderKey)

	if !m.limiter.Allow(key) {
		http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
		// w.Header().Add("Retry-After", strconv.FormatFloat(m.limiter.RetryAfter().Seconds(), 'f', -1, 64))
		// w.Header().Add("X-RateLimit-Limit", strconv.Itoa(m.limiter.Limit()))
		return
	}

	m.handler.ServeHTTP(w, r)
}

// StdLib creates a new instance of the stdLib middleware.
// It takes an http.Handler as the main handler and Options as the rate limiting middleware options to initialize the rate limiter.
// It returns an http.Handler that applies rate limiting to incoming requests, compatible with standard lib and frameworks that accept the same interface.
func StdLib(handler http.Handler, options Options) http.Handler {
	limiter := ratelimiter.New(ratelimiter.Options{
		MaxRatePerSecond: options.MaxRatePerSecond,
		MaxBurst:         options.MaxBurst,
	})
	return stdLib{
		handler:         handler,
		limiter:         limiter,
		sourceHeaderKey: options.SourceHeaderKey,
	}
}
