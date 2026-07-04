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

func TestTemplateRepository_Create(t *testing.T) {
	t.Parallel()

	client, mock, _ := newMockPgClient(t)
	repo := NewTemplateRepository(client)
	now := time.Date(2026, 7, 4, 12, 0, 0, 0, time.UTC)

	mock.ExpectQuery(`INSERT INTO templates`).
		WithArgs("ws-1", "Welcome", "welcome", "email", nil, nil, "<p>hi</p>", true, false).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).AddRow("tpl-1", now, now))

	tpl := &domain.Template{
		WorkspaceID: "ws-1",
		Name:        "Welcome",
		UniqueKey:   "welcome",
		ChannelType: "email",
		ContentHTML: "<p>hi</p>",
		IsActive:    true,
	}
	if err := repo.Create(context.Background(), tpl); err != nil {
		t.Fatalf("Create: %v", err)
	}
	if tpl.ID != "tpl-1" {
		t.Fatalf("unexpected template: %+v", tpl)
	}
}

func TestTemplateRepository_GetByWorkspaceAndKey_notFound(t *testing.T) {
	t.Parallel()

	client, mock, _ := newMockPgClient(t)
	repo := NewTemplateRepository(client)

	mock.ExpectQuery(`FROM templates`).
		WithArgs("ws-1", "missing").
		WillReturnError(sql.ErrNoRows)

	_, err := repo.GetByWorkspaceAndKey(context.Background(), "ws-1", "missing")
	if err == nil || !errors.Is(err, port.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestTemplateRepository_Delete_notFound(t *testing.T) {
	t.Parallel()

	client, mock, _ := newMockPgClient(t)
	repo := NewTemplateRepository(client)

	mock.ExpectExec(`DELETE FROM templates WHERE id`).
		WithArgs("tpl-missing").
		WillReturnResult(sqlmock.NewResult(0, 0))

	err := repo.Delete(context.Background(), "tpl-missing")
	if err == nil || !errors.Is(err, port.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}
