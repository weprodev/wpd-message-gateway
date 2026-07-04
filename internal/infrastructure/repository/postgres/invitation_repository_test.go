package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/weprodev/wpd-message-gateway/internal/core/domain"
)

func TestInvitationRepository_Create(t *testing.T) {
	t.Parallel()

	client, mock, _ := newMockPgClient(t)
	repo := NewInvitationRepository(client)
	createdAt := time.Date(2026, 7, 4, 12, 0, 0, 0, time.UTC)
	expiresAt := createdAt.Add(7 * 24 * time.Hour)

	mock.ExpectQuery(`SELECT id FROM roles WHERE name`).
		WithArgs("member", domain.RBACGuardName).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("role-member-id"))
	mock.ExpectQuery(`INSERT INTO invitations`).
		WithArgs("ws-1", "invitee@example.com", "role-member-id", "hash", expiresAt, "pending").
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at"}).AddRow("inv-1", createdAt))

	inv := &domain.Invitation{
		WorkspaceID: "ws-1",
		Email:       "invitee@example.com",
		Role:        domain.RoleMember,
		TokenHash:   "hash",
		ExpiresAt:   expiresAt,
		Status:      "pending",
	}
	if err := repo.Create(context.Background(), inv); err != nil {
		t.Fatalf("Create: %v", err)
	}
	if inv.ID != "inv-1" || !inv.CreatedAt.Equal(createdAt) {
		t.Fatalf("unexpected invitation after create: %+v", inv)
	}
}

func TestInvitationRepository_ListPendingByWorkspace(t *testing.T) {
	t.Parallel()

	client, mock, _ := newMockPgClient(t)
	repo := NewInvitationRepository(client)
	createdAt := time.Date(2026, 7, 4, 12, 0, 0, 0, time.UTC)
	expiresAt := createdAt.Add(24 * time.Hour)

	rows := sqlmock.NewRows([]string{
		"id", "workspace_id", "email", "name", "token_hash", "expires_at", "status", "created_at",
	}).AddRow("inv-1", "ws-1", "a@b.com", "member", "hash", expiresAt, "pending", createdAt)

	mock.ExpectQuery(`FROM invitations i`).
		WithArgs("ws-1").
		WillReturnRows(rows)

	list, err := repo.ListPendingByWorkspace(context.Background(), "ws-1")
	if err != nil {
		t.Fatalf("ListPendingByWorkspace: %v", err)
	}
	if len(list) != 1 || list[0].Email != "a@b.com" || list[0].Role != domain.RoleMember {
		t.Fatalf("unexpected list: %+v", list)
	}
}

func TestInvitationRepository_PendingInvitationExistsByEmail(t *testing.T) {
	t.Parallel()

	client, mock, _ := newMockPgClient(t)
	repo := NewInvitationRepository(client)

	mock.ExpectQuery(`SELECT EXISTS`).
		WithArgs("ws-1", "invitee@example.com").
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

	exists, err := repo.PendingInvitationExistsByEmail(context.Background(), "ws-1", "invitee@example.com")
	if err != nil {
		t.Fatalf("PendingInvitationExistsByEmail: %v", err)
	}
	if !exists {
		t.Fatal("expected pending invitation to exist")
	}
}
