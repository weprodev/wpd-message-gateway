package memory

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/weprodev/wpd-message-gateway/contracts"
)

// SMSProvider wraps the Store to implement contracts.SMSSender
type SMSProvider struct {
	store *Provider
}

// Store returns the underlying memory store
func (s *SMSProvider) Store() *Provider {
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
	s.store.addSMS(stored)

	return &contracts.SendResult{
		ID:         id,
		StatusCode: 200,
		Message:    "Stored SMS in memory",
	}, nil
}
