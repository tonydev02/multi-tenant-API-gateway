package auth

import "time"

// AdminUser represents an authenticated admin account.
type AdminUser struct {
	ID           int64     `json:"id"`
	TenantID     int64     `json:"tenant_id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
}

// APIKeyRecord stores key metadata and hashed secret.
type APIKeyRecord struct {
	ID        int64      `json:"id"`
	TenantID  int64      `json:"tenant_id"`
	Name      string     `json:"name"`
	Prefix    string     `json:"prefix"`
	KeyHash   string     `json:"-"`
	RevokedAt *time.Time `json:"revoked_at,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
}

// Claims is the JWT payload used by admin APIs.
type Claims struct {
	Subject  int64  `json:"sub"`
	TenantID int64  `json:"tenant_id"`
	Email    string `json:"email"`
	Issuer   string `json:"iss"`
	IssuedAt int64  `json:"iat"`
	Expiry   int64  `json:"exp"`
}
