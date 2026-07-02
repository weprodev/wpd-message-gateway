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

type stubSettingsRepo struct {
	values map[string]string
}

func (s *stubSettingsRepo) Get(ctx context.Context, workspaceID, key string) (string, error) {
	if s.values == nil {
		return "", nil
	}
	return s.values[key], nil
}

func (s *stubSettingsRepo) Set(ctx context.Context, workspaceID, key, value string) error {
	return nil
}

func (s *stubSettingsRepo) GetAll(ctx context.Context, workspaceID string) (map[string]string, error) {
	return nil, nil
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
		if log.Retained {
			t.Error("expected retained false with default memory-only mode")
		}
	})

	t.Run("retained true for provider_database setting", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/v1/email", strings.NewReader(`{"to":["test@example.com"]}`))
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		settings := &stubSettingsRepo{values: map[string]string{
			domain.SettingKeyMessageDispatchMode: string(domain.APIProviderAndDatabase),
		}}
		repo := &mockLogRepo{}
		svc := service.NewGatewayService(nil, nil, settings, nil, repo)
		helper := NewSendHelper(svc)

		var dst contracts.Email
		err := helper.DispatchAndLog(c, "email", "ws-123", "key-456", "/v1/email", &dst, func(ctx context.Context) (*contracts.SendResult, error) {
			return &contracts.SendResult{ID: "msg-111"}, nil
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(repo.entries) != 1 {
			t.Fatalf("expected 1 log entry, got %d", len(repo.entries))
		}
		if !repo.entries[0].Retained {
			t.Error("expected retained true for provider_database")
		}
	})

	t.Run("retained true for memory_and_database setting", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/v1/email", strings.NewReader(`{"to":["test@example.com"]}`))
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		settings := &stubSettingsRepo{values: map[string]string{
			domain.SettingKeyMessageDispatchMode: string(domain.APIMemoryAndDatabase),
		}}
		repo := &mockLogRepo{}
		svc := service.NewGatewayService(nil, nil, settings, nil, repo)
		helper := NewSendHelper(svc)

		var dst contracts.Email
		err := helper.DispatchAndLog(c, "email", "ws-123", "key-456", "/v1/email", &dst, func(ctx context.Context) (*contracts.SendResult, error) {
			return &contracts.SendResult{ID: "msg-111"}, nil
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(repo.entries) != 1 {
			t.Fatalf("expected 1 log entry, got %d", len(repo.entries))
		}
		if !repo.entries[0].Retained {
			t.Error("expected retained true for memory_and_database")
		}
	})

	t.Run("retained false for provider_only setting", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/v1/email", strings.NewReader(`{"to":["test@example.com"]}`))
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		settings := &stubSettingsRepo{values: map[string]string{
			domain.SettingKeyMessageDispatchMode: string(domain.APIProviderOnly),
		}}
		repo := &mockLogRepo{}
		svc := service.NewGatewayService(nil, nil, settings, nil, repo)
		helper := NewSendHelper(svc)

		var dst contracts.Email
		err := helper.DispatchAndLog(c, "email", "ws-123", "key-456", "/v1/email", &dst, func(ctx context.Context) (*contracts.SendResult, error) {
			return &contracts.SendResult{ID: "msg-111"}, nil
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(repo.entries) != 1 {
			t.Fatalf("expected 1 log entry, got %d", len(repo.entries))
		}
		if repo.entries[0].Retained {
			t.Error("expected retained false for provider_only")
		}
	})

	t.Run("repository error logged without failing dispatch", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/v1/email", strings.NewReader(`{"to":["test@example.com"]}`))
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		repo := &mockLogRepo{createErr: errors.New("db error")}
		svc := service.NewGatewayService(nil, nil, nil, nil, repo)
		helper := NewSendHelper(svc)

		var dst contracts.Email
		err := helper.DispatchAndLog(c, "email", "ws-123", "key-456", "/v1/email", &dst, func(ctx context.Context) (*contracts.SendResult, error) {
			return &contracts.SendResult{ID: "msg-111"}, nil
		})

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if rec.Code != http.StatusOK {
			t.Errorf("expected status OK, got %d", rec.Code)
		}
	})
}
