package service

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/weprodev/wpd-message-gateway/internal/core/domain"
	"github.com/weprodev/wpd-message-gateway/internal/core/port"
	"github.com/weprodev/wpd-message-gateway/pkg/contracts"
)

const testJWTSecret = "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"

type fakeUserRepo struct {
	byEmail *domain.User
	err     error
	created *domain.User
}

func (f *fakeUserRepo) Create(ctx context.Context, u *domain.User) error {
	f.created = u
	u.ID = uuid.New()
	return nil
}

func (f *fakeUserRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	return f.byEmail, f.err
}

func (f *fakeUserRepo) GetByID(ctx context.Context, id string) (*domain.User, error) {
	return nil, port.ErrNotFound
}

func (f *fakeUserRepo) SetEmailVerified(ctx context.Context, id string) error {
	return nil
}

type fakeEmailSender struct{}

func (f *fakeEmailSender) Send(ctx context.Context, email contracts.Email) (*contracts.SendResult, error) {
	return &contracts.SendResult{ID: "test-id"}, nil
}

func (f *fakeEmailSender) Name() string {
	return "fake"
}

func TestAuthService_RegisterUser_rejectsDuplicateEmail(t *testing.T) {
	repo := &fakeUserRepo{byEmail: &domain.User{Email: "a@b.com"}}
	sender := &fakeEmailSender{}
	svc := NewAuthService(repo, sender, testJWTSecret, false, time.Hour, time.Hour, "")

	_, err := svc.RegisterUser(context.Background(), "First", "Last", "a@b.com", "secret")
	if err == nil || err.Error() != "user with this email already exists" {
		t.Fatalf("expected duplicate error, got %v", err)
	}
}

func TestAuthService_RegisterUser_propagatesLookupErrors(t *testing.T) {
	repo := &fakeUserRepo{err: errors.New("database unavailable")}
	sender := &fakeEmailSender{}
	svc := NewAuthService(repo, sender, testJWTSecret, false, time.Hour, time.Hour, "")

	_, err := svc.RegisterUser(context.Background(), "First", "Last", "new@b.com", "secret")
	if err == nil || !errors.Is(err, repo.err) && err.Error() != "database unavailable" {
		t.Fatalf("expected lookup error, got %v", err)
	}
}

func TestAuthService_RegisterUser_allowsWhenUserMissing(t *testing.T) {
	repo := &fakeUserRepo{err: fmt.Errorf("lookup: %w", port.ErrNotFound)}
	sender := &fakeEmailSender{}
	svc := NewAuthService(repo, sender, testJWTSecret, false, time.Hour, time.Hour, "")

	user, err := svc.RegisterUser(context.Background(), "First", "Last", "fresh@b.com", "secret123456")
	if err != nil {
		t.Fatalf("RegisterUser: %v", err)
	}
	if user.Email != "fresh@b.com" {
		t.Fatalf("expected fresh@b.com, got %s", user.Email)
	}
	if repo.created == nil || repo.created.Email != "fresh@b.com" {
		t.Fatalf("user not created: %+v", repo.created)
	}
}

func TestAuthService_Login_invalidCredentialsWhenNotFound(t *testing.T) {
	repo := &fakeUserRepo{err: fmt.Errorf("wrapped: %w", port.ErrNotFound)}
	sender := &fakeEmailSender{}
	svc := NewAuthService(repo, sender, testJWTSecret, false, time.Hour, time.Hour, "")

	_, _, err := svc.Login(context.Background(), "a@b.com", "pw")
	if err == nil || err.Error() != "invalid email or password" {
		t.Fatalf("expected invalid email or password error, got %v", err)
	}
}

func TestAuthService_Login_masksInfrastructureErrors(t *testing.T) {
	repo := &fakeUserRepo{err: errors.New("connection reset")}
	sender := &fakeEmailSender{}
	svc := NewAuthService(repo, sender, testJWTSecret, false, time.Hour, time.Hour, "")

	_, _, err := svc.Login(context.Background(), "a@b.com", "pw")
	if err == nil || err.Error() != "invalid email or password" {
		t.Fatalf("expected invalid email or password error, got %v", err)
	}
}

type fakeWorkspaceRepo struct {
	byID       *domain.Workspace
	hashedPin  string
	updatedPin string
}

func (f *fakeWorkspaceRepo) Create(ctx context.Context, workspace *domain.Workspace) error {
	return nil
}

func (f *fakeWorkspaceRepo) GetByID(ctx context.Context, id string) (*domain.Workspace, error) {
	if f.byID == nil {
		return nil, port.ErrNotFound
	}
	w := *f.byID
	w.HashedPin = f.hashedPin
	return &w, nil
}

func (f *fakeWorkspaceRepo) GetBySlug(ctx context.Context, slug string) (*domain.Workspace, error) {
	if f.byID == nil || f.byID.Slug != slug {
		return nil, port.ErrNotFound
	}
	w := *f.byID
	w.HashedPin = f.hashedPin
	return &w, nil
}

func (f *fakeWorkspaceRepo) Update(ctx context.Context, workspace *domain.Workspace) error {
	return nil
}

func (f *fakeWorkspaceRepo) UpdateHashedPin(ctx context.Context, workspaceID, hashedPin string) error {
	f.updatedPin = hashedPin
	f.hashedPin = hashedPin
	return nil
}

func (f *fakeWorkspaceRepo) SetStatus(ctx context.Context, id, status string) error {
	return nil
}

