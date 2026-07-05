package middleware

import (
	"context"
	"errors"
	"testing"

	"github.com/weprodev/wpd-message-gateway/internal/core/domain"
	"github.com/weprodev/wpd-message-gateway/internal/core/port"
)

type stubWorkspaceRepo struct {
	byID   map[string]*domain.Workspace
	bySlug map[string]*domain.Workspace
}

func (s *stubWorkspaceRepo) Create(context.Context, *domain.Workspace, string) error { return nil }
func (s *stubWorkspaceRepo) Delete(context.Context, string) error                    { return nil }
func (s *stubWorkspaceRepo) Update(context.Context, *domain.Workspace) error         { return nil }
func (s *stubWorkspaceRepo) SetStatus(context.Context, string, string) error         { return nil }
func (s *stubWorkspaceRepo) ListForUser(context.Context, string) ([]domain.Workspace, error) {
	return nil, nil
}

func (s *stubWorkspaceRepo) GetByID(_ context.Context, id string) (*domain.Workspace, error) {
	ws, ok := s.byID[id]
	if !ok {
		return nil, port.ErrNotFound
	}
	return ws, nil
}

func (s *stubWorkspaceRepo) GetBySlug(_ context.Context, slug string) (*domain.Workspace, error) {
	ws, ok := s.bySlug[slug]
	if !ok {
		return nil, port.ErrNotFound
	}
	return ws, nil
}

func TestResolveWorkspaceByKey(t *testing.T) {
	repo := &stubWorkspaceRepo{
		byID: map[string]*domain.Workspace{
			"00000000-0000-0000-0000-000000000001": {ID: "00000000-0000-0000-0000-000000000001", Slug: "demo"},
		},
		bySlug: map[string]*domain.Workspace{
			"demo": {ID: "00000000-0000-0000-0000-000000000001", Slug: "demo"},
		},
	}

	t.Run("resolves workspace id", func(t *testing.T) {
		ws, err := resolveWorkspaceByKey(context.Background(), repo, "00000000-0000-0000-0000-000000000001")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if ws.ID != "00000000-0000-0000-0000-000000000001" {
			t.Fatalf("got id %q", ws.ID)
		}
	})

	t.Run("resolves workspace slug", func(t *testing.T) {
		ws, err := resolveWorkspaceByKey(context.Background(), repo, "demo")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if ws.Slug != "demo" {
			t.Fatalf("got slug %q", ws.Slug)
		}
	})

	t.Run("unknown key returns not found", func(t *testing.T) {
		_, err := resolveWorkspaceByKey(context.Background(), repo, "missing")
		if !errors.Is(err, port.ErrNotFound) {
			t.Fatalf("expected ErrNotFound, got %v", err)
		}
	})
}
