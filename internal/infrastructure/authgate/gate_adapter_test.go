package authgate

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	gogate "github.com/weprodev/wpd-gogate"

	"github.com/weprodev/go-pkg/pgsql"
)

type staticDBProvider struct {
	db pgsql.DBTX
}

func (s staticDBProvider) GetDB(context.Context) pgsql.DBTX {
	return s.db
}

func TestGoGateAdapter_GetPermissionsForTeams_emptyInput(t *testing.T) {
	t.Parallel()

	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })

	adapter := NewGoGateAdapter(nil, staticDBProvider{db: db}, gogate.DefaultConfig())
	out, err := adapter.GetPermissionsForTeams(context.Background(), "User", "user-1", nil)
	if err != nil {
		t.Fatalf("GetPermissionsForTeams: %v", err)
	}
	if len(out) != 0 {
		t.Fatalf("expected empty map, got %+v", out)
	}
}

func TestGoGateAdapter_GetPermissionsForTeams_dedupesPermissions(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })

	rows := sqlmock.NewRows([]string{"team_id", "permission"}).
		AddRow("team-1", "logs.read").
		AddRow("team-1", "logs.read").
		AddRow("team-1", "settings.write")

	mock.ExpectQuery(`SELECT team_id::text, permission`).
		WithArgs("User", "user-1", sqlmock.AnyArg()).
		WillReturnRows(rows)

	adapter := NewGoGateAdapter(nil, staticDBProvider{db: db}, gogate.DefaultConfig())
	out, err := adapter.GetPermissionsForTeams(context.Background(), "User", "user-1", []string{"team-1"})
	if err != nil {
		t.Fatalf("GetPermissionsForTeams: %v", err)
	}
	perms := out["team-1"]
	if len(perms) != 2 || perms[0] != "logs.read" || perms[1] != "settings.write" {
		t.Fatalf("unexpected permissions: %+v", perms)
	}
}

func TestDedupeSortedStrings(t *testing.T) {
	t.Parallel()

	got := dedupeSortedStrings([]string{"a", "a", "b", "c", "c"})
	want := []string{"a", "b", "c"}
	if len(got) != len(want) {
		t.Fatalf("got %v want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("got %v want %v", got, want)
		}
	}
}
