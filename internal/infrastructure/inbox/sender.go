package inbox

import (
	"context"

	"github.com/weprodev/wpd-message-gateway/internal/core/port"
	"github.com/weprodev/wpd-message-gateway/pkg/contracts"
)

// Compile-time assertion
var _ contracts.EmailSender = (*InboxEmailSender)(nil)

// InboxEmailSender wraps an underlying EmailSender and writes any sent email
// to the in-process capture store before sending.
type InboxEmailSender struct {
	writer port.InboxWriter
	sender contracts.EmailSender
}

// NewInboxEmailSender creates a new InboxEmailSender.
func NewInboxEmailSender(writer port.InboxWriter, sender contracts.EmailSender) *InboxEmailSender {
	return &InboxEmailSender{
		writer: writer,
		sender: sender,
	}
}

// Send captures the email to the in-process store, then forwards to the wrapped sender.
func (s *InboxEmailSender) Send(ctx context.Context, email contracts.Email) (*contracts.SendResult, error) {
	if s.writer != nil {
		_, _ = s.writer.WriteEmail(ctx, "", email)
	}
	return s.sender.Send(ctx, email)
}

// Name returns the name of the underlying sender.
func (s *InboxEmailSender) Name() string {
	return s.sender.Name()
}
