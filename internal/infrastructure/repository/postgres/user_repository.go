package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

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
		INSERT INTO users (first_name, last_name, email, password_hash)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at
	`
	err := r.client.GetDB(ctx).QueryRowContext(ctx, query,
		u.FirstName, u.LastName, u.Email, u.PasswordHash,
	).Scan(&u.ID, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		slog.ErrorContext(ctx, "database error: failed to create user", "error", err, "email", u.Email)
	}
	return err
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `
		SELECT id, first_name, last_name, email, password_hash, email_verified, created_at, updated_at
		FROM users WHERE email = $1
	`
	var u domain.User
	err := r.client.GetDB(ctx).QueryRowContext(ctx, query, email).
		Scan(&u.ID, &u.FirstName, &u.LastName, &u.Email, &u.PasswordHash, &u.EmailVerified, &u.CreatedAt, &u.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("user %s: %w", email, port.ErrNotFound)
	}
	if err != nil {
		slog.ErrorContext(ctx, "database error: failed to get user by email", "error", err, "email", email)
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	query := `
		SELECT id, first_name, last_name, email, password_hash, email_verified, created_at, updated_at
		FROM users WHERE id = $1
	`
	var u domain.User
	err := r.client.GetDB(ctx).QueryRowContext(ctx, query, id).
		Scan(&u.ID, &u.FirstName, &u.LastName, &u.Email, &u.PasswordHash, &u.EmailVerified, &u.CreatedAt, &u.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("user %s: %w", id, port.ErrNotFound)
	}
	if err != nil {
		slog.ErrorContext(ctx, "database error: failed to get user by id", "error", err, "user_id", id)
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) SetEmailVerified(ctx context.Context, id string) error {
	query := `
		UPDATE users
		SET email_verified = TRUE, updated_at = NOW()
		WHERE id = $1
	`
	_, err := r.client.GetDB(ctx).ExecContext(ctx, query, id)
	if err != nil {
		slog.ErrorContext(ctx, "database error: failed to set email verified", "error", err, "user_id", id)
	}
	return err
}
