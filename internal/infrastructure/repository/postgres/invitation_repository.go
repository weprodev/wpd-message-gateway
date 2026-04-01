package postgres

import (
	"context"

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
	query := `
		INSERT INTO invitations (workspace_id, email, role, token_hash, expires_at, status)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at
	`
	return r.client.GetDB(ctx).QueryRowContext(ctx, query,
		inv.WorkspaceID, inv.Email, inv.Role, inv.TokenHash, inv.ExpiresAt, inv.Status,
	).Scan(&inv.ID, &inv.CreatedAt)
}

func (r *InvitationRepository) ListPendingByWorkspace(ctx context.Context, workspaceID string) ([]domain.Invitation, error) {
	rows, err := r.client.GetDB(ctx).QueryContext(ctx, `
		SELECT id, workspace_id, email, role, token_hash, expires_at, status, created_at
		FROM invitations
		WHERE workspace_id = $1 AND status = 'pending' AND expires_at > NOW()
		ORDER BY created_at DESC
	`, workspaceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close() //nolint:errcheck

	var out []domain.Invitation
	for rows.Next() {
		var inv domain.Invitation
		if err := rows.Scan(&inv.ID, &inv.WorkspaceID, &inv.Email, &inv.Role, &inv.TokenHash, &inv.ExpiresAt, &inv.Status, &inv.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, inv)
	}
	return out, rows.Err()
}
