# Go Rate Limiter

Go Rate Limiter is a Go package that provides rate limiting functionality with middleware implementations for standard lib and common frameworks.

## Features

- Simple and easy-to-use rate limiter for Go applications.
- Middleware implementations for standard lib and popular frameworks.
- Configurable rate limiting options.

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

    if !rateLimiter.Allow() {
        // over rate limit, deny action and stop execution
    }

    // proceed normaly
```
