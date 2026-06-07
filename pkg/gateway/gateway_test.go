package gateway

import (
	"context"
	"testing"

	"github.com/weprodev/wpd-message-gateway/pkg/contracts"
	"github.com/weprodev/wpd-message-gateway/pkg/registry"
)

type mockEmailSender struct{ called bool }

func (m *mockEmailSender) Send(ctx context.Context, e contracts.Email) (*contracts.SendResult, error) {
	m.called = true
	return &contracts.SendResult{ID: "mock"}, nil
}
func (m *mockEmailSender) Name() string { return "mock" }

type mockSMSSender struct{ called bool }

func (m *mockSMSSender) Send(ctx context.Context, e contracts.SMS) (*contracts.SendResult, error) {
	m.called = true
	return &contracts.SendResult{ID: "mock"}, nil
}
func (m *mockSMSSender) Name() string { return "mock" }

type mockPushSender struct{ called bool }

func (m *mockPushSender) Send(ctx context.Context, e contracts.PushNotification) (*contracts.SendResult, error) {
	m.called = true
	return &contracts.SendResult{ID: "mock"}, nil
}
func (m *mockPushSender) Name() string { return "mock" }

type mockChatSender struct{ called bool }

func (m *mockChatSender) Send(ctx context.Context, e contracts.ChatMessage) (*contracts.SendResult, error) {
	m.called = true
	return &contracts.SendResult{ID: "mock"}, nil
}
func (m *mockChatSender) Name() string { return "mock" }

func init() {
	registry.RegisterEmailProvider("mock", func(cfg registry.EmailConfig) (contracts.EmailSender, error) {
		return &mockEmailSender{}, nil
	})
	registry.RegisterSMSProvider("mock", func(cfg registry.SMSConfig) (contracts.SMSSender, error) {
		return &mockSMSSender{}, nil
	})
	registry.RegisterPushProvider("mock", func(cfg registry.PushConfig) (contracts.PushSender, error) {
		return &mockPushSender{}, nil
	})
	registry.RegisterChatProvider("mock", func(cfg registry.ChatConfig) (contracts.ChatSender, error) {
		return &mockChatSender{}, nil
	})
}

func TestGatewayNew(t *testing.T) {
	tests := []struct {
		name    string
		cfg     Config
		wantErr bool
	}{
		{
			name: "valid config",
			cfg: Config{
				DefaultEmailProvider: "mock",
				EmailProviders:       map[string]EmailConfig{"mock": {}},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := New(tt.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGatewaySend(t *testing.T) {
	cfg := Config{
		DefaultEmailProvider: "mock",
		DefaultSMSProvider:   "mock",
		DefaultPushProvider:  "mock",
		DefaultChatProvider:  "mock",
		EmailProviders:       map[string]EmailConfig{"mock": {}},
		SMSProviders:         map[string]SMSConfig{"mock": {}},
		PushProviders:        map[string]PushConfig{"mock": {}},
		ChatProviders:        map[string]ChatConfig{"mock": {}},
	}
	gw, err := New(cfg)
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}

	res, err := gw.SendEmail(context.Background(), contracts.Email{})
	if err != nil {
		t.Errorf("SendEmail() error = %v", err)
	}
	if res.ID != "mock" {
		t.Errorf("expected mock ID, got %s", res.ID)
	}

	res, err = gw.SendSMS(context.Background(), contracts.SMS{})
	if err != nil {
		t.Errorf("SendSMS() error = %v", err)
	}
	if res.ID != "mock" {
		t.Errorf("expected mock ID, got %s", res.ID)
	}

	res, err = gw.SendPush(context.Background(), contracts.PushNotification{})
	if err != nil {
		t.Errorf("SendPush() error = %v", err)
	}
	if res.ID != "mock" {
		t.Errorf("expected mock ID, got %s", res.ID)
	}

	res, err = gw.SendChat(context.Background(), contracts.ChatMessage{})
	if err != nil {
		t.Errorf("SendChat() error = %v", err)
	}
	if res.ID != "mock" {
		t.Errorf("expected mock ID, got %s", res.ID)
	}
}
