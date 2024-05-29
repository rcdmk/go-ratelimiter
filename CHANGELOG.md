# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Changed

- `ratelimitermiddleware.StdLib` now takes cache options to enable replacing the cache implementation.

## [0.1.0]

### Added

- `ratelimiter` package brings a generic rate limiter implementation that can be used as base for any use case.
- `ratelimitermiddleware.StdLib` implements the rate limiter for the standard lib `http` package.
- `ratelimitermiddleware.StdLib` adds standard `RateLimit-Limit`, `RateLimit-Remaining`, `RateLimit-Reset` and `Retry-After` headers to the response when blocking a request.
- `cache.GetterSetter` interface provices a mean for specifying a cache implementation for the rate limiter.
- `cache.InMemory` is a default implementation for caching that supports TTL and uses standard maps as backing storage. This implementation is used by default if no implementation is provided.

[Unreleased]: https://github.com/rcdmk/go-ratelimiter/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/rcdmk/go-ratelimiter/commits/v0.1.0/
