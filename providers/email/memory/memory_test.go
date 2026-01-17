package memory

import (
	"context"
	"testing"

	"github.com/weprodev/wpd-message-gateway/contracts"
)

func TestProvider_Send_And_Retrieval(t *testing.T) {
	p := New()

	if p.Name() != ProviderName {
		t.Errorf("Name() = %s, want %s", p.Name(), ProviderName)
	}

	email := &contracts.Email{
		To:      []string{"test@example.com"},
		Subject: "Test",
		HTML:    "<p>Hello</p>",
	}

	// Test Send
	result, err := p.Send(context.Background(), email)
	if err != nil {
		t.Fatalf("Send() error: %v", err)
	}

	if result.ID == "" {
		t.Error("Expected ID in SendResult")
	}

	// Test Retrieval
	if p.Count() != 1 {
		t.Errorf("Count() = %d, want 1", p.Count())
	}

	msgs := p.Messages()
	if len(msgs) != 1 {
		t.Fatalf("Messages() length = %d, want 1", len(msgs))
	}

	if msgs[0] != email {
		t.Error("Stored message does not match sent message")
	}
}

func TestProvider_Clear(t *testing.T) {
	p := New()
	_, _ = p.Send(context.Background(), &contracts.Email{To: []string{"a"}})
	_, _ = p.Send(context.Background(), &contracts.Email{To: []string{"b"}})

	if p.Count() != 2 {
		t.Errorf("Count() = %d, want 2", p.Count())
	}

	p.Clear()

	if p.Count() != 0 {
		t.Errorf("Count() = %d, want 0 after Clear()", p.Count())
	}
}

func TestProvider_Concurrency(t *testing.T) {
	p := New()
	count := 100
	done := make(chan bool)

	// Launch concurrent senders
	for i := 0; i < count; i++ {
		go func() {
			_, _ = p.Send(context.Background(), &contracts.Email{To: []string{"test"}})
			done <- true
		}()
	}

	// Wait for all
	for i := 0; i < count; i++ {
		<-done
	}

	if p.Count() != count {
		t.Errorf("Count() = %d, want %d", p.Count(), count)
	}
}
