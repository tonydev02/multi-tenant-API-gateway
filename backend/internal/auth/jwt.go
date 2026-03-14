package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

var ErrInvalidToken = errors.New("invalid token")

type JWTManager struct {
	secret []byte
	issuer string
	expiry time.Duration
	nowFn  func() time.Time
}

func NewJWTManager(secret, issuer string, expiry time.Duration) *JWTManager {
	return &JWTManager{
		secret: []byte(secret),
		issuer: issuer,
		expiry: expiry,
		nowFn:  time.Now,
	}
}

func (m *JWTManager) Issue(user AdminUser) (string, error) {
	now := m.nowFn().UTC().Unix()
	claims := Claims{
		Subject:  user.ID,
		TenantID: user.TenantID,
		Email:    user.Email,
		Issuer:   m.issuer,
		IssuedAt: now,
		Expiry:   now + int64(m.expiry.Seconds()),
	}

	header := map[string]string{"alg": "HS256", "typ": "JWT"}
	headerRaw, err := json.Marshal(header)
	if err != nil {
		return "", fmt.Errorf("marshal token header: %w", err)
	}
	claimsRaw, err := json.Marshal(claims)
	if err != nil {
		return "", fmt.Errorf("marshal token claims: %w", err)
	}

	headerEnc := base64.RawURLEncoding.EncodeToString(headerRaw)
	claimsEnc := base64.RawURLEncoding.EncodeToString(claimsRaw)
	unsigned := headerEnc + "." + claimsEnc
	sig := m.sign(unsigned)

	return unsigned + "." + sig, nil
}

func (m *JWTManager) Parse(token string) (Claims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return Claims{}, ErrInvalidToken
	}

	unsigned := parts[0] + "." + parts[1]
	expectedSig := m.sign(unsigned)
	if !hmac.Equal([]byte(parts[2]), []byte(expectedSig)) {
		return Claims{}, ErrInvalidToken
	}

	claimsBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return Claims{}, ErrInvalidToken
	}

	var claims Claims
	if err := json.Unmarshal(claimsBytes, &claims); err != nil {
		return Claims{}, ErrInvalidToken
	}

	now := m.nowFn().UTC().Unix()
	if claims.Issuer != m.issuer || claims.Expiry <= now || claims.Subject == 0 || claims.TenantID == 0 {
		return Claims{}, ErrInvalidToken
	}

	return claims, nil
}

func (m *JWTManager) sign(value string) string {
	mac := hmac.New(sha256.New, m.secret)
	_, _ = mac.Write([]byte(value))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}
