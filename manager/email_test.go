package manager

import (
	"context"
	"testing"

	"github.com/weprodev/wpd-message-gateway/config"
	"github.com/weprodev/wpd-message-gateway/contracts"
)

// mockEmailSender is a mock implementation of contracts.EmailSender for testing
type mockEmailSender struct {
	name       string
	sendCalled bool
	lastEmail  *contracts.Email
	sendResult *contracts.SendResult
	sendErr    error
}

func (m *mockEmailSender) Send(ctx context.Context, email *contracts.Email) (*contracts.SendResult, error) {
	m.sendCalled = true
	m.lastEmail = email
	return m.sendResult, m.sendErr
}

func (m *mockEmailSender) Name() string {
	return m.name
}

func TestManager_SendEmail(t *testing.T) {
	cfg := &config.Config{
		DefaultEmailProvider: "mock",
		EmailProviders:       make(map[string]config.EmailConfig),
	}
	mgr, err := New(cfg)
	if err != nil {
		t.Fatalf("Failed to create manager: %v", err)
	}

	mock := &mockEmailSender{
		name: "mock",
		sendResult: &contracts.SendResult{
			ID:         "test-id-123",
			StatusCode: 200,
			Message:    "Email sent successfully",
		},
	}
	mgr.RegisterEmailProvider("mock", mock)

	email := &contracts.Email{
		To:      []string{"test@example.com"},
		Subject: "Test",
		HTML:    "<h1>Test</h1>",
	}

	result, err := mgr.SendEmail(context.Background(), email)
	if err != nil {
		t.Fatalf("SendEmail() error: %v", err)
	}

	// Verify result
	if result.ID != "test-id-123" {
		t.Errorf("Expected ID 'test-id-123', got %s", result.ID)
	}

	// Verify mock was called with correct email
	if !mock.sendCalled {
		t.Error("Expected mock.Send() to be called")
	}
	if mock.lastEmail != email {
		t.Error("Expected email to be passed to mock")
	}
}

func TestManager_Email_NoDefault(t *testing.T) {
	cfg := &config.Config{
		DefaultEmailProvider: "", // No default set
		EmailProviders:       make(map[string]config.EmailConfig),
	}
	mgr, err := New(cfg)
	if err != nil {
		t.Fatalf("Failed to create manager: %v", err)
	}

	// Should error when no default provider configured
	_, err = mgr.Email()
	if err == nil {
		t.Error("Email() expected error when no default configured")
	}
}

func TestManager_EmailProvider_NotFound(t *testing.T) {
	cfg := &config.Config{
		EmailProviders: make(map[string]config.EmailConfig),
	}
	mgr, err := New(cfg)
	if err != nil {
		t.Fatalf("Failed to create manager: %v", err)
	}

	// Should error for nonexistent provider
	_, err = mgr.EmailProvider("nonexistent")
	if err == nil {
		t.Error("EmailProvider() expected error for nonexistent provider")
	}
}
