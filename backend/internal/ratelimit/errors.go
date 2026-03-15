package ratelimit

import "errors"

var (
	ErrInvalidPolicy      = errors.New("invalid rate limit policy")
	ErrLimiterUnavailable = errors.New("rate limiter unavailable")
)
