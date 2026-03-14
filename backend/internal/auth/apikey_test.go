package auth

import "testing"

func TestParseAPIKeyPrefix(t *testing.T) {
	tests := []struct {
		name    string
		raw     string
		want    string
		wantErr bool
	}{
		{name: "valid", raw: "gk_abcd1234.deadbeef", want: "abcd1234"},
		{name: "missing marker", raw: "abcd1234.deadbeef", wantErr: true},
		{name: "missing secret", raw: "gk_abcd1234.", wantErr: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := ParseAPIKeyPrefix(tc.raw)
			if tc.wantErr {
				if err == nil {
					t.Fatal("expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tc.want {
				t.Fatalf("prefix = %q, want %q", got, tc.want)
			}
		})
	}
}
