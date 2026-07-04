package dto

import (
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/weprodev/wpd-message-gateway/internal/core/domain"
)

func TestUserFromDomain(t *testing.T) {
	t.Parallel()

	uid := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	now := time.Date(2026, 7, 4, 10, 0, 0, 0, time.UTC)
	got := UserFromDomain(&domain.User{
		ID:            uid,
		FirstName:     "Ada",
		LastName:      "Lovelace",
		Email:         "ada@example.com",
		PasswordHash:  "must-not-leak",
		EmailVerified: true,
		CreatedAt:     now,
		UpdatedAt:     now,
	})

	if got.ID != uid.String() || got.Email != "ada@example.com" {
		t.Fatalf("unexpected user mapping: %+v", got)
	}
}

func TestUserProfileResponseFromDomain(t *testing.T) {
	t.Parallel()

	uid := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	resp := UserProfileResponseFromDomain(&domain.User{
		ID:        uid,
		FirstName: "Ada",
		Email:     "ada@example.com",
	}, []domain.Workspace{{
		ID:   "ws-1",
		Name: "Demo",
		Slug: "demo",
		Role: "owner",
	}})

	if resp.ID != uid.String() || resp.FirstName != "Ada" {
		t.Fatalf("embedded user fields not promoted: %+v", resp)
	}
	if len(resp.Workspaces) != 1 || resp.Workspaces[0].Role != "owner" {
		t.Fatalf("unexpected workspaces: %+v", resp.Workspaces)
	}
}
