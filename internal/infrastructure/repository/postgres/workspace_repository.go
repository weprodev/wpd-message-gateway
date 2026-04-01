package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/weprodev/wpd-packages/common"

	"github.com/weprodev/wpd-message-gateway/internal/core/domain"
	"github.com/weprodev/wpd-message-gateway/internal/core/port"
)

type WorkspaceRepository struct {
	client *common.PgClient
}

func NewWorkspaceRepository(client *common.PgClient) port.WorkspaceRepository {
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
		INSERT INTO workspaces (name, unique_key, admin_email, status, visibility, hashed_pin, icon_key)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
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
	return err
}

func (r *WorkspaceRepository) GetByID(ctx context.Context, id string) (*domain.Workspace, error) {
	query := `
		SELECT id, name, unique_key, admin_email, status, visibility, hashed_pin, icon_key, created_at, updated_at
		FROM workspaces WHERE id = $1
	`
	w, err := scanWorkspace(r.client.GetDB(ctx).QueryRowContext(ctx, query, id))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("workspace %s: %w", id, port.ErrNotFound)
	}
	if err != nil {
		return nil, err
	}
	return &w, nil
}

func (r *WorkspaceRepository) GetByUniqueKey(ctx context.Context, uniqueKey string) (*domain.Workspace, error) {
	query := `
		SELECT id, name, unique_key, admin_email, status, visibility, hashed_pin, icon_key, created_at, updated_at
		FROM workspaces WHERE unique_key = $1
	`
	w, err := scanWorkspace(r.client.GetDB(ctx).QueryRowContext(ctx, query, uniqueKey))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("workspace key=%s: %w", uniqueKey, port.ErrNotFound)
	}
	if err != nil {
		return nil, err
	}
	return &w, nil
}

func (r *WorkspaceRepository) Update(ctx context.Context, workspace *domain.Workspace) error {
	query := `
		UPDATE workspaces SET
			name = $2, admin_email = $3, visibility = $4, icon_key = $5, updated_at = NOW()
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
	return err
}

func (r *WorkspaceRepository) SetStatus(ctx context.Context, id, status string) error {
	_, err := r.client.GetDB(ctx).ExecContext(ctx,
		`UPDATE workspaces SET status = $2, updated_at = NOW() WHERE id = $1`, id, status)
	return err
}

func (r *WorkspaceRepository) ListForUser(ctx context.Context, userID string) ([]domain.Workspace, error) {
	query := `
		SELECT w.id, w.name, w.unique_key, w.admin_email, w.status, w.visibility, w.hashed_pin, w.icon_key, w.created_at, w.updated_at
		FROM workspaces w
		INNER JOIN workspace_members wm ON wm.workspace_id = w.id
		WHERE wm.user_id = $1 AND w.status = 'active'
		ORDER BY w.name ASC
	`
	rows, err := r.client.GetDB(ctx).QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close() //nolint:errcheck

	var out []domain.Workspace
	for rows.Next() {
		w, err := scanWorkspace(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, w)
	}
	return out, rows.Err()
}
