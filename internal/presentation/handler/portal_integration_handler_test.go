package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"

	"github.com/weprodev/wpd-message-gateway/internal/core/domain"
	"github.com/weprodev/wpd-message-gateway/internal/core/service"
)

type fakeIntegrationRepository struct {
	list       []domain.Integration
	upserted   *domain.Integration
	deletedID  string
	getByID    *domain.Integration
	listErr    error
	upsertErr  error
	deleteErr  error
	getByIDErr error
	fields     []domain.ProviderConfigField
	fieldsErr  error
}

func (f *fakeIntegrationRepository) Create(ctx context.Context, integration *domain.Integration) error {
	return nil
}

func (f *fakeIntegrationRepository) GetActiveByWorkspaceAndChannel(ctx context.Context, workspaceID, channelType string) (*domain.Integration, error) {
	return nil, nil
}

func (f *fakeIntegrationRepository) ListByWorkspace(ctx context.Context, workspaceID string) ([]domain.Integration, error) {
	return f.list, f.listErr
}

func (f *fakeIntegrationRepository) GetByID(ctx context.Context, id string) (*domain.Integration, error) {
	return f.getByID, f.getByIDErr
}

func (f *fakeIntegrationRepository) Delete(ctx context.Context, id string) error {
	f.deletedID = id
	return f.deleteErr
}

func (f *fakeIntegrationRepository) Upsert(ctx context.Context, integration *domain.Integration) error {
	f.upserted = integration
	return f.upsertErr
}

func (f *fakeIntegrationRepository) GetProviderFields(ctx context.Context, providerName string) ([]domain.ProviderConfigField, error) {
	return f.fields, f.fieldsErr
}

func (f *fakeIntegrationRepository) ListProviders(ctx context.Context) ([]domain.Provider, error) {
	return nil, nil
}

func TestPortalIntegrationHandler_ListIntegrations(t *testing.T) {
	e := echo.New()

	t.Run("lists active integrations and filters out memory provider", func(t *testing.T) {
		repo := &fakeIntegrationRepository{
			list: []domain.Integration{
				{
					ID:           "intg-1",
					WorkspaceID:  "ws-123",
					ChannelType:  "email",
					ProviderName: "mailgun",
					Config:       []byte(`{"api_key":"123"}`),
					Status:       "connected",
				},
				{
					ID:           "intg-2",
					WorkspaceID:  "ws-123",
					ChannelType:  "email",
					ProviderName: "memory",
					Config:       []byte(`{}`),
					Status:       "connected",
				},
			},
		}

		svc := service.NewPortalService(service.PortalDeps{
			Integrations: repo,
		})
		h := NewPortalIntegrationHandler(svc)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/workspaces/ws-123/integrations", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("wid")
		c.SetParamValues("ws-123")

		if err := h.ListIntegrations(c); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if rec.Code != http.StatusOK {
			t.Errorf("expected 200 OK, got %d", rec.Code)
		}

		var res []map[string]any
		if err := json.Unmarshal(rec.Body.Bytes(), &res); err != nil {
			t.Fatalf("failed to decode JSON response: %v", err)
		}

		// The memory provider integration must be filtered out
		if len(res) != 1 {
			t.Fatalf("expected exactly 1 integration, got %d", len(res))
		}

		if res[0]["provider_name"] != "mailgun" {
			t.Errorf("expected provider mailgun, got %v", res[0]["provider_name"])
		}
	})

	t.Run("handles repository list error", func(t *testing.T) {
		repo := &fakeIntegrationRepository{
			listErr: errors.New("db error"),
		}

		svc := service.NewPortalService(service.PortalDeps{
			Integrations: repo,
		})
		h := NewPortalIntegrationHandler(svc)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/workspaces/ws-123/integrations", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("wid")
		c.SetParamValues("ws-123")

		err := h.ListIntegrations(c)
		if err == nil {
			t.Fatal("expected handler error")
		}

		var echoErr *echo.HTTPError
		if errors.As(err, &echoErr) {
			if echoErr.Code != http.StatusInternalServerError {
				t.Errorf("expected 500 error, got %d", echoErr.Code)
			}
		} else {
			t.Errorf("expected Echo HTTPError, got %v", err)
		}
	})
}

