package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"

	"github.com/weprodev/wpd-message-gateway/internal/infrastructure/inbox"
	"github.com/weprodev/wpd-message-gateway/pkg/contracts"
)

func TestPortalInboxHandler_HandleGetEmails(t *testing.T) {
	t.Parallel()
	e := echo.New()

	store := inbox.NewStore()
	ctx := context.Background()
	id, err := store.WriteEmail(ctx, "ws-1", contracts.Email{
		To:      []string{"a@b.com"},
		Subject: "Hello",
		HTML:    "<p>hi</p>",
	})
	if err != nil {
		t.Fatalf("WriteEmail: %v", err)
	}

	h := NewPortalInboxHandler(store, store)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/workspaces/ws-1/inbox/emails?limit=10", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("wid")
	c.SetParamValues("ws-1")

	if err := h.HandleGetEmails(c); err != nil {
		t.Fatalf("HandleGetEmails: %v", err)
	}
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var page struct {
		Items []struct {
			ID string `json:"id"`
		} `json:"items"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &page); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if len(page.Items) != 1 || page.Items[0].ID != id {
		t.Fatalf("unexpected page: %+v", page)
	}
}

func TestPortalInboxHandler_HandleGetSMS(t *testing.T) {
	t.Parallel()
	e := echo.New()

	store := inbox.NewStore()
	ctx := context.Background()
	id, err := store.WriteSMS(ctx, "ws-1", contracts.SMS{
		To:      []string{"+15550001111"},
		Message: "ping",
	})
	if err != nil {
		t.Fatalf("WriteSMS: %v", err)
	}

	h := NewPortalInboxHandler(store, store)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/workspaces/ws-1/inbox/sms?limit=10", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("wid")
	c.SetParamValues("ws-1")

	if err := h.HandleGetSMS(c); err != nil {
		t.Fatalf("HandleGetSMS: %v", err)
	}

	var page struct {
		Items []struct {
			ID string `json:"id"`
		} `json:"items"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &page); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if len(page.Items) != 1 || page.Items[0].ID != id {
		t.Fatalf("unexpected page: %+v", page)
	}
}

func TestPortalInboxHandler_HandleGetEmailByID_notFound(t *testing.T) {
	t.Parallel()
	e := echo.New()

	h := NewPortalInboxHandler(inbox.NewStore(), inbox.NewStore())

	req := httptest.NewRequest(http.MethodGet, "/api/v1/workspaces/ws-1/inbox/emails/missing", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("wid", "id")
	c.SetParamValues("ws-1", "missing")

	if err := h.HandleGetEmailByID(c); err != nil {
		t.Fatalf("HandleGetEmailByID: %v", err)
	}
	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}
}

func TestPortalInboxHandler_HandleStats(t *testing.T) {
	t.Parallel()
	e := echo.New()

	store := inbox.NewStore()
	ctx := context.Background()
	if _, err := store.WriteEmail(ctx, "ws-1", contracts.Email{To: []string{"a@b.com"}, Subject: "s"}); err != nil {
		t.Fatalf("WriteEmail: %v", err)
	}

	h := NewPortalInboxHandler(store, store)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/workspaces/ws-1/inbox/stats", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("wid")
	c.SetParamValues("ws-1")

	if err := h.HandleStats(c); err != nil {
		t.Fatalf("HandleStats: %v", err)
	}

	var stats map[string]int
	if err := json.Unmarshal(rec.Body.Bytes(), &stats); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if stats["emails"] < 1 {
		t.Fatalf("expected at least one email in stats, got %+v", stats)
	}
}
