package memory

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/weprodev/wpd-message-gateway/pkg/contracts"
)

// EmailProvider implements contracts.EmailSender using an in-memory store.
type EmailProvider struct {
	store *Store
}

// NewEmailProvider creates a new memory email provider.
func NewEmailProvider(store *Store) *EmailProvider {
	return &EmailProvider{
		store: store,
	}
}

// Name returns the provider name.
func (e *EmailProvider) Name() string {
	return ProviderName
}

// Send stores the email in memory.
func (e *EmailProvider) Send(ctx context.Context, email contracts.Email) (*contracts.SendResult, error) {
	id := uuid.New().String()

	stored := SentEmail{
		ID:        id,
		CreatedAt: time.Now(),
		Email:     email,
	}
	e.store.AddEmail(stored)

	return &contracts.SendResult{
		ID:         id,
		StatusCode: 200,
		Message:    "Stored email in memory",
	}, nil
}
