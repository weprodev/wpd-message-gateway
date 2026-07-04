package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"

	"github.com/weprodev/wpd-message-gateway/internal/core/domain"
	"github.com/weprodev/wpd-message-gateway/internal/core/port"
	"github.com/weprodev/wpd-message-gateway/internal/core/service"
	"github.com/weprodev/wpd-message-gateway/internal/infrastructure/logger"
	"github.com/weprodev/wpd-message-gateway/pkg/contracts"
)

type mockLogRepo struct {
	createErr error
	entries   []*domain.MessageRequestLog
}

func (m *mockLogRepo) Create(ctx context.Context, entry *domain.MessageRequestLog) error {
	m.entries = append(m.entries, entry)
	return m.createErr
}

func (m *mockLogRepo) ListWithSource(ctx context.Context, q port.MessageLogQuery) ([]domain.MessageRequestLogWithSource, int, error) {
	return nil, 0, nil
}

func TestSendHelper_DispatchAndLog(t *testing.T) {
	e := echo.New()

	t.Run("successful dispatch", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/v1/email", strings.NewReader(`{"to":["test@example.com"]}`))
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Set request ID in context
		ctx := req.Context()
		ctx = logger.WithRequestID(ctx, "req-xyz")
		c.SetRequest(req.WithContext(ctx))

		repo := &mockLogRepo{}
		svc := service.NewGatewayService(nil, nil, nil, nil, repo)
		helper := NewSendHelper(svc)

		var dst contracts.Email
		err := helper.DispatchAndLog(c, "email", "ws-123", "key-456", "/v1/email", &dst, func(ctx context.Context) (*contracts.SendResult, error) {
			return &contracts.SendResult{
				ID: "msg-111",
				Meta: map[string]string{
					"provider_name": "mailgun",
				},
			}, nil
		})

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if rec.Code != http.StatusOK {
			t.Errorf("expected status OK, got %d", rec.Code)
		}

		if len(repo.entries) != 1 {
			t.Fatalf("expected 1 log entry, got %d", len(repo.entries))
		}

		log := repo.entries[0]
		if log.WorkspaceID != "ws-123" {
			t.Errorf("expected workspace ID ws-123, got %s", log.WorkspaceID)
		}
		if log.APIKeyID != "key-456" {
			t.Errorf("expected API key ID key-456, got %s", log.APIKeyID)
		}
		if log.RequestID != "req-xyz" {
			t.Errorf("expected request ID req-xyz, got %s", log.RequestID)
		}
		if log.ProviderName != "mailgun" {
			t.Errorf("expected provider name mailgun, got %s", log.ProviderName)
		}
		if log.StatusCode != http.StatusOK {
			t.Errorf("expected log status code 200, got %d", log.StatusCode)
		}
	})

	t.Run("successful dispatch stores inbox message id in memory mode", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/v1/email", strings.NewReader(`{"to":["test@example.com"]}`))
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		repo := &mockLogRepo{}
		svc := service.NewGatewayService(nil, nil, nil, nil, repo)
		helper := NewSendHelper(svc)

		var dst contracts.Email
		err := helper.DispatchAndLog(c, "email", "ws-123", "key-456", "/v1/email", &dst, func(ctx context.Context) (*contracts.SendResult, error) {
			return &contracts.SendResult{
				ID: "inbox-msg-1",
				Meta: map[string]string{
					contracts.MetaKeyProviderName: "memory",
					contracts.MetaKeyDispatchMode: contracts.DispatchModeMemory,
				},
			}, nil
		})

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(repo.entries) != 1 {
			t.Fatalf("expected 1 log entry, got %d", len(repo.entries))
		}
		if repo.entries[0].InboxMessageID != "inbox-msg-1" {
			t.Errorf("expected inbox message id inbox-msg-1, got %s", repo.entries[0].InboxMessageID)
		}
	})

	t.Run("provider dispatch with inbox capture stores inbox_message_id from meta", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/v1/email", strings.NewReader(`{"to":["test@example.com"]}`))
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		repo := &mockLogRepo{}
		svc := service.NewGatewayService(nil, nil, nil, nil, repo)
		helper := NewSendHelper(svc)

		var dst contracts.Email
		err := helper.DispatchAndLog(c, "email", "ws-123", "key-456", "/v1/email", &dst, func(ctx context.Context) (*contracts.SendResult, error) {
			return &contracts.SendResult{
				ID: "provider-msg-1",
				Meta: map[string]string{
					contracts.MetaKeyProviderName:   "mailgun",
					contracts.MetaKeyDispatchMode:   contracts.DispatchModeProvider,
					contracts.MetaKeyInboxMessageID: "inbox-msg-2",
				},
			}, nil
		})

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if repo.entries[0].InboxMessageID != "inbox-msg-2" {
			t.Errorf("expected inbox message id inbox-msg-2, got %s", repo.entries[0].InboxMessageID)
		}
	})

	t.Run("store_content persists request and response payloads", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/v1/email", strings.NewReader(`{"to":["test@example.com"],"subject":"hi"}`))
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		repo := &mockLogRepo{}
		svc := service.NewGatewayService(nil, nil, nil, nil, repo)
		helper := NewSendHelper(svc)

		var dst contracts.Email
		err := helper.DispatchAndLog(c, "email", "ws-123", "key-456", "/v1/email", &dst, func(ctx context.Context) (*contracts.SendResult, error) {
			r := &contracts.SendResult{ID: "msg-111", StatusCode: 200, Message: "sent"}
			contracts.SetStoreContentMeta(r, true)
			r.Meta[contracts.MetaKeyProviderName] = "mailgun"
			return r, nil
		})

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(repo.entries) != 1 {
			t.Fatalf("expected 1 log entry, got %d", len(repo.entries))
		}

		log := repo.entries[0]
		if log.Payload == nil {
			t.Fatal("expected payload on log entry")
		}
		if !strings.Contains(log.Payload.RequestBody, "test@example.com") {
			t.Errorf("request body missing recipient: %q", log.Payload.RequestBody)
		}
		if !strings.Contains(log.Payload.ResponseBody, "msg-111") {
			t.Errorf("response body missing message id: %q", log.Payload.ResponseBody)
		}
	})

	t.Run("store_content off skips payload", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/v1/email", strings.NewReader(`{"to":["test@example.com"]}`))
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		repo := &mockLogRepo{}
		svc := service.NewGatewayService(nil, nil, nil, nil, repo)
		helper := NewSendHelper(svc)

		var dst contracts.Email
		err := helper.DispatchAndLog(c, "email", "ws-123", "key-456", "/v1/email", &dst, func(ctx context.Context) (*contracts.SendResult, error) {
			r := &contracts.SendResult{ID: "msg-111"}
			contracts.SetStoreContentMeta(r, false)
			return r, nil
		})

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(repo.entries) != 1 {
			t.Fatalf("expected 1 log entry, got %d", len(repo.entries))
		}
		if repo.entries[0].Payload != nil {
			t.Fatalf("expected no payload when store_content is false, got %+v", repo.entries[0].Payload)
		}
	})

	t.Run("missing workspace returns unauthorized", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/v1/email", strings.NewReader(`{"to":["test@example.com"]}`))
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		repo := &mockLogRepo{}
		svc := service.NewGatewayService(nil, nil, nil, nil, repo)
		helper := NewSendHelper(svc)

		var dst contracts.Email
		err := helper.DispatchAndLog(c, "email", "", "key-456", "/v1/email", &dst, func(ctx context.Context) (*contracts.SendResult, error) {
			return &contracts.SendResult{ID: "msg-111"}, nil
		})

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if rec.Code != http.StatusUnauthorized {
			t.Errorf("expected status 401, got %d", rec.Code)
		}
		if len(repo.entries) != 0 {
			t.Fatalf("expected no log entry without workspace, got %d", len(repo.entries))
		}
	})

	t.Run("invalid JSON returns bad request", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/v1/email", strings.NewReader(`{invalid`))
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		repo := &mockLogRepo{}
		svc := service.NewGatewayService(nil, nil, nil, nil, repo)
		helper := NewSendHelper(svc)

		var dst contracts.Email
		err := helper.DispatchAndLog(c, "email", "ws-123", "key-456", "/v1/email", &dst, func(ctx context.Context) (*contracts.SendResult, error) {
			return &contracts.SendResult{ID: "msg-111"}, nil
		})

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d", rec.Code)
		}
	})

	t.Run("send failure returns 500 and logs error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/v1/email", strings.NewReader(`{"to":["test@example.com"]}`))
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		repo := &mockLogRepo{}
		svc := service.NewGatewayService(nil, nil, nil, nil, repo)
		helper := NewSendHelper(svc)

		var dst contracts.Email
		err := helper.DispatchAndLog(c, "email", "ws-123", "key-456", "/v1/email", &dst, func(ctx context.Context) (*contracts.SendResult, error) {
			return nil, errors.New("dispatch failed")
		})

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if rec.Code != http.StatusInternalServerError {
			t.Errorf("expected status 500, got %d", rec.Code)
		}
		if len(repo.entries) != 1 || repo.entries[0].StatusCode != http.StatusInternalServerError {
			t.Fatalf("expected error log entry, got %+v", repo.entries)
		}
	})
}
