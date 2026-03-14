package tenant

import (
	"context"
	"database/sql"
	"fmt"
)

// Store provides tenant persistence against PostgreSQL.
type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) Create(ctx context.Context, name, slug string) (Tenant, error) {
	var t Tenant
	err := s.db.QueryRowContext(ctx, `
		INSERT INTO tenants (name, slug)
		VALUES ($1, $2)
		RETURNING id, name, slug, created_at, updated_at
	`, name, slug).Scan(&t.ID, &t.Name, &t.Slug, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		return Tenant{}, fmt.Errorf("create tenant: %w", err)
	}
	return t, nil
}

func (s *Store) GetByID(ctx context.Context, id int64) (Tenant, error) {
	var t Tenant
	err := s.db.QueryRowContext(ctx, `
		SELECT id, name, slug, created_at, updated_at
		FROM tenants
		WHERE id = $1
	`, id).Scan(&t.ID, &t.Name, &t.Slug, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		return Tenant{}, fmt.Errorf("get tenant by id: %w", err)
	}
	return t, nil
}

func (s *Store) UpdateName(ctx context.Context, id int64, name string) (Tenant, error) {
	var t Tenant
	err := s.db.QueryRowContext(ctx, `
		UPDATE tenants
		SET name = $2
		WHERE id = $1
		RETURNING id, name, slug, created_at, updated_at
	`, id, name).Scan(&t.ID, &t.Name, &t.Slug, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		return Tenant{}, fmt.Errorf("update tenant: %w", err)
	}
	return t, nil
}

func (s *Store) Delete(ctx context.Context, id int64) error {
	result, err := s.db.ExecContext(ctx, `DELETE FROM tenants WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete tenant: %w", err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("delete tenant rows affected: %w", err)
	}
	if affected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (s *Store) GetBySlug(ctx context.Context, slug string) (Tenant, error) {
	var t Tenant
	err := s.db.QueryRowContext(ctx, `
		SELECT id, name, slug, created_at, updated_at
		FROM tenants
		WHERE slug = $1
	`, slug).Scan(&t.ID, &t.Name, &t.Slug, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		return Tenant{}, fmt.Errorf("get tenant by slug: %w", err)
	}
	return t, nil
}
