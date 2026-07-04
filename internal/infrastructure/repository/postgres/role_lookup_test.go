package postgres

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/weprodev/wpd-message-gateway/internal/core/domain"
)

func TestLookupRoleIDByName(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })

	mock.ExpectQuery(`SELECT id FROM roles WHERE name = \$1 AND guard_name = \$2`).
		WithArgs("admin", domain.RBACGuardName).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("role-admin-id"))

	roleID, err := lookupRoleIDByName(context.Background(), db, "admin")
	if err != nil {
		t.Fatalf("lookupRoleIDByName: %v", err)
	}
	if roleID != "role-admin-id" {
		t.Fatalf("roleID: got %q", roleID)
	}
}
