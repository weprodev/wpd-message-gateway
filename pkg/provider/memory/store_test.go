package memory

import (
	"context"
	"testing"

	"github.com/weprodev/wpd-message-gateway/pkg/contracts"
)

func TestMemoryStore(t *testing.T) {
	s := GetStore()
	if s == nil {
		t.Fatal("expected store")
	}

	// Email
	provider := NewEmailProvider(s)
	if provider.Name() != ProviderName {
		t.Errorf("expected %s, got %s", ProviderName, provider.Name())
	}
	res, err := provider.Send(context.Background(), contracts.Email{Subject: "test"})
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if res.ID == "" {
		t.Error("expected ID")
	}

	emails := s.Emails()
	if len(emails) != 1 {
		t.Errorf("expected 1 email, got %d", len(emails))
	}

	s.Clear()
	if len(s.Emails()) != 0 {
		t.Error("expected 0 emails after clear")
	}
}
