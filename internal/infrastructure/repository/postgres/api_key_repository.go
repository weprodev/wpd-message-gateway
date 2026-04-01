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

type APIKeyRepository struct {
	client *pgsql.PgClient
}

func NewAPIKeyRepository(client *pgsql.PgClient) port.APIKeyRepository {
	return &APIKeyRepository{client: client}
}

func (r *APIKeyRepository) Create(ctx context.Context, apiKey *domain.APIKey) error {
	query := `
		INSERT INTO api_keys (workspace_id, client_id, client_secret_hash, name, is_active, expires_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at
	`
	err := r.client.GetDB(ctx).QueryRowContext(ctx, query,
		apiKey.WorkspaceID, apiKey.ClientID, apiKey.ClientSecretHash, apiKey.Name, apiKey.IsActive, apiKey.ExpiresAt,
	).Scan(&apiKey.ID, &apiKey.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (r *APIKeyRepository) GetByClientID(ctx context.Context, clientID string) (*domain.APIKey, error) {
	query := `
		SELECT id, workspace_id, client_id, client_secret_hash, name, is_active, last_used_at, created_at, expires_at
		FROM api_keys
		WHERE client_id = $1
	`
	var key domain.APIKey
	var lastUsed sql.NullTime
	err := r.client.GetDB(ctx).QueryRowContext(ctx, query, clientID).
		Scan(&key.ID, &key.WorkspaceID, &key.ClientID, &key.ClientSecretHash, &key.Name, &key.IsActive,
			&lastUsed, &key.CreatedAt, &key.ExpiresAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("api key %s: %w", clientID, port.ErrNotFound)
	}
	if err != nil {
		return nil, err
	}
	if lastUsed.Valid {
		t := lastUsed.Time
		key.LastUsedAt = &t
	}
	return &key, nil
}

func (r *APIKeyRepository) UpdateLastUsedAt(ctx context.Context, id string) error {
	_, err := r.client.GetDB(ctx).ExecContext(ctx,
		`UPDATE api_keys SET last_used_at = NOW() WHERE id = $1`,
		id,
	)
	return err
}

func (r *APIKeyRepository) GetByID(ctx context.Context, id string) (*domain.APIKey, error) {
	query := `
		SELECT id, workspace_id, client_id, client_secret_hash, name, is_active, last_used_at, created_at, expires_at
		FROM api_keys WHERE id = $1
	`
	var key domain.APIKey
	var lastUsed sql.NullTime
	var exp sql.NullTime
	err := r.client.GetDB(ctx).QueryRowContext(ctx, query, id).
		Scan(&key.ID, &key.WorkspaceID, &key.ClientID, &key.ClientSecretHash, &key.Name, &key.IsActive,
			&lastUsed, &key.CreatedAt, &exp)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("api key %s: %w", id, port.ErrNotFound)
	}
	if err != nil {
		return nil, err
	}
	if lastUsed.Valid {
		t := lastUsed.Time
		key.LastUsedAt = &t
	}
	if exp.Valid {
		t := exp.Time
		key.ExpiresAt = &t
	}
	return &key, nil
}

func (r *APIKeyRepository) ListByWorkspace(ctx context.Context, workspaceID string) ([]domain.APIKey, error) {
	rows, err := r.client.GetDB(ctx).QueryContext(ctx, `
		SELECT id, workspace_id, client_id, client_secret_hash, name, is_active, last_used_at, created_at, expires_at
		FROM api_keys
		WHERE workspace_id = $1
		ORDER BY created_at DESC
	`, workspaceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close() //nolint:errcheck

	var out []domain.APIKey
	for rows.Next() {
		var key domain.APIKey
		var lastUsed, exp sql.NullTime
		if err := rows.Scan(&key.ID, &key.WorkspaceID, &key.ClientID, &key.ClientSecretHash, &key.Name, &key.IsActive,
			&lastUsed, &key.CreatedAt, &exp); err != nil {
			return nil, err
		}
		if lastUsed.Valid {
			t := lastUsed.Time
			key.LastUsedAt = &t
		}
		if exp.Valid {
			t := exp.Time
			key.ExpiresAt = &t
		}
		out = append(out, key)
	}
	return out, rows.Err()
}

func (r *APIKeyRepository) Delete(ctx context.Context, id string) error {
	res, err := r.client.GetDB(ctx).ExecContext(ctx, `DELETE FROM api_keys WHERE id = $1`, id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("api key %s: %w", id, port.ErrNotFound)
	}
	return nil
}

func (r *APIKeyRepository) UpdateSecret(ctx context.Context, id, clientID, secretHash string) error {
	res, err := r.client.GetDB(ctx).ExecContext(ctx, `
		UPDATE api_keys SET client_id = $2, client_secret_hash = $3 WHERE id = $1
	`, id, clientID, secretHash)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("api key %s: %w", id, port.ErrNotFound)
	}
	return nil
}
