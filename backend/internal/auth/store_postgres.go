package auth

import (
	"context"
	"database/sql"
	"fmt"
)

// Store provides auth persistence in PostgreSQL.
type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) CreateAdminUser(ctx context.Context, tenantID int64, email, passwordHash string) (AdminUser, error) {
	var user AdminUser
	err := s.db.QueryRowContext(ctx, `
		INSERT INTO admin_users (tenant_id, email, password_hash)
		VALUES ($1, $2, $3)
		RETURNING id, tenant_id, email, password_hash, created_at
	`, tenantID, email, passwordHash).Scan(
		&user.ID,
		&user.TenantID,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
	)
	if err != nil {
		return AdminUser{}, fmt.Errorf("create admin user: %w", err)
	}
	return user, nil
}

func (s *Store) GetAdminByEmail(ctx context.Context, email string) (AdminUser, error) {
	var user AdminUser
	err := s.db.QueryRowContext(ctx, `
		SELECT id, tenant_id, email, password_hash, created_at
		FROM admin_users
		WHERE email = $1
	`, email).Scan(
		&user.ID,
		&user.TenantID,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
	)
	if err != nil {
		return AdminUser{}, fmt.Errorf("get admin by email: %w", err)
	}
	return user, nil
}

func (s *Store) CreateAPIKey(ctx context.Context, tenantID int64, name, prefix, hash string) (APIKeyRecord, error) {
	var key APIKeyRecord
	err := s.db.QueryRowContext(ctx, `
		INSERT INTO api_keys (tenant_id, name, key_prefix, key_hash)
		VALUES ($1, $2, $3, $4)
		RETURNING id, tenant_id, name, key_prefix, key_hash, revoked_at, created_at
	`, tenantID, name, prefix, hash).Scan(
		&key.ID,
		&key.TenantID,
		&key.Name,
		&key.Prefix,
		&key.KeyHash,
		&key.RevokedAt,
		&key.CreatedAt,
	)
	if err != nil {
		return APIKeyRecord{}, fmt.Errorf("create api key: %w", err)
	}
	return key, nil
}

func (s *Store) ListAPIKeysByTenant(ctx context.Context, tenantID int64) ([]APIKeyRecord, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, tenant_id, name, key_prefix, key_hash, revoked_at, created_at
		FROM api_keys
		WHERE tenant_id = $1
		ORDER BY id DESC
	`, tenantID)
	if err != nil {
		return nil, fmt.Errorf("list api keys: %w", err)
	}
	defer rows.Close()

	keys := make([]APIKeyRecord, 0)
	for rows.Next() {
		var key APIKeyRecord
		if err := rows.Scan(&key.ID, &key.TenantID, &key.Name, &key.Prefix, &key.KeyHash, &key.RevokedAt, &key.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan api key: %w", err)
		}
		keys = append(keys, key)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate api keys: %w", err)
	}
	return keys, nil
}

func (s *Store) RevokeAPIKey(ctx context.Context, tenantID, keyID int64) error {
	result, err := s.db.ExecContext(ctx, `
		UPDATE api_keys
		SET revoked_at = NOW()
		WHERE id = $1 AND tenant_id = $2 AND revoked_at IS NULL
	`, keyID, tenantID)
	if err != nil {
		return fmt.Errorf("revoke api key: %w", err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("revoke api key rows affected: %w", err)
	}
	if affected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (s *Store) GetAPIKeyByPrefix(ctx context.Context, prefix string) (APIKeyRecord, error) {
	var key APIKeyRecord
	err := s.db.QueryRowContext(ctx, `
		SELECT id, tenant_id, name, key_prefix, key_hash, revoked_at, created_at
		FROM api_keys
		WHERE key_prefix = $1
	`, prefix).Scan(
		&key.ID,
		&key.TenantID,
		&key.Name,
		&key.Prefix,
		&key.KeyHash,
		&key.RevokedAt,
		&key.CreatedAt,
	)
	if err != nil {
		return APIKeyRecord{}, fmt.Errorf("get api key by prefix: %w", err)
	}
	return key, nil
}
