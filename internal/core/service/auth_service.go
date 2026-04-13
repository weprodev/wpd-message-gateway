package service

import (
	"context"
	"errors"
	"fmt"
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
	emailSender               port.EmailSender
	secret                    string
	emailVerificationEnabled  bool
	emailVerificationTokenTTL time.Duration
}

// NewAuthService creates a new AuthService.
func NewAuthService(
	userRepo port.UserRepository,
	emailSender port.EmailSender,
	secret string,
	emailVerificationEnabled bool,
	emailVerificationTokenTTL time.Duration,
) *AuthService {
	return &AuthService{
		userRepo:                  userRepo,
		emailSender:               emailSender,
		secret:                    secret,
		emailVerificationEnabled:  emailVerificationEnabled,
		emailVerificationTokenTTL: emailVerificationTokenTTL,
	}
}

// RegisterUser registers a new user.
func (s *AuthService) RegisterUser(ctx context.Context, firstName, lastName, email, password string) (*domain.User, error) {
	if _, err := s.userRepo.GetByEmail(ctx, email); !errors.Is(err, port.ErrNotFound) {
		return nil, errors.New("user with this email already exists")
	}

	passwordHash, err := crypto.HashSecret(password)
	if err != nil {
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
		return nil, fmt.Errorf("could not create user: %w", err)
	}

	if s.emailVerificationEnabled {
		if err := s.sendVerificationEmail(ctx, user); err != nil {
			return nil, fmt.Errorf("could not send verification email: %w", err)
		}
	}

	return user, nil
}

func (s *AuthService) sendVerificationEmail(ctx context.Context, user *domain.User) error {
	token, err := authjwt.Sign(user.ID.String(), user.Email, s.secret, s.emailVerificationTokenTTL)
	if err != nil {
		return err
	}

	// TODO: Make the verification URL configurable
	verificationURL := fmt.Sprintf("https://portal.weprodev.com/verify-email?token=%s", token)

	_, err = s.emailSender.Send(ctx, &contracts.Email{
		To:      []string{user.Email},
		Subject: "Verify your email address",
		HTML:    fmt.Sprintf("Please click the following link to verify your email address: %s", verificationURL),
	})
	return err
}

// VerifyEmail verifies a user's email address.
func (s *AuthService) VerifyEmail(ctx context.Context, t string) error {
	claims, err := authjwt.Parse(t, s.secret)
	if err != nil {
		return errors.New("invalid or expired verification token")
	}

	user, err := s.userRepo.GetByID(ctx, claims.UserID)
	if err != nil {
		return errors.New("user not found")
	}

	if user.EmailVerified {
		return errors.New("email already verified")
	}

	return s.userRepo.SetEmailVerified(ctx, claims.UserID)
}

// Login logs in a user.
func (s *AuthService) Login(ctx context.Context, email, password string) (*domain.User, string, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, "", errors.New("invalid email or password")
	}

	if !user.EmailVerified {
		return nil, "", errors.New("email not verified")
	}

	if !crypto.CheckSecretHash(password, user.PasswordHash) {
		return nil, "", errors.New("invalid email or password")
	}

	token, err := authjwt.Sign(user.ID.String(), user.Email, s.secret, s.emailVerificationTokenTTL)
	if err != nil {
		return nil, "", fmt.Errorf("could not sign token: %w", err)
	}

	return user, token, nil
}
