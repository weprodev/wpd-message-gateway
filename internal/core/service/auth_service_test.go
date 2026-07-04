package service

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/weprodev/go-pkg/crypto"

	"github.com/weprodev/wpd-message-gateway/internal/core/authjwt"
	"github.com/weprodev/wpd-message-gateway/internal/core/domain"
	"github.com/weprodev/wpd-message-gateway/internal/core/port"
	"github.com/weprodev/wpd-message-gateway/pkg/contracts"
)

type failingEmailSender struct{}

func (f *failingEmailSender) Send(context.Context, contracts.Email) (*contracts.SendResult, error) {
	return nil, errors.New("smtp down")
}

func (f *failingEmailSender) Name() string { return "failing" }

func TestAuthService_VerifyEmail_success(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	repo := &fakeUserRepo{
		byID: &domain.User{
			ID:            userID,
			Email:         "user@example.com",
			EmailVerified: false,
		},
	}
	svc := NewAuthService(repo, &fakeEmailSender{}, testJWTSecret, true, time.Hour, time.Hour, "https://portal.example.com")

	token, err := authjwt.Sign(userID.String(), "user@example.com", testJWTSecret, time.Hour)
	if err != nil {
		t.Fatalf("Sign: %v", err)
	}

	if err := svc.VerifyEmail(context.Background(), token); err != nil {
		t.Fatalf("VerifyEmail: %v", err)
	}
	if len(repo.verifiedUserIDs) != 1 || repo.verifiedUserIDs[0] != userID.String() {
		t.Fatalf("expected SetEmailVerified for %s, got %v", userID, repo.verifiedUserIDs)
	}
}

func TestAuthService_VerifyEmail_rejectsAlreadyVerified(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	repo := &fakeUserRepo{
		byID: &domain.User{ID: userID, Email: "user@example.com", EmailVerified: true},
	}
	svc := NewAuthService(repo, &fakeEmailSender{}, testJWTSecret, true, time.Hour, time.Hour, "")

	token, err := authjwt.Sign(userID.String(), "user@example.com", testJWTSecret, time.Hour)
	if err != nil {
		t.Fatalf("Sign: %v", err)
	}

	err = svc.VerifyEmail(context.Background(), token)
	if err == nil || err.Error() != "email already verified" {
		t.Fatalf("expected already verified error, got %v", err)
	}
}

func TestAuthService_VerifyEmail_rejectsInvalidToken(t *testing.T) {
	t.Parallel()

	svc := NewAuthService(&fakeUserRepo{}, &fakeEmailSender{}, testJWTSecret, true, time.Hour, time.Hour, "")

	err := svc.VerifyEmail(context.Background(), "not-a-jwt")
	if err == nil || err.Error() != "invalid or expired verification token" {
		t.Fatalf("expected invalid token error, got %v", err)
	}
}

func TestAuthService_Login_success(t *testing.T) {
	t.Parallel()

	hash, err := crypto.HashSecret("correct-password")
	if err != nil {
		t.Fatalf("HashSecret: %v", err)
	}
	userID := uuid.New()
	repo := &fakeUserRepo{
		byEmail: &domain.User{
			ID:            userID,
			Email:         "user@example.com",
			PasswordHash:  hash,
			EmailVerified: true,
		},
	}
	svc := NewAuthService(repo, &fakeEmailSender{}, testJWTSecret, false, time.Hour, time.Hour, "")

	user, token, err := svc.Login(context.Background(), "user@example.com", "correct-password")
	if err != nil {
		t.Fatalf("Login: %v", err)
	}
	if user.ID != userID || token == "" {
		t.Fatalf("unexpected login result: user=%+v token=%q", user, token)
	}

	claims, err := authjwt.Parse(token, testJWTSecret)
	if err != nil {
		t.Fatalf("Parse token: %v", err)
	}
	if claims.UserID != userID.String() {
		t.Fatalf("token subject: got %q want %q", claims.UserID, userID)
	}
}

func TestAuthService_Login_rejectsUnverifiedEmail(t *testing.T) {
	t.Parallel()

	repo := &fakeUserRepo{
		byEmail: &domain.User{
			ID:            uuid.New(),
			Email:         "user@example.com",
			PasswordHash:  "hash",
			EmailVerified: false,
		},
	}
	svc := NewAuthService(repo, &fakeEmailSender{}, testJWTSecret, true, time.Hour, time.Hour, "")

	_, _, err := svc.Login(context.Background(), "user@example.com", "pw")
	if err == nil || err.Error() != "email not verified" {
		t.Fatalf("expected email not verified error, got %v", err)
	}
}

func TestAuthService_RegisterUser_sendsVerificationEmail(t *testing.T) {
	t.Parallel()

	repo := &fakeUserRepo{err: fmt.Errorf("lookup: %w", port.ErrNotFound)}
	sender := &fakeEmailSender{}
	svc := NewAuthService(repo, sender, testJWTSecret, true, time.Hour, time.Hour, "https://portal.example.com")

	user, err := svc.RegisterUser(context.Background(), "First", "Last", "new@example.com", "secret123456")
	if err != nil {
		t.Fatalf("RegisterUser: %v", err)
	}
	if user.EmailVerified {
		t.Fatal("expected EmailVerified false when verification enabled")
	}
}

func TestAuthService_RegisterUser_failsWhenVerificationEmailFails(t *testing.T) {
	t.Parallel()

	repo := &fakeUserRepo{err: fmt.Errorf("lookup: %w", port.ErrNotFound)}
	svc := NewAuthService(repo, &failingEmailSender{}, testJWTSecret, true, time.Hour, time.Hour, "https://portal.example.com")

	_, err := svc.RegisterUser(context.Background(), "First", "Last", "new@example.com", "secret123456")
	if err == nil || err.Error() != "could not send verification email: smtp down" {
		t.Fatalf("expected verification email failure, got %v", err)
	}
}
