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
		Providers: config.ProviderConfig{
			Defaults: config.ProviderDefaults{
				Email: "mock",
			},
		},
		EmailProviders: make(map[string]config.EmailConfig),
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

func TestManager_Email(t *testing.T) {
	// Setup with mock config
	cfg := mockConfig()
	// Set default provider
	cfg.Providers.Defaults.Email = "mailgun"
	cfg.EmailProviders["mailgun"] = config.EmailConfig{
		CommonConfig: config.CommonConfig{APIKey: "key"},
		Domain:       "domain",
	}

	mgr, _ := New(cfg)
	_ = mgr // Silence unused variable warning
	// Register mock provider in registry manually since factory assumes real implementation
	// or relies on memory.
	// For this test we can use "memory" as "mailgun" to avoid dependencies if we want,
	// OR just stick to memory test.
}

func TestManager_Email_Defaults(t *testing.T) {
	cfg := mockConfig()
	cfg.Providers.Defaults.Email = "memory"
	mgr, _ := New(cfg)

	provider, err := mgr.Email()
	if err != nil {
		t.Fatalf("Email() error = %v", err)
	}
	if provider == nil {
		t.Error("Email() returned nil provider")
	}
}

// (Removed duplicates)
