package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/weprodev/wpd-message-gateway/internal/core/domain"
	"github.com/weprodev/wpd-message-gateway/internal/core/port"
	"github.com/weprodev/wpd-message-gateway/internal/core/service"
)

type handlerFakeAPIKeyRepo struct {
	keys    map[string]*domain.APIKey
	byWS    map[string][]domain.APIKey
	created *domain.APIKey
	deleted string
}

func (f *handlerFakeAPIKeyRepo) Create(_ context.Context, k *domain.APIKey) error {
	if f.keys == nil {
		f.keys = make(map[string]*domain.APIKey)
	}
	if f.byWS == nil {
		f.byWS = make(map[string][]domain.APIKey)
	}
	k.ID = "key-" + uuid.NewString()
	f.created = k
	f.keys[k.ID] = k
	f.byWS[k.WorkspaceID] = append(f.byWS[k.WorkspaceID], *k)
	return nil
}

func (f *handlerFakeAPIKeyRepo) GetByClientID(context.Context, string) (*domain.APIKey, error) {
	return nil, port.ErrNotFound
}

func (f *handlerFakeAPIKeyRepo) GetByID(_ context.Context, id string) (*domain.APIKey, error) {
	k, ok := f.keys[id]
	if !ok {
		return nil, port.ErrNotFound
	}
	return k, nil
}

func (f *handlerFakeAPIKeyRepo) ListByWorkspace(_ context.Context, workspaceID string) ([]domain.APIKey, error) {
	return f.byWS[workspaceID], nil
}

func (f *handlerFakeAPIKeyRepo) Delete(_ context.Context, id string) error {
	f.deleted = id
	delete(f.keys, id)
	return nil
}

func (f *handlerFakeAPIKeyRepo) UpdateSecret(context.Context, string, string, string) error {
	return nil
}

func (f *handlerFakeAPIKeyRepo) UpdateLastUsedAt(context.Context, string) error {
	return nil
}

type handlerFakeSettingsRepo struct {
	values map[string]string
}

func (f *handlerFakeSettingsRepo) Get(_ context.Context, _, key string) (string, error) {
	if f.values == nil {
		return "", nil
	}
	return f.values[key], nil
}

func (f *handlerFakeSettingsRepo) Set(_ context.Context, _, key, value string) error {
	if f.values == nil {
		f.values = make(map[string]string)
	}
	f.values[key] = value
	return nil
}

func (f *handlerFakeSettingsRepo) GetAll(_ context.Context, _ string) (map[string]string, error) {
	return f.values, nil
}

func TestPortalHandler_CreateAPIKey(t *testing.T) {
	t.Parallel()
	e := echo.New()

	repo := &handlerFakeAPIKeyRepo{}
	svc := service.NewPortalService(service.PortalDeps{APIKeys: repo})
	h := NewPortalHandler(svc)

	body := `{"name":"Portal Key"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/workspaces/ws-1/api-keys", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("wid")
	c.SetParamValues("ws-1")

	if err := h.CreateAPIKey(c); err != nil {
		t.Fatalf("CreateAPIKey: %v", err)
	}
	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d body=%s", rec.Code, rec.Body.String())
	}

	var res map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &res); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if res["name"] != "Portal Key" {
		t.Fatalf("unexpected response: %+v", res)
	}
	secret, _ := res["client_secret"].(string)
	if secret == "" {
		t.Fatal("expected client_secret in response")
	}
}

func TestPortalHandler_CreateAPIKey_rejectsMissingName(t *testing.T) {
	t.Parallel()
	e := echo.New()

	h := NewPortalHandler(service.NewPortalService(service.PortalDeps{}))

	req := httptest.NewRequest(http.MethodPost, "/api/v1/workspaces/ws-1/api-keys", strings.NewReader(`{}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("wid")
	c.SetParamValues("ws-1")

	err := h.CreateAPIKey(c)
	if err == nil {
		t.Fatal("expected validation error")
	}
	httpErr, ok := err.(*echo.HTTPError)
	if !ok || httpErr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %v", err)
	}
}

func TestPortalHandler_ListAPIKeys(t *testing.T) {
	t.Parallel()
	e := echo.New()

	ts := time.Now()
	repo := &handlerFakeAPIKeyRepo{
		byWS: map[string][]domain.APIKey{
			"ws-1": {{
				ID:          "key-1",
				WorkspaceID: "ws-1",
				ClientID:    "client-1",
				Name:        "Demo",
				IsActive:    true,
				CreatedAt:   ts,
			}},
		},
	}
	svc := service.NewPortalService(service.PortalDeps{APIKeys: repo})
	h := NewPortalHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/workspaces/ws-1/api-keys", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("wid")
	c.SetParamValues("ws-1")

	if err := h.ListAPIKeys(c); err != nil {
		t.Fatalf("ListAPIKeys: %v", err)
	}
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var keys []map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &keys); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if len(keys) != 1 || keys[0]["name"] != "Demo" {
		t.Fatalf("unexpected keys: %+v", keys)
	}
}

func TestPortalHandler_GetSettings(t *testing.T) {
	t.Parallel()
	e := echo.New()

	settings := &handlerFakeSettingsRepo{
		values: map[string]string{
			domain.SettingKeyMessageDispatchMode: string(domain.DispatchMemory),
		},
	}
	svc := service.NewPortalService(service.PortalDeps{Settings: settings})
	h := NewPortalHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/workspaces/ws-1/settings", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("wid")
	c.SetParamValues("ws-1")

	if err := h.GetSettings(c); err != nil {
		t.Fatalf("GetSettings: %v", err)
	}

	var res map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &res); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if res[domain.SettingKeyMessageDispatchMode] != string(domain.DispatchMemory) {
		t.Fatalf("unexpected settings: %+v", res)
	}
}
