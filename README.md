# Go Rate Limiter

Go Rate Limiter is a Go package that provides rate limiting functionality with middleware implementations for standard lib and common frameworks.

## Features

- Simple and easy-to-use rate limiter for Go applications.
- Middleware implementations for standard lib and popular frameworks.
- Configurable rate limiting options.
- In-memory caching for rate-limiter can be replaced by your own implementation.

## Installation

Use `go get` to install the package:

```sh
go get github.com/rcdmk/go-ratelimiter
```

## Usage

Import the package and use the `New` method to create a new rate limiter instance for use:

```go
import (
    // ...
    "github.com/rcdmk/go-ratelimiter"
    // ...
)

// ...
    rateLimiter := ratelimiter.New(ratelimiter.Options{
        MaxRatePerSecond: 15,
        MaxBurst: 15,
    })

    if !rateLimiter.Allow("my-operation-name") {
        // over rate limit, deny action and stop execution
        return
    }

    // proceed normaly
// ...
```

### Middleware

**`StdLib`** is a standard lib compatible middleware implementation for limitting requests served through an HTTP server.

It is also compatible with all frameworks that can use standard lib middleware (eg. [chi](https://github.com/go-chi/chi), [Gorilla](https://github.com/gorilla/mux), etc.).

```go
import (
    // ...
    "github.com/rcdmk/go-ratelimiter/ratelimitermiddleware"
    // ...
)

// ...

options := ratelimitermiddleware.Options{
    MaxRatePerSecond: 15,
    MaxBurst:         10,
    SourceHeaderKey:  "Authorization",
}

// wrap your handler with the middleware
rateLimitedHandler := ratelimitermiddleware.StdLib(handler, options)

http.Handle("/my-resource", rateLimitedHandler)

// ...
```
