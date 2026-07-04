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

type TemplateRepository struct {
	client *pgsql.PgClient
}

func NewTemplateRepository(client *pgsql.PgClient) port.TemplateRepository {
	return &TemplateRepository{client: client}
}

func (r *TemplateRepository) Create(ctx context.Context, template *domain.Template) error {
	query := `
		INSERT INTO templates (workspace_id, name, unique_key, channel_type, category, subject, content, is_active, is_default)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, created_at, updated_at
	`
	err := r.client.GetDB(ctx).QueryRowContext(ctx, query,
		template.WorkspaceID, template.Name, template.UniqueKey, template.ChannelType,
		nullIfEmpty(template.Category), nullIfEmpty(template.Subject), template.ContentHTML, template.IsActive, template.IsDefault,
	).Scan(&template.ID, &template.CreatedAt, &template.UpdatedAt)
	if err != nil {
		slog.ErrorContext(ctx, "database error: failed to create template", "error", err, "workspace_id", template.WorkspaceID, "unique_key", template.UniqueKey)
		return err
	}
	return nil
}

func (r *TemplateRepository) GetByWorkspaceAndKey(ctx context.Context, workspaceID, uniqueKey string) (*domain.Template, error) {
	query := `
		SELECT id, workspace_id, name, unique_key, channel_type, category, subject, content, is_active, is_default, created_at, updated_at
		FROM templates
		WHERE workspace_id = $1 AND unique_key = $2
	`
	var t domain.Template
	var subject, category sql.NullString
	err := r.client.GetDB(ctx).QueryRowContext(ctx, query, workspaceID, uniqueKey).
		Scan(&t.ID, &t.WorkspaceID, &t.Name, &t.UniqueKey, &t.ChannelType, &category, &subject,
			&t.ContentHTML, &t.IsActive, &t.IsDefault, &t.CreatedAt, &t.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("template key=%s: %w", uniqueKey, port.ErrNotFound)
	}
	if err != nil {
		slog.ErrorContext(ctx, "database error: failed to get template by workspace and key", "error", err, "workspace_id", workspaceID, "unique_key", uniqueKey)
		return nil, err
	}
	if subject.Valid {
		t.Subject = subject.String
	}
	if category.Valid {
		t.Category = category.String
	}
	return &t, nil
}

func (r *TemplateRepository) GetByID(ctx context.Context, id string) (*domain.Template, error) {
	query := `
		SELECT id, workspace_id, name, unique_key, channel_type, category, subject, content, is_active, is_default, created_at, updated_at
		FROM templates WHERE id = $1
	`
	var t domain.Template
	var subject, category sql.NullString
	err := r.client.GetDB(ctx).QueryRowContext(ctx, query, id).
		Scan(&t.ID, &t.WorkspaceID, &t.Name, &t.UniqueKey, &t.ChannelType, &category, &subject,
			&t.ContentHTML, &t.IsActive, &t.IsDefault, &t.CreatedAt, &t.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("template %s: %w", id, port.ErrNotFound)
	}
	if err != nil {
		slog.ErrorContext(ctx, "database error: failed to get template by id", "error", err, "id", id)
		return nil, err
	}
	if subject.Valid {
		t.Subject = subject.String
	}
	if category.Valid {
		t.Category = category.String
	}
	return &t, nil
}

func (r *TemplateRepository) ListByWorkspace(ctx context.Context, workspaceID string) ([]domain.Template, error) {
	rows, err := r.client.GetDB(ctx).QueryContext(ctx, `
		SELECT id, workspace_id, name, unique_key, channel_type, category, subject, content, is_active, is_default, created_at, updated_at
		FROM templates
		WHERE workspace_id = $1
		ORDER BY name ASC
	`, workspaceID)
	if err != nil {
		slog.ErrorContext(ctx, "database error: failed to list templates for workspace", "error", err, "workspace_id", workspaceID)
		return nil, err
	}
	defer rows.Close() //nolint:errcheck

	var out []domain.Template
	for rows.Next() {
		var t domain.Template
		var subject, category sql.NullString
		if err := rows.Scan(&t.ID, &t.WorkspaceID, &t.Name, &t.UniqueKey, &t.ChannelType, &category, &subject,
			&t.ContentHTML, &t.IsActive, &t.IsDefault, &t.CreatedAt, &t.UpdatedAt); err != nil {
			slog.ErrorContext(ctx, "database error: failed to scan template in list", "error", err, "workspace_id", workspaceID)
			return nil, err
		}
		if subject.Valid {
			t.Subject = subject.String
		}
		if category.Valid {
			t.Category = category.String
		}
		out = append(out, t)
	}
	if err = rows.Err(); err != nil {
		slog.ErrorContext(ctx, "database error: rows iteration failed for templates", "error", err, "workspace_id", workspaceID)
		return nil, err
	}
	return out, nil
}

func (r *TemplateRepository) Update(ctx context.Context, template *domain.Template) error {
	_, err := r.client.GetDB(ctx).ExecContext(ctx, `
		UPDATE templates SET
			name = $2,
			unique_key = $3,
			channel_type = $4,
			category = $5,
			subject = $6,
			content = $7,
			is_active = $8,
			is_default = $9,
			updated_at = NOW()
		WHERE id = $1 AND workspace_id = $10
	`, template.ID, template.Name, template.UniqueKey, template.ChannelType,
		nullIfEmpty(template.Category), nullIfEmpty(template.Subject), template.ContentHTML,
		template.IsActive, template.IsDefault, template.WorkspaceID)
	if err != nil {
		slog.ErrorContext(ctx, "database error: failed to update template", "error", err, "id", template.ID, "workspace_id", template.WorkspaceID)
		return err
	}
	return nil
}

func (r *TemplateRepository) Delete(ctx context.Context, id string) error {
	res, err := r.client.GetDB(ctx).ExecContext(ctx, `DELETE FROM templates WHERE id = $1`, id)
	if err != nil {
		slog.ErrorContext(ctx, "database error: failed to delete template", "error", err, "id", id)
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("template %s: %w", id, port.ErrNotFound)
	}
	return nil
}
