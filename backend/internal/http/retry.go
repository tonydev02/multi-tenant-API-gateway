package gatewayhttp

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"strings"

	"github.com/lib/pq"
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
	if errors.Is(err, driver.ErrBadConn) || errors.Is(err, context.DeadlineExceeded) || errors.Is(err, sql.ErrConnDone) {
		return true
	}
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		// 08P01 = protocol_violation. In managed postgres/proxy environments this can
		// surface as intermittent bind message/result format mismatches and is often
		// recoverable on immediate retry using a fresh connection.
		if string(pqErr.Code) == "08P01" {
			return true
		}
	}

	msg := strings.ToLower(err.Error())
	transientSubstrings := []string{
		"bad connection",
		"bind message has",
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
