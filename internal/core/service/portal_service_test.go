package service

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/weprodev/wpd-message-gateway/internal/core/domain"
	"github.com/weprodev/wpd-message-gateway/internal/core/port"
)

const testJWTSecret = "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"

type fakeUserRepo struct {
	byEmail *domain.User
	err     error
	created *domain.User
}

func (f *fakeUserRepo) Create(ctx context.Context, u *domain.User) error {
	f.created = u
	u.ID = "new-user-id"
	return nil
}

func (f *fakeUserRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	return f.byEmail, f.err
}

func (f *fakeUserRepo) GetByID(ctx context.Context, id string) (*domain.User, error) {
	return nil, port.ErrNotFound
}

func TestPortalService_Register_rejectsDuplicateEmail(t *testing.T) {
	repo := &fakeUserRepo{byEmail: &domain.User{Email: "a@b.com"}}
	svc := NewPortalService(PortalDeps{Users: repo, JWTSecret: testJWTSecret, JWTTTL: time.Hour})

	_, _, err := svc.Register(context.Background(), "a@b.com", "secret", "n")
	if err == nil || err.Error() != "email already registered" {
		t.Fatalf("expected duplicate error, got %v", err)
	}
}

func TestPortalService_Register_propagatesLookupErrors(t *testing.T) {
	repo := &fakeUserRepo{err: errors.New("database unavailable")}
	svc := NewPortalService(PortalDeps{Users: repo, JWTSecret: testJWTSecret})

	_, _, err := svc.Register(context.Background(), "new@b.com", "secret", "n")
	if !errors.Is(err, repo.err) {
		t.Fatalf("expected lookup error, got %v", err)
	}
}

func TestPortalService_Register_allowsWhenUserMissing(t *testing.T) {
	repo := &fakeUserRepo{err: fmt.Errorf("lookup: %w", port.ErrNotFound)}
	svc := NewPortalService(PortalDeps{Users: repo, JWTSecret: testJWTSecret, JWTTTL: time.Hour})

	_, token, err := svc.Register(context.Background(), "fresh@b.com", "secret123456", "Name")
	if err != nil {
		t.Fatalf("Register: %v", err)
	}
	if token == "" {
		t.Fatal("expected JWT token")
	}
	if repo.created == nil || repo.created.Email != "fresh@b.com" {
		t.Fatalf("user not created: %+v", repo.created)
	}
}

func TestPortalService_Login_invalidCredentialsWhenNotFound(t *testing.T) {
	repo := &fakeUserRepo{err: fmt.Errorf("wrapped: %w", port.ErrNotFound)}
	svc := NewPortalService(PortalDeps{Users: repo, JWTSecret: testJWTSecret})

	_, _, err := svc.Login(context.Background(), "a@b.com", "pw")
	if !errors.Is(err, port.ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}
}

func TestPortalService_Login_propagatesInfrastructureErrors(t *testing.T) {
	repo := &fakeUserRepo{err: errors.New("connection reset")}
	svc := NewPortalService(PortalDeps{Users: repo, JWTSecret: testJWTSecret})

	_, _, err := svc.Login(context.Background(), "a@b.com", "pw")
	if !errors.Is(err, repo.err) {
		t.Fatalf("expected infrastructure error, got %v", err)
	}
}
