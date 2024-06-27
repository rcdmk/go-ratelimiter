# Go Rate Limiter Redis Cache

[Go Rate Limiter](https://github.com/rcdmk/go-ratelimiter) is a Go package that provides rate limiting functionality with middleware implementations for standard lib and common frameworks.

This package provides an implementation of the cache provider for that library using [Redis](https://redis.io), with the [go-redis](https://github.com/redis/go-redis) pacakge.

## Installation

Use `go get` to install the package:

```sh
go get github.com/rcdmk/go-ratelimiter
go get github.com/rcdmk/go-ratelimiter/cache/rediscache
```

## Usage

Import the package and use the `New` method to create a new rate limiter instance for use:

```go
import (
    // ...
    "github.com/rcdmk/go-ratelimiter"
    "github.com/rcdmk/go-ratelimiter/rediscache"
    "github.com/redis/go-redis/v9"
    // ...
)

// ...

    redisClient := redis.NewClient(&redis.Options{
		// ...
	})

    redisCache := rediscache.New(redisClient)

    rateLimiter := ratelimiter.New(ratelimiter.Options{
        MaxRatePerSecond: 15,
        MaxBurst:         15,
        Cache:            redisCache,
    })

    if !rateLimiter.Allow("my-operation-name") {
        // over rate limit, deny action and stop execution
        return
    }

    // proceed normaly
// ...
```

### Middleware

[**`StdLib`**](https://github.com/rcdmk/go-ratelimiter/tree/master/ratelimitermiddleware) is a standard lib compatible middleware implementation for limitting requests served through an HTTP server and supports this cache provider.

```go
import (
    // ...
    "github.com/rcdmk/go-ratelimiter/ratelimitermiddleware"
    "github.com/rcdmk/go-ratelimiter/rediscache"
    "github.com/redis/go-redis/v9"
    // ...
)

// ...

redisClient := redis.NewClient(&redis.Options{
    // ...
})

redisCache := rediscache.New(redisClient)

options := ratelimitermiddleware.Options{
    MaxRatePerSecond: 15,
    MaxBurst:         10,
    SourceHeaderKey:  "Authorization",
    Cache:            redisCache,
}

// wrap your handler with the middleware
rateLimitedHandler := ratelimitermiddleware.StdLib(handler, options)

http.Handle("/my-resource", rateLimitedHandler)

// ...
```
