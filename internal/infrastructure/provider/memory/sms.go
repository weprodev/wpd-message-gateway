package memory

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/weprodev/wpd-message-gateway/pkg/contracts"
)

// SMSProvider implements port.SMSSender using an in-memory store.
type SMSProvider struct {
	store *Store
}

// NewSMSProvider creates a new memory SMS provider.
func NewSMSProvider(store *Store) *SMSProvider {
	return &SMSProvider{store: store}
}

// Store returns the underlying memory store.
func (s *SMSProvider) Store() *Store {
	return s.store
}

// Name returns the provider name.
func (s *SMSProvider) Name() string {
	return ProviderName
}

// Send stores the SMS in memory and returns a success result.
func (s *SMSProvider) Send(ctx context.Context, sms *contracts.SMS) (*contracts.SendResult, error) {
	id := uuid.New().String()

	stored := &StoredSMS{
		ID:        id,
		CreatedAt: time.Now(),
		SMS:       sms,
	}
	s.store.AddSMS(stored)

	return &contracts.SendResult{
		ID:         id,
		StatusCode: 200,
		Message:    "Stored SMS in memory",
	}, nil
}
