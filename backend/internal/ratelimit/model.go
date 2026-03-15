package ratelimit

import "time"

// Policy defines a fixed-window rate limit.
type Policy struct {
	Requests int64
	Window   time.Duration
}

// Decision is the evaluated outcome for a request.
type Decision struct {
	Allowed   bool
	Limit     int64
	Remaining int64
	ResetAt   time.Time
}
