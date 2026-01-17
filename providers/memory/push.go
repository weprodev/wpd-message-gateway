package memory

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/weprodev/wpd-message-gateway/contracts"
)

// PushProvider wraps the Store to implement contracts.PushSender
type PushProvider struct {
	store *Provider
}

// Store returns the underlying memory store
func (p *PushProvider) Store() *Provider {
	return p.store
}

// Name returns the provider name.
func (p *PushProvider) Name() string {
	return ProviderName
}

// Send stores the push notification in memory and returns a success result.
func (p *PushProvider) Send(ctx context.Context, push *contracts.PushNotification) (*contracts.SendResult, error) {
	id := uuid.New().String()

	stored := &StoredPush{
		ID:        id,
		CreatedAt: time.Now(),
		Push:      push,
	}
	p.store.addPush(stored)

	return &contracts.SendResult{
		ID:         id,
		StatusCode: 200,
		Message:    "Stored push notification in memory",
	}, nil
}
