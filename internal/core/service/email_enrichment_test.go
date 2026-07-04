package service

import (
	"context"
	"testing"

	"github.com/weprodev/wpd-message-gateway/internal/core/domain"
	"github.com/weprodev/wpd-message-gateway/pkg/contracts"
	"github.com/weprodev/wpd-message-gateway/pkg/registry"
)

func TestApplyEmailSenderDefaults(t *testing.T) {
	t.Parallel()

	cfg := registry.EmailConfig{
		FromEmail: "noreply@weprodev.com",
		FromName:  "WeProDev",
	}

	t.Run("fills empty from and from_name", func(t *testing.T) {
		t.Parallel()
		got := applyEmailSenderDefaults(contracts.Email{To: []string{"user@example.com"}}, cfg)
		if got.From != "noreply@weprodev.com" {
			t.Fatalf("From: %q", got.From)
		}
		if got.FromName != "WeProDev" {
			t.Fatalf("FromName: %q", got.FromName)
		}
	})

	t.Run("preserves explicit from", func(t *testing.T) {
		t.Parallel()
		got := applyEmailSenderDefaults(contracts.Email{
			From: "custom@example.com",
			To:   []string{"user@example.com"},
		}, cfg)
		if got.From != "custom@example.com" {
			t.Fatalf("From: %q", got.From)
		}
		if got.FromName != "WeProDev" {
			t.Fatalf("FromName: %q", got.FromName)
		}
	})

	t.Run("preserves explicit from_name", func(t *testing.T) {
		t.Parallel()
		got := applyEmailSenderDefaults(contracts.Email{
			FromName: "Custom Name",
			To:       []string{"user@example.com"},
		}, cfg)
		if got.From != "noreply@weprodev.com" {
			t.Fatalf("From: %q", got.From)
		}
		if got.FromName != "Custom Name" {
			t.Fatalf("FromName: %q", got.FromName)
		}
	})
}

func TestGatewayService_SendEmail_enrichesSenderFromIntegration(t *testing.T) {
	inbox := &capturingInbox{}
	intg := &domain.Integration{
		ID:           "intg-1",
		ProviderName: "mailgun",
		Config:       []byte(`{"from_email":"noreply@weprodev.com","from_name":"WeProDev"}`),
		Status:       "connected",
	}
	svc := NewGatewayService(&stubIntegrationRepo{active: intg}, nil, &stubSettingsRepo{values: map[string]string{
		domain.SettingKeyStoreMessageContent: "true",
	}}, inbox, nil)

	_, err := svc.SendEmail(t.Context(), "ws-1", contracts.Email{
		To:      []string{"user@example.com"},
		Subject: "hi",
		HTML:    "<p>x</p>",
	})
	if err != nil {
		t.Fatalf("SendEmail: %v", err)
	}
	if inbox.lastEmail.From != "noreply@weprodev.com" {
		t.Fatalf("From: %q", inbox.lastEmail.From)
	}
	if inbox.lastEmail.FromName != "WeProDev" {
		t.Fatalf("FromName: %q", inbox.lastEmail.FromName)
	}
}

type capturingInbox struct {
	lastEmail contracts.Email
}

func (s *capturingInbox) WriteEmail(_ context.Context, _ string, email contracts.Email) (string, error) {
	s.lastEmail = email
	return "inbox-msg-1", nil
}

func (s *capturingInbox) WriteSMS(context.Context, string, contracts.SMS) (string, error) {
	return "inbox-sms-1", nil
}

func (s *capturingInbox) WritePush(context.Context, string, contracts.PushNotification) (string, error) {
	return "inbox-push-1", nil
}

func (s *capturingInbox) WriteChat(context.Context, string, contracts.ChatMessage) (string, error) {
	return "inbox-chat-1", nil
}
