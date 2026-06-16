package service

import (
	"context"
	"testing"
	"time"

	"github.com/weprodev/wpd-message-gateway/internal/core/domain"
	"github.com/weprodev/wpd-message-gateway/internal/core/port"
	"github.com/weprodev/wpd-message-gateway/pkg/contracts"
)

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

type stubIntegrationRepo struct {
	active *domain.Integration
	err    error
}

func (s *stubIntegrationRepo) Create(ctx context.Context, integration *domain.Integration) error {
	return nil
}

func (s *stubIntegrationRepo) GetActiveByWorkspaceAndChannel(ctx context.Context, workspaceID, channel string) (*domain.Integration, error) {
	if s.err != nil {
		return nil, s.err
	}
	return s.active, nil
}

func (s *stubIntegrationRepo) ListByWorkspace(ctx context.Context, workspaceID string) ([]domain.Integration, error) {
	return nil, nil
}

func (s *stubIntegrationRepo) GetByID(ctx context.Context, id string) (*domain.Integration, error) {
	return nil, port.ErrNotFound
}

func (s *stubIntegrationRepo) Delete(ctx context.Context, id string) error {
	return nil
}

func (s *stubIntegrationRepo) Upsert(ctx context.Context, integration *domain.Integration) error {
	return nil
}

func (s *stubIntegrationRepo) GetProviderFields(ctx context.Context, providerName string) ([]domain.ProviderConfigField, error) {
	return nil, nil
}

func (s *stubIntegrationRepo) ListProviders(ctx context.Context) ([]domain.Provider, error) {
	return nil, nil
}

type stubInbox struct {
	emailID string
}

func (s *stubInbox) WriteEmail(ctx context.Context, workspaceID string, email contracts.Email) (string, error) {
	if s.emailID != "" {
		return s.emailID, nil
	}
	return "inbox-msg-1", nil
}

func (s *stubInbox) WriteSMS(ctx context.Context, workspaceID string, sms contracts.SMS) (string, error) {
	return "inbox-sms-1", nil
}

func (s *stubInbox) WritePush(ctx context.Context, workspaceID string, push contracts.PushNotification) (string, error) {
	return "inbox-push-1", nil
}

func (s *stubInbox) WriteChat(ctx context.Context, workspaceID string, chat contracts.ChatMessage) (string, error) {
	return "inbox-chat-1", nil
}

type stubStoredMessages struct {
	emailID string
}

func (s *stubStoredMessages) WriteEmail(ctx context.Context, workspaceID string, email contracts.Email) (string, error) {
	if s.emailID != "" {
		return s.emailID, nil
	}
	return "stored-msg-1", nil
}

func (s *stubStoredMessages) WriteSMS(ctx context.Context, workspaceID string, sms contracts.SMS) (string, error) {
	return "stored-sms-1", nil
}

func (s *stubStoredMessages) WritePush(ctx context.Context, workspaceID string, push contracts.PushNotification) (string, error) {
	return "stored-push-1", nil
}

func (s *stubStoredMessages) WriteChat(ctx context.Context, workspaceID string, chat contracts.ChatMessage) (string, error) {
	return "stored-chat-1", nil
}

func TestGatewayService_SendEmail_memoryOnly(t *testing.T) {
	inbox := &stubInbox{emailID: "mem-1"}
	svc := NewGatewayService(&stubIntegrationRepo{}, nil, nil, inbox, nil, nil)

	res, err := svc.SendEmail(context.Background(), "ws-1", contracts.Email{
		To:      []string{"a@b.com"},
		Subject: "hi",
		HTML:    "<p>x</p>",
	})
	if err != nil {
		t.Fatalf("SendEmail: %v", err)
	}
	if res.ID != "mem-1" {
		t.Fatalf("got ID %q", res.ID)
	}
	if res.Meta["dispatch_mode"] != string(domain.DispatchMemoryOnly) {
		t.Fatalf("dispatch_mode: %v", res.Meta["dispatch_mode"])
	}
	if res.Meta["channel"] != "email" {
		t.Fatalf("channel: %v", res.Meta["channel"])
	}
	if res.Meta["provider_name"] != memoryProviderName {
		t.Fatalf("provider_name: %v", res.Meta["provider_name"])
	}
}

func TestGatewayService_SendEmail_memoryOnly_inboxNil(t *testing.T) {
	svc := NewGatewayService(&stubIntegrationRepo{}, nil, nil, nil, nil, nil)

	_, err := svc.SendEmail(context.Background(), "ws-1", contracts.Email{To: []string{"a@b.com"}, Subject: "s"})
	if err == nil {
		t.Fatal("expected error when inbox is nil")
	}
}

func TestGatewayService_SendEmail_providerOnly_memoryIntegration(t *testing.T) {
	ts := time.Now()
	intg := &domain.Integration{
		ID:           "int-1",
		WorkspaceID:  "ws-1",
		ChannelType:  "email",
		ProviderName: memoryProviderName,
		Config:       []byte(`{}`),
		Status:       "connected",
		CreatedAt:    ts,
		UpdatedAt:    ts,
	}
	settings := &stubSettingsRepo{values: map[string]string{
		domain.SettingKeyMessageDispatchMode: string(domain.DispatchProviderOnly),
	}}
	inbox := &stubInbox{emailID: "cap-1"}
	svc := NewGatewayService(&stubIntegrationRepo{active: intg}, nil, settings, inbox, nil, nil)

	res, err := svc.SendEmail(context.Background(), "ws-1", contracts.Email{
		To: []string{"a@b.com"}, Subject: "s", HTML: "h",
	})
	if err != nil {
		t.Fatalf("SendEmail: %v", err)
	}
	if res.ID != "cap-1" {
		t.Fatalf("got ID %q", res.ID)
	}
	if res.Meta["dispatch_mode"] != string(domain.DispatchProviderOnly) {
		t.Fatalf("dispatch_mode: %v", res.Meta["dispatch_mode"])
	}
	if res.Meta["integration_id"] != "int-1" {
		t.Fatalf("integration_id: %v", res.Meta["integration_id"])
	}
}

func TestGatewayService_SendEmail_providerAndDatabase_memoryIntegration(t *testing.T) {
	ts := time.Now()
	intg := &domain.Integration{
		ID:           "int-1",
		WorkspaceID:  "ws-1",
		ChannelType:  "email",
		ProviderName: memoryProviderName,
		Config:       []byte(`{}`),
		Status:       "connected",
		CreatedAt:    ts,
		UpdatedAt:    ts,
	}
	settings := &stubSettingsRepo{values: map[string]string{
		domain.SettingKeyMessageDispatchMode: string(domain.DispatchProviderAndDatabase),
	}}
	inbox := &stubInbox{emailID: "cap-1"}
	svc := NewGatewayService(&stubIntegrationRepo{active: intg}, nil, settings, inbox, &stubStoredMessages{}, nil)

	res, err := svc.SendEmail(context.Background(), "ws-1", contracts.Email{
		To: []string{"a@b.com"}, Subject: "s", HTML: "h",
	})
	if err != nil {
		t.Fatalf("SendEmail: %v", err)
	}
	if res.ID != "cap-1" {
		t.Fatalf("got ID %q", res.ID)
	}
	if res.Meta["dispatch_mode"] != string(domain.DispatchProviderAndDatabase) {
		t.Fatalf("dispatch_mode: %v", res.Meta["dispatch_mode"])
	}
}
