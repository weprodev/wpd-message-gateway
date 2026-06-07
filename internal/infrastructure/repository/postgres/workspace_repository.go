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

type WorkspaceRepository struct {
	client *pgsql.PgClient
}

func NewWorkspaceRepository(client *pgsql.PgClient) port.WorkspaceRepository {
	return &WorkspaceRepository{client: client}
}

func scanWorkspace(scanner interface {
	Scan(dest ...any) error
}) (domain.Workspace, error) {
	var w domain.Workspace
	var hashedPin, iconKey sql.NullString
	err := scanner.Scan(&w.ID, &w.Name, &w.UniqueKey, &w.AdminEmail, &w.Status, &w.Visibility,
		&hashedPin, &iconKey, &w.CreatedAt, &w.UpdatedAt)
	if err != nil {
		return w, err
	}
	if hashedPin.Valid {
		w.HashedPin = hashedPin.String
	}
	if iconKey.Valid {
		w.IconKey = iconKey.String
	}
	return w, nil
}

func (r *WorkspaceRepository) Create(ctx context.Context, workspace *domain.Workspace) error {
	query := `
		INSERT INTO workspaces (name, slug, owner_id, status, is_private, hashed_pin_code, icon_key)
		VALUES ($1, $2, (SELECT id FROM users WHERE email = $3), $4, $5 = 'private', $6, $7)
		RETURNING id, created_at, updated_at
	`
	var hashedPin, icon interface{}
	if workspace.HashedPin != "" {
		hashedPin = workspace.HashedPin
	}
	if workspace.IconKey != "" {
		icon = workspace.IconKey
	}
	err := r.client.GetDB(ctx).QueryRowContext(ctx, query,
		workspace.Name, workspace.UniqueKey, workspace.AdminEmail, workspace.Status, workspace.Visibility, hashedPin, icon,
	).Scan(&workspace.ID, &workspace.CreatedAt, &workspace.UpdatedAt)
	if err != nil {
		slog.ErrorContext(ctx, "database error: failed to create workspace", "error", err, "name", workspace.Name, "unique_key", workspace.UniqueKey)
	}
	return err
}

func (r *WorkspaceRepository) GetByID(ctx context.Context, id string) (*domain.Workspace, error) {
	query := `
		SELECT w.id, w.name, w.slug, u.email, w.status,
		       CASE WHEN w.is_private THEN 'private' ELSE 'public' END,
		       w.hashed_pin_code, w.icon_key, w.created_at, w.updated_at
		FROM workspaces w
		INNER JOIN users u ON u.id = w.owner_id
		WHERE w.id = $1
	`
	w, err := scanWorkspace(r.client.GetDB(ctx).QueryRowContext(ctx, query, id))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("workspace %s: %w", id, port.ErrNotFound)
	}
	if err != nil {
		slog.ErrorContext(ctx, "database error: failed to get workspace by id", "error", err, "id", id)
		return nil, err
	}
	return &w, nil
}

func (r *WorkspaceRepository) GetByUniqueKey(ctx context.Context, uniqueKey string) (*domain.Workspace, error) {
	query := `
		SELECT w.id, w.name, w.slug, u.email, w.status,
		       CASE WHEN w.is_private THEN 'private' ELSE 'public' END,
		       w.hashed_pin_code, w.icon_key, w.created_at, w.updated_at
		FROM workspaces w
		INNER JOIN users u ON u.id = w.owner_id
		WHERE w.slug = $1
	`
	w, err := scanWorkspace(r.client.GetDB(ctx).QueryRowContext(ctx, query, uniqueKey))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("workspace key=%s: %w", uniqueKey, port.ErrNotFound)
	}
	if err != nil {
		slog.ErrorContext(ctx, "database error: failed to get workspace by unique key", "error", err, "unique_key", uniqueKey)
		return nil, err
	}
	return &w, nil
}

func (r *WorkspaceRepository) Update(ctx context.Context, workspace *domain.Workspace) error {
	query := `
		UPDATE workspaces SET
			name = $2,
			owner_id = (SELECT id FROM users WHERE email = $3),
			is_private = ($4 = 'private'),
			icon_key = $5,
			updated_at = NOW()
		WHERE id = $1
	`
	var icon interface{}
	if workspace.IconKey != "" {
		icon = workspace.IconKey
	} else {
		icon = nil
	}
	_, err := r.client.GetDB(ctx).ExecContext(ctx, query,
		workspace.ID, workspace.Name, workspace.AdminEmail, workspace.Visibility, icon,
	)
	if err != nil {
		slog.ErrorContext(ctx, "database error: failed to update workspace", "error", err, "id", workspace.ID)
	}
	return err
}

func (r *WorkspaceRepository) SetStatus(ctx context.Context, id, status string) error {
	_, err := r.client.GetDB(ctx).ExecContext(ctx,
		`UPDATE workspaces SET status = $2, updated_at = NOW() WHERE id = $1`, id, status)
	if err != nil {
		slog.ErrorContext(ctx, "database error: failed to set workspace status", "error", err, "id", id, "status", status)
	}
	return err
}

func (r *WorkspaceRepository) ListForUser(ctx context.Context, userID string) ([]domain.Workspace, error) {
	query := `
		SELECT w.id, w.name, w.slug, u.email, w.status,
		       CASE WHEN w.is_private THEN 'private' ELSE 'public' END,
		       w.hashed_pin_code, w.icon_key, w.created_at, w.updated_at
		FROM workspaces w
		INNER JOIN users u ON u.id = w.owner_id
		INNER JOIN workspace_members wm ON wm.workspace_id = w.id
		WHERE wm.user_id = $1 AND w.status = 'active'
		ORDER BY w.name ASC
	`
	rows, err := r.client.GetDB(ctx).QueryContext(ctx, query, userID)
	if err != nil {
		slog.ErrorContext(ctx, "database error: failed to list workspaces for user", "error", err, "user_id", userID)
		return nil, err
	}
	defer rows.Close() //nolint:errcheck

	var out []domain.Workspace
	for rows.Next() {
		w, err := scanWorkspace(rows)
		if err != nil {
			slog.ErrorContext(ctx, "database error: failed to scan workspace in list", "error", err, "user_id", userID)
			return nil, err
		}
		out = append(out, w)
	}
	if err = rows.Err(); err != nil {
		slog.ErrorContext(ctx, "database error: rows iteration failed for user workspaces", "error", err, "user_id", userID)
		return nil, err
	}
	return out, nil
}
