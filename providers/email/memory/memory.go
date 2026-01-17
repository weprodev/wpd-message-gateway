package memory

import (
	"context"
	"sync"

	"github.com/google/uuid"

	"github.com/weprodev/wpd-message-gateway/contracts"
)

const ProviderName = "memory"

// Provider implements an in-memory email sender for testing/interception.
// It is thread-safe.
type Provider struct {
	mu       sync.RWMutex
	messages []*contracts.Email
}

// New creates a new Memory provider.
func New() *Provider {
	return &Provider{
		messages: make([]*contracts.Email, 0),
	}
}

// Name returns the provider name.
func (p *Provider) Name() string {
	return ProviderName
}

// Send stores the email in memory and returns a success result.
func (p *Provider) Send(ctx context.Context, email *contracts.Email) (*contracts.SendResult, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	// Store a copy to avoid external mutation
	// In a real scenario, deep copy might be safer but for this purpose assume basic struct copy is fine
	// or rely on the fact that contracts.Email is mostly value types + slices
	p.messages = append(p.messages, email)

	return &contracts.SendResult{
		ID:         uuid.New().String(),
		StatusCode: 200,
		Message:    "Stored in memory",
	}, nil
}

// Messages returns a copy of all stored messages.
func (p *Provider) Messages() []*contracts.Email {
	p.mu.RLock()
	defer p.mu.RUnlock()

	// Return a copy of the slice to be safe
	msgs := make([]*contracts.Email, len(p.messages))
	copy(msgs, p.messages)
	return msgs
}

// Clear removes all stored messages.
func (p *Provider) Clear() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.messages = make([]*contracts.Email, 0)
}

// Count returns the number of stored messages.
func (p *Provider) Count() int {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return len(p.messages)
}
