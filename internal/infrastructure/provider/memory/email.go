package memory

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/weprodev/wpd-message-gateway/pkg/contracts"
)

// EmailProvider implements port.EmailSender using an in-memory store.
type EmailProvider struct {
	store       *Store
	workspaceID string
}

// NewEmailProvider creates a new memory email provider (no workspace scope; legacy / global inbox).
func NewEmailProvider(store *Store) *EmailProvider {
	return NewEmailProviderForWorkspace(store, "")
}

// NewEmailProviderForWorkspace tags stored emails with workspaceID for multi-tenant portal inbox.
func NewEmailProviderForWorkspace(store *Store, workspaceID string) *EmailProvider {
	return &EmailProvider{
		store:       store,
		workspaceID: workspaceID,
	}
}

// Name returns the provider name.
func (e *EmailProvider) Name() string {
	return ProviderName
}

// Send stores the email in memory and optionally forwards to Mailpit.
func (e *EmailProvider) Send(ctx context.Context, email *contracts.Email) (*contracts.SendResult, error) {
	id := uuid.New().String()

	stored := &StoredEmail{
		ID:          id,
		WorkspaceID: e.workspaceID,
		CreatedAt:   time.Now(),
		Email:       email,
	}
	e.store.AddEmail(stored)

	return &contracts.SendResult{
		ID:         id,
		StatusCode: 200,
		Message:    "Stored email in memory",
	}, nil
}
