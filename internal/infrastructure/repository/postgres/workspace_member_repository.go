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

type WorkspaceMemberRepository struct {
	client *pgsql.PgClient
}

func NewWorkspaceMemberRepository(client *pgsql.PgClient) port.WorkspaceMemberRepository {
	return &WorkspaceMemberRepository{client: client}
}

func (r *WorkspaceMemberRepository) Add(ctx context.Context, workspaceID, userID, role string) error {
	_, err := r.client.GetDB(ctx).ExecContext(ctx, `
		INSERT INTO workspace_members (workspace_id, user_id, role)
		VALUES ($1, $2, $3)
		ON CONFLICT (workspace_id, user_id) DO UPDATE SET role = EXCLUDED.role
	`, workspaceID, userID, role)
	return err
}

func (r *WorkspaceMemberRepository) Remove(ctx context.Context, workspaceID, userID string) error {
	res, err := r.client.GetDB(ctx).ExecContext(ctx,
		`DELETE FROM workspace_members WHERE workspace_id = $1 AND user_id = $2`,
		workspaceID, userID,
	)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("workspace member workspace=%s user=%s: %w", workspaceID, userID, port.ErrNotFound)
	}
	return nil
}

func (r *WorkspaceMemberRepository) GetRole(ctx context.Context, workspaceID, userID string) (string, error) {
	var role string
	err := r.client.GetDB(ctx).QueryRowContext(ctx,
		`SELECT role FROM workspace_members WHERE workspace_id = $1 AND user_id = $2`,
		workspaceID, userID,
	).Scan(&role)
	if errors.Is(err, sql.ErrNoRows) {
		return "", fmt.Errorf("workspace member workspace=%s user=%s: %w", workspaceID, userID, port.ErrNotFound)
	}
	return role, err
}

func (r *WorkspaceMemberRepository) ListMembers(ctx context.Context, workspaceID string) ([]domain.WorkspaceMember, error) {
	rows, err := r.client.GetDB(ctx).QueryContext(ctx, `
		SELECT wm.workspace_id, wm.user_id, wm.role, wm.joined_at, u.email, COALESCE(u.display_name, '')
		FROM workspace_members wm
		INNER JOIN users u ON u.id = wm.user_id
		WHERE wm.workspace_id = $1
		ORDER BY wm.joined_at ASC
	`, workspaceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close() //nolint:errcheck

	var out []domain.WorkspaceMember
	for rows.Next() {
		var m domain.WorkspaceMember
		if err := rows.Scan(&m.WorkspaceID, &m.UserID, &m.Role, &m.JoinedAt, &m.UserEmail, &m.DisplayName); err != nil {
			return nil, err
		}
		out = append(out, m)
	}
	return out, rows.Err()
}
