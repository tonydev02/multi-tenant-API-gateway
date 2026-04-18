package gatewayhttp

import (
	"context"
	"database/sql/driver"
	"errors"
	"fmt"
	"testing"
)

func TestShouldRetryReadError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		err  error
		want bool
	}{
		{name: "nil", err: nil, want: false},
		{name: "context canceled", err: context.Canceled, want: false},
		{name: "wrapped context canceled", err: fmt.Errorf("query failed: %w", context.Canceled), want: false},
		{name: "bad conn", err: driver.ErrBadConn, want: true},
		{name: "wrapped bad conn", err: fmt.Errorf("query failed: %w", driver.ErrBadConn), want: true},
		{name: "deadline exceeded", err: context.DeadlineExceeded, want: true},
		{name: "reset by peer", err: errors.New("read tcp: connection reset by peer"), want: true},
		{name: "broken pipe", err: errors.New("write: broken pipe"), want: true},
		{name: "unexpected eof", err: errors.New("unexpected EOF"), want: true},
		{name: "other", err: errors.New("permission denied"), want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := shouldRetryReadError(tt.err); got != tt.want {
				t.Fatalf("shouldRetryReadError() = %v, want %v", got, tt.want)
			}
		})
	}
}
