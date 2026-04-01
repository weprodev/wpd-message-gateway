package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/weprodev/go-pkg/pgsql"

	"github.com/weprodev/wpd-message-gateway/internal/core/domain"
	"github.com/weprodev/wpd-message-gateway/internal/core/port"
)

type UserRepository struct {
	client *pgsql.PgClient
}

func NewUserRepository(client *pgsql.PgClient) port.UserRepository {
	return &UserRepository{client: client}
}

func (r *UserRepository) Create(ctx context.Context, u *domain.User) error {
	query := `
		INSERT INTO users (email, password_hash, display_name)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at
	`
	return r.client.GetDB(ctx).QueryRowContext(ctx, query,
		u.Email, u.PasswordHash, nullStr(u.DisplayName),
	).Scan(&u.ID, &u.CreatedAt, &u.UpdatedAt)
}

func nullStr(s string) interface{} {
	if s == "" {
		return nil
	}
	return s
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `
		SELECT id, email, password_hash, display_name, created_at, updated_at
		FROM users WHERE email = $1
	`
	var u domain.User
	var dn sql.NullString
	err := r.client.GetDB(ctx).QueryRowContext(ctx, query, email).
		Scan(&u.ID, &u.Email, &u.PasswordHash, &dn, &u.CreatedAt, &u.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("user %s: %w", email, port.ErrNotFound)
	}
	if err != nil {
		return nil, err
	}
	if dn.Valid {
		u.DisplayName = dn.String
	}
	return &u, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	query := `
		SELECT id, email, password_hash, display_name, created_at, updated_at
		FROM users WHERE id = $1
	`
	var u domain.User
	var dn sql.NullString
	err := r.client.GetDB(ctx).QueryRowContext(ctx, query, id).
		Scan(&u.ID, &u.Email, &u.PasswordHash, &dn, &u.CreatedAt, &u.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("user %s: %w", id, port.ErrNotFound)
	}
	if err != nil {
		return nil, err
	}
	if dn.Valid {
		u.DisplayName = dn.String
	}
	return &u, nil
}
