package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
)

const keyPrefixMarker = "gk_"

var ErrInvalidAPIKey = errors.New("invalid api key")

func GenerateAPIKey() (raw, prefix, hash string, err error) {
	prefixBytes := make([]byte, 4)
	secretBytes := make([]byte, 24)
	if _, err = rand.Read(prefixBytes); err != nil {
		return "", "", "", fmt.Errorf("generate key prefix: %w", err)
	}
	if _, err = rand.Read(secretBytes); err != nil {
		return "", "", "", fmt.Errorf("generate key secret: %w", err)
	}

	prefix = hex.EncodeToString(prefixBytes)
	secret := hex.EncodeToString(secretBytes)
	raw = keyPrefixMarker + prefix + "." + secret
	hash = HashAPIKey(raw)
	return raw, prefix, hash, nil
}

func HashAPIKey(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}

func ParseAPIKeyPrefix(raw string) (string, error) {
	if !strings.HasPrefix(raw, keyPrefixMarker) {
		return "", ErrInvalidAPIKey
	}
	parts := strings.SplitN(strings.TrimPrefix(raw, keyPrefixMarker), ".", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", ErrInvalidAPIKey
	}
	return parts[0], nil
}

type apiKeyReader interface {
	GetAPIKeyByPrefix(ctx context.Context, prefix string) (APIKeyRecord, error)
}

// APIKeyAuthenticator validates consumer API keys and resolves tenant context.
type APIKeyAuthenticator struct {
	store apiKeyReader
}

func NewAPIKeyAuthenticator(store apiKeyReader) *APIKeyAuthenticator {
	return &APIKeyAuthenticator{store: store}
}

func (a *APIKeyAuthenticator) Authenticate(ctx context.Context, raw string) (APIKeyRecord, error) {
	prefix, err := ParseAPIKeyPrefix(raw)
	if err != nil {
		return APIKeyRecord{}, ErrInvalidAPIKey
	}

	record, err := a.store.GetAPIKeyByPrefix(ctx, prefix)
	if err != nil {
		return APIKeyRecord{}, ErrInvalidAPIKey
	}

	hash := HashAPIKey(raw)
	if subtle.ConstantTimeCompare([]byte(record.KeyHash), []byte(hash)) != 1 {
		return APIKeyRecord{}, ErrInvalidAPIKey
	}
	if record.RevokedAt != nil {
		return APIKeyRecord{}, ErrInvalidAPIKey
	}

	return record, nil
}
