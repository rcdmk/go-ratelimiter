# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- `ratelimiter` package brings a generic rate limiter implementation that can be used as base for any use case.
- `ratelimitermiddleware.StdLib` implements the rate limiter for the standard lib `http` package.
- `ratelimitermiddleware.StdLib` adds standard `RateLimit-Limit`, `RateLimit-Remaining`, `RateLimit-Reset` and `Retry-After` headers to the response when blocking a request.

[Unreleased]: https://github.com/rcdmk/go-ratelimiter/compare/d8b2554...HEAD
