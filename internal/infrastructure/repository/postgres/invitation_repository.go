package postgres

import (
	"context"
	"log/slog"

	"github.com/weprodev/go-pkg/pgsql"

	"github.com/weprodev/wpd-message-gateway/internal/core/domain"
	"github.com/weprodev/wpd-message-gateway/internal/core/port"
)

type InvitationRepository struct {
	client *pgsql.PgClient
}

func NewInvitationRepository(client *pgsql.PgClient) port.InvitationRepository {
	return &InvitationRepository{client: client}
}

func (r *InvitationRepository) Create(ctx context.Context, inv *domain.Invitation) error {
	roleID, err := lookupRoleIDByName(ctx, r.client.GetDB(ctx), inv.Role)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO invitations (workspace_id, email, role_id, token_hash, expires_at, status)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at
	`
	err = r.client.GetDB(ctx).QueryRowContext(ctx, query,
		inv.WorkspaceID, inv.Email, roleID, inv.TokenHash, inv.ExpiresAt, inv.Status,
	).Scan(&inv.ID, &inv.CreatedAt)
	if err != nil {
		slog.ErrorContext(ctx, "database error: failed to create invitation", "error", err, "workspace_id", inv.WorkspaceID, "email", inv.Email, "role", inv.Role)
	}
	return err
}

func (r *InvitationRepository) ListPendingByWorkspace(ctx context.Context, workspaceID string) ([]domain.Invitation, error) {
	rows, err := r.client.GetDB(ctx).QueryContext(ctx, `
		SELECT i.id, i.workspace_id, i.email, r.name, i.token_hash, i.expires_at, i.status, i.created_at
		FROM invitations i
		INNER JOIN roles r ON r.id = i.role_id
		WHERE i.workspace_id = $1 AND i.status = 'pending' AND i.expires_at > NOW()
		ORDER BY i.created_at DESC
	`, workspaceID)
	if err != nil {
		slog.ErrorContext(ctx, "database error: failed to list pending invitations", "error", err, "workspace_id", workspaceID)
		return nil, err
	}
	defer rows.Close() //nolint:errcheck

	var out []domain.Invitation
	for rows.Next() {
		var inv domain.Invitation
		if err := rows.Scan(&inv.ID, &inv.WorkspaceID, &inv.Email, &inv.Role, &inv.TokenHash, &inv.ExpiresAt, &inv.Status, &inv.CreatedAt); err != nil {
			slog.ErrorContext(ctx, "database error: failed to scan pending invitation", "error", err, "workspace_id", workspaceID)
			return nil, err
		}
		out = append(out, inv)
	}
	if err := rows.Err(); err != nil {
		slog.ErrorContext(ctx, "database error: rows iteration failed for pending invitations", "error", err, "workspace_id", workspaceID)
		return nil, err
	}
	return out, nil
}
