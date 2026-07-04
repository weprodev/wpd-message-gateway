package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/weprodev/wpd-message-gateway/internal/core/domain"
)

type queryRower interface {
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

func lookupRoleIDByName(ctx context.Context, db queryRower, roleName string) (string, error) {
	var roleID string
	err := db.QueryRowContext(ctx, `
		SELECT id FROM roles WHERE name = $1 AND guard_name = $2
	`, roleName, domain.RBACGuardName).Scan(&roleID)
	if err != nil {
		slog.ErrorContext(ctx, "database error: lookup role failed", "error", err, "role", roleName)
		return "", fmt.Errorf("wpd-message-gateway: lookup role %q: %w", roleName, err)
	}
	return roleID, nil
}
