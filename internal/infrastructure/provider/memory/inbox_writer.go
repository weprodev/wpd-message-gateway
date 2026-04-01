package memory

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/weprodev/wpd-message-gateway/internal/core/port"
	"github.com/weprodev/wpd-message-gateway/pkg/contracts"
)

// Compile-time assertion: InboxWriterAdapter implements port.InboxWriter.
var _ port.InboxWriter = (*InboxWriterAdapter)(nil)

// InboxWriterAdapter wraps the in-process Store and optional Mailpit forwarder
// to implement port.InboxWriter. Construct via NewInboxWriter.
type InboxWriterAdapter struct {
	store   *Store
	mailpit *smtpForwarder
}

// NewInboxWriter returns a port.InboxWriter backed by the given Store.
// If mailpitCfg.Enabled is true, emails are also forwarded to Mailpit via SMTP.
func NewInboxWriter(store *Store, mailpitCfg MailpitConfig) port.InboxWriter {
	return &InboxWriterAdapter{
		store:   store,
		mailpit: newSMTPForwarder(mailpitCfg),
	}
}

// WriteEmail stores an email in memory and optionally forwards to Mailpit.
func (a *InboxWriterAdapter) WriteEmail(_ context.Context, workspaceID string, email *contracts.Email) (string, error) {
	id := uuid.New().String()
	a.store.AddEmail(&StoredEmail{
		ID:          id,
		WorkspaceID: workspaceID,
		CreatedAt:   time.Now(),
		Email:       email,
	})
	a.mailpit.forward(email)
	return id, nil
}

// WriteSMS stores an SMS in memory.
func (a *InboxWriterAdapter) WriteSMS(_ context.Context, _ string, sms *contracts.SMS) (string, error) {
	id := uuid.New().String()
	a.store.AddSMS(&StoredSMS{
		ID:        id,
		CreatedAt: time.Now(),
		SMS:       sms,
	})
	return id, nil
}

// WritePush stores a push notification in memory.
func (a *InboxWriterAdapter) WritePush(_ context.Context, _ string, push *contracts.PushNotification) (string, error) {
	id := uuid.New().String()
	a.store.AddPush(&StoredPush{
		ID:        id,
		CreatedAt: time.Now(),
		Push:      push,
	})
	return id, nil
}

// WriteChat stores a chat message in memory.
func (a *InboxWriterAdapter) WriteChat(_ context.Context, _ string, chat *contracts.ChatMessage) (string, error) {
	id := uuid.New().String()
	a.store.AddChat(&StoredChat{
		ID:        id,
		CreatedAt: time.Now(),
		Chat:      chat,
	})
	return id, nil
}