func TestPortalIntegrationHandler_UpsertIntegration(t *testing.T) {
	e := echo.New()

	t.Run("upserts integration successfully", func(t *testing.T) {
		repo := &fakeIntegrationRepository{}
		svc := service.NewPortalService(service.PortalDeps{
			Integrations: repo,
		})
		h := NewPortalIntegrationHandler(svc)

		body := `{"channel_type":"email","provider_name":"mailgun","config":{"api_key":"secret-key"},"status":"connected","is_default":true}`
		req := httptest.NewRequest(http.MethodPost, "/api/v1/workspaces/ws-123/integrations", strings.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("wid")
		c.SetParamValues("ws-123")

		if err := h.UpsertIntegration(c); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if rec.Code != http.StatusOK {
			t.Errorf("expected 200 OK, got %d", rec.Code)
		}

		if repo.upserted == nil {
			t.Fatal("expected upserted integration to be captured in repo")
		}

		if repo.upserted.WorkspaceID != "ws-123" {
			t.Errorf("expected workspace ws-123, got %s", repo.upserted.WorkspaceID)
		}

		if repo.upserted.ProviderName != "mailgun" {
			t.Errorf("expected provider mailgun, got %s", repo.upserted.ProviderName)
		}

		if !repo.upserted.IsDefault {
			t.Error("expected integration to be marked as default")
		}
	})
}

func TestPortalIntegrationHandler_DeleteIntegration(t *testing.T) {
	e := echo.New()

	t.Run("deletes integration successfully", func(t *testing.T) {
		repo := &fakeIntegrationRepository{
			getByID: &domain.Integration{
				ID:          "intg-789",
				WorkspaceID: "ws-123",
			},
		}
		svc := service.NewPortalService(service.PortalDeps{
			Integrations: repo,
		})
		h := NewPortalIntegrationHandler(svc)

		req := httptest.NewRequest(http.MethodDelete, "/api/v1/workspaces/ws-123/integrations/intg-789", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("wid", "iid")
		c.SetParamValues("ws-123", "intg-789")

		if err := h.DeleteIntegration(c); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if rec.Code != http.StatusNoContent {
			t.Errorf("expected 240 No Content, got %d", rec.Code)
		}

		if repo.deletedID != "intg-789" {
			t.Errorf("expected deleted ID intg-789, got %s", repo.deletedID)
		}
	})

	t.Run("rejects deletion if integration not in workspace", func(t *testing.T) {
		repo := &fakeIntegrationRepository{
			getByID: &domain.Integration{
				ID:          "intg-789",
				WorkspaceID: "ws-different",
			},
		}
		svc := service.NewPortalService(service.PortalDeps{
			Integrations: repo,
		})
		h := NewPortalIntegrationHandler(svc)

		req := httptest.NewRequest(http.MethodDelete, "/api/v1/workspaces/ws-123/integrations/intg-789", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("wid", "iid")
		c.SetParamValues("ws-123", "intg-789")

		err := h.DeleteIntegration(c)
		if err == nil {
			t.Fatal("expected error")
		}
	})
}

func TestPortalIntegrationHandler_GetProviderConfigFields(t *testing.T) {
	e := echo.New()

	t.Run("gets provider config fields successfully", func(t *testing.T) {
		repo := &fakeIntegrationRepository{
			fields: []domain.ProviderConfigField{
				{
					Key:       "api_key",
					Label:     "API Key",
					FieldType: "password",
					Required:  true,
				},
			},
		}
		svc := service.NewPortalService(service.PortalDeps{
			Integrations: repo,
		})
		h := NewPortalIntegrationHandler(svc)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/workspaces/ws-123/providers/mailgun/config", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("wid", "name")
		c.SetParamValues("ws-123", "mailgun")

		if err := h.GetProviderConfigFields(c); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if rec.Code != http.StatusOK {
			t.Errorf("expected 200 OK, got %d", rec.Code)
		}

		var res []domain.ProviderConfigField
		if err := json.Unmarshal(rec.Body.Bytes(), &res); err != nil {
			t.Fatalf("failed to decode JSON response: %v", err)
		}

		if len(res) != 1 {
			t.Fatalf("expected 1 field, got %d", len(res))
		}

		if res[0].Key != "api_key" {
			t.Errorf("expected key api_key, got %s", res[0].Key)
		}
	})
}
