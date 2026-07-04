package postgres

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/weprodev/wpd-message-gateway/internal/core/domain"
	"github.com/weprodev/wpd-message-gateway/internal/core/port"
)

func TestWorkspaceMemberRepository_Add(t *testing.T) {
	t.Parallel()

	client, mock, _ := newMockPgClient(t)
	repo := NewWorkspaceMemberRepository(client)

	mock.ExpectQuery(`SELECT id FROM roles WHERE name`).
		WithArgs("admin", domain.RBACGuardName).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("role-admin-id"))
	mock.ExpectExec(`INSERT INTO workspace_members`).
		WithArgs("ws-1", "user-1", "role-admin-id").
		WillReturnResult(sqlmock.NewResult(1, 1))

	if err := repo.Add(context.Background(), "ws-1", "user-1", domain.RoleAdmin); err != nil {
		t.Fatalf("Add: %v", err)
	}
}

func TestWorkspaceMemberRepository_GetRole_notFound(t *testing.T) {
	t.Parallel()

	client, mock, _ := newMockPgClient(t)
	repo := NewWorkspaceMemberRepository(client)

	mock.ExpectQuery(`FROM workspace_members wm`).
		WithArgs("ws-1", "user-missing").
		WillReturnError(sql.ErrNoRows)

	_, err := repo.GetRole(context.Background(), "ws-1", "user-missing")
	if err == nil || !errors.Is(err, port.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestWorkspaceMemberRepository_ListMembers(t *testing.T) {
	t.Parallel()

	client, mock, _ := newMockPgClient(t)
	repo := NewWorkspaceMemberRepository(client)
	joinedAt := time.Date(2026, 7, 4, 12, 0, 0, 0, time.UTC)

	rows := sqlmock.NewRows([]string{
		"workspace_id", "user_id", "name", "joined_at", "email", "display_name",
	}).AddRow("ws-1", "user-1", "admin", joinedAt, "owner@example.com", "Owner User")

	mock.ExpectQuery(`FROM workspace_members wm`).
		WithArgs("ws-1").
		WillReturnRows(rows)

	members, err := repo.ListMembers(context.Background(), "ws-1")
	if err != nil {
		t.Fatalf("ListMembers: %v", err)
	}
	if len(members) != 1 || members[0].Role != domain.RoleAdmin || members[0].DisplayName != "Owner User" {
		t.Fatalf("unexpected members: %+v", members)
	}
}
