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

type fakeSettingsRepo struct {
	values map[string]string
}

func (f *fakeSettingsRepo) Get(ctx context.Context, workspaceID, key string) (string, error) {
	return f.values[key], nil
}

func (f *fakeSettingsRepo) Set(ctx context.Context, workspaceID, key, value string) error {
	if f.values == nil {
		f.values = make(map[string]string)
	}
	f.values[key] = value
	return nil
}

func (f *fakeSettingsRepo) GetAll(ctx context.Context, workspaceID string) (map[string]string, error) {
	return f.values, nil
}

type fakeIntegrationRepo struct {
	integrations []domain.Integration
}

func (f *fakeIntegrationRepo) Create(ctx context.Context, intg *domain.Integration) error { return nil }
func (f *fakeIntegrationRepo) ListByWorkspace(ctx context.Context, workspaceID string) ([]domain.Integration, error) {
	return f.integrations, nil
}
func (f *fakeIntegrationRepo) Upsert(ctx context.Context, intg *domain.Integration) error { return nil }
func (f *fakeIntegrationRepo) GetByID(ctx context.Context, id string) (*domain.Integration, error) {
	return nil, port.ErrNotFound
}
func (f *fakeIntegrationRepo) Delete(ctx context.Context, id string) error { return nil }
func (f *fakeIntegrationRepo) ListProviders(ctx context.Context) ([]domain.Provider, error) {
	return nil, nil
}
func (f *fakeIntegrationRepo) GetProviderFields(ctx context.Context, name string) ([]domain.ProviderConfigField, error) {
	return nil, nil
}
func (f *fakeIntegrationRepo) GetActiveByWorkspaceAndChannel(ctx context.Context, workspaceID, channel string) (*domain.Integration, error) {
	return nil, nil
}

func TestPortalService_PatchSettings_rejectsInvalidDispatchMode(t *testing.T) {
	svc := NewPortalService(PortalDeps{Settings: &fakeSettingsRepo{}})

	err := svc.PatchSettings(context.Background(), "ws-1", map[string]string{
		domain.SettingKeyMessageDispatchMode: "invalid",
	})
	if !errors.Is(err, port.ErrInvalidInput) {
		t.Fatalf("expected ErrInvalidInput, got %v", err)
	}
}

func TestPortalService_PatchSettings_rejectsProviderWithoutIntegration(t *testing.T) {
	repo := &fakeSettingsRepo{}
	intgRepo := &fakeIntegrationRepo{
		integrations: []domain.Integration{},
	}
	svc := NewPortalService(PortalDeps{
		Settings:     repo,
		Integrations: intgRepo,
	})

	err := svc.PatchSettings(context.Background(), "ws-1", map[string]string{
		domain.SettingKeyMessageDispatchMode: string(domain.DispatchProvider),
	})
	if !errors.Is(err, port.ErrInvalidInput) {
		t.Fatalf("expected ErrInvalidInput, got %v", err)
	}
}

func TestPortalService_PatchSettings_persistsValidDispatchSettings(t *testing.T) {
	repo := &fakeSettingsRepo{}
	intgRepo := &fakeIntegrationRepo{
		integrations: []domain.Integration{
			{ProviderName: "mailgun", Status: domain.IntegrationStatusConnected},
		},
	}
	svc := NewPortalService(PortalDeps{
		Settings:     repo,
		Integrations: intgRepo,
	})

	err := svc.PatchSettings(context.Background(), "ws-1", map[string]string{
		domain.SettingKeyMessageDispatchMode: string(domain.DispatchProvider),
		domain.SettingKeyStoreMessageContent: "true",
	})
	if err != nil {
		t.Fatalf("PatchSettings: %v", err)
	}
	if repo.values[domain.SettingKeyMessageDispatchMode] != string(domain.DispatchProvider) {
		t.Fatalf("mode not saved: %+v", repo.values)
	}
}
