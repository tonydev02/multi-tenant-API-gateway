package gatewayhttp

import (
	"context"
	"database/sql/driver"
	"errors"
	"strings"
)

// shouldRetryReadError identifies transient read failures where a fast one-time retry
// can recover without changing endpoint semantics.
func shouldRetryReadError(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, context.Canceled) {
		return false
	}
	if errors.Is(err, driver.ErrBadConn) || errors.Is(err, context.DeadlineExceeded) {
		return true
	}

	msg := strings.ToLower(err.Error())
	transientSubstrings := []string{
		"bad connection",
		"connection reset by peer",
		"broken pipe",
		"unexpected eof",
		"i/o timeout",
		"server closed the connection unexpectedly",
	}
	for _, needle := range transientSubstrings {
		if strings.Contains(msg, needle) {
			return true
		}
	}

	return false
}
