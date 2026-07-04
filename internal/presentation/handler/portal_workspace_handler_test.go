package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/weprodev/wpd-message-gateway/internal/core/domain"
	"github.com/weprodev/wpd-message-gateway/internal/core/service"
)

type fakeInvitationRepository struct {
	created *domain.Invitation
}

func (f *fakeInvitationRepository) Create(_ context.Context, inv *domain.Invitation) error {
	f.created = inv
	inv.ID = "inv-" + uuid.NewString()
	return nil
}

func (f *fakeInvitationRepository) ListPendingByWorkspace(context.Context, string) ([]domain.Invitation, error) {
	return nil, nil
}

func TestPortalWorkspaceHandler_CreateInvitation(t *testing.T) {
	t.Parallel()
	e := echo.New()

	repo := &fakeInvitationRepository{}
	svc := service.NewPortalService(service.PortalDeps{Invitations: repo})
	h := NewPortalWorkspaceHandler(svc)

	body := `{"email":"invitee@example.com","role":"member"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/workspaces/ws-1/invitations", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("wid")
	c.SetParamValues("ws-1")

	if err := h.CreateInvitation(c); err != nil {
		t.Fatalf("CreateInvitation: %v", err)
	}
	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d body=%s", rec.Code, rec.Body.String())
	}

	var res map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &res); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if res["email"] != "invitee@example.com" || res["role"] != "member" {
		t.Fatalf("unexpected response: %+v", res)
	}
	token, _ := res["token"].(string)
	if token == "" {
		t.Fatal("expected token in response")
	}
	if repo.created == nil || repo.created.TokenHash == "" {
		t.Fatal("expected invitation persisted with token hash")
	}
	if repo.created.Email != "invitee@example.com" || repo.created.Status != "pending" {
		t.Fatalf("unexpected persisted invitation: %+v", repo.created)
	}
}

func TestPortalWorkspaceHandler_CreateInvitation_rejectsInvalidJSON(t *testing.T) {
	t.Parallel()
	e := echo.New()

	h := NewPortalWorkspaceHandler(service.NewPortalService(service.PortalDeps{}))

	req := httptest.NewRequest(http.MethodPost, "/api/v1/workspaces/ws-1/invitations", strings.NewReader("{bad"))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("wid")
	c.SetParamValues("ws-1")

	err := h.CreateInvitation(c)
	if err == nil {
		t.Fatal("expected error")
	}
	var httpErr *echo.HTTPError
	if !errors.As(err, &httpErr) || httpErr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %v", err)
	}
}
