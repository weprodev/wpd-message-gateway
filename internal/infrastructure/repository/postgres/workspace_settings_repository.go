package postgres

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/weprodev/go-pkg/pgsql"

	"github.com/weprodev/wpd-message-gateway/internal/core/port"
)

type WorkspaceSettingsRepository struct {
	client *pgsql.PgClient
}

func NewWorkspaceSettingsRepository(client *pgsql.PgClient) port.WorkspaceSettingsRepository {
	return &WorkspaceSettingsRepository{client: client}
}

func (r *WorkspaceSettingsRepository) Get(ctx context.Context, workspaceID, key string) (string, error) {
	var v string
	err := r.client.GetDB(ctx).QueryRowContext(ctx,
		`SELECT value FROM workspace_settings WHERE workspace_id = $1 AND key = $2`,
		workspaceID, key,
	).Scan(&v)
	if err == sql.ErrNoRows {
		return "", nil
	}
	if err != nil {
		slog.ErrorContext(ctx, "database error: failed to get workspace setting", "error", err, "workspace_id", workspaceID, "key", key)
	}
	return v, err
}

func (r *WorkspaceSettingsRepository) Delete(ctx context.Context, workspaceID, key string) error {
	_, err := r.client.GetDB(ctx).ExecContext(ctx,
		`DELETE FROM workspace_settings WHERE workspace_id = $1 AND key = $2`,
		workspaceID, key,
	)
	if err != nil {
		slog.ErrorContext(ctx, "database error: failed to delete workspace setting", "error", err, "workspace_id", workspaceID, "key", key)
	}
	return err
}

func (r *WorkspaceSettingsRepository) Set(ctx context.Context, workspaceID, key, value string) error {
	_, err := r.client.GetDB(ctx).ExecContext(ctx, `
		INSERT INTO workspace_settings (workspace_id, key, value)
		VALUES ($1, $2, $3)
		ON CONFLICT (workspace_id, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW()
	`, workspaceID, key, value)
	if err != nil {
		slog.ErrorContext(ctx, "database error: failed to set workspace setting", "error", err, "workspace_id", workspaceID, "key", key)
	}
	return err
}

func (r *WorkspaceSettingsRepository) GetAll(ctx context.Context, workspaceID string) (map[string]string, error) {
	rows, err := r.client.GetDB(ctx).QueryContext(ctx,
		`SELECT key, value FROM workspace_settings WHERE workspace_id = $1`,
		workspaceID,
	)
	if err != nil {
		slog.ErrorContext(ctx, "database error: failed to query all workspace settings", "error", err, "workspace_id", workspaceID)
		return nil, err
	}
	defer rows.Close() //nolint:errcheck

	out := make(map[string]string)
	for rows.Next() {
		var k, v string
		if err := rows.Scan(&k, &v); err != nil {
			slog.ErrorContext(ctx, "database error: failed to scan workspace setting", "error", err, "workspace_id", workspaceID)
			return nil, err
		}
		out[k] = v
	}
	if err := rows.Err(); err != nil {
		slog.ErrorContext(ctx, "database error: rows iteration failed for workspace settings", "error", err, "workspace_id", workspaceID)
		return nil, err
	}
	return out, nil
}
