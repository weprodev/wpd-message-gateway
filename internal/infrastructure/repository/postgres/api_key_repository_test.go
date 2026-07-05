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

func TestAPIKeyRepository_Create(t *testing.T) {
	t.Parallel()

	client, mock, _ := newMockPgClient(t)
	repo := NewAPIKeyRepository(client)
	createdAt := time.Date(2026, 7, 4, 10, 0, 0, 0, time.UTC)

	mock.ExpectQuery(`INSERT INTO api_keys`).
		WithArgs("ws-1", "client-1", "hash", "Demo", true, nil).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at"}).AddRow("key-1", createdAt))

	key := &domain.APIKey{
		WorkspaceID:      "ws-1",
		ClientID:         "client-1",
		ClientSecretHash: "hash",
		Name:             "Demo",
		IsActive:         true,
	}
	if err := repo.Create(context.Background(), key); err != nil {
		t.Fatalf("Create: %v", err)
	}
	if key.ID != "key-1" || !key.CreatedAt.Equal(createdAt) {
		t.Fatalf("unexpected key after create: %+v", key)
	}
}

func TestAPIKeyRepository_GetByClientID_notFound(t *testing.T) {
	t.Parallel()

	client, mock, _ := newMockPgClient(t)
	repo := NewAPIKeyRepository(client)

	mock.ExpectQuery(`SELECT id, workspace_id`).
		WithArgs("missing").
		WillReturnError(sql.ErrNoRows)

	_, err := repo.GetByClientID(context.Background(), "missing")
	if err == nil || !errors.Is(err, port.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestAPIKeyRepository_Delete_notFound(t *testing.T) {
	t.Parallel()

	client, mock, _ := newMockPgClient(t)
	repo := NewAPIKeyRepository(client)

	mock.ExpectExec(`DELETE FROM api_keys WHERE id = \$1`).
		WithArgs("key-missing").
		WillReturnResult(sqlmock.NewResult(0, 0))

	err := repo.Delete(context.Background(), "key-missing")
	if err == nil || !errors.Is(err, port.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}
