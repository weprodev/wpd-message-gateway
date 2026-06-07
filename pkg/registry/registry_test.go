package registry

import (
	"context"
	"testing"

	"github.com/weprodev/wpd-message-gateway/pkg/contracts"
)

type mockEmailSender struct{}

func (m *mockEmailSender) Send(ctx context.Context, e contracts.Email) (*contracts.SendResult, error) {
	return nil, nil
}
func (m *mockEmailSender) Name() string { return "mock" }

func TestRegistry(t *testing.T) {
	// Register a mock
	RegisterEmailProvider("mock", func(cfg EmailConfig) (contracts.EmailSender, error) {
		return &mockEmailSender{}, nil
	})

	// Test success
	factory, err := GetEmailFactory("mock")
	if err != nil {
		t.Errorf("GetEmailFactory failed: %v", err)
	}
	if factory == nil {
		t.Errorf("expected non-nil factory")
	}

	// Test missing
	_, err = GetEmailFactory("missing")
	if err == nil {
		t.Errorf("expected error for missing provider, got nil")
	}
}
