package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/weprodev/go-pkg/crypto"

	"github.com/weprodev/wpd-message-gateway/internal/core/authjwt"
	"github.com/weprodev/wpd-message-gateway/internal/core/domain"
	"github.com/weprodev/wpd-message-gateway/internal/core/port"
	"github.com/weprodev/wpd-message-gateway/pkg/contracts"
)

// AuthService provides user authentication services.
type AuthService struct {
	userRepo                  port.UserRepository
	emailSender               contracts.EmailSender
	secret                    string
	emailVerificationEnabled  bool
	sessionTTL                time.Duration
	emailVerificationTokenTTL time.Duration
	// verificationBaseURL is the URL prefix for email verification links.
	verificationBaseURL string
}

// NewAuthService creates a new AuthService.
func NewAuthService(
	userRepo port.UserRepository,
	emailSender contracts.EmailSender,
	secret string,
	emailVerificationEnabled bool,
	sessionTTL time.Duration,
	emailVerificationTokenTTL time.Duration,
	verificationBaseURL string,
) *AuthService {
	if sessionTTL <= 0 {
		sessionTTL = 24 * time.Hour
	}
	if emailVerificationTokenTTL <= 0 {
		emailVerificationTokenTTL = 24 * time.Hour
	}
	return &AuthService{
		userRepo:                  userRepo,
		emailSender:               emailSender,
		secret:                    secret,
		emailVerificationEnabled:  emailVerificationEnabled,
		sessionTTL:                sessionTTL,
		emailVerificationTokenTTL: emailVerificationTokenTTL,
		verificationBaseURL:       verificationBaseURL,
	}
}

// RegisterUser registers a new user.
func (s *AuthService) RegisterUser(ctx context.Context, firstName, lastName, email, password string) (*domain.User, error) {
	slog.InfoContext(ctx, "registering new user", "email", email)
	_, err := s.userRepo.GetByEmail(ctx, email)
	if err == nil {
		slog.WarnContext(ctx, "user registration rejected: email already exists", "email", email)
		return nil, errors.New("user with this email already exists")
	}
	if !errors.Is(err, port.ErrNotFound) {
		slog.ErrorContext(ctx, "user registration failed: database error", "error", err, "email", email)
		return nil, err
	}

	passwordHash, err := crypto.HashSecret(password)
	if err != nil {
		slog.ErrorContext(ctx, "user registration failed: hashing error", "error", err, "email", email)
		return nil, fmt.Errorf("could not hash password: %w", err)
	}

	user := &domain.User{
		FirstName:     firstName,
		LastName:      lastName,
		Email:         email,
		PasswordHash:  passwordHash,
		EmailVerified: !s.emailVerificationEnabled,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		slog.ErrorContext(ctx, "user registration failed: failed to persist user", "error", err, "email", email)
		return nil, fmt.Errorf("could not create user: %w", err)
	}

	if s.emailVerificationEnabled {
		if err := s.sendVerificationEmail(ctx, user); err != nil {
			slog.ErrorContext(ctx, "user registration completed but verification email failed", "error", err, "email", email)
			return nil, fmt.Errorf("could not send verification email: %w", err)
		}
	}

	slog.InfoContext(ctx, "user registered successfully", "user_id", user.ID, "email", email)
	return user, nil
}

func (s *AuthService) sendVerificationEmail(ctx context.Context, user *domain.User) error {
	slog.InfoContext(ctx, "sending verification email", "user_id", user.ID, "email", user.Email)
	token, err := authjwt.Sign(user.ID.String(), user.Email, s.secret, s.emailVerificationTokenTTL)
	if err != nil {
		return err
	}

	verificationURL := fmt.Sprintf("%s/verify-email?token=%s", s.verificationBaseURL, token)

	_, err = s.emailSender.Send(ctx, contracts.Email{
		To:      []string{user.Email},
		Subject: "Verify your email address",
		HTML:    fmt.Sprintf("Please click the following link to verify your email address: %s", verificationURL),
	})
	if err != nil {
		slog.ErrorContext(ctx, "failed to send verification email", "error", err, "user_id", user.ID, "email", user.Email)
	} else {
		slog.InfoContext(ctx, "verification email sent successfully", "user_id", user.ID, "email", user.Email)
	}
	return err
}

// VerifyEmail verifies a user's email address.
func (s *AuthService) VerifyEmail(ctx context.Context, t string) error {
	slog.InfoContext(ctx, "verifying email token")
	claims, err := authjwt.Parse(t, s.secret)
	if err != nil {
		slog.WarnContext(ctx, "email verification failed: invalid or expired token")
		return errors.New("invalid or expired verification token")
	}

	user, err := s.userRepo.GetByID(ctx, claims.UserID)
	if err != nil {
		slog.ErrorContext(ctx, "email verification failed: user not found", "error", err, "user_id", claims.UserID)
		return errors.New("user not found")
	}

	if user.EmailVerified {
		slog.WarnContext(ctx, "email verification rejected: already verified", "user_id", claims.UserID)
		return errors.New("email already verified")
	}

	if err := s.userRepo.SetEmailVerified(ctx, claims.UserID); err != nil {
		slog.ErrorContext(ctx, "email verification failed: database error", "error", err, "user_id", claims.UserID)
		return err
	}
	slog.InfoContext(ctx, "email verified successfully", "user_id", claims.UserID)
	return nil
}

// Login logs in a user.
func (s *AuthService) Login(ctx context.Context, email, password string) (*domain.User, string, error) {
	slog.InfoContext(ctx, "user login attempt", "email", email)
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		slog.WarnContext(ctx, "login failed: invalid email", "email", email)
		return nil, "", errors.New("invalid email or password")
	}

	if !user.EmailVerified {
		slog.WarnContext(ctx, "login failed: email not verified", "user_id", user.ID, "email", email)
		return nil, "", errors.New("email not verified")
	}

	if !crypto.CheckSecretHash(password, user.PasswordHash) {
		slog.WarnContext(ctx, "login failed: invalid password", "user_id", user.ID, "email", email)
		return nil, "", errors.New("invalid email or password")
	}

	token, err := authjwt.Sign(user.ID.String(), user.Email, s.secret, s.sessionTTL)
	if err != nil {
		slog.ErrorContext(ctx, "login failed: token signing error", "error", err, "user_id", user.ID, "email", email)
		return nil, "", fmt.Errorf("could not sign token: %w", err)
	}

	slog.InfoContext(ctx, "user logged in successfully", "user_id", user.ID, "email", email)
	return user, token, nil
}
