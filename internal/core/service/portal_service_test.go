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

func (f *fakeEmailSender) Send(ctx context.Context, email *contracts.Email) (*contracts.SendResult, error) {
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
