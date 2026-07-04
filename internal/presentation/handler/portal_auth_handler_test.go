package handler

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/weprodev/wpd-message-gateway/internal/core/domain"
	"github.com/weprodev/wpd-message-gateway/internal/core/port"
	"github.com/weprodev/wpd-message-gateway/internal/core/service"
	"github.com/weprodev/wpd-message-gateway/pkg/contracts"
)

const handlerTestJWTSecret = "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"

type authFakeUserRepo struct {
	byEmail *domain.User
	err     error
}

func (f *authFakeUserRepo) Create(context.Context, *domain.User) error { return nil }

func (f *authFakeUserRepo) GetByEmail(_ context.Context, email string) (*domain.User, error) {
	return f.byEmail, f.err
}

func (f *authFakeUserRepo) GetByID(context.Context, string) (*domain.User, error) {
	return nil, port.ErrNotFound
}

func (f *authFakeUserRepo) SetEmailVerified(context.Context, string) error { return nil }

type authFakeEmailSender struct{}

func (authFakeEmailSender) Send(context.Context, contracts.Email) (*contracts.SendResult, error) {
	return &contracts.SendResult{ID: "test-id"}, nil
}

func (authFakeEmailSender) Name() string { return "fake" }

func newPortalAuthHandlerForTest(userRepo *authFakeUserRepo) *PortalAuthHandler {
	authSvc := service.NewAuthService(userRepo, authFakeEmailSender{}, handlerTestJWTSecret, false, time.Hour, time.Hour, "")
	portalSvc := service.NewPortalService(service.PortalDeps{})
	return NewPortalAuthHandler(portalSvc, authSvc)
}

func TestPortalAuthHandler_Register_rejectsInvalidJSON(t *testing.T) {
	t.Parallel()
	e := echo.New()
	h := newPortalAuthHandlerForTest(&authFakeUserRepo{})

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", strings.NewReader("{not-json"))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := h.Register(c)
	if err == nil {
		t.Fatal("expected error")
	}
	var httpErr *echo.HTTPError
	if !errors.As(err, &httpErr) || httpErr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 invalid JSON, got %v", err)
	}
}

func TestPortalAuthHandler_Login_returns401OnBadCredentials(t *testing.T) {
	t.Parallel()
	e := echo.New()
	repo := &authFakeUserRepo{err: fmt.Errorf("lookup: %w", port.ErrNotFound)}
	h := newPortalAuthHandlerForTest(repo)

	body := `{"email":"nobody@example.com","password":"wrong-password"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := h.Login(c)
	if err == nil {
		t.Fatal("expected error")
	}
	var httpErr *echo.HTTPError
	if !errors.As(err, &httpErr) || httpErr.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %v", err)
	}
}
