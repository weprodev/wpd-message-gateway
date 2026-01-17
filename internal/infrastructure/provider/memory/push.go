package memory

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/weprodev/wpd-message-gateway/pkg/contracts"
)

// PushProvider implements port.PushSender using an in-memory store.
type PushProvider struct {
	store *Store
}

// NewPushProvider creates a new memory push provider.
func NewPushProvider(store *Store) *PushProvider {
	return &PushProvider{store: store}
}

// Store returns the underlying memory store.
func (p *PushProvider) Store() *Store {
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
	p.store.AddPush(stored)

	return &contracts.SendResult{
		ID:         id,
		StatusCode: 200,
		Message:    "Stored push notification in memory",
	}, nil
}
