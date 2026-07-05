package postgres

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/weprodev/wpd-message-gateway/internal/core/port"
)

func TestWorkspaceRepository_GetByID(t *testing.T) {
	t.Parallel()

	client, mock, _ := newMockPgClient(t)
	repo := NewWorkspaceRepository(client)
	now := time.Date(2026, 7, 4, 12, 0, 0, 0, time.UTC)

	mock.ExpectQuery(`FROM workspaces w`).
		WithArgs("ws-1").
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "name", "slug", "email", "status", "visibility", "hashed_pin_code", "icon_key", "created_at", "updated_at",
		}).AddRow("ws-1", "Demo", "demo", "owner@example.com", "active", "public", nil, nil, now, now))

	ws, err := repo.GetByID(context.Background(), "ws-1")
	if err != nil {
		t.Fatalf("GetByID: %v", err)
	}
	if ws.Name != "Demo" || ws.Slug != "demo" || ws.AdminEmail != "owner@example.com" {
		t.Fatalf("unexpected workspace: %+v", ws)
	}
}

func TestWorkspaceRepository_GetByID_notFound(t *testing.T) {
	t.Parallel()

	client, mock, _ := newMockPgClient(t)
	repo := NewWorkspaceRepository(client)

	mock.ExpectQuery(`FROM workspaces w`).
		WithArgs("missing").
		WillReturnError(sql.ErrNoRows)

	_, err := repo.GetByID(context.Background(), "missing")
	if err == nil || !errors.Is(err, port.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestWorkspaceRepository_GetBySlug_notFound(t *testing.T) {
	t.Parallel()

	client, mock, _ := newMockPgClient(t)
	repo := NewWorkspaceRepository(client)

	mock.ExpectQuery(`FROM workspaces w`).
		WithArgs("missing-slug").
		WillReturnError(sql.ErrNoRows)

	_, err := repo.GetBySlug(context.Background(), "missing-slug")
	if err == nil || !errors.Is(err, port.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}
