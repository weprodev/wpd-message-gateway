package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/weprodev/go-pkg/crypto"

	"github.com/weprodev/wpd-message-gateway/internal/core/domain"
	"github.com/weprodev/wpd-message-gateway/internal/core/port"
)

type fakeInvitationRepo struct {
	created *domain.Invitation
	list    []domain.Invitation
}

func (f *fakeInvitationRepo) Create(_ context.Context, inv *domain.Invitation) error {
	f.created = inv
	inv.ID = "inv-" + uuid.NewString()
	return nil
}

func (f *fakeInvitationRepo) ListPendingByWorkspace(_ context.Context, workspaceID string) ([]domain.Invitation, error) {
	return f.list, nil
}

type fakeAPIKeyRepo struct {
	keys    map[string]*domain.APIKey
	byWS    map[string][]domain.APIKey
	created *domain.APIKey
	deleted string
}

func (f *fakeAPIKeyRepo) Create(_ context.Context, k *domain.APIKey) error {
	if f.keys == nil {
		f.keys = make(map[string]*domain.APIKey)
	}
	if f.byWS == nil {
		f.byWS = make(map[string][]domain.APIKey)
	}
	k.ID = "key-" + uuid.NewString()
	f.created = k
	f.keys[k.ID] = k
	f.byWS[k.WorkspaceID] = append(f.byWS[k.WorkspaceID], *k)
	return nil
}

func (f *fakeAPIKeyRepo) GetByClientID(_ context.Context, clientID string) (*domain.APIKey, error) {
	for _, k := range f.keys {
		if k.ClientID == clientID {
			return k, nil
		}
	}
	return nil, port.ErrNotFound
}

func (f *fakeAPIKeyRepo) GetByID(_ context.Context, id string) (*domain.APIKey, error) {
	k, ok := f.keys[id]
	if !ok {
		return nil, port.ErrNotFound
	}
	return k, nil
}

func (f *fakeAPIKeyRepo) ListByWorkspace(_ context.Context, workspaceID string) ([]domain.APIKey, error) {
	return f.byWS[workspaceID], nil
}

func (f *fakeAPIKeyRepo) Delete(_ context.Context, id string) error {
	f.deleted = id
	delete(f.keys, id)
	return nil
}

func (f *fakeAPIKeyRepo) UpdateSecret(_ context.Context, id, clientID, secretHash string) error {
	k, ok := f.keys[id]
	if !ok {
		return port.ErrNotFound
	}
	k.ClientID = clientID
	k.ClientSecretHash = secretHash
	return nil
}

func (f *fakeAPIKeyRepo) UpdateLastUsedAt(_ context.Context, id string) error {
	return nil
}

func TestPortalService_NewPendingInvitation(t *testing.T) {
	t.Parallel()

	svc := NewPortalService(PortalDeps{})
	before := time.Now()
	inv := svc.NewPendingInvitation("ws-1", "user@example.com", domain.RoleMember)
	after := time.Now()

	if inv.WorkspaceID != "ws-1" || inv.Email != "user@example.com" || inv.Role != domain.RoleMember {
		t.Fatalf("unexpected invitation fields: %+v", inv)
	}
	if inv.Status != "pending" {
		t.Fatalf("expected pending status, got %q", inv.Status)
	}
	wantMin := before.Add(7 * 24 * time.Hour)
	wantMax := after.Add(7 * 24 * time.Hour)
	if inv.ExpiresAt.Before(wantMin) || inv.ExpiresAt.After(wantMax) {
		t.Fatalf("expected ~7-day expiry, got %v (window %v–%v)", inv.ExpiresAt, wantMin, wantMax)
	}
}

func TestPortalService_CreateInvitation(t *testing.T) {
	t.Parallel()

	repo := &fakeInvitationRepo{}
	svc := NewPortalService(PortalDeps{Invitations: repo})
	inv := svc.NewPendingInvitation("ws-1", "user@example.com", domain.RoleMember)

	rawToken, err := svc.CreateInvitation(context.Background(), inv)
	if err != nil {
		t.Fatalf("CreateInvitation: %v", err)
	}
	if rawToken == "" {
		t.Fatal("expected non-empty token")
	}
	if repo.created == nil || repo.created.TokenHash == "" {
		t.Fatal("expected hashed token persisted")
	}
	if rawToken == repo.created.TokenHash {
		t.Fatal("stored value must be hash, not plaintext token")
	}
	sum := sha256.Sum256([]byte(rawToken))
	if repo.created.TokenHash != hex.EncodeToString(sum[:]) {
		t.Fatalf("token hash mismatch: got %q", repo.created.TokenHash)
	}
}

func TestPortalService_APIKeys(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	repo := &fakeAPIKeyRepo{}
	svc := NewPortalService(PortalDeps{APIKeys: repo})

	key, secret, err := svc.CreateAPIKey(ctx, "ws-1", "  My Key  ")
	if err != nil {
		t.Fatalf("CreateAPIKey: %v", err)
	}
	if secret == "" {
		t.Fatal("expected plaintext secret")
	}
	if key.Name != "My Key" {
		t.Fatalf("expected trimmed name, got %q", key.Name)
	}
	if !crypto.CheckSecretHash(secret, repo.created.ClientSecretHash) {
		t.Fatal("stored hash does not match returned secret")
	}

	list, err := svc.ListAPIKeys(ctx, "ws-1")
	if err != nil {
		t.Fatalf("ListAPIKeys: %v", err)
	}
	if len(list) != 1 || list[0].ID != key.ID {
		t.Fatalf("expected one key in list, got %+v", list)
	}

	if err := svc.DeleteAPIKey(ctx, "ws-1", key.ID); err != nil {
		t.Fatalf("DeleteAPIKey: %v", err)
	}
	if repo.deleted != key.ID {
		t.Fatalf("expected delete id %q, got %q", key.ID, repo.deleted)
	}
}

func TestPortalService_DeleteAPIKey_rejectsWrongWorkspace(t *testing.T) {
	t.Parallel()

	repo := &fakeAPIKeyRepo{
		keys: map[string]*domain.APIKey{
			"key-1": {ID: "key-1", WorkspaceID: "ws-other"},
		},
	}
	svc := NewPortalService(PortalDeps{APIKeys: repo})

	err := svc.DeleteAPIKey(context.Background(), "ws-1", "key-1")
	if err == nil || err.Error() != "api key not in workspace" {
		t.Fatalf("expected workspace mismatch error, got %v", err)
	}
	if repo.deleted != "" {
		t.Fatal("delete should not run for wrong workspace")
	}
}