func (f *fakeWorkspaceRepo) ListForUser(ctx context.Context, userID string) ([]domain.Workspace, error) {
	return nil, nil
}

type fakeSettingsRepo struct {
	values  map[string]string
	deleted []string
}

func (f *fakeSettingsRepo) Get(ctx context.Context, workspaceID, key string) (string, error) {
	if f.values == nil {
		return "", nil
	}
	return f.values[key], nil
}

func (f *fakeSettingsRepo) Set(ctx context.Context, workspaceID, key, value string) error {
	if f.values == nil {
		f.values = make(map[string]string)
	}
	f.values[key] = value
	return nil
}

func (f *fakeSettingsRepo) Delete(ctx context.Context, workspaceID, key string) error {
	delete(f.values, key)
	f.deleted = append(f.deleted, key)
	return nil
}

func (f *fakeSettingsRepo) GetAll(ctx context.Context, workspaceID string) (map[string]string, error) {
	out := make(map[string]string, len(f.values))
	for k, v := range f.values {
		out[k] = v
	}
	return out, nil
}

func TestPortalService_PatchSettings_hashesPinOnWorkspaceNotSettings(t *testing.T) {
	wsRepo := &fakeWorkspaceRepo{byID: &domain.Workspace{ID: "ws-1"}}
	settingsRepo := &fakeSettingsRepo{values: map[string]string{"owner_email": "a@b.com"}}
	svc := NewPortalService(PortalDeps{Workspaces: wsRepo, Settings: settingsRepo})

	err := svc.PatchSettings(context.Background(), "ws-1", map[string]string{
		"pin_code":    "567890",
		"owner_email": "owner@b.com",
	})
	if err != nil {
		t.Fatalf("PatchSettings: %v", err)
	}
	if wsRepo.updatedPin == "" {
		t.Fatal("expected workspace hashed PIN to be updated")
	}
	if settingsRepo.values["pin_code"] != "" {
		t.Fatalf("pin_code must not be stored in settings, got %q", settingsRepo.values["pin_code"])
	}
	if settingsRepo.values["owner_email"] != "owner@b.com" {
		t.Fatalf("expected owner_email updated, got %q", settingsRepo.values["owner_email"])
	}
}

func TestPortalService_GetSettings_reportsPinConfiguredWithoutPlaintext(t *testing.T) {
	wsRepo := &fakeWorkspaceRepo{
		byID:      &domain.Workspace{ID: "ws-1"},
		hashedPin: "$2a$14$abcdefghijklmnopqrstuv",
	}
	settingsRepo := &fakeSettingsRepo{values: map[string]string{
		"owner_email": "owner@b.com",
		"pin_code":    "567890",
	}}
	svc := NewPortalService(PortalDeps{Workspaces: wsRepo, Settings: settingsRepo})

	got, err := svc.GetSettings(context.Background(), "ws-1")
	if err != nil {
		t.Fatalf("GetSettings: %v", err)
	}
	if got["pin_code"] != "" {
		t.Fatalf("pin_code must not be returned, got %q", got["pin_code"])
	}
	if got["pin_configured"] != "true" {
		t.Fatalf("expected pin_configured=true, got %q", got["pin_configured"])
	}
}

func TestPortalService_JoinWorkspaceWithPIN_migratesLegacySettingsPin(t *testing.T) {
	wsRepo := &fakeWorkspaceRepo{
		byID: &domain.Workspace{ID: "ws-1", Slug: "demo"},
	}
	settingsRepo := &fakeSettingsRepo{values: map[string]string{"pin_code": "567890"}}
	members := &fakeMemberRepo{}
	gate := &fakeAuthGate{}
	svc := NewPortalService(PortalDeps{
		Workspaces: wsRepo,
		Settings:   settingsRepo,
		Members:    members,
		Gate:       gate,
	})

	err := svc.JoinWorkspaceWithPIN(context.Background(), "user-1", "demo", "567890")
	if err != nil {
		t.Fatalf("JoinWorkspaceWithPIN: %v", err)
	}
	if wsRepo.updatedPin == "" {
		t.Fatal("expected legacy pin migrated to workspace hash")
	}
	if settingsRepo.values["pin_code"] != "" {
		t.Fatal("legacy pin_code setting should be deleted")
	}
	if !members.added {
		t.Fatal("expected user added as member")
	}
}

type fakeMemberRepo struct {
	added bool
}

func (f *fakeMemberRepo) Add(ctx context.Context, workspaceID, userID, role string) error {
	f.added = true
	return nil
}

func (f *fakeMemberRepo) Remove(ctx context.Context, workspaceID, userID string) error {
	return nil
}

func (f *fakeMemberRepo) GetRole(ctx context.Context, workspaceID, userID string) (string, error) {
	return "", port.ErrNotFound
}

func (f *fakeMemberRepo) ListMembers(ctx context.Context, workspaceID string) ([]domain.WorkspaceMember, error) {
	return nil, nil
}

type fakeAuthGate struct{}

func (f *fakeAuthGate) AssignRole(ctx context.Context, modelType, modelID, teamID, roleName string) error {
	return nil
}

func (f *fakeAuthGate) RemoveRole(ctx context.Context, modelType, modelID, teamID, roleName string) error {
	return nil
}

func (f *fakeAuthGate) GetRoleNames(ctx context.Context, modelType, modelID, teamID string) ([]string, error) {
	return nil, nil
}

func (f *fakeAuthGate) GetAllPermissions(ctx context.Context, modelType, modelID, teamID string) ([]string, error) {
	return nil, nil
}
